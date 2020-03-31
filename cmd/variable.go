package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcl"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// variablesCmd represents the variables command
var variableCmd = &cobra.Command{
	Use:   "variable",
	Short: "Manage TFE variables",
	Long:  `Manage TFE variables.`,
}

var variableCreateCmd = &cobra.Command{
	Use:   "create [WORKSPACE]",
	Short: "Create TFE variables",
	Long:  `Create TFE variables.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read the flags.
		name := args[0]
		vars, _ := cmd.Flags().GetStringArray("var")
		svars, _ := cmd.Flags().GetStringArray("svar")
		HCLvars, _ := cmd.Flags().GetStringArray("hvar")
		sHCLVars, _ := cmd.Flags().GetStringArray("shvar")
		EnvVars, _ := cmd.Flags().GetStringArray("evar")
		sEnvVars, _ := cmd.Flags().GetStringArray("sevar")
		force, _ := cmd.Flags().GetBool("force")
		varFile, _ := cmd.Flags().GetString("var-file")

		// Setup the command.
		organization, client, err := setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// Retrieve the workspace.
		workspace, err := readWorkspace(client, organization, name)
		if err != nil {
			log.Fatalf("Cannot retrieve workspace %q: %s.", name, err)
		}

		// Prepare the variables.
		varOptions := []tfe.VariableCreateOptions{}
		varOptions = append(varOptions, createVariableOptions(vars, tfe.CategoryTerraform, false, false)...)
		varOptions = append(varOptions, createVariableOptions(svars, tfe.CategoryTerraform, false, true)...)
		varOptions = append(varOptions, createVariableOptions(HCLvars, tfe.CategoryTerraform, true, false)...)
		varOptions = append(varOptions, createVariableOptions(sHCLVars, tfe.CategoryTerraform, true, true)...)
		varOptions = append(varOptions, createVariableOptions(EnvVars, tfe.CategoryEnv, false, false)...)
		varOptions = append(varOptions, createVariableOptions(sEnvVars, tfe.CategoryEnv, false, true)...)

		// Read variables from file.
		if varFile != "" {
			// Read the content of the file.
			FileContent, err := ioutil.ReadFile(varFile)
			if err != nil {
				log.Fatalf("Cannot read the file %q: %s.", varFile, err)
			}

			// Parse it.
			var v interface{}
			err = hcl.Unmarshal(FileContent, &v)
			if err != nil {
				log.Fatalf("Cannot parse the file %q: %s.", varFile, err)
			}

			// Convert the content to `key=value` format.
			regularVarFile := []string{}
			HCLVarFile := []string{}
			s := reflect.ValueOf(v)
			if s.Kind() == reflect.Map {
				for _, key := range s.MapKeys() {
					strct := s.MapIndex(key)
					k := fmt.Sprintf("%s", key.Interface())
					value := reflect.ValueOf(strct.Interface())
					// If the type is Slice, we consider it HCL.
					if value.Kind() == reflect.Slice {
						// Use reflection to extract the values of the slice.
						b := make([]string, value.Len())
						for i := 0; i < value.Len(); i++ {
							b[i] = fmt.Sprintf("%s", value.Index(i))
						}
						HCLVarFile = append(HCLVarFile, fmt.Sprintf("%s=[\"%s\"]", k, strings.Join(b, "\", \"")))
					} else {
						// Otherwise it is always a regular variable.
						regularVarFile = append(regularVarFile, fmt.Sprintf("%s=%s", k, strct.Interface()))
					}
				}
			}

			// Add it to the list of variables to create.
			varOptions = append(varOptions, createVariableOptions(regularVarFile, tfe.CategoryTerraform, false, false)...)
			varOptions = append(varOptions, createVariableOptions(HCLVarFile, tfe.CategoryTerraform, true, false)...)
		}

		// List existing variables and index them by key.
		indexedVars, err := indexVariables(client, workspace.ID, organization)
		if err != nil {
			log.Fatalf("Cannot index variables: %s.", err)
		}

		// Go through all the variables.
		for _, opts := range varOptions {
			//Check if the variable already exists.
			v, exists := indexedVars[*(opts.Key)]

			// If the variable exists.
			if exists {
				// Update it.
				if force {
					options := tfe.VariableUpdateOptions{
						Key:       opts.Key,
						Value:     opts.Value,
						HCL:       opts.HCL,
						Sensitive: opts.Sensitive,
					}
					if _, err := client.Variables.Update(context.Background(), workspace.ID, v.ID, options); err != nil {
						log.Fatalf("Cannot update variable %q: %s.", *(opts.Key), err)
					}
					log.Debugf("Variable %q updated successfully.", *(opts.Key))
					continue
				}

				// Raise an error..
				log.Fatalf("Cannot create %q: variable already exists.", *(opts.Key))
			}

			// Otherwise create it.
			if _, err = createVariable(client, workspace.ID, opts); err != nil {
				log.Fatalf("Cannot create variable %q: %s.", *(opts.Key), err)
			}
			log.Infof("Variable %q created successfully.", *(opts.Key))
		}
	},
}

var variableListCmd = &cobra.Command{
	Use:   "list [WORKSPACE]",
	Short: "List TFE variables for a specific workspace",
	Long:  `List TFE variables for a specific workspace.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read the flags.
		name := args[0]

		// Setup the command.
		organization, client, err := setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// Retrieve the workspace exist.
		workspace, err := readWorkspace(client, organization, name)
		if err != nil {
			log.Fatalf("Cannot retrieve workspace %q: %s.", name, err)
		}

		// List variables.
		variables, err := listVariables(client, workspace.ID)
		if err != nil {
			log.Fatalf("Cannot list the variables for  %q: %s.", organization, err)
		}

		// Print the variables.
		for _, variable := range variables {
			fmt.Printf("%s=%s\n", variable.Key, variable.Value)
		}

	},
}

var variableDeleteCmd = &cobra.Command{
	Use:   "delete [WORKSPACE] [VARIABLE]",
	Short: "Delete a TFE variable for a specific workspace",
	Long:  `Delete a TFE variable for a specific workspace.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Read the flags.
		wsName := args[0]
		varName := args[1]

		// Setup the command.
		organization, client, err := setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// Retrieve the workspace.
		workspace, err := readWorkspace(client, organization, wsName)
		if err != nil {
			log.Fatalf("Cannot retrieve workspace %q: %s.", wsName, err)
		}

		// List existing variables and index them by key.
		indexedVars, err := indexVariables(client, workspace.ID, organization)
		if err != nil {
			log.Fatalf("Cannot index variables: %s.", err)
		}

		// Check if it exists.
		v, exists := indexedVars[varName]
		if !exists {
			log.Warningf("Cannot delete variable %q: it does not exist.", varName)
			return
		}

		// And delete it if it does.
		if err := deleteVariable(client, workspace.ID, v.ID); err != nil {
			log.Fatalf("Cannot delete variable %q: %s.", varName, err)
		}
		log.Infof("Variable %q deleted successfully.", varName)
	},
}

func init() {
	rootCmd.AddCommand(variableCmd)
	variableCmd.AddCommand(variableCreateCmd)
	variableCmd.AddCommand(variableDeleteCmd)
	variableCmd.AddCommand(variableListCmd)

	variableCreateCmd.Flags().StringArray("var", []string{}, "Create a regular variable")
	variableCreateCmd.Flags().StringArray("svar", []string{}, "Create a regular sensitive variable")
	variableCreateCmd.Flags().StringArray("hvar", []string{}, "Create an HCL variable")
	variableCreateCmd.Flags().StringArray("shvar", []string{}, "Create a sensitive HCL variable")
	variableCreateCmd.Flags().StringArray("evar", []string{}, "Create an environment variable")
	variableCreateCmd.Flags().StringArray("sevar", []string{}, "Create a sensitive environment variable")
	variableCreateCmd.Flags().String("var-file", "", "Create non-sensitive regular and HCL variables from a file")
	variableCreateCmd.Flags().BoolP("force", "f", false, "Overwrite a variable if it exists")
}

func createVariable(client *tfe.Client, workspaceID string, options tfe.VariableCreateOptions) (*tfe.Variable, error) {
	v, err := client.Variables.Create(context.Background(), workspaceID, options)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func listVariables(client *tfe.Client, workspaceID string) ([]*tfe.Variable, error) {
	results := []*tfe.Variable{}
	currentPage := 1

	// Go through the pages of results until there is no more pages.
	for {
		log.Debugf("Processing page %d.\n", currentPage)
		options := tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: currentPage,
			}}
		v, err := client.Variables.List(context.Background(), workspaceID, options)
		if err != nil {
			return nil, err
		}
		results = append(results, v.Items...)

		// Check if there is another poage to retrieve.
		if v.Pagination.NextPage == 0 {
			break
		}

		// Incremment the page number.
		currentPage++
	}

	return results, nil
}

func createVariableOptions(vars []string, category tfe.CategoryType, hcl, sensitive bool) []tfe.VariableCreateOptions {
	optionList := []tfe.VariableCreateOptions{}

	for _, v := range vars {
		splitV := strings.Split(v, "=")
		if len(splitV) != 2 {
			log.Fatalf("Invalid variable %q: the format must be key=value.", v)
		}

		options := tfe.VariableCreateOptions{
			Key:       tfe.String(splitV[0]),
			Value:     tfe.String(splitV[1]),
			Category:  tfe.Category(category),
			HCL:       tfe.Bool(hcl),
			Sensitive: tfe.Bool(sensitive),
		}

		optionList = append(optionList, options)
	}
	return optionList
}

func indexVariables(client *tfe.Client, workspaceID, organization string) (map[string]*tfe.Variable, error) {
	// List existing variables.
	variables, err := listVariables(client, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("cannot list the variables for  %q: %s", organization, err)
	}

	// Index them by key.
	indexedVars := map[string]*tfe.Variable{}
	for _, v := range variables {
		indexedVars[v.Key] = v
	}

	return indexedVars, nil
}

func deleteVariable(client *tfe.Client, workspaceID, variableID string) error {
	return client.Variables.Delete(context.Background(), workspaceID, variableID)

}

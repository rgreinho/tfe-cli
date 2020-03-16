package cmd

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
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
		HCLvars, _ := cmd.Flags().GetStringArray("hcl-var")
		sHCLVars, _ := cmd.Flags().GetStringArray("shcl-var")
		EnvVars, _ := cmd.Flags().GetStringArray("svar")
		sEnvVars, _ := cmd.Flags().GetStringArray("svar")
		force, _ := cmd.Flags().GetBool("force")

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

		indexedVars := map[string]*tfe.Variable{}
		if force {
			// List variables.
			variables, err := listVariables(client, workspace.ID)
			if err != nil {
				log.Fatalf("Cannot list the variables for  %q: %s.", organization, err)
			}

			// Index them by key.
			for _, v := range variables {
				indexedVars[v.Key] = v
			}
		}

		// Go through all the variables.
		for _, opts := range varOptions {

			// Update it if needed.
			if force {
				// Ensure the variable exists.
				if v, ok := indexedVars[*(opts.Key)]; ok {
					// If so, update it.
					options := tfe.VariableUpdateOptions{
						Key:       opts.Key,
						Value:     opts.Value,
						HCL:       opts.HCL,
						Sensitive: opts.Sensitive,
					}
					if _, err := client.Variables.Update(context.Background(), workspace.ID, v.ID, options); err != nil {
						log.Fatalf("Cannot update variable %q: %s.", v.Key, err)
					}
					break
				}
			}

			// Otherwise create it.
			if _, err = createVariable(client, workspace.ID, opts); err != nil {
				log.Fatalf("Cannot create variable %q: %s.", *(opts.Key), err)
			}
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

func init() {
	rootCmd.AddCommand(variableCmd)
	variableCmd.AddCommand(variableCreateCmd)
	variableCmd.AddCommand(variableListCmd)

	variableCreateCmd.Flags().StringArray("var", []string{}, "Create a regular variable")
	variableCreateCmd.Flags().StringArray("svar", []string{}, "Create a regular sensitive variable")
	variableCreateCmd.Flags().StringArray("hcl-var", []string{}, "Create an HCL variable")
	variableCreateCmd.Flags().StringArray("shcl-var", []string{}, "Create a sensitive HCL variable")
	variableCreateCmd.Flags().StringArray("env-var", []string{}, "Create an environment variable")
	variableCreateCmd.Flags().StringArray("senv-var", []string{}, "Create a sensitive environment variable")
	// variableCreateCmd.Flags().String("var-file", "", "Create HCL non-sensitive variables from a file")
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

package cmd

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/rgreinho/tfe-cli/tfecli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace command.
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage TFE workspaces",
	Long:  `Manage TFE workspaces.`,
}

// createCmd represents the create command.
var workspaceCreateCmd = &cobra.Command{
	Use:   "create [WORKSPACE]",
	Short: "Create a TFE workspace",
	Long:  `Create a TFE workspace.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the command.
		organization, client, err := tfecli.Setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// Read the flags.
		name := args[0]
		autoapply, _ := cmd.Flags().GetBool("autoapply")
		filetriggers, _ := cmd.Flags().GetBool("filetriggers")
		terraformversion, _ := cmd.Flags().GetString("terraformversion")
		force, _ := cmd.Flags().GetBool("force")
		workingdirectory, _ := cmd.Flags().GetString("workingdirectory")
		vcsrepository, _ := cmd.Flags().GetString("vcsrepository")
		splitVCS := strings.Split(vcsrepository, ":")

		// Check whether the workspace exists.
		w, err := readWorkspace(client, organization, name)
		if err != nil {
			if !strings.Contains(err.Error(), "resource not found") {
				log.Fatalf("Cannot retrieve workspace %q: %s.", name, err)
			}
		}

		if w != nil {
			// Update the workspace if needed.
			if force {
				// Prepare the new workspace options.
				options := tfe.WorkspaceUpdateOptions{
					AutoApply:           tfe.Bool(autoapply),
					FileTriggersEnabled: tfe.Bool(filetriggers),
					Name:                tfe.String(name),
					TerraformVersion:    tfe.String(terraformversion),
					WorkingDirectory:    tfe.String(workingdirectory),
				}
				if len(splitVCS) == 3 {
					options.VCSRepo = &tfe.VCSRepoOptions{
						Branch:       tfe.String(splitVCS[2]),
						Identifier:   tfe.String(splitVCS[1]),
						OAuthTokenID: tfe.String(splitVCS[0]),
					}
				}

				if _, err := updateWorkspaceByID(client, w.ID, options); err != nil {
					log.Fatalf("Cannot update workspace %q: %s.", name, err)
				}

				log.Infof("Workspace %q updated successfully.", name)
				return
			}

			// Otherwise do nothing.
			log.Infof("Workspace %q already exists.", name)
			return
		}

		// Prepare the new workspace options.
		options := tfe.WorkspaceCreateOptions{
			AutoApply:           tfe.Bool(autoapply),
			FileTriggersEnabled: tfe.Bool(filetriggers),
			Name:                tfe.String(name),
			TerraformVersion:    tfe.String(terraformversion),
			WorkingDirectory:    tfe.String(workingdirectory),
		}
		if len(splitVCS) == 3 {
			options.VCSRepo = &tfe.VCSRepoOptions{
				Branch:       tfe.String(splitVCS[2]),
				Identifier:   tfe.String(splitVCS[1]),
				OAuthTokenID: tfe.String(splitVCS[0]),
			}
		}

		// Create the workspace.
		if _, err = createWorkspace(client, organization, options); err != nil {
			log.Fatalf("Cannot create workspace %q: %s.", name, err)
		}
	},
}

// workspaceDeleteCmd represents the delete command.
var workspaceDeleteCmd = &cobra.Command{
	Use:   "delete [WORKSPACE]",
	Short: "Delete a TFE workspace",
	Long:  `Delete a TFE workspace.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the command.
		organization, client, err := tfecli.Setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// Read the flags.
		name := args[0]

		// Delete the workspace.
		if err := deleteWorkspace(client, organization, name); err != nil {
			log.Fatalf("Cannot delete workspace %q: %s.", name, err)
		}

		log.Infof("Workspace %q deleted successfully.", name)
	},
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List TFE workspaces",
	Long:  `List TFE workspaces.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the command.
		organization, client, err := tfecli.Setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// List workspaces.
		workspaces, err := listWorkspaces(client, organization)
		if err != nil {
			log.Fatalf("Cannot list the workspaces for  %q: %s.", organization, err)
		}

		// Print the workspace names.
		for _, workspace := range workspaces {
			fmt.Println(workspace.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceCreateCmd)
	workspaceCmd.AddCommand(workspaceDeleteCmd)
	workspaceCmd.AddCommand(workspaceListCmd)

	workspaceCreateCmd.Flags().Bool("autoapply", false, "Apply changes automatically")
	workspaceCreateCmd.Flags().Bool("filetriggers", false, "Filter runs based on the changed files in a VCS push")
	workspaceCreateCmd.Flags().String("terraformversion", "", "Specify the Terraform version")
	workspaceCreateCmd.Flags().String("workingdirectory", "", "Specify a relative path that Terraform will execute within")
	// colon sperated values: <OAuthTokenID>:<repository>:<branch>
	// example: ot-8Xc1NTYpjIQZIwIh:organization/repository:master
	workspaceCreateCmd.Flags().String("vcsrepository", "", "Specify a workspace's VCS repository")
	workspaceCreateCmd.Flags().BoolP("force", "f", false, "Update workspace if it exists")
}

func readWorkspace(client *tfe.Client, organization, workspace string) (*tfe.Workspace, error) {
	w, err := client.Workspaces.Read(context.Background(), organization, workspace)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func createWorkspace(client *tfe.Client, organization string, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	w, err := client.Workspaces.Create(context.Background(), organization, options)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func updateWorkspaceByID(client *tfe.Client, workspaceID string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	w, err := client.Workspaces.UpdateByID(context.Background(), workspaceID, options)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func listWorkspaces(client *tfe.Client, organization string) ([]*tfe.Workspace, error) {
	results := []*tfe.Workspace{}
	currentPage := 1

	// Go through the pages of results until there is no more pages.
	for {
		log.Debugf("Processing page %d.\n", currentPage)
		options := tfe.WorkspaceListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: currentPage,
			}}
		w, err := client.Workspaces.List(context.Background(), organization, options)
		if err != nil {
			return nil, err
		}
		results = append(results, w.Items...)

		// Check if there is another poage to retrieve.
		if w.Pagination.NextPage == 0 {
			break
		}

		// Incremment the page number.
		currentPage++
	}

	return results, nil
}

func deleteWorkspace(client *tfe.Client, organization, workspace string) error {
	if err := client.Workspaces.Delete(context.Background(), organization, workspace); err != nil {
		return err
	}
	return nil
}

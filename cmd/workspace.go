package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage TFE workspaces",
	Long:  `Manage TFE workspaces.`,
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [WORKSPACE]",
	Short: "Create a TFE workspace",
	Long:  `Create a TFE workspace.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Get organization.
		organization, err := getOrganization(cmd)
		if err != nil {
			log.Fatalf("No organization specified.")
		}

		// Get create command options.
		name := args[0]
		autoapply, _ := cmd.Flags().GetBool("autoapply")
		filetriggers, _ := cmd.Flags().GetBool("filetriggers")
		terraformversion, _ := cmd.Flags().GetString("terraformversion")
		force, _ := cmd.Flags().GetBool("force")
		workingdirectory, _ := cmd.Flags().GetString("workingdirectory")
		vcsrepository, _ := cmd.Flags().GetString("vcsrepository")
		splitVCS := strings.Split(vcsrepository, ":")

		// Create the TFE client.
		client, err := newClient(cmd)
		if err != nil {
			log.Fatalf("Cannot create TFE client: %s.", err)
		}

		// Check whether the workspace exists.
		w, err := readWorkspace(client, organization, name)
		if err != nil {
			log.Fatalf("Cannot retrieve workspace %q: %s.", name, err)
		}

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
		if w != nil {
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List TFE workspaces",
	Long:  `List TFE workspaces.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get organization.
		organization, err := getOrganization(cmd)
		if err != nil {
			log.Fatalf("No organization specified.")
		}

		// Create the TFE client.
		client, err := newClient(cmd)
		if err != nil {
			log.Fatalf("Cannot create TFE client: %s.", err)
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
	workspaceCmd.AddCommand(createCmd)
	workspaceCmd.AddCommand(listCmd)

	createCmd.Flags().Bool("autoapply", false, "Apply changes automatically")
	createCmd.Flags().Bool("filetriggers", false, "Filter runs based on the changed files in a VCS push")
	createCmd.Flags().String("terraformversion", "", "Specify the Terraform version")
	createCmd.Flags().String("workingdirectory", "", "Specify a relative path that Terraform will execute within")
	// colon sperated values: <OAuthTokenID>:<repository>:<branch>
	// example: ot-8Xc1NTYpjIQZIwIh:shipstation/mercury-appstack:master
	createCmd.Flags().String("vcsrepository", "", "Specify a workspace's VCS repository")
	createCmd.Flags().BoolP("force", "f", false, "Update workspace if it exists")
}

func newClient(cmd *cobra.Command) (*tfe.Client, error) {
	// Get token.
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		// Read the environment variable as a fallback.
		token = os.Getenv("TFE_TOKEN")
	}
	if token == "" {
		return nil, fmt.Errorf("no token specified")
	}

	// Prepare TFE config.
	config := &tfe.Config{
		Token: token,
	}

	// Create TFE client.
	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client, err
}

func getOrganization(cmd *cobra.Command) (string, error) {
	// Get organization from CLI flag.
	organization, _ := cmd.Flags().GetString("organization")
	if organization == "" {
		// Read the environment variable as a fallback.
		organization = os.Getenv("TFE_ORG")
	}
	if organization == "" {
		return "", fmt.Errorf("no organization specified")
	}
	return organization, nil
}

func readWorkspace(client *tfe.Client, organization string, workspace string) (*tfe.Workspace, error) {
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

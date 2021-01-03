package cmd

import (
	"context"
	"fmt"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/rgreinho/tfe-cli/tfecli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// notificationCmd represents the notifications command
var notificationCmd = &cobra.Command{
	Use:   "notification",
	Short: "Manage TFE notifications",
	Long:  `Manage TFE notifications.`,
}

func init() {
	rootCmd.AddCommand(notificationCmd)
	notificationCmd.AddCommand(notificationListCmd)
	notificationCmd.AddCommand(notificationCreateCmd)
	notificationCmd.AddCommand(notificationDeleteCmd)

	notificationCreateCmd.Flags().String("type", "", "Specify the destination type")
	notificationCreateCmd.Flags().Bool("disabled", false, "Disable the notification")
	notificationCreateCmd.Flags().String("token", "", "Specify the notification token")
	notificationCreateCmd.Flags().StringArray("triggers", []string{}, "Specify the list of run events that will trigger notifications")
	notificationCreateCmd.Flags().String("url", "", "Specify the notification url")
	notificationCreateCmd.Flags().StringArray("emailaddresses", []string{}, "Specify the list of email addresses that will receive notification emails")
	notificationCreateCmd.Flags().BoolP("force", "f", false, "Update notification if it exists")
}

var notificationListCmd = &cobra.Command{
	Use:   "list [WORKSPACE]",
	Short: "List TFE notifications for a specific workspace",
	Long:  `List TFE notifications for a specific workspace.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read the flags.
		name := args[0]

		// Setup the command.
		organization, client, err := tfecli.Setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// Retrieve the workspace.
		workspace, err := readWorkspace(client, organization, name)
		if err != nil {
			log.Fatalf("Cannot retrieve workspace %q: %s.", name, err)
		}

		// List workspace notifications.
		notifications, err := listNotifications(client, workspace.ID)
		if err != nil {
			log.Fatalf("Cannot list the notifications for  %q: %s.", organization, err)
		}

		// Print the variables.
		for _, notification := range notifications {
			fmt.Printf("%s: %s\n", notification.Name, notification.DestinationType)
		}
	},
}

var notificationCreateCmd = &cobra.Command{
	Use:   "create [WORKSPACE] [NOTIFICATION_NAME]",
	Short: "Create TFE notification for a specific workspace",
	Long:  `Create TFE notification for a specific workspace.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Read the flags.
		workspaceName := args[0]
		notificationName := args[1]
		destinationType, _ := cmd.Flags().GetString("type")
		disabled, _ := cmd.Flags().GetBool("disabled")
		token, _ := cmd.Flags().GetString("token")
		triggers, _ := cmd.Flags().GetStringArray("triggers")
		url, _ := cmd.Flags().GetString("url")
		emailAddresses, _ := cmd.Flags().GetStringArray("emailaddresses")
		force, _ := cmd.Flags().GetBool("force")

		// Setup the command.
		organization, client, err := tfecli.Setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// Retrieve the workspace.
		workspace, err := readWorkspace(client, organization, workspaceName)
		if err != nil {
			log.Fatalf("Cannot retrieve workspace %q: %s.", workspaceName, err)
		}

		// List existing variables and index them by key.
		indexedNotifications, err := indexNotifications(client, workspace.ID, organization)
		if err != nil {
			log.Fatalf("Cannot index notifications: %s.", err)
		}

		// Check if it exists.
		notification, exists := indexedNotifications[notificationName]

		if exists {
			if force {
				// Create the notification options.
				options := tfe.NotificationConfigurationUpdateOptions{
					Enabled:        tfe.Bool(!disabled),
					Name:           tfe.String(notificationName),
					Token:          tfe.String(token),
					Triggers:       triggers,
					URL:            tfe.String(url),
					EmailAddresses: emailAddresses,
				}

				// Update
				if _, err := updateNotification(client, notification.ID, options); err != nil {
					log.Fatalf("Cannot update notification %q: %s.", notificationName, err)
				}
				log.Infof("Notification %q of type %q updated.", notificationName, destinationType)
			} else {
				log.Fatalf("cannot create %q: notification already exists", notificationName)
			}
		} else {
			// Create the notification options.
			options := tfe.NotificationConfigurationCreateOptions{
				DestinationType: (*tfe.NotificationDestinationType)(&destinationType),
				Enabled:         tfe.Bool(!disabled),
				Name:            tfe.String(notificationName),
				Token:           tfe.String(token),
				Triggers:        triggers,
				URL:             tfe.String(url),
				EmailAddresses:  emailAddresses,
			}

			// Create the notification.
			if _, err = createNotification(client, workspace.ID, options); err != nil {
				log.Fatalf("Cannot create notification %q: %s.", notificationName, err)
			}
			log.Infof("Notification %q of type %q created.", notificationName, destinationType)

		}
	},
}

var notificationDeleteCmd = &cobra.Command{
	Use:   "delete [WORKSPACE] [NOTIFICATION_NAME]",
	Short: "Delete a TFE notification",
	Long:  `Delete a TFE notification.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Setup the command.
		organization, client, err := tfecli.Setup(cmd)
		if err != nil {
			log.Fatalf("Cannot execute the command: %s.", err)
		}

		// Read the flags.
		workspaceName := args[0]
		notificationName := args[1]

		// Retrieve the workspace.
		workspace, err := readWorkspace(client, organization, workspaceName)
		if err != nil {
			log.Fatalf("Cannot retrieve workspace %q: %s.", workspaceName, err)
		}

		// List existing variables and index them by key.
		indexedNotifications, err := indexNotifications(client, workspace.ID, organization)
		if err != nil {
			log.Fatalf("Cannot index notifications: %s.", err)
		}

		// Check if it exists.
		notification, exists := indexedNotifications[notificationName]

		if !exists {
			log.Warningf("Cannot delete notification %q: it does not exist.", notificationName)
			return
		}

		// Delete the workspace.
		if err := deleteNotification(client, notification.ID); err != nil {
			log.Fatalf("Cannot delete notification %q: %s.", notificationName, err)
		}

		log.Infof("Notification %q deleted successfully.", notificationName)
	},
}

func listNotifications(client *tfe.Client, workspaceID string) ([]*tfe.NotificationConfiguration, error) {
	results := []*tfe.NotificationConfiguration{}
	currentPage := 1

	// Go through the pages of results until there is no more pages.
	for {
		log.Debugf("Processing page %d.\n", currentPage)
		options := tfe.NotificationConfigurationListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: currentPage,
			}}
		v, err := client.NotificationConfigurations.List(context.Background(), workspaceID, options)
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

func createNotification(client *tfe.Client, workspaceID string, options tfe.NotificationConfigurationCreateOptions) (*tfe.NotificationConfiguration, error) {
	v, err := client.NotificationConfigurations.Create(context.Background(), workspaceID, options)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func indexNotifications(client *tfe.Client, workspaceID, organization string) (map[string]*tfe.NotificationConfiguration, error) {
	// List existing variables.
	notifications, err := listNotifications(client, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("cannot list the notifications for  %q: %s", organization, err)
	}

	// Index them by key.
	indexedNotifications := map[string]*tfe.NotificationConfiguration{}
	for _, n := range notifications {
		indexedNotifications[n.Name] = n
	}

	return indexedNotifications, nil
}

func updateNotification(client *tfe.Client, notificationID string, options tfe.NotificationConfigurationUpdateOptions) (*tfe.NotificationConfiguration, error) {
	n, err := client.NotificationConfigurations.Update(context.Background(), notificationID, options)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func deleteNotification(client *tfe.Client, notificationID string) error {
	if err := client.NotificationConfigurations.Delete(context.Background(), notificationID); err != nil {
		return err
	}
	return nil
}

package cmd

import (
	"fmt"
	"os"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

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

func getToken(cmd *cobra.Command) (string, error) {
	// Get token.
	token, _ := cmd.Flags().GetString("token")
	if token == "" {
		// Read the environment variable as a fallback.
		token = os.Getenv("TFE_TOKEN")
	}
	if token == "" {
		return "", fmt.Errorf("no token specified")
	}
	return token, nil
}

func newClient(token string) (*tfe.Client, error) {
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

func setup(cmd *cobra.Command) (organization string, client *tfe.Client, err error) {
	// Get organization.
	organization, err = getOrganization(cmd)
	if err != nil {
		fmt.Errorf("no organization specified")
	}

	// Get token.
	token, err := getToken(cmd)
	if err != nil {
		fmt.Errorf("no organization specified")
	}

	// Create the TFE client.
	client, err = newClient(token)
	if err != nil {
		fmt.Errorf("cannot create TFE client: %s", err)
	}

	return
}

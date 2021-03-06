package service

import (
	"fmt"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func newCreateCommand(dockerCli *client.DockerCli) *cobra.Command {
	opts := newServiceOptions()

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] IMAGE [COMMAND] [ARG...]",
		Short: "Create a new service",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.image = args[0]
			if len(args) > 1 {
				opts.args = args[1:]
			}
			return runCreate(dockerCli, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.mode, flagMode, "replicated", "Service mode (replicated or global)")
	addServiceFlags(cmd, opts)
	cmd.Flags().SetInterspersed(false)
	return cmd
}

func runCreate(dockerCli *client.DockerCli, opts *serviceOptions) error {
	apiClient := dockerCli.Client()
	headers := map[string][]string{}

	service, err := opts.ToService()
	if err != nil {
		return err
	}

	ctx := context.Background()

	// only send auth if flag was set
	if opts.registryAuth {
		// Retrieve encoded auth token from the image reference
		encodedAuth, err := dockerCli.RetrieveAuthTokenFromImage(ctx, opts.image)
		if err != nil {
			return err
		}
		headers["X-Registry-Auth"] = []string{encodedAuth}
	}

	response, err := apiClient.ServiceCreate(ctx, service, headers)
	if err != nil {
		return err
	}

	fmt.Fprintf(dockerCli.Out(), "%s\n", response.ID)
	return nil
}

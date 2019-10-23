package cmd

import (
	"log"

	"github.com/docker/mayday/pkg/mayday"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// MaydayConfig -
type MaydayConfig struct {
	Host string
}

// NewRootCommand -
func NewRootCommand() *cobra.Command {
	clientProvider := mayday.NewClientProvider()
	config := MaydayConfig{}

	command := &cobra.Command{
		Use:   "mayday",
		Short: "track various metadata",
		Long:  `track various metadata`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(config.Host, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}

			clientProvider.Set(mayday.NewClient(conn))
		},
	}

	command.PersistentFlags().StringVar(&config.Host, "host", "mayday:8050", "host of mayday server")
	command.AddCommand(NewTypesCommand(config, clientProvider))
	command.AddCommand(NewObservationsCommand(config, clientProvider))

	return command
}

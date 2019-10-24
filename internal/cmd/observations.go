package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/docker/mayday/pkg/mayday"
	"github.com/spf13/cobra"
)

// NewObservationsCommand -
func NewObservationsCommand(config MaydayConfig, clientProvider mayday.ClientProvider) *cobra.Command {
	typesCommand := &cobra.Command{
		Use:   "observations [string to echo]",
		Short: "used to created and list observations",
		Long:  `does things with observations`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("observations: " + strings.Join(args, " "))
		},
	}

	typesCommand.AddCommand(newObservationsCreateCommand(config, clientProvider))
	typesCommand.AddCommand(newObservationListCommand(config, clientProvider))

	return typesCommand
}

func newObservationsCreateCommand(config MaydayConfig, clientProvider mayday.ClientProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "create [thing to create]",
		Short: "create observation",
		Long:  `create observation`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			client := clientProvider.Get()

			response, err := client.CreateObservation(context.TODO(), []byte(args[0]))
			if err != nil {
				println(err.Error())
				return
			}

			t := response.GetObservation()

			w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "ID\tName\tPayload\tCreated\tUpdated")
			fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\n", t.GetId().GetValue(), t.GetName(), t.GetPayload(), t.GetCreated(), t.GetUpdated())
			w.Flush()
		},
	}
}

func newObservationListCommand(config MaydayConfig, clientProvider mayday.ClientProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list observations",
		Long:  `list observations`,
		Run: func(c *cobra.Command, args []string) {
			client := clientProvider.Get()

			response, err := client.ListObservations(context.TODO())
			if err != nil {
				println(err.Error())
				return
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "ID\tName\tPayload\tCreated\tUpdated")
			for _, t := range response.GetObservations() {
				fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\n", t.GetId().GetValue(), t.GetName(), t.GetPayload(), t.GetCreated(), t.GetUpdated())
			}
			w.Flush()
		},
	}
}

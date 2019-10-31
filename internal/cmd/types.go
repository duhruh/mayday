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

// NewTypesCommand -
func NewTypesCommand(config MaydayConfig, clientProvider mayday.ClientProvider) *cobra.Command {
	typesCommand := &cobra.Command{
		Use:   "types [command]",
		Short: "manipulate type entities",
		Long:  `basic CRUD on type entities`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("types: " + strings.Join(args, " "))
		},
	}

	typesCommand.AddCommand(newTypesCreateCommand(config, clientProvider))
	typesCommand.AddCommand(newTypesListCommand(config, clientProvider))
	typesCommand.AddCommand(newTypesDeleteCommand(config, clientProvider))

	return typesCommand
}

func newTypesCreateCommand(config MaydayConfig, clientProvider mayday.ClientProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "create [json entitiy]",
		Short: "create type",
		Long:  `create type`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			client := clientProvider.Get()

			response, err := client.CreateType(context.TODO(), []byte(args[0]))
			if err != nil {
				println(err.Error())
			}

			t := response.GetType()
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "ID\tName\tSchema\tCreated\tUpdated")
			fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\n", t.GetId().GetValue(), t.GetName(), t.GetSchema(), t.GetCreated(), t.GetUpdated())
			w.Flush()
		},
	}
}

func newTypesDeleteCommand(config MaydayConfig, clientProvider mayday.ClientProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "Delete [json entitiy]",
		Short: "Delete type",
		Long:  `Delete type`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(c *cobra.Command, args []string) {
			client := clientProvider.Get()

			response, err := client.DeleteType(context.TODO(), []byte(args[0]))
			if err != nil {
				println(err.Error())
			}

			t := response.GetType()
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "ID\tName\tSchema\tCreated\tUpdated")
			fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\n", t.GetId().GetValue(), t.GetName(), t.GetSchema(), t.GetCreated(), t.GetUpdated())
			w.Flush()
		},
	}
}

func newTypesListCommand(config MaydayConfig, clientProvider mayday.ClientProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "list types",
		Long:  `list types`,
		Run: func(c *cobra.Command, args []string) {
			client := clientProvider.Get()

			response, err := client.ListTypes(context.TODO())
			if err != nil {
				println(err.Error())
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
			fmt.Fprintln(w, "ID\tName\tSchema\tCreated\tUpdated")
			for _, t := range response.GetTypes() {
				fmt.Fprintf(w, "%s\t%s\t%v\t%v\t%v\n", t.GetId().GetValue(), t.GetName(), t.GetSchema(), t.GetCreated(), t.GetUpdated())
			}
			w.Flush()
		},
	}
}

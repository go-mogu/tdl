package cmd

import (
	"github.com/iyear/tdl/app/contacts"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/spf13/cobra"
)

func NewContacts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cont",
		Short: "A set of contacts tools",
	}

	cmd.AddCommand(NewContactsList())

	return cmd
}

func NewContactsList() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List your contacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return contacts.List(logger.Named(cmd.Context(), "ls"))
		},
	}

	return cmd
}

package cmd

import (
	"github.com/iyear/tdl/app/test"
	"github.com/spf13/cobra"
)

func NewTest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "A set of chat tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			return test.Test(cmd.Context())
		},
	}

	return cmd
}

package cmd

import (
	"github.com/iyear/tdl/app/test"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/spf13/cobra"
)

func NewTest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "A set of chat tools",
	}
	cmd.AddCommand(NewTest1(), NewTest2(), NewTest3(), NewTest4(), NewTest5())
	return cmd
}

func NewTest1() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "test1",
		Short: "List your chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return test.Test1(logger.Named(cmd.Context(), "test1"))
		},
	}

	return cmd
}

func NewTest2() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "test2",
		Short: "List your chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return test.Test2(logger.Named(cmd.Context(), "test2"))
		},
	}

	return cmd
}

func NewTest3() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "test3",
		Short: "List your chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return test.Test3(logger.Named(cmd.Context(), "test3"))
		},
	}

	return cmd
}

func NewTest4() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "test4",
		Short: "List your chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return test.Test4(logger.Named(cmd.Context(), "test4"))
		},
	}

	return cmd
}

func NewTest5() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "test5",
		Short: "List your chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return test.Test5(logger.Named(cmd.Context(), "test5"))
		},
	}

	return cmd
}

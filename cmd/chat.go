package cmd

import (
	"github.com/spf13/cobra"

	"github.com/iyear/tdl/app/chat"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/iyear/tdl/pkg/utils"
)

func NewChat() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chat",
		Short: "A set of chat tools",
	}

	cmd.AddCommand(NewChatList(), NewChatSendMsg(), NewChatSearch(), NewChatMessage(), NewChatUsers())

	return cmd
}

func NewChatList() *cobra.Command {
	var opts chat.ListOptions

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List your chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return chat.List(logger.Named(cmd.Context(), "ls"), opts)
		},
	}

	utils.Cmd.StringEnumFlag(cmd, &opts.Output, "output", "o", string(chat.OutputTable), []string{string(chat.OutputTable), string(chat.OutputJSON)}, "output format")
	cmd.Flags().StringVarP(&opts.Filter, "filter", "f", "true", "filter chats by expression")

	return cmd
}

func NewChatSendMsg() *cobra.Command {
	var opts chat.SendOptions

	cmd := &cobra.Command{
		Use:   "send",
		Short: "send msg to contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			return chat.Send(logger.Named(cmd.Context(), "send"), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Username, "username", "u", "true", "username")
	cmd.Flags().StringVarP(&opts.Msg, "msg", "m", "true", "msg")

	return cmd
}

func NewChatSearch() *cobra.Command {
	var opts chat.AddOptions

	cmd := &cobra.Command{
		Use:   "add",
		Short: "add contacts to user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return chat.Add(logger.Named(cmd.Context(), "add"), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Username, "username", "u", "true", "username")
	return cmd
}

func NewChatMessage() *cobra.Command {
	var opts chat.MsgOptions
	cmd := &cobra.Command{
		Use:   "msg",
		Short: "msg",
		RunE: func(cmd *cobra.Command, args []string) error {
			return chat.NewMessage(logger.Named(cmd.Context(), "msg"), opts)
		},
	}
	cmd.Flags().BoolVar(&opts.Store, "store", false, "store")
	return cmd
}

func NewChatUsers() *cobra.Command {
	var opts chat.UsersOptions

	cmd := &cobra.Command{
		Use:   "users",
		Short: "export users from (protected) channels",
		RunE: func(cmd *cobra.Command, args []string) error {
			return chat.Users(logger.Named(cmd.Context(), "users"), opts)
		},
	}

	cmd.Flags().StringVarP(&opts.Output, "output", "o", "tdl-users.json", "output JSON file path")
	cmd.Flags().StringVarP(&opts.Chat, "chat", "c", "", "domain id (channels, supergroups, etc.)")
	cmd.Flags().BoolVar(&opts.Raw, "raw", false, "export raw message struct of Telegram MTProto API, useful for debugging")
	return cmd
}

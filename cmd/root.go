package cmd

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/iyear/tdl/pkg/consts"
	"github.com/iyear/tdl/pkg/kv"
	"github.com/iyear/tdl/pkg/logger"
	"github.com/iyear/tdl/pkg/utils"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "tdl",
		Short:             "Telegram Downloader, but more than a downloader",
		DisableAutoGenTag: true,
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// init logger
			debug, level := viper.GetBool(consts.FlagDebug), zap.InfoLevel
			if debug {
				level = zap.DebugLevel
			}
			cmd.SetContext(logger.With(cmd.Context(),
				logger.New(level, filepath.Join(consts.LogPath, "latest.log"))))

			ns := viper.GetString(consts.FlagNamespace)
			if ns != "" {
				logger.From(cmd.Context()).Info("Namespace",
					zap.String("namespace", ns))
			}
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return logger.From(cmd.Context()).Sync()
		},
	}

	cmd.AddCommand(NewVersion(), NewLogin(), NewChat(), NewContacts())

	cmd.PersistentFlags().String(consts.FlagProxy, "", "proxy address, only socks5 is supported, format: protocol://username:password@host:port")
	cmd.PersistentFlags().StringP(consts.FlagNamespace, "n", "", "namespace for Telegram session")
	cmd.PersistentFlags().Bool(consts.FlagDebug, false, "enable debug mode")

	cmd.PersistentFlags().IntP(consts.FlagPartSize, "s", 512*1024, "part size for transfer, max is 512*1024")
	cmd.PersistentFlags().IntP(consts.FlagThreads, "t", 4, "max threads for transfer one item")
	cmd.PersistentFlags().IntP(consts.FlagLimit, "l", 2, "max number of concurrent tasks")

	cmd.PersistentFlags().String(consts.FlagNTP, "", "ntp server host, if not set, use system time")
	cmd.PersistentFlags().Duration(consts.FlagReconnectTimeout, 2*time.Minute, "Telegram client reconnection backoff timeout, infinite if set to 0") // #158

	cmd.PersistentFlags().String(consts.FlagTest, "", "use test Telegram client, only for developer")

	// completion
	_ = cmd.RegisterFlagCompletionFunc(consts.FlagNamespace, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ns, err := kv.Namespaces(consts.KVPath)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return ns, cobra.ShellCompDirectiveNoFileComp
	})

	_ = viper.BindPFlags(cmd.PersistentFlags())

	viper.SetEnvPrefix("tdl")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	generateCommandDocs(cmd)

	return cmd
}

func generateCommandDocs(cmd *cobra.Command) {
	docs := filepath.Join(consts.DocsPath, "command")
	if utils.FS.PathExists(docs) {
		if err := doc.GenMarkdownTree(cmd, docs); err != nil {
			color.Red("generate cmd docs failed: %v", err)
		}
	}
}

type completeFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

func completeExtFiles(ext ...string) completeFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		files := make([]string, 0)
		for _, e := range ext {
			f, err := filepath.Glob(toComplete + "*." + e)
			if err != nil {
				return nil, cobra.ShellCompDirectiveDefault
			}
			files = append(files, f...)
		}

		return files, cobra.ShellCompDirectiveFilterDirs
	}
}

package cmd

import (
	"github.com/alipourhabibi/urlshortener/config"
	handlers "github.com/alipourhabibi/urlshortener/internal/handlers/urlhdl"
	"github.com/spf13/cobra"
)

func init() {
	RootCMD.AddCommand(runCMD)
}

var runCMD = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Long:  "Run the application",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().String("config", "config/prod_config.yaml", "config file path")
		err := cmd.ParseFlags(args)
		if err != nil {
			return err
		}

		configFilePath := getConfigfilePath(cmd)
		if configFilePath != "" {
			config.Confs.Load(configFilePath)
		}
		return nil
	},
	RunE: runCmdE,
}

func runCmdE(cmd *cobra.Command, args []string) error {
	return handlers.Launch()
}

package cmd

import (
	"github.com/alipourhabibi/urlshortener/config"
	postgresdb "github.com/alipourhabibi/urlshortener/internal/repository/postgres"
	"github.com/spf13/cobra"
)

func init() {
	RootCMD.AddCommand(migrateCMD)
}

var migrateCMD = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate all models",
	Long:  "Migrate all models",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cmd.Flags().String("config", "config/local_config.yaml", "config file path")
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
	RunE: migrateCmdE,
}

func migrateCmdE(cmd *cobra.Command, args []string) error {
	p, err := postgresdb.New(config.Confs.PostgresDB)
	if err != nil {
		return err
	}
	return p.MigrateAll()
}

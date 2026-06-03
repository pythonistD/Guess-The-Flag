package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pythonistD/Guess-The-Flag/internal/config"
	"github.com/pythonistD/Guess-The-Flag/internal/db"
	"github.com/pythonistD/Guess-The-Flag/internal/db/cleardb"
	"github.com/pythonistD/Guess-The-Flag/internal/db/filldb"
	"github.com/urfave/cli/v3"
)

func configFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    "config",
		Usage:   "path to config file",
		Value:   "./config.yml",
		Aliases: []string{"c"},
	}
}

func getDBFromConfig(yamlConfigPath string) (*sqlx.DB, error) {
	if yamlConfigPath == "" {
		yamlConfigPath = "../config.yml"
	}
	cfg, err := config.LoadConfi
func getDBFromConfig(yamlConfigPath string) (*sqlx.DB, error) {
	if yamlConfigPath == "" {
		yamlConfigPath = "../config"
	}
	cfg, err := config.LoadConfigFromFile(yamlConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load data from config: %s", yamlConfigPath)
	}
	database, err := db.NewPostgres(cfg.DBConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create db instance from cfg: %v", cfg)
	}
	var existing int
	err = database.Get(&existing, "SELECT COUNT(*) FROM countries")
	if err != nil {
		return nil, fmt.Errorf("failed to execute select query for countries table")
}

func main() {
	var cmd *cli.Command
	// Database subcommands
	databaseCommands := []*cli.Command{
		{
			Name:  "fill",
			Usage: "fill db with countries, country names and images",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				cfgPath := cmd.String("config")
				database, err := getDBFromConfig(cfgPath)
				if err != nil {
					return fmt.Errorf("fill db error: %w", err)
				}
				err = filldb.FillCountriesInDb(database)
				if err != nil {
					return fmt.Errorf("fill db error: %w", err)
				}
				return nil
			},
		},
		{
			Name:  "clear",
			Usage: "clears db tables",
			Action: func(ctx context.Context, cmd *cli.Command) error {
				cfgPath := cmd.String("config")
				database, err := getDBFromConfig(cfgPath)
				if err != nil {
					return fmt.Errorf("clear db tables error: %w", err)
				}
				err = cleardb.ClearDB(database)
				if err != nil {
					return fmt.Errorf("clear db tables error: %w", err)
				}
	databaseCmd := &cli.Command{
		Name:     "database",
		Usage:    "Interaction with DB: clear db, fill db and etc.",
		Aliases:  []string{"db"},
		Flags:    []cli.Flag{configFlag()},
				Usage:   "path to config file",
				Value:   "./config.yml",
				Aliases: []string{"c"},
	composeEnvCmd := &cli.Command{
		Name:  "compose-env",
		Usage: "print docker compose variables from config.yml (for --env-file)",
		Flags: []cli.Flag{configFlag()},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			env, err := config.LoadComposeEnv(cmd.String("config"))
			if err != nil {
				return err
			}
			_, err = os.Stdout.WriteString(env.FormatDotenv())
			return err
		},
	}

	root := &cli.Command{
		Name:     "guess-the-flag",
		Usage:    "Guess The Flag utilities",
		Commands: []*cli.Command{databaseCmd, composeEnvCmd},
	}

	if err := root.Run(context.Background(), os.Args); err != nil {
		},
		Commands: databaseCommands,

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

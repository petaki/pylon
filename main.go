package main

import (
	"github.com/joho/godotenv"
	"github.com/petaki/pylon/internal/cmd"
	"github.com/petaki/support-go/cli"
)

func init() {
	configFile, err := cmd.ConfigFile()
	if err == nil {
		godotenv.Load(configFile)
	}
}

func main() {
	(&cli.App{
		Name:       "Pylon",
		Version:    "1.0.0",
		TryDefault: true,
		Groups: []*cli.Group{
			{
				Name:  "link",
				Usage: "Handle the links",
				Commands: []*cli.Command{
					{
						Name:       "add",
						Usage:      "Add a link",
						Arguments:  []string{"url"},
						HandleFunc: cmd.LinkAdd,
					},
					{
						Name:       "delete",
						Usage:      "Delete the link",
						Arguments:  []string{"id"},
						HandleFunc: cmd.LinkDelete,
					},
				},
			},
			{
				Name:  "config",
				Usage: "Handle the configs",
				Commands: []*cli.Command{
					{
						Name:       "init",
						Usage:      "Initialize the config file",
						HandleFunc: cmd.ConfigInit,
					},
				},
			},
		},
	}).Execute()
}

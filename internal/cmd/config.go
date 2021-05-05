package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path"

	"github.com/petaki/support-go/cli"
)

// ConfigFile function.
func ConfigFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	return path.Join(usr.HomeDir, ".pylonfile"), nil
}

// ConfigInit command.
func ConfigInit(group *cli.Group, command *cli.Command, arguments []string) int {
	configFile, err := ConfigFile()
	if err != nil {
		return command.PrintError(err)
	}

	fmt.Println("=> Create a " + cli.Green("config file"))

	f, err := os.Create(configFile)
	if err != nil {
		return command.PrintError(err)
	}

	defer f.Close()

	fmt.Println("---> Write values")

	_, err = f.WriteString(`HEADLESS_SHELL_HOST=http://127.0.0.1:9222
MEILISEARCH_HOST=http://127.0.0.1:7700
MEILISEARCH_API_KEY=
MEILISEARCH_INDEX=pylon
`)
	if err != nil {
		return command.PrintError(err)
	}

	fmt.Println("---> Sync content")

	err = f.Sync()
	if err != nil {
		return command.PrintError(err)
	}

	fmt.Println("=> Config file created: " + cli.Green(f.Name()))

	return cli.Success
}

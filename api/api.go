package api

import (
	application "cf-appmonit/application"
	term "cf-appmonit/terminal"
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
)

func ask(prompt string, args ...interface{}) (answer string) {
	fmt.Println("")
	fmt.Printf(prompt+" ", args...)
	fmt.Scanln(&answer)
	return
}

func New() (app *cli.App) {
	app = cli.NewApp()
	app.Name = "cf-appmonit"
	app.Usage = "A command line tool to interact with Cloud Foundry"
	app.Version = "1.0.0.alpha"
	app.Commands = []cli.Command{
		{
			Name:        "monit",
			ShortName:   "m",
			Description: "Monitor specific application CPU, memory and disk usages",
			Action: func(c *cli.Context) {
				domain := ask("Domain [mandatory]%s", term.Cyan(">"))
				if domain == "" {
					log.Println("DOMAIN is a mandatory")
					os.Exit(0)
				}

				username := ask("Email [mandatory]%s", term.Cyan(">"))
				if username == "" {
					log.Println("Email is a mandatory")
					os.Exit(0)
				}

				password := ask("Password [mandatory]%s", term.Cyan(">"))
				if password == "" {
					log.Println("Password is a mandatory")
					os.Exit(0)
				}

				appName := ask("App Name [mandatory]%s", term.Cyan(">"))
				if appName == "" {
					log.Println("App name is a mandatory")
					os.Exit(0)
				}
				application.Monitor(domain, username, password, appName)
			},
		},
	}
	return
}

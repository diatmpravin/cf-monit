package main

import (
	"cf-appmonit/api"
	"os"
)

func main() {
	app := api.New()
	app.Run(os.Args)
}

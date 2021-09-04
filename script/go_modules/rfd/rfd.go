package main

import (
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"os"
)

var logger Trace = TraceLog{}
var config *configuration

type configuration struct {
	RFDRootDirectory string `yaml:"rfd-root-directory"`
	InstallDirectory string `yaml:"install-directory"`
	RFDRelativeDirectory string `yaml:"rfd-relative-directory"`
}



func main() {

	config = populateConfig()
	app := createCommandLineApp()
	err := app.Run(os.Args)
	CheckFatal(err)

}

func createCommandLineApp() *cli.App {
	app := &cli.App{
		Name:  "rfd",
		Usage: "Create new rfd's, index and output their status, and manage their .",
		Commands: []*cli.Command{
			{
				Name:     "update-status",
				Category: "Information",
				Usage:    "Output the status of all rfd's to `FILE` in markdown format.",
				Action: func(c *cli.Context) error {
					return CreateEntries()
				},
			},
			{
				Name:  "new",
				Usage: "Create a new rfd",
				Action: func(c *cli.Context) error {
					NewRFD()
					return nil
				},
				Subcommands: []*cli.Command{
					{
						Name: "repository",
						Action: func(c *cli.Context) error {
							return nil
						},
					},
				},
			},
		},
		/*Flags: []cli.Flag {
			&cli.StringFlag{
				Name: "status",
				Aliases: []string{"s"},
				Usage: "Output the status of all rfd's to `FILE` in markdown format.",
				Value: "status.md",
				DefaultText: "status.md",
			},
			&cli.StringFlag{
				Name: "create",
				Aliases: []string{"c"},
				Usage: "Create a new rfd, with an optionally specified `RFD ID` in nnnn format",
			},
		},
		Action: func(c *cli.Context) error {
			//fmt.Println("rfd %q", c.Args().Get(0))
			if c.NArg() > 0 {

			} else {
				fmt.Println("")
			}
			return nil
		},*/
	}

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`
	return app
}

func populateConfig() *configuration {
	// Create config structure
	config := &configuration{}

	// Open config file
	file, err := os.Open("./config.yml")
	CheckFatal(err)

	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	err = d.Decode(&config)
	CheckFatal(err)

	return config
}

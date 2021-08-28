package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var logger Trace = TraceLog{}

func main() {

	app := &cli.App{
		Name:  "rfd",
		Usage: "Create new rfd's, index and output their status, and manage their .",
		Commands: []*cli.Command{
			{
				Name:     "test-branch",
				Category: "Test",
				Usage:    "Test",
				Action: func(c *cli.Context) error {
					return CreateBranch("TEST-BRANCH")
				},
			},
			{
				Name:     "write-status",
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
						Name: "rfd",
						Action: func(c *cli.Context) error {
							return nil
						},
					},
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

	/*cli.AppHelpTemplate = fmt.Sprintf(`%s

	WEBSITE: http://awesometown.example.com

	SUPPORT: support@awesometown.example.com

	`, cli.AppHelpTemplate)*/

	// EXAMPLE: Override a template
	/*cli.AppHelpTemplate = `NAME:
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
	`*/
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

	// EXAMPLE: Replace the `HelpPrinter` func
	//cli.HelpPrinter = func(w io.Writer, templ string, data interface{}) {
	//	fmt.Println("Ha HA.  I pwnd the help!!1")
	//}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

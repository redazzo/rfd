package main

import (
	"encoding/base64"
	"fmt"
	"github.com/redazzo/rfd/cmd/rfd/internal/config"
	"github.com/redazzo/rfd/cmd/rfd/internal/global"
	"github.com/redazzo/rfd/cmd/rfd/internal/util"
	"github.com/urfave/cli/v2"
	"os"
	"runtime"
)

func main() {
	app := createCommandLineApp()
	err := app.Run(os.Args)
	util.CheckFatal(err)
}

func createCommandLineApp() *cli.App {
	app := &cli.App{
		Name:  "rfd",
		Usage: "Create new rfd's, index, and manage their status.",
		Commands: []*cli.Command{
			{
				Name:  "check",
				Usage: "Check environment is suitable to ensure a clean run when creating a new RFD.",
				Action: func(c *cli.Context) error {
					config.PreConfigure()
					config.Configure()
					config.PostConfigure()
					config.CheckAndReportOnRepositoryState()
					return nil
				},
			},
			{
				Name:  "index",
				Usage: "Output the status of all rfd's to index.md in markdown format.",
				Action: func(c *cli.Context) error {
					config.PreConfigure()
					config.Configure()
					config.PostConfigure()
					Index()
					return nil
				},
			},
			{
				Name:  "new",
				Usage: "Create a new rfd",
				Action: func(c *cli.Context) error {
					if config.CheckAndReportOnRepositoryState() {
						config.PreConfigure()
						config.Configure()
						config.PostConfigure()
						new()
					} else {
						fmt.Println("Creating a new RFD creates and switches to new branch. Commit (or otherwise) unstaged and/or uncommitted work first.")
						fmt.Println()
					}
					return nil
				},
			},

			{
				Name:  "init",
				Usage: "Initialise an RFD repository.",
				Action: func(c *cli.Context) error {
					config.PreConfigure()
					config.InitialiseRepo()

					return nil
				},
			},
			{
				Name:  "environment",
				Usage: "Displays configuration settings and relevant operating system environment variables.",
				Action: func(c *cli.Context) error {
					config.PreConfigure()
					config.Configure()
					config.PostConfigure()
					displayEnvironment()
					return nil
				},
			},
			{
				Name:  "merge",
				Usage: "Transitions an RFD's status to Accepted, captures the discussion link from the user, and merges it into the main branch (NOT IMPLEMENTED)",
				Action: func(c *cli.Context) error {
					config.PreConfigure()
					config.Configure()
					config.PostConfigure()
					doMerge()
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "Displays the status of the current RFD (as per branch). (NOT IMPLEMENTED)",
				Action: func(c *cli.Context) error {

					config.PreConfigure()
					config.Configure()
					config.PostConfigure()

					getDefaultStatus()
					return nil
				},
			},
		},
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

func displayEnvironment() {
	operatingSystem := runtime.GOOS
	fmt.Println("OS: " + operatingSystem)
	switch operatingSystem {
	case "windows":
		fmt.Println("HOMEDRIVE=" + os.Getenv(global.HOMEDRIVE))
		fmt.Println("HOMEPATH =" + os.Getenv(global.HOMEPATH))
	case "linux":
		fmt.Println("HOME=" + os.Getenv(global.HOME))
	}
	fmt.Println("RFD root directory=" + global.APP_CONFIG.RootDirectory)
	//fmt.Println("RFD relative directory=" + appConfig.RFDRelativeDirectory)
	fmt.Println("Installation directory=" + global.APP_CONFIG.TemplatesDirectory)
	fmt.Println("SSH public key directory=" + util.GetSSHPath())

	publicKey, err := util.GetPublicKey()
	util.CheckFatal(err)

	bytes := publicKey.Signer.PublicKey().Marshal()

	sPublicKey := base64.StdEncoding.EncodeToString(bytes) + " " + publicKey.User
	util.CheckFatal(err)

	fmt.Println("SSH Public Key=" + sPublicKey)

	println()
	println()

	fmt.Printf("%+v", global.APP_CONFIG)
	fmt.Println()
}

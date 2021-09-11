package main

import (
	"encoding/base64"
	"fmt"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
)

const HOME string = "HOME"
const HOMEDRIVE string = "HOMEDRIVE"
const HOMEPATH string = "HOMEPATH"

var logger Trace
var appConfig *configuration
var sshDir string

type configuration struct {
	RFDRootDirectory     string `yaml:"rfd-root-directory"`
	InstallDirectory     string `yaml:"install-directory"`
	RFDRelativeDirectory string `yaml:"rfd-relative-directory"`
	PrivateKeyFileName   string `yaml:"private-key-file-name"`
}

func init() {
	logger = TraceLog{}
	appConfig = populateConfig()
	initSSHDIR()
}

func main() {
	app := createCommandLineApp()
	err := app.Run(os.Args)
	CheckFatal(err)
}

//
func createCommandLineApp() *cli.App {
	app := &cli.App{
		Name:  "rfd",
		Usage: "Create new rfd's, index and output their status, and manage their .",
		Commands: []*cli.Command{
			{
				Name:  "create-index",
				Usage: "Output the status of all rfd's to index.md in markdown format.",
				Action: func(c *cli.Context) error {
					Index()
					return nil
				},
			},
			{
				Name:  "new",
				Usage: "Create a new rfd",
				Action: func(c *cli.Context) error {
					New()
					return nil
				},
			},
			{
				Name:  "show-status",
				Usage: "Displays the status of <nnnn>. Will output the status of every RFD if it isn't provided an RFD ID.",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  "init",
				Usage: "Displays the status of <nnnn>. Will output the status of every RFD if it isn't provided an RFD ID.",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  "environment",
				Usage: "Displays configuration settings and relevant operating system environment variables.",
				Action: func(c *cli.Context) error {
					displayEnvironment()
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

func initSSHDIR() {
	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		sshDir = os.Getenv(HOMEDRIVE) + os.Getenv(HOMEPATH)
	case "linux":
		sshDir = os.Getenv(HOME)
	}
}

func displayEnvironment() {
	operatingSystem := runtime.GOOS
	fmt.Println(operatingSystem)
	switch operatingSystem {
	case "windows":
		fmt.Println("HOMEDRIVE=" + os.Getenv(HOMEDRIVE))
		fmt.Println("HOMEPATH =" + os.Getenv(HOMEPATH))
	case "linux":
		fmt.Println("HOME=" + os.Getenv(HOME))
	}
	fmt.Println("RFD root directory=" + appConfig.RFDRootDirectory)
	fmt.Println("RFD relative directory=" + appConfig.RFDRelativeDirectory)
	fmt.Println("Installation directory=" + appConfig.InstallDirectory)
	fmt.Println("SSH public key directory=" + getSSHPath())

	publicKey, err := getPublicKey()
	CheckFatal(err)

	bytes := publicKey.Signer.PublicKey().Marshal()

	sPublicKey := base64.StdEncoding.EncodeToString(bytes) + " " + publicKey.User
	CheckFatal(err)

	fmt.Println("SSH Public Key=" + sPublicKey)
}

func populateConfig() *configuration {
	// Create appConfig structure
	config := &configuration{}

	// Open appConfig file
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

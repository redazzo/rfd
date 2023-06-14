package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

var APP_CONFIG *Configuration

func Configure() {
	configuration := &Configuration{}
	configuration.Populate()
	APP_CONFIG = configuration
}

type Configuration struct {
	RootDirectory      string `yaml:"root-directory"`
	TemplatesDirectory string `yaml:"templates-directory"`
	PrivateKeyFileName string `yaml:"private-key-file-name"`
	InitialAuthor      string `yaml:"initial-author"`
	Organisation       string `yaml:"organisation"`
	InstigationDate    string `yaml:"instigation-date"`
}

func (c *Configuration) Get001ReadmeFileLocation() string {
	return c.TemplatesDirectory + PATH_SEPARATOR + "0001" + PATH_SEPARATOR + "readme.md"
}

func (c *Configuration) GetReadmeTemplateLocation() string {
	return c.TemplatesDirectory + PATH_SEPARATOR + "readme.md"
}

func (c *Configuration) Populate() {

	err := CheckConfigurationFilePresence()
	CheckFatal(err)
	err = c.populateConfig()
	if err != nil {
		fmt.Println("Error populating configuration")
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *Configuration) populateConfig() error {

	err := CheckConfigurationFilePresence()
	if err != nil {
		return err
	}

	// Open appConfig file
	file, err := os.Open("./config.yml")
	CheckFatal(err)

	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	err = d.Decode(&c)
	CheckFatal(err)

	return err
}

func CheckConfigurationFilePresence() error {
	// Check to ensure there is a config file present
	_, err := os.Stat("./config.yml")
	if os.IsNotExist(err) {
		fmt.Print(
			"\n  There doesn't appear to be a configuration file present.\n" +
				"  Either run 'rfd init', or if you have, make sure you are\n" +
				"  in the root directory of your rfd repository.\n\n")
		return err
	}
	return nil
}

type States struct {
	RFDStates []map[string]map[string]string `yaml:"rfd-states"`
}

type RFDMetadata struct {
	RFDID     string
	Title     string
	Authors   string
	State     string
	Link      string
	RFDStates []map[string]map[string]string
}

const HOME string = "HOME"
const HOMEDRIVE string = "HOMEDRIVE"
const HOMEPATH string = "HOMEPATH"
const PATH_SEPARATOR string = string(os.PathSeparator)

var APP_STATES *States

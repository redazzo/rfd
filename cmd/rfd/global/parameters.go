package global

type Configuration struct {
	RootDirectory      string `yaml:"root-directory"`
	TemplatesDirectory string `yaml:"templates-directory"`
	//RFDRelativeDirectory string                         `yaml:"rfd-relative-directory"`
	PrivateKeyFileName string `yaml:"private-key-file-name"`
	InitialAuthor      string `yaml:"initial-author"`
	Organisation       string `yaml:"organisation"`
	InstigationDate    string `yaml:"instigation-date"`
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

var SSHDIR string
var REPO_TEMPLATE_FILE_LOCATION string
var TEMPLATE_FILE_LOCATION string
var PATH_SEPARATOR string

var APP_CONFIG *Configuration
var APP_STATES *States

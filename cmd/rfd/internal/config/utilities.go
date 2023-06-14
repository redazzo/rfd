package config

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	//"github.com/redazzo/rfd/cmd/rfd/internal/config"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
)

type Trace interface {
	TraceLog(msg string)
}

type TraceLog struct {
	// In future, log this to a file

}

var Logger Trace = TraceLog{}

func CheckFatalWithMessage(e error, msg string) {

	if e != nil {
		if msg != "" {
			println(msg)
		}
		debug.PrintStack()
		log.Fatal(e)

	}
}

func CheckFatal(e error) {

	CheckFatalWithMessage(e, "")
}

func (t TraceLog) TraceLog(msg string) {
	log.Print(msg)
}

func IsRFDIDFormat(name string) (bool, error) {
	entryIsBranchID, err := regexp.MatchString(`(^\d{4}).*`, name)
	return entryIsBranchID, err
}

func Exists(sFile string) bool {
	_, err := os.Stat(sFile)

	exists := true

	if err != nil {
		if os.IsNotExist(err) {
			exists = false
		}
	}

	return exists
}

func GetUserInput(txt string) string {

	print(txt + " ")
	reader := bufio.NewReader(os.Stdin)

	// Hack, but it'll do. Too lazy to find a better way ...
	responseTxt, err := reader.ReadString('\n')
	responseTxt = strings.TrimSuffix(responseTxt, "\n")
	CheckFatal(err)
	return responseTxt
}

func CopyToRoot(source string, target string, force bool) {

	bytesRead, err := ioutil.ReadFile(source)

	CheckFatal(err)

	if Exists(APP_CONFIG.RootDirectory + PATH_SEPARATOR + target) {
		if force {
			err = os.Remove(APP_CONFIG.RootDirectory + PATH_SEPARATOR + target)
			CheckFatal(err)
		} else {
			log.Fatal("Error: Attempted to overwrite " + source + " to RFD root.")
		}
	}

	err = os.WriteFile(APP_CONFIG.RootDirectory+PATH_SEPARATOR+target, bytesRead, 0744)

	if err != nil {
		log.Fatal(err)
	}

}

func PushToOrigin(r *git.Repository) error {

	publicKey, err := GetPublicKey()

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       publicKey,
		Force:      true,
	})
	CheckFatal(err)

	return err
}

func GetRFDDirectory(sRfdNumber string) string {
	return APP_CONFIG.RootDirectory + "/" + sRfdNumber
}

func FetchTemplateDirectory() {

	// Fetch the template directory from the remote repo
	// and copy it to the local RFD root directory.

	// This is a one-time operation, so we don't need to
	// check if the directory already exists.

	repoURL := "https://github.com/redazzo/rfd"
	targetDirectory := APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template"

	// Clone the repo to the target directory
	_, err := git.PlainClone(targetDirectory, false, &git.CloneOptions{
		URL:               repoURL,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		fmt.Printf("Failed to clone repository: %v\n", err)
		return
	}

	// Repository cloned successfully
	fmt.Println("Repository cloned successfully.")

}

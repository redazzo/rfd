package config

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"io"
	"net/http"
	"path/filepath"

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

func fetchFileFromURL(url string) ([]byte, error) {

	// Make an HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error while fetching the URL: %s", err)
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error while reading the response body: %s", err)
		return nil, err
	}

	return body, nil

}

func WriteTemplates() (string, error) {

	type FileToWrite struct {
		URL         string
		Destination string
	}

	targetDirectory := APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR

	filesToWrite := []FileToWrite{
		{
			URL:         "https://sea-turtle-app-ufxrk.ondigitalocean.app/readme.md",
			Destination: targetDirectory + "readme.md",
		},
		{
			URL:         "https://sea-turtle-app-ufxrk.ondigitalocean.app/states.yml",
			Destination: targetDirectory + "states.yml",
		},
		{
			URL:         "https://sea-turtle-app-ufxrk.ondigitalocean.app/0001/readme.md",
			Destination: targetDirectory + "0001" + PATH_SEPARATOR + "readme.md",
		},
	}

	for _, file := range filesToWrite {
		fileContent, err := fetchFileFromURL(file.URL)
		if err != nil {
			return targetDirectory, err
		}

		err = os.MkdirAll(filepath.Dir(file.Destination), 0744)
		if err != nil {
			return targetDirectory, err
		}

		err = os.WriteFile(file.Destination, fileContent, 0744)
		if err != nil {
			return targetDirectory, err
		}

		log.Printf("Wrote %s template to %s ...\n", filepath.Base(file.URL), filepath.Dir(file.Destination))
	}

	return targetDirectory, nil

}

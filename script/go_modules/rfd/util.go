package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
)

type Trace interface {
	traceLog(msg string)
}

type TraceLog struct {
	// In future, log this to a file

}

func CheckFatal(e error) {

	if e != nil {
		debug.PrintStack()

		fmt.Println("Rolling back ...")

		for _, rollbackFunction := range rollbackFunctions {
			rollbackFunction()
		}
		log.Fatal(e)

	}
}

func (t TraceLog) traceLog(msg string) {
	log.Print(msg)
}

func isRFDIDFormat(name string) (bool, error) {
	entryIsBranchID, err := regexp.MatchString(`(^\d{4}).*`, name)
	return entryIsBranchID, err
}

func getPublicKey() (*ssh.PublicKeys, error) {
	sshPath := getSSHPath()
	publicKey, err := ssh.NewPublicKeysFromFile("git", sshPath, "")
	return publicKey, err
}

func getSSHPath() string {
	sPathseparator := string(os.PathSeparator)
	sshPath := sshDir + sPathseparator + ".ssh" + sPathseparator + appConfig.PrivateKeyFileName
	return sshPath
}

func isFileExists(sFile string) bool {
	_ , err := os.Stat(sFile)

	exists := true

	if err != nil {
		if os.IsNotExist(err) {
			exists = false
		}
	}

	return exists
}

func getUserInput(txt string) string {

	print(txt + " ")
	reader := bufio.NewReader(os.Stdin)

	// Hack, but it'll do. Too lazy to find a better way ...
	responseTxt, err := reader.ReadString('\n')
	responseTxt = strings.TrimSuffix(responseTxt, "\n")
	CheckFatal(err)
	return responseTxt
}

func copyToRoot(source string, target string, force bool) {

	bytesRead, err := ioutil.ReadFile(source)

	CheckFatal(err)

	if isFileExists(appConfig.RFDRootDirectory + sPathseparator + target) {
		if force {
			err = os.Remove(appConfig.RFDRootDirectory + sPathseparator + target)
			CheckFatal(err)
		} else {
			log.Fatal("Error: Attempted to overwrite " + source + " to RFD root.")
		}
	}

	err = ioutil.WriteFile(appConfig.RFDRootDirectory + sPathseparator + target, bytesRead, 0744)

	if err != nil {
		log.Fatal(err)
	}

}

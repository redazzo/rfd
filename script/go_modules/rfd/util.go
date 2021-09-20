package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
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

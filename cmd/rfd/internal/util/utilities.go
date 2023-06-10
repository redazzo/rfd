package util

import (
	"bufio"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/redazzo/rfd/cmd/rfd/internal/global"
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

func GetPublicKey() (*ssh.PublicKeys, error) {
	sshPath := GetSSHPath()
	publicKey, err := ssh.NewPublicKeysFromFile("git", sshPath, "")
	return publicKey, err
}

func GetSSHPath() string {
	sPathseparator := string(os.PathSeparator)
	sshPath := global.SSHDIR + sPathseparator + ".ssh" + sPathseparator + global.APP_CONFIG.PrivateKeyFileName
	return sshPath
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

	if Exists(global.APP_CONFIG.RootDirectory + global.PATH_SEPARATOR + target) {
		if force {
			err = os.Remove(global.APP_CONFIG.RootDirectory + global.PATH_SEPARATOR + target)
			CheckFatal(err)
		} else {
			log.Fatal("Error: Attempted to overwrite " + source + " to RFD root.")
		}
	}

	err = os.WriteFile(global.APP_CONFIG.RootDirectory+global.PATH_SEPARATOR+target, bytesRead, 0744)

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
	return global.APP_CONFIG.RootDirectory + "/" + sRfdNumber
}

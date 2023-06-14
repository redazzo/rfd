package config

import (
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
)

var SSHDIR string

func GetPublicKey() (*ssh.PublicKeys, error) {
	sshPath := GetSSHPath()
	publicKey, err := ssh.NewPublicKeysFromFile("git", sshPath, "")
	return publicKey, err
}

func GetSSHPath() string {
	sPathseparator := string(os.PathSeparator)
	sshPath := SSHDIR + sPathseparator + ".ssh" + sPathseparator + APP_CONFIG.PrivateKeyFileName
	return sshPath
}

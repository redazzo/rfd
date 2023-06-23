package config

import (
	"log"
	"os"
	"testing"
)

func TestWriteTemplates(t *testing.T) {

	APP_CONFIG = &Configuration{}

	APP_CONFIG.RootDirectory = "/tmp"

	err := WriteTemplates()
	if err != nil {
		t.Errorf("Error writing templates: %s", err)
	}

	// Test for existence of readme.md
	_, err = os.Stat(APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR + "readme.md")
	if err != nil {
		t.Errorf("Error reading readme.md: %s", err)
	}
	log.Println("Readme.md template present in root directory ...")

	// Test for existence of states.yml
	_, err = os.Stat(APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR + "states.yml")
	if err != nil {
		t.Errorf("Error reading states.yml: %s", err)
	}
	log.Println("States.yml template present in root directory ...")

	// Test for existence of 0001 directory
	_, err = os.Stat(APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR + "0001")
	if err != nil {
		t.Errorf("Error reading 0001 directory: %s", err)
	}
	log.Println("0001 directory present in root directory ...")

	// Test for existence of readme.md
	_, err = os.Stat(APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR + "0001" + PATH_SEPARATOR + "readme.md")
	if err != nil {
		t.Errorf("Error reading readme.md: %s", err)
	}
	log.Println("Readme.md template present in 0001 directory ...")

	// Clean up
	log.Println("Cleaning up ...")

	// We only have rights to remove the files we created
	if err := os.Remove(APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR + "readme.md"); err != nil {
		t.Errorf("Error removing readme.md: %s", err)
	}
	log.Println("Removed readme.md ...")

	if err := os.Remove(APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR + "states.yml"); err != nil {
		t.Errorf("Error removing states.yml: %s", err)
	}
	log.Println("Removed states.yml ...")

	if err := os.Remove(APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR + "0001" + PATH_SEPARATOR + "readme.md"); err != nil {
		t.Errorf("Error removing readme.md: %s", err)
	}
	log.Println("Removed readme.md from 0001 directory ...")

	if err := os.RemoveAll(APP_CONFIG.RootDirectory + PATH_SEPARATOR + "template" + PATH_SEPARATOR + "0001"); err != nil {
		t.Errorf("Error removing 0001 directory contents: %s", err)
	}

	log.Println("Cleaned up ...")

}

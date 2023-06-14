package main

import (
	"bytes"
	"fmt"
	"github.com/redazzo/rfd/cmd/rfd/internal/config"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func Index() error {

	config.Logger.TraceLog("Creating index file ...")

	// The entries to be read and written to the
	// Markdown table file.

	// The file that will be updated with the
	// rfd entries and their status.
	mdTableFile := openMetadataTableFile()

	defer mdTableFile.Close()

	entries := getDirectories()
	for _, entry := range entries {

		entryIsBranchID, err := config.IsRFDIDFormat(entry.Name())
		config.CheckFatal(err)

		if entryIsBranchID {

			config.Logger.TraceLog("Matched " + entry.Name())

			branchID := entry.Name()

			if entry.IsDir() {

				subEntries, err := ioutil.ReadDir(config.APP_CONFIG.RootDirectory + "/" + entry.Name())
				config.CheckFatal(err)

				for _, subEntry := range subEntries {

					if !subEntry.IsDir() {

						isReadmeFile, err := regexp.MatchString(`(?i)^readme.md`, subEntry.Name())
						config.CheckFatal(err)

						if isReadmeFile {

							metaData := readMetadataFromReadmeFile(subEntry, entry)

							writeMetadataToTableFile(metaData, mdTableFile, branchID)

						}

					}
				}

			}

		}

	}

	return nil
}

func getDirectories() []os.FileInfo {
	entries, err := ioutil.ReadDir(config.APP_CONFIG.RootDirectory)
	config.CheckFatal(err)
	return entries
}

func openMetadataTableFile() *os.File {
	mdTableFile, err := os.Create(config.APP_CONFIG.RootDirectory + "/index.md")
	config.CheckFatal(err)

	_, err = mdTableFile.WriteString("**Index of Requests for Discussion**\n\n")
	config.CheckFatal(err)
	_, err = mdTableFile.WriteString("| **RFD Id** | **Title** | **State** | **Author(s)** |\n")
	config.CheckFatal(err)
	_, err = mdTableFile.WriteString("|------------|-----------|-----------|------------------------|\n")
	config.CheckFatal(err)
	return mdTableFile
}

func writeMetadataToTableFile(metaData map[string]interface{}, mdTableFile *os.File, branchID string) {
	title := fmt.Sprintf("%v", metaData["title"])
	authors := fmt.Sprintf("%v", metaData["authors"])
	state := fmt.Sprintf("%v", metaData["state"])

	config.Logger.TraceLog(title + ":" + authors + ":" + state)

	_, err := mdTableFile.WriteString("|[" + branchID + "](./" + branchID + "/readme.md)|" + title + "|" + state + "|" + authors + "|\n")
	config.CheckFatal(err)

	config.Logger.TraceLog("recorded: " + branchID)
	config.Logger.TraceLog("----------------------------------------------")
}

func readMetadataFromReadmeFile(subEntry os.FileInfo, entry os.FileInfo) map[string]interface{} {
	config.Logger.TraceLog("Found " + config.APP_CONFIG.RootDirectory + "/" + entry.Name() + "/" + subEntry.Name())
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	file, err := os.ReadFile(config.APP_CONFIG.RootDirectory + "/" + entry.Name() + "/" + subEntry.Name())
	config.CheckFatal(err)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(file, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	metaData := meta.Get(context)
	return metaData
}

package main

import (
	"bytes"
	"fmt"
	"github.com/redazzo/rfd/cmd/rfd/global"
	"github.com/redazzo/rfd/cmd/rfd/util"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func Index() error {

	util.Logger.TraceLog("Creating index file ...")

	// The entries to be read and written to the
	// Markdown table file.

	// The file that will be updated with the
	// rfd entries and their status.
	mdTableFile := openMetadataTableFile()

	defer mdTableFile.Close()

	entries := getDirectories()
	for _, entry := range entries {

		entryIsBranchID, err := util.IsRFDIDFormat(entry.Name())
		util.CheckFatal(err)

		if entryIsBranchID {

			util.Logger.TraceLog("Matched " + entry.Name())

			branchID := entry.Name()

			if entry.IsDir() {

				subEntries, err := ioutil.ReadDir(global.APP_CONFIG.RootDirectory + "/" + entry.Name())
				util.CheckFatal(err)

				for _, subEntry := range subEntries {

					if !subEntry.IsDir() {

						isReadmeFile, err := regexp.MatchString(`(?i)^readme.md`, subEntry.Name())
						util.CheckFatal(err)

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
	entries, err := ioutil.ReadDir(global.APP_CONFIG.RootDirectory)
	util.CheckFatal(err)
	return entries
}

func openMetadataTableFile() *os.File {
	mdTableFile, err := os.Create(global.APP_CONFIG.RootDirectory + "/index.md")
	util.CheckFatal(err)

	_, err = mdTableFile.WriteString("**Index of Requests for Discussion**\n\n")
	util.CheckFatal(err)
	_, err = mdTableFile.WriteString("| **RFD Id** | **Title** | **State** | **Author(s)** |\n")
	util.CheckFatal(err)
	_, err = mdTableFile.WriteString("|------------|-----------|-----------|------------------------|\n")
	util.CheckFatal(err)
	return mdTableFile
}

func writeMetadataToTableFile(metaData map[string]interface{}, mdTableFile *os.File, branchID string) {
	title := fmt.Sprintf("%v", metaData["title"])
	authors := fmt.Sprintf("%v", metaData["authors"])
	state := fmt.Sprintf("%v", metaData["state"])

	util.Logger.TraceLog(title + ":" + authors + ":" + state)

	_, err := mdTableFile.WriteString("|[" + branchID + "](./" + branchID + "/readme.md)|" + title + "|" + state + "|" + authors + "|\n")
	util.CheckFatal(err)

	util.Logger.TraceLog("recorded: " + branchID)
	util.Logger.TraceLog("----------------------------------------------")
}

func readMetadataFromReadmeFile(subEntry os.FileInfo, entry os.FileInfo) map[string]interface{} {
	util.Logger.TraceLog("Found " + global.APP_CONFIG.RootDirectory + "/" + entry.Name() + "/" + subEntry.Name())
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	file, err := os.ReadFile(global.APP_CONFIG.RootDirectory + "/" + entry.Name() + "/" + subEntry.Name())
	util.CheckFatal(err)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(file, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	metaData := meta.Get(context)
	return metaData
}

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func Index() error {

	logger.traceLog("Creating index file ...")

	// The entries to be read and written to the
	// markdown table file.

	// The file that will be updated with the
	// rfd entries and their status.
	mdTableFile := openMetadataTableFile()

	defer mdTableFile.Close()

	entries := getDirectories()
	for _, entry := range entries {

		entryIsBranchID, err := isRFDIDFormat(entry.Name())
		CheckFatal(err)

		if entryIsBranchID {

			logger.traceLog("Matched " + entry.Name())

			branchID := entry.Name()

			if entry.IsDir() {

				subEntries, err := ioutil.ReadDir(appConfig.RFDRootDirectory + "/" + entry.Name())
				CheckFatal(err)

				for _, subEntry := range subEntries {

					if !subEntry.IsDir() {

						isReadmeFile, err := regexp.MatchString(`(?i)^readme.md`, subEntry.Name())
						CheckFatal(err)

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
	entries, err := ioutil.ReadDir(appConfig.RFDRootDirectory)
	CheckFatal(err)
	return entries
}

func openMetadataTableFile() *os.File {
	mdTableFile, err := os.Create(appConfig.RFDRootDirectory + "/index.md")
	CheckFatal(err)

	_, err = mdTableFile.WriteString("**Index of Requests for Discussion**\n\n")
	CheckFatal(err)
	_, err = mdTableFile.WriteString("| **RFD Id** | **Title** | **State** | **Author(s)** |\n")
	CheckFatal(err)
	_, err = mdTableFile.WriteString("|------------|-----------|-----------|------------------------|\n")
	CheckFatal(err)
	return mdTableFile
}

func writeMetadataToTableFile(metaData map[string]interface{}, mdTableFile *os.File, branchID string) {
	title := fmt.Sprintf("%v", metaData["title"])
	authors := fmt.Sprintf("%v", metaData["authors"])
	state := fmt.Sprintf("%v", metaData["state"])

	logger.traceLog(title + ":" + authors + ":" + state)

	_, err := mdTableFile.WriteString("|" + branchID + "|" + title + "|" + state + "|" + authors + "|\n")
	CheckFatal(err)

	logger.traceLog("recorded: " + branchID)
	logger.traceLog("----------------------------------------------")
}

func readMetadataFromReadmeFile(subEntry os.FileInfo, entry os.FileInfo) map[string]interface{} {
	logger.traceLog("Found " + appConfig.RFDRootDirectory + "/" + entry.Name() + "/" + subEntry.Name())
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	file, err := ioutil.ReadFile(appConfig.RFDRootDirectory + "/" + entry.Name() + "/" + subEntry.Name())
	CheckFatal(err)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(file, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	metaData := meta.Get(context)
	return metaData
}

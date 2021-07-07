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

/*
	An RFD can be in one of two branch states:
		1. Newly created, or still being updated and not yet ready for mainlining into the trunk. In this case
		   there won't yet be an separate RFD directory in the trunk.
		2. Trunk - Merged into the trunk.

		Steps:
		1.
*/

func CreateEntries() error {

	logger.traceLog("Creating new entries ...")

	// The entries to be read and written to the
	// markdown table file.
	entries := getDirectories()

	// The file that will be updated with the
	// rfd entries and their status.
	mdTableFile := openMetadataTableFile()
	CheckFatal(mdTableFile.Close())

	for _, entry := range entries {

		entryIsBranchID, err := regexp.MatchString(`(^\d{4}).*`, entry.Name())
		CheckFatal(err)

		if entryIsBranchID {

			logger.traceLog("Matched " + entry.Name())

			branchID := entry.Name()

			if entry.IsDir() {

				subEntries, err := ioutil.ReadDir("./" + entry.Name())
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
	entries, err := ioutil.ReadDir(".")
	CheckFatal(err)
	return entries
}

func openMetadataTableFile() *os.File {
	mdTableFile, err := os.Create("record.md")
	CheckFatal(err)

	_, err = mdTableFile.WriteString("**Record of Merged Requests for Discussion**\n\n")
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

	_, err := mdTableFile.WriteString("|" + branchID + "|" + title + "|" + state + "|" + authors + "|\n")
	CheckFatal(err)

	logger.traceLog("recorded: " + branchID)
	logger.traceLog("----------------------------------------------")
}

func readMetadataFromReadmeFile(subEntry os.FileInfo, entry os.FileInfo) map[string]interface{} {
	logger.traceLog("Found " + subEntry.Name())
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	file, err := ioutil.ReadFile("./" + entry.Name() + "/" + subEntry.Name())
	CheckFatal(err)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(file, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	metaData := meta.Get(context)
	return metaData
}

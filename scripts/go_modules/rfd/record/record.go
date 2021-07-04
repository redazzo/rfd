package record

import (
	"bytes"
	"fmt"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"io/ioutil"
	"os"
	"regexp"
	"rfd.kessellhaak.dev/process/rfd/util"
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

	util.TraceLog("Creating new entries ...")

	// The entries to be read and written to the
	// markdown table file.
	entries := getDirectories()

	// The file that will be updated with the
	// rfd entries and their status.
	mdTableFile := openMetadataTableFile()

	defer util.CheckFatal(mdTableFile.Close())

	for _, entry := range entries {

		entryIsBranchID, err := regexp.MatchString(`(^\d{4}).*`, entry.Name())
		util.CheckFatal(err)

		if entryIsBranchID {

			util.TraceLog("Matched " + entry.Name())

			branchID := entry.Name()

			if entry.IsDir() {

				subEntries, err := ioutil.ReadDir("./" + entry.Name())
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
	entries, err := ioutil.ReadDir(".")
	util.CheckFatal(err)
	return entries
}

func openMetadataTableFile() *os.File {
	mdTableFile, err := os.Create("record.md")
	util.CheckFatal(err)

	_, err = mdTableFile.WriteString("**Record of Merged Requests for Discussion**\n\n")
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

	_, err := mdTableFile.WriteString("|" + branchID + "|" + title + "|" + state + "|" + authors + "|\n")
	util.CheckFatal(err)

	util.TraceLog("recorded: " + branchID)
	util.TraceLog("----------------------------------------------")
}

func readMetadataFromReadmeFile(subEntry os.FileInfo, entry os.FileInfo) map[string]interface{} {
	util.TraceLog("Found " + subEntry.Name())
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	file, err := ioutil.ReadFile("./" + entry.Name() + "/" + subEntry.Name())
	util.CheckFatal(err)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(file, &buf, parser.WithContext(context)); err != nil {
		panic(err)
	}

	metaData := meta.Get(context)
	return metaData
}

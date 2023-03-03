package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	s "strings"
)

func checkIfValidFile(filename string) (bool, error) {
	// Check if file is yaml
	if fileExtension := filepath.Ext(filename); !(fileExtension == ".yaml" || fileExtension == ".yml") {
		return false, fmt.Errorf("File %s is not YAML", filename)
	}

	// Check if file does exist
	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("File %s does not exist", filename)
	}

	return true, nil
}

func filterFingerprintsFromFile(filePath string) []string {
	// open file
	file, err := os.Open(filePath)
	check(err)
	defer func(file *os.File) {
		err := file.Close()
		check(err)
	}(file)

	// parse input into string array
	fileReader := bufio.NewScanner(file)
	fileReader.Split(bufio.ScanLines)

	var content []string
	for fileReader.Scan() {
		content = append(content, fileReader.Text())
	}

	debug(fmt.Sprintf("Read lines from file: %d", len(content)))

	result := make([]string, 0)

	// filter content for fingerprints
	isInPgpBracket := false
	for i := 0; i < len(content); i++ {
		line := content[i]

		// check for pgp bracket start
		if isInPgpBracket == false && s.Contains(line, "pgp:") {
			debug("found pgp bracket")
			isInPgpBracket = true
			continue
		}

		// check for pgp bracket end
		if isInPgpBracket == true && s.Contains(line, ": ") {
			debug("closing pgp bracket")
			isInPgpBracket = false
			continue
		}

		// filter out comment lines
		if isInPgpBracket == true && s.HasPrefix(s.TrimSpace(line), "#") {
			continue
		}

		if isInPgpBracket == true {
			trimmedLine := s.Trim(line, " -")           // remove list operator
			trimmedLine = s.Split(trimmedLine, " #")[0] // remove possible comments after fingerprint
			debug("adding line to file result: " + trimmedLine)
			result = append(result, trimmedLine)
		}
	}

	debug(fmt.Sprintf("Fingerprints in file: %d", len(result)))
	result = removeDuplicates(result)

	debug(fmt.Sprintf("Unique fingerprints: %d", len(result)))
	return result
}

func removeDuplicates[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func createPublicKeyFile(fingerprint string, key string) {
	f, err := os.Create(fingerprint)
	check(err)
	_, err = f.WriteString(key)
	check(err)
	err = f.Close()
	check(err)
}

func removeFile(filePath string) {
	err := os.Remove(filePath)
	check(err)
}

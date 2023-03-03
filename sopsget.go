package main

import (
	"flag"
	"fmt"
)

const DefaultPath string = ".sops.yaml"

var debugEnabled bool

func main() {
	flag.BoolVar(&debugEnabled, "v", false, "enable verbose logging")
	flag.Parse()
	debug("logging on debug level enabled")

	filePath := DefaultPath
	if len(flag.Args()) < 1 {
		fmt.Println("Using default filepath: ", filePath)
	} else {
		filePath = flag.Arg(0)
		fmt.Println("Using specified filepath: ", filePath)
	}

	// read fingerprints from file
	if _, err := checkIfValidFile(filePath); err != nil {
		exitGracefully(err)
	}
	fileFingerprints := filterFingerprintsFromFile(filePath)

	log(fmt.Sprintf("fingerprints from file: %d", len(fileFingerprints)))

	// query gpg
	gpgFingerprints := getLocalFingerprints()

	log(fmt.Sprintf("fingerprints from gpg: %d", len(gpgFingerprints)))

	missingFingerprints := detectUnknownFingerprints(fileFingerprints, gpgFingerprints)

	log(fmt.Sprintf("missing fingerprints: %d", len(missingFingerprints)))

	log("Starting to import missing keys...")
	for i := 0; i < len(missingFingerprints); i++ {
		currentFp := missingFingerprints[i]
		content := fetchFingerprint(currentFp)
		createPublicKeyFile(currentFp, content)
		importKey(currentFp)
		removeFile(currentFp)
		debug(fmt.Sprintf("Imported key %s.", currentFp))
	}

	// check again
	gpgFingerprints = getLocalFingerprints()
	missingFingerprints = detectUnknownFingerprints(fileFingerprints, gpgFingerprints)

	log(fmt.Sprintf("missing fingerprints: %d", len(missingFingerprints)))
}

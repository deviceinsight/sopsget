package main

import (
	"fmt"
	"os/exec"
	"runtime"
	s "strings"
)

func getLocalFingerprints() []string {
	result := make([]string, 0)

	var out []byte
	if runtime.GOOS == "windows" {
		out = gpgListWindows()
	} else {
		out = gpgListLinux()
	}
	content := string(out)
	lines := s.Split(content, "\n")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if s.HasPrefix(line, " ") {
			trimmedLine := s.TrimSpace(line)
			debug("adding line to gpg result: " + trimmedLine)
			result = append(result, trimmedLine)
		}
	}

	debug(fmt.Sprintf("Fingerprints in gpg: %d", len(result)))
	return result
}

func gpgListWindows() []byte {
	out, err := exec.Command("cmd", "/c", "gpg", "--list-keys").Output()
	check(err)
	return out
}

func gpgListLinux() []byte {
	out, err := exec.Command("bash", "-c", "gpg", "--list-keys").Output()
	check(err)
	return out
}

func detectUnknownFingerprints(requestedFp []string, presentFp []string) []string {
	result := make([]string, 0)

	for i := 0; i < len(requestedFp); i++ {
		found := false
		for k := 0; k < len(presentFp); k++ {
			if requestedFp[i] == presentFp[k] {
				found = true
				break
			}
		}
		if !found {
			result = append(result, requestedFp[i])
		}
	}

	return result
}

func importKey(fileName string) {
	var output []byte
	if runtime.GOOS == "windows" {
		output = gpgImportWindows(fileName)
	} else {
		output = gpgImportLinux(fileName)
	}
	debug(string(output))
}

func gpgImportWindows(fileName string) []byte {
	out, err := exec.Command("cmd", "/c", "gpg", "--import", fileName).Output()
	check(err)
	return out
}

func gpgImportLinux(fileName string) []byte {
	out, err := exec.Command("bash", "-c", "gpg", "--import", fileName).Output()
	check(err)
	return out
}

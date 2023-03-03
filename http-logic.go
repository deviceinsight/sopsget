package main

import (
	"fmt"
	"io"
	"net/http"
)

func fetchFingerprint(fingerprint string) string {
	resp, err := http.Get(fmt.Sprintf("https://keys.openpgp.org/vks/v1/by-fingerprint/%s", fingerprint))
	check(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	check(err)

	return string(body)
}

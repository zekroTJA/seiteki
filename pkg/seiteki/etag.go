package seiteki

import (
	"crypto/sha1"
	"fmt"
)

// getETag returns an ETag header content for the
// passed body content in form of a SHA1
// hex string.
// If weak is set to true, the ETag will be
// prefixed with W/ to identify a weak ETag.
func getETag(body []byte, weak bool) string {
	hash := sha1.Sum(body)

	weakTag := ""
	if weak {
		weakTag = "W/"
	}

	tag := fmt.Sprintf("%s\"%x\"", weakTag, hash)

	return tag
}

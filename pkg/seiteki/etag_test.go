package seiteki

import "testing"

func TestGetETag(t *testing.T) {
	testBody := []byte("This is a test body.")
	const controlETagStrong = `"2c1a253771877942941e797d4907366001204c29"`
	const controlETagWeak = `W/"2c1a253771877942941e797d4907366001204c29"`

	etag := getETag(testBody, false)
	if etag != controlETagStrong {
		t.Errorf("strong etag was '%s' (expected: '%s')",
			etag, controlETagStrong)
	}

	etag = getETag(testBody, true)
	if etag != controlETagWeak {
		t.Errorf("weak etag was '%s' (expected: '%s')",
			etag, controlETagWeak)
	}
}

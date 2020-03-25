package seiteki

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func getConfig() *Config {
	return &Config{
		Addr:          "localhost:8080",
		CacheDuration: "720h",
		CertFile:      "",
		KeyFile:       "",
		IndexFile:     "index.html",
		StaticDir:     "testdata",
	}
}

func TestNew(t *testing.T) {
	cfg := getConfig()
	s, err := New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if s == nil {
		t.Fatal("new instance was nil")
	}
	if s.cacheHeader != "max-age=2592000, public" {
		t.Errorf("cache-control header had unexpected value: %s",
			s.cacheHeader)
	}

	cfg.CacheDuration = "invalid"
	s, err = New(cfg)
	if err == nil {
		t.Errorf("invalid CacheDuration in config returned no error")
	}
}

func TestListenAndServe(t *testing.T) {
	cfg := getConfig()
	cfg.StaticDir = os.Getenv("STK_TEST_WEBDIR")

	contIndex, err := ioutil.ReadFile(cfg.StaticDir + "/index.html")
	if err != nil {
		t.Error(err)
	}

	contJs, err := ioutil.ReadFile(cfg.StaticDir + "/test.js")
	if err != nil {
		t.Error(err)
	}

	s, err := New(cfg)
	if err != nil {
		t.Error(err)
	}

	go func() {
		err := s.ListenAndServeBlocking()
		if err != nil {
			t.Fatal(err)
		}
	}()

	res, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Errorf("req http://localhost:8080 failed: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("req http://localhost:8080 failed: status code was %d",
			res.StatusCode)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(body, contIndex) != 0 {
		t.Errorf("req http://localhost:8080 failed: wrong body content")
	}

	res, err = http.Get("http://localhost:8080/test")
	if err != nil {
		t.Errorf("req http://localhost:8080/test failed: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("req http://localhost:8080/test failed: status code was %d",
			res.StatusCode)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(body, contIndex) != 0 {
		t.Errorf("req http://localhost:8080/test failed: wrong body content")
	}

	res, err = http.Get("http://localhost:8080/test.js")
	if err != nil {
		t.Errorf("req http://localhost:8080/test.js failed: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("req http://localhost:8080/test.js failed: status code was %d",
			res.StatusCode)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(body, contJs) != 0 {
		t.Errorf("req http://localhost:8080/test.js failed: wrong body content")
	}
}

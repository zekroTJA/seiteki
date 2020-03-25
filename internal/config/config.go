package config

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/zekroTJA/seiteki/pkg/seiteki"
)

// configFile defines the default config file location
const configFile = "/etc/seiteki/config.json"

// Get returns a merged seiteki.Config instance
// with configuration values collected from
// flags, environment variables and default
// configuration file.
func Get() (*seiteki.Config, error) {
	c := getFromFlags()

	cEnv, err := getFromEnv()
	if err != nil {
		return nil, err
	}
	merge(cEnv, c)

	cFile, err := getFromFile()
	if err != nil {
		return nil, err
	}
	merge(cFile, c)

	return c, nil
}

// getFromFlags crteates a seiteki.Config instance
// with data passed by flags.
// This is used as base config instance because defaults
// are set by non-set flag values.
func getFromFlags() *seiteki.Config {
	addr := flag.String("addr", "localhost:80", "expose address and port")
	cacheDuration := flag.String("cd", "720h", "cache duration (for time format see https://golang.org/pkg/time/#ParseDuration)")
	certFile := flag.String("cert", "", "ssl cert file location")
	compress := flag.Bool("compress", false, "whether or not to gzip compress static files")
	indexFile := flag.String("index", "index.html", "default index file location")
	keyFile := flag.String("key", "", "ssl key file location")
	staticDir := flag.String("dir", "web", "static file location")

	flag.Parse()

	return &seiteki.Config{
		Addr:          *addr,
		CacheDuration: *cacheDuration,
		CertFile:      *certFile,
		Compress:      *compress,
		IndexFile:     *indexFile,
		KeyFile:       *keyFile,
		StaticDir:     *staticDir,
	}
}

// getFromEnv creates a seiteki.Config instance
// from environment variable configuration
// prefixed with 'STK_'.
func getFromEnv() (*seiteki.Config, error) {
	c := new(seiteki.Config)
	return c, envconfig.Process("STK", c)
}

// getFromFile tries to read the content of the
// default config file and creates a seiteki.Config
// instance if the file exists.
// If IO errors occur, this will be returned as
// error.
func getFromFile() (*seiteki.Config, error) {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	f, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := new(seiteki.Config)
	dec := json.NewDecoder(f)
	return c, dec.Decode(c)
}

// merge sets or replaces data in target fields
// if the fields are not empty in source.
func merge(source, target *seiteki.Config) {
	if source == nil {
		return
	}

	if source.Addr != "" {
		target.Addr = source.Addr
	}
	if source.CacheDuration != "" {
		target.CacheDuration = source.CacheDuration
	}
	if source.CertFile != "" {
		target.CertFile = source.CertFile
	}
	if source.Compress {
		target.Compress = source.Compress
	}
	if source.IndexFile != "" {
		target.IndexFile = source.IndexFile
	}
	if source.KeyFile != "" {
		target.KeyFile = source.KeyFile
	}
	if source.StaticDir != "" {
		target.StaticDir = source.StaticDir
	}
}

// seiteki wraps around fasthttp.Server and fasthttp.FS
// to provide a simple to use static host serve for
// routed web applications like Angular or VueJS apps.
//
// Specified source files are served via fasthttp.FS
// and every other request will be directed directly
// onto the specified index file.
package seiteki

import (
	"path"
	"regexp"
	"time"

	"github.com/valyala/fasthttp"
)

// FileRx defines the regular expression used to
// match static files by file extension. If you really
// want to, just overwrite this variable to replace it.
var FileRx = regexp.MustCompile(`^.*\.(ico|css|js|svg|gif|jpe?g|png)$`)

// Config includes everything which defines
// preferences and parameters to configure
// the server.
type Config struct {
	// Expose address and port
	Addr string `json:"addr"`
	// Duration how long static files will be cached
	// by setting cache control headers.
	// See here how to format this:
	// https://golang.org/pkg/time/#ParseDuration
	CacheDuration string `json:"cacheduration"`
	// Whether or not to compress static files
	Compress bool `json:"compress"`
	// Static files location to serve from
	StaticDir string `json:"staticdir"`
	// Default index file name
	IndexFile string `json:"indexfile"`

	// SSL key file directory
	KeyFile string `json:"keyfile"`
	// SSL cert file directory
	CertFile string `json:"certfile"`
}

// Seiteki web server
type Seiteki struct {
	config *Config

	fs        *fasthttp.FS
	s         *fasthttp.Server
	fsHandler fasthttp.RequestHandler
}

// New creates a new instance of Seiteki with the
// passed configuration.
func New(config *Config) (*Seiteki, error) {
	server := &Seiteki{
		config: config,
	}

	server.s = &fasthttp.Server{
		Handler: server.requestHandler,
	}

	cacheDur, err := time.ParseDuration(config.CacheDuration)
	if err != nil {
		return nil, err
	}

	server.fs = &fasthttp.FS{
		CacheDuration: cacheDur,
		Compress:      config.Compress,
		Root:          config.StaticDir,
		IndexNames:    []string{config.IndexFile},
	}

	server.fsHandler = server.fs.NewRequestHandler()

	return server, nil
}

// ListenAndServeBlocking blocks the current go routine
// and starts the listening and serving routines.
// Depending on if the SSL config was passed properly,
// the server will be started using SSL, else, it
// will automatically use non-SSL setup.
func (server *Seiteki) ListenAndServeBlocking() error {
	useSSL := server.config.CertFile != "" &&
		server.config.KeyFile != ""

	if useSSL {
		return server.s.ListenAndServeTLS(
			server.config.Addr,
			server.config.CertFile,
			server.config.KeyFile)
	}

	return server.s.ListenAndServe(server.config.Addr)
}

// requestHandler checks if the request destination is a
// file or a web route. If it is a file, serve the file
// via FS handler, else serve the "index.html" file.
func (server *Seiteki) requestHandler(ctx *fasthttp.RequestCtx) {
	if FileRx.Match(ctx.Path()) {
		server.fsHandler(ctx)
	} else {
		ctx.SendFile(
			path.Join(server.fs.Root, server.config.IndexFile))
	}
}

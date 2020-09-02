// Package seiteki wraps around fasthttp.Server and
// fasthttp.FS to provide a simple to use static host
// serve for routed web applications like Angular or
// VueJS apps.
//
// Specified source files are served via fasthttp.FS
// and every other request will be directed directly
// onto the specified index file.
package seiteki

import (
	"fmt"
	"math"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/valyala/fasthttp"
)

// RouteMode is the mode type used to determine if the
// accessed request path is a file or SPA route.
type RouteMode string

const (
	// Route mode using regular expression on req path.
	//
	// This is useful when you exacly know which file extensions
	// are served as static files and that they are included in
	// the match regex. This mode is especially verry fast and
	// has a less performance footprint because no file ops are
	// executed on path check.
	RouteModeRegex RouteMode = "regex"

	// Route mode using fs stat on accessed file.
	//
	// This mode can be used if you don't know exacly which file
	// extensions are served as static files or if these extensions
	// are not included in the match regex. Also this is useful if
	// you may have path missmatches for some reasons. This mode only
	// serves a file when it is actually existent on the fs, otherwise
	// the index file will be served. But keep in mind that this mode
	// executes a fs stat operation on EVERY request, which might lead
	// to poor performance.
	RouteModeStat RouteMode = "stat"

	// Route mode which bypasses the SPA matcher and tries to serve
	// all request paths as static file.
	//
	// This is especially useful if you want to use seiteki as pure
	// static file server without SPA routing.
	RouteModeStatic RouteMode = "static"
)

// FileRx defines the regular expression used to
// match static files by file extension. If you really
// want to, just overwrite this variable to replace it.
var FileRx = regexp.MustCompile(`^.*\.(ico|css|js|svg|gif|jpe?g|png|html?)$`)

const (
	cacheControlHeader = "cache-control"
	etagHeader         = "etag"
)

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

	// RouteMode used for determining if the
	// requested path is a static file or
	// SPA route. Default is "regex".
	RouteMode RouteMode `json:"routemode"`

	// SSL key file directory
	KeyFile string `json:"keyfile"`
	// SSL cert file directory
	CertFile string `json:"certfile"`
}

// Seiteki web server
type Seiteki struct {
	config      *Config
	logger      Logger
	cacheHeader string

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
		Handler: server.RequestHandler,
	}

	cacheDur, err := time.ParseDuration(config.CacheDuration)
	if err != nil {
		return nil, err
	}

	server.cacheHeader = fmt.Sprintf("max-age=%.0f, public",
		math.Floor(cacheDur.Seconds()))

	server.fs = &fasthttp.FS{
		Compress:   config.Compress,
		Root:       config.StaticDir,
		IndexNames: []string{config.IndexFile},
	}

	server.fsHandler = server.fs.NewRequestHandler()

	server.logger = newLogegrWrapper(nil)

	return server, nil
}

// SetLogger sets a logger interface as
// request logger
func (server *Seiteki) SetLogger(logger Logger) {
	server.logger = newLogegrWrapper(logger)
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

func (server *Seiteki) IsStaticFile(path []byte) bool {
	switch server.config.RouteMode {
	case RouteModeRegex:
		return server.isStaticRegex(path)
	case RouteModeStat:
		return server.isStaticStat(path)
	case RouteModeStatic:
		return true
	default:
		return server.isStaticRegex(path)
	}
}

// RequestHandler checks if the request destination is a
// file or a web route. If it is a file, serve the file
// via FS handler, else serve the "index.html" file.
func (server *Seiteki) RequestHandler(ctx *fasthttp.RequestCtx) {
	const serverHeader = "seiteki/" + Version

	reqPath := ctx.Path()

	ctx.Response.Header.Set(cacheControlHeader, server.cacheHeader)
	ctx.Response.Header.SetServer(serverHeader)

	if server.IsStaticFile(reqPath) {
		server.fsHandler(ctx)
		server.logRequest(ctx, reqPath, reqPath, server.fs.Root)
	} else {
		ctx.SendFile(path.Join(server.fs.Root, server.config.IndexFile))
		server.logRequest(ctx, reqPath, server.config.IndexFile, server.fs.Root)
	}

	if ctx.Response.StatusCode() == fasthttp.StatusOK {
		etag := getETag(ctx.Response.Body(), false)
		ctx.Response.Header.Set(etagHeader, etag)
	}
}

// isStaticRegex uses the FileRx regular expression to determine
// if the passed reuqest path is a static file.
func (server *Seiteki) isStaticRegex(reqPath []byte) bool {
	return FileRx.Match(reqPath)
}

// isStaticStat uses an fs stat operation on the given path to
// determine if the passed reuqest path is a static file.
func (server *Seiteki) isStaticStat(reqPath []byte) bool {
	dir := path.Join(server.config.StaticDir, string(reqPath))
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		server.logger.Error("file stat failed: ", err)
		return false
	}

	return !info.IsDir()
}

// logRequest loggs an incomming requests remote address,
// request path and the file path which will be sent in the
// body of the response as same as the response status code.
func (server *Seiteki) logRequest(ctx *fasthttp.RequestCtx, reqPath []byte, resPath interface{}, fsRoot string) {
	server.logger.Infof("REQ [%s] %s -> %s%s [%d]",
		ctx.RemoteAddr().String(), reqPath, fsRoot,
		resPath, ctx.Response.StatusCode())
}

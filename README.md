# seiteki
[![](https://godoc.org/github.com/zekroTJA/seiteki/pkg/seiteki?status.svg)](https://godoc.org/github.com/zekroTJA/seiteki/pkg/seiteki) &nbsp; [![](https://github.com/zekroTJA/seiteki/workflows/Main%20CI/badge.svg)](https://github.com/zekroTJA/seiteki/actions) &nbsp; [![Go Report Card](https://goreportcard.com/badge/github.com/zekroTJA/seiteki)](https://goreportcard.com/report/github.com/zekroTJA/seiteki)

seiteki *(静的, japanese for 'static')* is a package and application which wraps around `fasthttp.Server` and `fasthttp.FS` to provide a static delivery server for routed web applications like Angular or VueJS apps. This application is easy to configure via command flags, environment variables or via a config file. So, this application is especially concipated to be used in container images.

## How does it work?

If you want to host a routed application like created with Angular or VueJS, you want to have all routes *(like `example.com/about`)* redirected to the `index.html` and let the MVVM engine control the route. On the other side, all static files *(like `example.com/js/app.js` or `example.com/css/app.css`)* should be served from this exact path in the web file directory.

So, seiteki simply checks the path for the following regular expression to check if the request path is a route or a file by checking the existence of a defined file extension:
```
^.*\.(ico|css|js|svg|gif|jpe?g|png)$
```

Then, some chache control headers are added for better cahcing and yeah, that's all the magic.

## Using as package

You can also use seiteki as package for your own application. Just download the package from pkg/seiteki:

```
$ go get github.com/zekroTJA/seiteki/pkg/seiteki
```

Example usage:

```go
package main

import (
    "github.com/zekroTJA/seiteki/pkg/seiteki"
)

func main() {
    s, err := seiteki.New(&seiteki.Config{
        Addr:          "localhost:80",
		CacheDuration: "720h",
		IndexFile:     "index.html",
		StaticDir:     "web",
    })
    if err != nil {
        panic(err)
    }

    err = s.ListenAndServeBlocking()
    if err != nil {
        panic(err)
    }
}
```

Further usage examples can be found in the text file [`seiteki_test.go`](pkg/seiteki/seiteki_test.go). Another production example is [my web page](https://github.com/zekroTJA/zekro-de-rewrite) which actually uses seiteki for serving the VueJS web application inside a Docker container.

## Configuration

The seiteki application uses a self written [config wrapper](internal/config/config.go) which merges configuration via flags, environemnt variables and the config file located in `/etc/seiteki/config.json`. The config values are overwritten in this exact order:
- command flags
- environment variables
- config file

### Flag Configuration

You can enter `go run cmd/seiteki/main.go -h` to get a list of configuration flags:
```
Usage of seiteki:
  -addr string
        expose address and port (default "localhost:80")
  -cd string
        cache duration (for time format see https://golang.org/pkg/time/#ParseDuration) (default "720h")
  -cert string
        ssl cert file location
  -compress
        whether or not to gzip compress static files
  -dir string
        static file location (default "web")
  -index string
        default index file location (default "index.html")
  -key string
        ssl key file location
```

### Environment Variable Configuration

Config env vars need to have the prefix `STK_`.

```bash
export STK_ADDR="localhost:80"
export STK_CACHEDURATION="720h"
export STK_CERTFILE="/etc/cert/cert.pem"
export STK_KEYFILE="/etc/cert/key.pem"
export STK_COMPRESS="true"
export STK_INDEXFILE="index.html"
export STK_STATICDIR="$PWD/web"
```

### File Configuration

The app checks if the config file exists under the location `/etc/seiteki/config.json` which looks like following:
```json
{
    "addr": "localhost:8080",
    "cacheduration": "720h",
    "compress": true,
    "staticdir": "web",
    "indexfile": "index.html",
    "keyfile": "/etc/cert/key.pem",
    "certfile": "/etc/cert/cert.pem"
}
```

---

© 2020 Ringo Hoffmann (zekro Development)  
Covered by the MIT Licence.
# alchemy-go-service-base

Provides an opinionated set of default configurations for a web service written in golang.

## Features

* [bunyan](https://github.com/trentm/node-bunyan) style logging facilitated by [logrus](https://github.com/sirupsen/logrus)
* access logging
* custom request specific logger provided via middleware
* instrumentation for [prometheus](https://prometheus.io/) using [the golang client](https://github.com/prometheus/client_golang)
* HTTP error definition / middleware

## Installation

With gvt:

```
gvt fetch github.com:pgalchemy/alchemy-go-service-base
```

With go:

```
go get github.com:pgalchemy/alchemy-go-service-base
```

## Development

If using gvt for development, the simplest (yet inconvenient) way to test library changes in a service is via rsync:

```
rsync -azP --delete ../alchemy-go-service-base/* vendor/github.com/pgalchemy/alchemy-go-service-base
```

## Usage

```golang

package main

import (
  "github.com/pgalchemy/alchemy-go-service-base/instrumentation"
  "github.com/pgalchemy/alchemy-go-service-base/logging"
  "github.com/pgalchemy/alchemy-go-service-base/scope"
  "github.com/pgalchemy/alchemy-go-service-base/server"
)

func main() {
  // Setup logger
  formatter := logging.NewFormatter("my-service")
  logger := logging.New(formatter)

  // Create server
  r := server.New(&server.Config{Logger: logger})

  // Setup instrumentation
  instrumentation.Serve(r)

  r.Get("/", func(c *gin.Context) {
    rs := scope.GetRequestScope(c)
    rs.Logger.Debug("I can log things that will be linked to this specific request!")
    c.JSON(200, gin.H{"success": true})
  })

  r.Run(":8000")
}
```

# http server


## How to start a simple http server
```
package main

import (
	"context"
	"log/slog"

	"github.com/devshansharma/tools/logger"
	"github.com/devshansharma/tools/server"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	log := logger.New(
		logger.WithJSON(true),
		logger.WithSource(true),
		logger.WithLevel("INFO"),
		logger.WithReplaceAttr(logger.WithShortFileNameAndErrorTrace),
	)
	slog.SetDefault(log)

	router := gin.New()
	srv := server.New(":8080", router)

	slog.InfoContext(ctx, "starting server")
	srv.Run(ctx)
}
```
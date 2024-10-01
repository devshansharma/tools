package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/devshansharma/tools/crypt"
	"github.com/devshansharma/tools/logger"
	"github.com/devshansharma/tools/server"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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

	router.GET("/", func(ctx *gin.Context) {
		slog.InfoContext(ctx.Request.Context(), "URL got hit")

		privateKey, err := crypt.GenerateES512PrivateKey()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		access, refresh, err := crypt.GenerateJWTTokens(privateKey, jwt.MapClaims{
			"jti": uuid.NewString(),
			"sub": uuid.NewString(),
			"aud": "test",
			"iss": "test",
			"exp": time.Now().Add(time.Hour * 24),
			"iat": time.Now(),
			"nbf": time.Now(),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access":  access,
			"refresh": refresh,
		})
	})

	srv := server.New(":8080", router,
		server.WithServerTimeout(11),
	)

	slog.InfoContext(ctx, "starting server on port: 8080")
	srv.Run(ctx)
}

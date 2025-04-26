package server

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	sloggin "github.com/samber/slog-gin"

	"github.com/etboye/calculator/api"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type GinServer struct {
}

func (ginServer *GinServer) StartServer(endpoints api.Endpoints) error {
	router := gin.Default()

	setupLogging(router)
	registerEndpoints(router, endpoints)

	return router.Run()
}

func setupLogging(router *gin.Engine) {
	config := sloggin.Config{
		WithRequestBody:    true,
		WithRequestHeader:  true,
		WithResponseBody:   true,
		WithResponseHeader: true,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	router.Use(sloggin.NewWithConfig(logger, config))
}

func registerEndpoints(router *gin.Engine, endpoints api.Endpoints) {
	router.POST("/sessions/:sessionid/compute", func(c *gin.Context) {
		computeRequest := api.ComputeRequest{}

		if unmarshallError := c.ShouldBindBodyWith(&computeRequest, binding.JSON); unmarshallError != nil {
			bodyAsString := getBodyAsString(c.Request.Body)
			log.Println("Could not marshall request", bodyAsString) // TODO: Test this
			c.Status(http.StatusBadRequest)
			return
		}

		sessionId := c.Param("sessionid") // TODO: Test empty
		computeResponse := endpoints.ComputationHandler.Compute(sessionId, computeRequest)
		sendSimpleHttpResponse(c, computeResponse)
	})

	router.GET("/sessions/:sessionid/history", func(c *gin.Context) {
		sessionId := c.Param("sessionid") // TODO: Test empty
		cursor := c.Query("cursor")
		response := endpoints.SessionHistoryHandler.GetSessionHistory(sessionId, cursor)
		sendSimpleHttpResponse(c, response)
	})
}

func getBodyAsString(body io.ReadCloser) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	return buf.String()
}

func sendSimpleHttpResponse[T any](c *gin.Context, httpResponse api.SimpleHttpResponse[T]) {
	c.JSON(httpResponse.Status, httpResponse.Response)
}

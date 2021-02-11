package httpServer

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"obfuscator/obfuscating"
)

const (
	processIdParam = "processId"
)

func obfuscatorRouter(router gin.RouterGroup) {
	router.POST("/schema-info", getSchemaInfo)

	router.POST("/obfuscate", obfuscate)

	router.GET("/status/:"+processIdParam, getProcessStatus)

	router.POST("/empty-progress-ctx", emptyProgressCtx)
}

func obfuscate(c *gin.Context) {
	var request obfuscating.ObfuscateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	err := obfuscating.ValidateObfuscationModel(request.Model, request.Origin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	processId, err := obfuscating.InitProcess(len(request.Model))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	go obfuscating.ObfuscateSchema(request.Model, processId, request.Origin, request.Destination)
	c.JSON(http.StatusOK, ObfuscationResponse{
		SuccessfulResponse: SuccessfulResponse{
			"Obfuscation was started.",
		},
		ProcessId: processId,
	})
}

func getSchemaInfo(c *gin.Context) {
	var request obfuscating.ConnectionInfo
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	result, err := obfuscating.GetSchemaInfo(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

func getProcessStatus(c *gin.Context) {
	processId := c.Param(processIdParam)
	result, exists := obfuscating.GetProcessCtx(processId)
	if !exists {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Entry with this process id doesn't exist",
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

func emptyProgressCtx(c *gin.Context) {
	obfuscating.EmptyProgressCtx()
	c.JSON(http.StatusOK, SuccessfulResponse{
		Status: "OK",
	})
}

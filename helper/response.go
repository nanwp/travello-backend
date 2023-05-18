package helper

import "github.com/gin-gonic/gin"

type Response struct {
	Status  string      `json:"status,omitempty"`
	Code    int32       `json:"code,omitempty"`
	Message interface{} `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func ResponseOutput(c *gin.Context, statusCode int32, statusMessage string, message interface{}, data interface{}) {
	resp := Response{
		Status:  statusMessage,
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	c.JSON(int(statusCode), resp)
}

package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nanwp/travello/helper"
	"github.com/nanwp/travello/models/ulasans"
	"github.com/nanwp/travello/pkg/middleware/auth"
	"github.com/nanwp/travello/service"
)

type ulasanHandler struct {
	urlApi        string
	ulasanService service.UlasanService
}

func NewUlasanHandler(service service.UlasanService) *ulasanHandler {
	return &ulasanHandler{"https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/ulasan", service}
}

func (h *ulasanHandler) AddUlasan(c *gin.Context) {
	var ulasanBody ulasans.RequestUlasan
	err := c.ShouldBindJSON(&ulasanBody)

	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	uuidGenerate := uuid.New().String()
	ulasan := ulasans.Ulasan{
		ID:            uuidGenerate,
		UserId:        auth.UserID,
		DestinationId: ulasanBody.DestinationId,
		Message:       ulasanBody.Message,
		Rating:        ulasanBody.Rating,
	}

	ulasanResp, err := h.ulasanService.Create(ulasan)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	respCreate := gin.H{
		"destination_id": ulasanResp.DestinationId,
		"user_id":        ulasanResp.UserId,
		"message":        ulasanResp.Message,
		"rating":         ulasanResp.Rating,
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", "berhasil menambahkan ulasan", respCreate)
}

func (h ulasanHandler) GetUlasanByDestination(c *gin.Context) {

	idDestination := c.Query("destination")

	data, err := h.ulasanService.GetUlasanByDestinationID(idDestination)
	log.Println(err)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err, nil)
		return
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", gin.H{"message": "berhasil mendapat ulasan", "count": len(data)}, data)

}

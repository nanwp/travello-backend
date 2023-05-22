package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	ulasan := ulasans.Ulasan{
		UserId:        auth.UserID,
		DestinationId: ulasanBody.DestinationId,
		Message:       ulasanBody.Message,
		Rating:        ulasanBody.Rating,
	}

	ulasanReq, err := json.Marshal(ulasan)

	//mengambil data ulasan yang sudah ada

	dataUlasan, err := http.Get(h.urlApi + "?destination=" + ulasanBody.DestinationId)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	defer dataUlasan.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(dataUlasan.Body)

	var arrayUlasan []ulasans.Ulasan
	json.Unmarshal(bodyBytes, &arrayUlasan)

	//validasi ketika sudah memberikan ulasan
	for _, u := range arrayUlasan {
		fmt.Println(u)
		if u.UserId == auth.UserID {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "hanya bisa 1 kali", nil)
			return
		}
	}

	//mengirim data ulasan
	ulasanResp, err := http.Post(h.urlApi, "application/json; charset=utf-8", bytes.NewBuffer(ulasanReq))
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	defer ulasanResp.Body.Close()

	//menghitung total ulasan
	var totalUlasan = ulasan.Rating

	//menambah total ulasan
	for _, u := range arrayUlasan {
		totalUlasan += u.Rating
	}

	//update ulasan destinasi

	type bodyUpdate struct {
		Rating float32 `json:"rating"`
	}

	updateDest := bodyUpdate{
		totalUlasan / float32(len(arrayUlasan)+1),
	}

	jsonReqDest, err := json.Marshal(updateDest)

	//mengupdate rating pada destinasi
	reqDestinationUpdate, err := http.NewRequest(http.MethodPut, "https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/destination?id="+ulasanBody.DestinationId, bytes.NewBuffer(jsonReqDest))
	reqDestinationUpdate.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(reqDestinationUpdate)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	helper.ResponseOutput(c, http.StatusOK, "OK", "berhasil menambahkan ulasan", nil)
}

func (h ulasanHandler) GetUlasanByDestination(c *gin.Context) {

	idDestination := c.Query("destination")

	data, count, err := h.ulasanService.Get(idDestination)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err, nil)
		return
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", gin.H{"message": "berhasil mendapat ulasan", "count": count}, data)

}

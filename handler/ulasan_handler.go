package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/helper"
	"github.com/nanwp/travello/middleware"
	"github.com/nanwp/travello/models/destinations"
	"github.com/nanwp/travello/models/ulasans"
	"github.com/nanwp/travello/repository"
	"github.com/nanwp/travello/service"
)

type ulasanHandler struct {
	urlApi string
}

func NewUlasanHandler() *ulasanHandler {
	return &ulasanHandler{"https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/ulasan"}
}

func (h *ulasanHandler) AddUlasan(c *gin.Context) {
	var ulasanBody ulasans.RequestUlasan
	err := c.ShouldBindJSON(&ulasanBody)

	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	ulasan := ulasans.Ulasan{
		UserId:        middleware.UserID,
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
		if u.UserId == middleware.UserID {
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

	resp, err := http.Get(h.urlApi + "?destination=" + idDestination)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyString := string(bodyBytes)

	var arrayUlasan []ulasans.Ulasan

	json.Unmarshal([]byte(bodyString), &arrayUlasan)

	if len(arrayUlasan) == 0 {
		helper.ResponseOutput(c, http.StatusNotFound, "NOT_FOUND", "belum ada ulasan", nil)
	}

	respDest, err := http.Get("https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/destination" + "?id=" + idDestination)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}
	defer respDest.Body.Close()

	dataDest, _ := ioutil.ReadAll(respDest.Body)

	var destination destinations.Destination
	json.Unmarshal(dataDest, &destination)

	if destination.ID == "" {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "data not found", nil)
		return
	}

	responseUlasans := []ulasans.ResponseUlasan{}

	for _, u := range arrayUlasan {
		userService := service.NewUserService(repository.NewUserRepository(config.ConnectDatabase()))
		user, err := userService.FindByID(u.UserId)
		if err != nil {
			panic(err)
		}
		responseUlasan := ulasans.ResponseUlasan{
			UserName: user.Name,
			Message:  u.Message,
			Rating:   u.Rating,
		}
		responseUlasans = append(responseUlasans, responseUlasan)
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", "success", gin.H{"destination": destination, "ulasan": responseUlasans})

}

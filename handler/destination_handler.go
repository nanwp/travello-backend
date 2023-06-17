package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/nanwp/travello/helper"
	"github.com/nanwp/travello/models/destinations"
	"github.com/nanwp/travello/pkg/middleware/auth"
	"github.com/nanwp/travello/service"
)

type destinatinHandler struct {
	urlApi      string
	uService    service.UlasanService
	userService service.UserService
}

func NewDestinationHandler(uService service.UlasanService, userService service.UserService) *destinatinHandler {
	return &destinatinHandler{"https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/destination", uService, userService}
}

func (h *destinatinHandler) Destination(c *gin.Context) {
	destinationId := c.Param("id")

	if destinationId == "" {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "Data not found", destinations.DestinationResponse{})
		return
	}

	response, err := http.Get(h.urlApi + "?id=" + destinationId)

	code := response.StatusCode

	if code != 200 {
		helper.ResponseOutput(c, int32(code), response.Status, nil, nil)
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	hasil := destinations.Destination{}
	json.Unmarshal(responseData, &hasil)

	ulas, err := h.uService.GetUlasanByDestinationID(hasil.ID)
	if err != nil {
		log.Printf("error message : %v", err.Error())
	}

	respData := convertDataToResponse(hasil)

	if ulas != nil {

		respData.Ulasan = ulas
		respData.CountUlasan = len(ulas)
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", respData)
}

func (h *destinatinHandler) Destinations(c *gin.Context) {

	category := c.Query("category")
	search := c.Query("search")
	limit := c.Query("limit")
	orderBy := c.Query("orderby")

	urlApi := h.urlApi

	if search != "" {
		urlApi += "?search=" + url.QueryEscape(search)
		if category != "" {
			urlApi += "&category=" + category
		}
	} else if category != "" {
		urlApi += "?category=" + category
	}

	response, err := http.Get(urlApi)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
		return
	}

	var hasil []destinations.Destination

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&hasil); err != nil {
		log.Printf("error get data from mongodb %v", err.Error())
	}

	var data []destinations.DestinationResponse

	for _, d := range hasil {
		data = append(data, convertDataToResponse(d))
	}

	var wg sync.WaitGroup
	for i := 0; i < len(data); i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			ulasanCount, err := h.uService.GetCountUlasanByDestinationID(data[index].ID)
			if err != nil {
				log.Printf("error message : %v", err.Error())
			}
			data[index].CountUlasan = int(ulasanCount)
		}(i)
	}
	wg.Wait()

	if orderBy != "" {
		switch orderBy {
		case "popular":
			sort.Slice(data, func(i, j int) bool {
				return data[j].CountUlasan < data[i].CountUlasan
			})
		case "highRating":
			sort.Slice(data, func(i, j int) bool {
				return data[j].Rating < data[i].Rating
			})
		default:
			sort.Slice(data, func(i, j int) bool {
				return data[j].UpdatedAt < data[i].UpdatedAt
			})
		}

	} else {
		sort.Slice(data, func(i, j int) bool {
			return data[j].UpdatedAt < data[i].UpdatedAt
		})
	}

	if limit != "" {
		limitInt, _ := strconv.Atoi(limit)
		if limitInt != 0 {
			limitHasil := data[:limitInt]
			helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", limitHasil)
			return
		}
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", data)

}

func (h *destinatinHandler) Create(c *gin.Context) {

	var createBody destinations.DestinationCreate
	err := c.ShouldBindJSON(&createBody)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	userLogin, _ := h.userService.FindByID(auth.UserID)

	if userLogin.Role != "admin" {
		helper.ResponseOutput(c, http.StatusUnauthorized, "Unauthorized", "anda bukan admin", nil)
		return
	}

	jsonReqDest, err := json.Marshal(createBody)

	reqDestinationCreate, err := http.NewRequest(http.MethodPost, h.urlApi, bytes.NewBuffer(jsonReqDest))
	reqDestinationCreate.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(reqDestinationCreate)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	helper.ResponseOutput(c, http.StatusOK, "OK", "success create data", createBody)
}

func (h *destinatinHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userLogin, _ := h.userService.FindByID(auth.UserID)

	if userLogin.Role != "admin" {
		helper.ResponseOutput(c, http.StatusUnauthorized, "Unauthorized", "anda bukan admin", nil)
		return
	}

	reqDeleteData, _ := http.NewRequest(http.MethodDelete, h.urlApi+"?id="+id, nil)
	reqDeleteData.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}

	resp, err := client.Do(reqDeleteData)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "data tidak ada", nil)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "data tidak ada", nil)
		return
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", "success delete data", nil)

}

func (h *destinatinHandler) Update(c *gin.Context) {
	var updateBody destinations.DestinationUpdate
	id := c.Param("id")
	err := c.ShouldBindJSON(&updateBody)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	userLogin, _ := h.userService.FindByID(auth.UserID)

	if userLogin.Role != "admin" {
		helper.ResponseOutput(c, http.StatusUnauthorized, "Unauthorized", "anda bukan admin", nil)
		return
	}

	jsonReqDest, err := json.Marshal(updateBody)

	reqDestinationUpdate, err := http.NewRequest(http.MethodPut, h.urlApi+"?id="+id, bytes.NewBuffer(jsonReqDest))
	reqDestinationUpdate.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(reqDestinationUpdate)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "data tidak ada", nil)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", "data tidak ada", nil)
		return
	}

	helper.ResponseOutput(c, http.StatusOK, "OK", "success update data", nil)

}
func convertDataToResponse(data destinations.Destination) destinations.DestinationResponse {
	resp := destinations.DestinationResponse{
		ID:          data.ID,
		Nama:        data.Nama,
		Description: data.Description,
		Location:    data.Location,
		Category:    data.Category,
		Image:       data.Image,
		Rating:      data.Rating,
		CountUlasan: data.CountUlasan,
		Ulasan:      data.Ulasan,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
	return resp
}

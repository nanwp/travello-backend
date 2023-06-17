package handler

import (
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
	"github.com/nanwp/travello/service"
)

type destinatinHandler struct {
	urlApi   string
	uService service.UlasanService
}

func NewDestinationHandler(uService service.UlasanService) *destinatinHandler {
	return &destinatinHandler{"https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/destination", uService}
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

	createDes := destinations.Destination{
		Nama:        c.Query("name"),
		Location:    c.Query("location"),
		Description: c.Query("description"),
	}

	urlParams := url.Values{}
	urlParams.Add("name", createDes.Nama)
	urlParams.Add("location", createDes.Location)
	urlParams.Add("description", createDes.Description)

	postData, err := http.PostForm(h.urlApi, urlParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	defer postData.Body.Close()
	body, err := ioutil.ReadAll(postData.Body)
	bodyString := string(body)

	c.JSON(http.StatusCreated, gin.H{
		"data": bodyString,
	})

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

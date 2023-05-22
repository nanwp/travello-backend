package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/nanwp/travello/helper"
	"github.com/nanwp/travello/models/destinations"
	"github.com/nanwp/travello/models/ulasans"
	"github.com/nanwp/travello/service"
)

type destinatinHandler struct {
	urlApi string
}

func NewDestinationHandler() *destinatinHandler {
	return &destinatinHandler{"https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/destination"}
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

	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	hasil := destinations.Destination{}
	json.Unmarshal(responseData, &hasil)

	ulas, jumlah, _ := service.NewUlasanService().Get(destinationId)
	var maxUlasan int

	if len(ulas) < 4 {
		maxUlasan = len(ulas)
	} else {
		maxUlasan = 4
	}

	if ulas == nil {
		ulas = []ulasans.ResponseUlasan{}
	}

	respData := convertDataToResponse(hasil, jumlah, ulas[:maxUlasan])

	helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", respData)
}

func (h *destinatinHandler) Destinations(c *gin.Context) {

	category := c.Query("category")
	search := c.Query("search")

	if search != "" && category != "" {
		response, err := http.Get(h.urlApi + "?search=" + url.QueryEscape(search) + "&category=" + category)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
			return
		}

		hasil := []destinations.Destination{}

		json.Unmarshal(responseData, &hasil)

		sort.Slice(hasil, func(i, j int) bool {
			return hasil[j].UpdatedAt < hasil[i].UpdatedAt
		})

		var data []destinations.DestinationResponse

		if len(hasil) != 0 {
			for _, d := range hasil {
				_, jumlah, _ := service.NewUlasanService().Get(d.ID)
				data = append(data, convertDataToResponse(d, jumlah, nil))
			}
			helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", data)
			return
		}

		helper.ResponseOutput(c, http.StatusNotFound, "NOT_FOUND", "data not found", []destinations.DestinationResponse{})
		return
	}

	if search != "" {
		response, err := http.Get(h.urlApi + "?search=" + url.QueryEscape(search))
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
			return
		}

		hasil := []destinations.Destination{}

		json.Unmarshal(responseData, &hasil)

		sort.Slice(hasil, func(i, j int) bool {
			return hasil[j].UpdatedAt < hasil[i].UpdatedAt
		})

		var data []destinations.DestinationResponse

		if len(hasil) != 0 {
			for _, d := range hasil {
				_, jumlah, _ := service.NewUlasanService().Get(d.ID)
				data = append(data, convertDataToResponse(d, jumlah, nil))
			}
			helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", data)
			return
		}

		helper.ResponseOutput(c, http.StatusNotFound, "NOT_FOUND", "data not found", []destinations.DestinationResponse{})
		return
	}

	if category != "" {
		response, err := http.Get(h.urlApi + "?category=" + category)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
			return
		}

		hasil := []destinations.Destination{}
		json.Unmarshal(responseData, &hasil)

		sort.Slice(hasil, func(i, j int) bool {
			return hasil[j].UpdatedAt < hasil[i].UpdatedAt
		})

		var data []destinations.DestinationResponse

		if len(hasil) != 0 {
			for _, d := range hasil {
				_, jumlah, _ := service.NewUlasanService().Get(d.ID)
				data = append(data, convertDataToResponse(d, jumlah, nil))
			}
			helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", data)
			return
		}

		helper.ResponseOutput(c, http.StatusNotFound, "NOT_FOUND", "data not found", []destinations.DestinationResponse{})
		return

	}

	response, err := http.Get(h.urlApi)

	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), []destinations.DestinationResponse{})
		return
	}

	hasil := []destinations.Destination{}

	json.Unmarshal(responseData, &hasil)

	sort.Slice(hasil, func(i, j int) bool {
		return hasil[j].UpdatedAt < hasil[i].UpdatedAt
	})

	var data []destinations.DestinationResponse

	for _, d := range hasil {
		_, jumlah, _ := service.NewUlasanService().Get(d.ID)
		data = append(data, convertDataToResponse(d, jumlah, nil))
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

func convertDataToResponse(data destinations.Destination, jumlah int, ulasan []ulasans.ResponseUlasan) destinations.DestinationResponse {
	resp := destinations.DestinationResponse{
		ID:          data.ID,
		Nama:        data.Nama,
		Description: data.Description,
		Location:    data.Location,
		Category:    data.Category,
		Image:       data.Image,
		Rating:      data.Rating,
		CountUlasan: jumlah,
		Ulasan:      ulasan,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
	return resp
}

// func (s *destinationService) Create(destinatin destinations.Destination) (http.Response, error) {
// 	params := url.Values{}
// 	params.Add("name", destinatin.Nama)
// 	params.Add("location", destinatin.Location)

// 	resp, err := http.PostForm(s.url, params)

// 	if err != nil {
// 		return *resp, err
// 	}

// 	return *resp, err

// }

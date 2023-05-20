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
)

type destinatinHandler struct {
	urlApi string
}

func NewDestinationHandler() *destinatinHandler {
	return &destinatinHandler{"https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/destination"}
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
				data = append(data, convertDataToResponse(d))
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
				data = append(data, convertDataToResponse(d))
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
				data = append(data, convertDataToResponse(d))
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
		data = append(data, convertDataToResponse(d))
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

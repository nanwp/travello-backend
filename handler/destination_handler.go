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
		response, err := http.Get(h.urlApi + "?search=" + search + "&category=" + category)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
			return
		}

		hasil := []destinations.Destination{}

		json.Unmarshal(responseData, &hasil)

		sort.Slice(hasil, func(i, j int) bool {
			return hasil[j].UpdatedAt < hasil[i].UpdatedAt
		})

		if len(hasil) != 0 {
			helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", hasil)
			return
		}

		helper.ResponseOutput(c, http.StatusNotFound, "NOT_FOUND", "data not found", nil)
		return
	}

	if search != "" {
		response, err := http.Get(h.urlApi + "?search=" + search)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
			return
		}

		hasil := []destinations.Destination{}

		json.Unmarshal(responseData, &hasil)

		sort.Slice(hasil, func(i, j int) bool {
			return hasil[j].UpdatedAt < hasil[i].UpdatedAt
		})

		if len(hasil) != 0 {
			helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", hasil)
			return
		}

		helper.ResponseOutput(c, http.StatusNotFound, "NOT_FOUND", "data not found", nil)
		return
	}

	if category != "" {
		response, err := http.Get(h.urlApi + "?category=" + category)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
			return
		}

		hasil := []destinations.Destination{}
		json.Unmarshal(responseData, &hasil)

		sort.Slice(hasil, func(i, j int) bool {
			return hasil[j].UpdatedAt < hasil[i].UpdatedAt
		})

		if len(hasil) != 0 {
			helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", hasil)
			return
		}

		helper.ResponseOutput(c, http.StatusNotFound, "NOT_FOUND", "data not found", nil)
		return

	}

	response, err := http.Get(h.urlApi)

	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		helper.ResponseOutput(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
		return
	}

	hasil := []destinations.Destination{}

	json.Unmarshal(responseData, &hasil)

	sort.Slice(hasil, func(i, j int) bool {
		return hasil[j].UpdatedAt < hasil[i].UpdatedAt
	})

	helper.ResponseOutput(c, http.StatusOK, "OK", "Success get data", hasil)
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

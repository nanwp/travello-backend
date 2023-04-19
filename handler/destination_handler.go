package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
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

	if search != "" {
		response, err := http.Get(h.urlApi + "?search=" + search)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		hasil := []destinations.Destination{}

		json.Unmarshal(responseData, &hasil)

		if len(hasil) != 0 {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"status":  "OK",
				"data":    hasil,
			})
			return
		}

		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"status":  "NOT_FOUND",
			"message": "tidak ditemukan",
		})
		return

	}

	if category != "" {
		response, err := http.Get(h.urlApi + "?category=" + category)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		hasil := []destinations.Destination{}
		json.Unmarshal(responseData, &hasil)

		if len(hasil) != 0 {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"status":  "OK",
				"data":    hasil,
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"data":    "tidak ditemukan",
		})
		return

	}

	response, err := http.Get(h.urlApi)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	hasil := []destinations.Destination{}

	json.Unmarshal(responseData, &hasil)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  "OK",
		"data":    hasil,
	})
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

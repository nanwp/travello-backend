package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/nanwp/travello/config"
	"github.com/nanwp/travello/models/ulasans"
	"github.com/nanwp/travello/repository"
)

type UlasanService interface {
	Get(idDestination string) ([]ulasans.ResponseUlasan, int, error)
}

type ulasanService struct {
	url string
}

func NewUlasanService() *ulasanService {
	return &ulasanService{"https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/ulasan"}
}

func (s *ulasanService) Get(idDestination string) ([]ulasans.ResponseUlasan, int, error) {

	resp, err := http.Get(s.url + "?destination=" + idDestination)

	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	bodyString := string(bodyBytes)

	var arrayUlasan []ulasans.Ulasan

	json.Unmarshal([]byte(bodyString), &arrayUlasan)

	if len(arrayUlasan) == 0 {
		return nil, 0, errors.New("ulasan kosong")
	}

	responseUlasans := []ulasans.ResponseUlasan{}

	for _, u := range arrayUlasan {
		userService := NewUserService(repository.NewUserRepository(config.ConnectDatabase()))
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

	return responseUlasans, len(responseUlasans), nil
}

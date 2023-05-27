package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/nanwp/travello/models/ulasans"
	"github.com/nanwp/travello/repository"
	"gorm.io/gorm"
)

type UlasanService interface {
	GetUlasanByDestinationID(idDestination string) ([]ulasans.ResponseUlasan, error)
	Create(ulasan ulasans.Ulasan) (ulasans.Ulasan, error)
	// GetAllUlasan() ([]ulasans.ResponseUlasan, error)
	GetCountUlasanByDestinationID(destinationId string) (int64, error)
}

type ulasanService struct {
	repository     repository.UlasanRepository
	urlDestination string
	usrService     userService
}

func NewUlasanService(repository repository.UlasanRepository, usrService userService) *ulasanService {
	return &ulasanService{repository, "https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/destination", usrService}
}

func (s *ulasanService) GetCountUlasanByDestinationID(destinationId string) (int64, error) {
	return s.repository.GetCountUlasanByDestinationID(destinationId)
}

func (s *ulasanService) Create(ulasan ulasans.Ulasan) (ulasans.Ulasan, error) {

	//mengambil data ulasan yang sudah ada

	dataUlasan, err := s.repository.GetByDestinationID(ulasan.DestinationId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ulasan, _ := s.repository.Create(ulasan)
			return ulasan, nil
		} else {
			return ulasans.Ulasan{}, err
		}
	}

	for _, u := range dataUlasan {
		if u.UserId == ulasan.UserId {
			error := errors.New("hanya bisa 1 kali")
			return ulasans.Ulasan{}, error
		}
	}

	//set data ulasan
	setUlasan, err := s.repository.Create(ulasan)
	if err != nil {
		return ulasans.Ulasan{}, err
	}

	//menghitung total ulasan
	var totalUlasan = ulasan.Rating

	//menambah total ulasan
	for _, u := range dataUlasan {
		totalUlasan += u.Rating
	}

	//update ulasan destinasi

	type bodyUpdate struct {
		Rating float32 `json:"rating"`
	}

	updateDest := bodyUpdate{
		totalUlasan / float32(len(dataUlasan)+1),
	}

	jsonReqDest, err := json.Marshal(updateDest)

	//mengupdate rating pada destinasi
	reqDestinationUpdate, err := http.NewRequest(http.MethodPut, s.urlDestination+"?id="+ulasan.DestinationId, bytes.NewBuffer(jsonReqDest))
	reqDestinationUpdate.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(reqDestinationUpdate)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	return setUlasan, nil
}

func (s *ulasanService) GetUlasanByDestinationID(idDestination string) ([]ulasans.ResponseUlasan, error) {

	ulasan, err := s.repository.GetByDestinationID(idDestination)
	if err != nil {
		return nil, err
	}

	responseUlasans := []ulasans.ResponseUlasan{}

	for _, u := range ulasan {
		response := ulasans.ResponseUlasan{
			UserName: u.User.Name,
			Message:  u.Message,
			Rating:   u.Rating,
		}
		responseUlasans = append(responseUlasans, response)
	}

	return responseUlasans, nil
}

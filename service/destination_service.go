package service

// import (
// 	"net/http"
// 	"net/url"

// 	"github.com/nanwp/travello/models/destinations"
// )

// type destinationService struct {
// 	url string
// }

// func NewDestinationService() *destinationService {
// 	return &destinationService{"https://ap-southeast-1.aws.data.mongodb-api.com/app/travello-sfoqh/endpoint/destination"}
// }

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

// func (s *destinationService) FindAll() (destinations.Destination, error){
// 	res, err := http.Get(s.url)
// }

// response, err := http.Get("https://ap-southeast-1.aws.data.mongodb-api.com/app/application-0-csufo/endpoint/pasien?status=sakit")

// if err != nil {
// 	fmt.Println(err)
// }

// responseData, err := ioutil.ReadAll(response.Body)
// if err != nil {
// 	log.Fatal(err)
// }

// hasil := []AutoGenerated{}

// json.Unmarshal(responseData, &hasil)

// for _, h := range hasil{
// 	fmt.Println("===============")
// 	fmt.Println("Nama : "+h.Nama)
// 	fmt.Println("Alamat : "+h.Alamat)
// 	fmt.Println("Penyakit : "+h.Penyakit)
// 	fmt.Println("Status : "+h.Status)
// 	fmt.Println("Id : "+h.ID)

// }

// fmt.Println(hasil)

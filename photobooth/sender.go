package PhotoBooth

import (
	"net/http"
	"mime/multipart"
	"bytes"
	"io"
	"log"
	"time"
	"fmt"
	"strconv"
)

type (
	Sender struct {
		serverURL string
	}
)

func (s *Sender) Init() {
	s.serverURL = "http://dev.vasypaulette.com/photomaton/wp-json/photo/v1/upload"
}

func (s *Sender) Upload(pictureData []byte, mail string) (int) {
	log.Println("-> Upload", len(pictureData), mail)
	extraParams := map[string]string{
		"email": mail,
	}

	request, err := s.newfileUploadRequest(extraParams, "project", pictureData)
	if err != nil {
		log.Println(err)
		return 503
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return 500
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			log.Println(err)
			return 501
		}

		log.Println("-> Uploaded", resp.StatusCode)
		resp.Body.Close()
		return resp.StatusCode;
	}
}


// Creates a new file upload http request with optional extra params
func (s *Sender) newfileUploadRequest(params map[string]string, paramName string, binary []byte) (*http.Request, error) {

	path := formattedTimestamp() + "_dirty_dancing.jpg"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(binary)
	_, err = io.Copy(part, r)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", s.serverURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

func formattedTimestamp() (string) {
	loc, _ := time.LoadLocation("Europe/Paris")
    t := time.Now().In(loc)
	return fmt.Sprintf("%d-%s-%d_%sh%sm%s", t.Year(), format(int(t.Month())), t.Day(), format(int(t.Hour())), format(int(t.Minute())), format(int(t.Second())))
} 


func format(hour int) (string) {
	formatted := strconv.Itoa(hour)
	if (hour < 10) {
		formatted = "0" + strconv.Itoa(hour)
	}
    return formatted
}
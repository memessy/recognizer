package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"memery-recognizer/recognizer/gosseract"
	"net/http"
)

type Controller struct {
	RecognizerFactory func() gosseract.Recognizer
	AuthToken         string
}

const (
	MaxSize = 10 << 20
	FileKey = "file"
)

type ErrResponse struct {
	Message string `json:"message"`
}

type SuccessResponse struct {
	Text string `json:"text"`
}

func (c *Controller) Upload(rw http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("token")
	if authToken != c.AuthToken {
		log.Warn().Msg("unauthorized access")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	responseInvalidInput := func(err error, message string) {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		log.Warn().Err(err).Msg("got error while getting user input")
		err = json.NewEncoder(rw).Encode(ErrResponse{Message: message})
		if err != nil {
			log.Error().Err(err).Msg("caught error while encoding error response")
		}
	}

	err := r.ParseMultipartForm(MaxSize)
	if err != nil {
		responseInvalidInput(err, "File size is too big.")
		return
	}
	file, _, err := r.FormFile(FileKey)
	if err != nil {
		responseInvalidInput(err, fmt.Sprintf("File must be stored under `%s` key.", FileKey))
		return
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		responseInvalidInput(err, "Could not read file from input.")
		return
	}
	rec := c.RecognizerFactory()
	defer rec.Close()
	text, err := rec.Recognize(bytes)
	if err != nil {
		log.Error().Err(err).Msg("caught error while recognizing text from image")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(SuccessResponse{Text: text})
	if err != nil {
		log.Error().Err(err).Msg("caught error while encoding success response")
	}
}

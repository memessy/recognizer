package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"memery-recognizer/pool"
	"net/http"
)

type Controller struct {
	Pool    pool.RecognizerPool
	MaxSize int64
	FileKey string
}

type errResponse struct {
	Message string `json:"message"`
}

type successResponse struct {
	Text string `json:"text"`
}

func (c *Controller) Upload(rw http.ResponseWriter, r *http.Request) {
	responseInvalidInput := func(err error, message string) {
		rw.WriteHeader(http.StatusUnprocessableEntity)
		log.Warn().Err(err).Msg("got error while getting user input")
		err = json.NewEncoder(rw).Encode(errResponse{Message: message})
		if err != nil {
			log.Error().Err(err).Msg("caught error while encoding error response")
		}
	}

	err := r.ParseMultipartForm(c.MaxSize)
	if err != nil {
		responseInvalidInput(err, "File size is too big.")
		return
	}
	file, _, err := r.FormFile(c.FileKey)
	if err != nil {
		responseInvalidInput(err, fmt.Sprintf("File must be stored under `%s` key.", c.FileKey))
		return
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		responseInvalidInput(err, "Could not read file from input.")
		return
	}
	text, err := c.Pool.Recognize(r.Context(), bytes)
	if err != nil {
		log.Error().Err(err).Msg("caught error while recognizing text from image")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	err = json.NewEncoder(rw).Encode(successResponse{Text: text})
	if err != nil {
		log.Error().Err(err).Msg("caught error while encoding success response")
	}
}

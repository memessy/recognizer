package api

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"memessy-recognizer/pkg/api/impl/multipart"
	"memessy-recognizer/pkg/pool"
	"net/http"
)

type Controller struct {
	Pool            pool.RecognizerPool
	MultipartReader multipart.Reader
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

	image, err := c.MultipartReader.Read(r)
	if err != nil {
		responseInvalidInput(err, "Could not parse request.")
		return
	}

	text, err := c.Pool.Recognize(r.Context(), image)
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

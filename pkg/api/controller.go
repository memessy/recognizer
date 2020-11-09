package api

import (
	"net/http"
)

type Controller interface {
	Upload(rw http.ResponseWriter, r *http.Request)
}

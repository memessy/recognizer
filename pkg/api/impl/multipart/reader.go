package multipart

import (
	"io/ioutil"
	"net/http"
)

type Reader interface {
	Read(*http.Request) ([]byte, error)
}

func NewReader(fileKey string, maxSize int64) Reader {
	return &reader{
		fileKey: fileKey,
		maxSize: maxSize,
	}
}

type reader struct {
	fileKey string
	maxSize int64
}

func (r *reader) Read(req *http.Request) ([]byte, error) {
	err := req.ParseMultipartForm(r.maxSize)
	if err != nil {
		return nil, err
	}
	file, _, err := req.FormFile(r.fileKey)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

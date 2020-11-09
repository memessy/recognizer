package gosseract

import (
	"github.com/otiai10/gosseract"
	"memessy-recognizer/pkg/recognizer"
	"sync"
)

type Recognizer struct {
	client *gosseract.Client
	mux    sync.Mutex
}

func NewRecognizer() recognizer.Recognizer {
	client := gosseract.NewClient()
	_ = client.SetLanguage("rus", "eng")
	return &Recognizer{
		client: client,
		mux:    sync.Mutex{},
	}
}

func (r *Recognizer) Close() error {
	return r.client.Close()
}

func (r *Recognizer) Recognize(data []byte) (string, error) {
	r.mux.Lock()
	defer r.mux.Unlock()
	err := r.client.SetImageFromBytes(data)
	if err != nil {
		return "", err
	}
	text, err := r.client.Text()
	if err != nil {
		return "", err
	}
	return text, nil
}

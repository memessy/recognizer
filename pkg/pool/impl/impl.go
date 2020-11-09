package impl

import (
	"context"
	"fmt"
	"memessy-recognizer/pkg/recognizer"
)

func NewPool(size int, recognizerFactory func() recognizer.Recognizer) *pool {
	p := &pool{
		recognizerFactory: recognizerFactory,
		sem:               make(chan struct{}, size),
	}
	return p
}

type pool struct {
	recognizerFactory func() recognizer.Recognizer
	sem               chan struct{}
}

func (p *pool) Recognize(ctx context.Context, data []byte) (string, error) {
	select {
	case p.sem <- struct{}{}:
		defer func() { <-p.sem }()
		rec := p.recognizerFactory()
		text, err := rec.Recognize(data)
		if err != nil {
			return "", err
		}
		return text, nil
	case <-ctx.Done():
		return "", fmt.Errorf("context is done, err: %v", ctx.Err())
	}
}

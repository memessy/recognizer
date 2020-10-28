package impl

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"memery-recognizer/recognizer"
)

func NewPool(ctx context.Context, size int, recognizerFactory func() recognizer.Recognizer) *pool {
	p := &pool{
		input:             make(chan message),
		recognizerFactory: recognizerFactory,
		ctx:               ctx,
	}
	p.init(size)
	go p.close()
	return p
}

type pool struct {
	input             chan message
	recognizerFactory func() recognizer.Recognizer
	ctx               context.Context
}

type message struct {
	ctx    context.Context
	data   []byte
	result chan string
	error  chan error
}

func (p *pool) Recognize(ctx context.Context, data []byte) (string, error) {
	msg := message{
		ctx:    ctx,
		data:   data,
		result: make(chan string),
		error:  make(chan error),
	}
	defer func() {
		close(msg.result)
		close(msg.error)
	}()
	p.input <- msg
	select {
	case text := <-msg.result:
		return text, nil
	case err := <-msg.error:
		return "", err
	case <-ctx.Done():
		return "", errors.New("context is done")
	}
}

func (p *pool) init(size int) {
	for i := 0; i < size; i++ {
		go p.worker()
	}
}

func (p *pool) close() {
	select {
	case <-p.ctx.Done():
		close(p.input)
	}
}

func (p *pool) worker() {
	rec := p.recognizerFactory()
	defer rec.Close()
	for {
		select {
		case msg := <-p.input:
			select {
			case <-msg.ctx.Done():
				break
			default:
				text, err := rec.Recognize(msg.data)
				if err != nil {
					log.Error().Err(err)
					msg.error <- err
					break
				}
				msg.result <- text
			}
		case <-p.ctx.Done():
			return
		}
	}
}

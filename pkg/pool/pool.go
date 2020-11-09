package pool

import (
	"context"
)

type RecognizerPool interface {
	Recognize(context.Context, []byte) (string, error)
}

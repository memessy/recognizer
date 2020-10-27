package recognizer

type Recognizer interface {
	Recognize(data []byte) (string, error)
	Close() error
}

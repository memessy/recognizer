package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	api "memessy-recognizer/pkg/api/impl"
	"memessy-recognizer/pkg/api/impl/multipart"
	pool "memessy-recognizer/pkg/pool/impl"
	recognizer "memessy-recognizer/pkg/recognizer/gosseract"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

func main() {
	port := os.Getenv("PORT")
	fileMaxSize, err := strconv.Atoi(os.Getenv("FILE_MAX_SIZE"))
	if err != nil {
		log.Fatal().Err(err)
	}
	fileKey := os.Getenv("FILE_KEY")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	multipartReader := multipart.NewReader(fileKey, int64(fileMaxSize))
	p := pool.NewPool(runtime.NumCPU(), recognizer.NewRecognizer)
	controller := api.Controller{
		Pool:            p,
		MultipartReader: multipartReader,
	}
	r.Post("/recognize", controller.Upload)
	log.Fatal().Err(http.ListenAndServe(":"+port, r))
}

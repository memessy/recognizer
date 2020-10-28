package main

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	api "memery-recognizer/api/impl"
	pool "memery-recognizer/pool/impl"
	recognizer "memery-recognizer/recognizer/gosseract"
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
	p := pool.NewPool(ctx, runtime.NumCPU()-1, recognizer.NewRecognizer)
	controller := api.Controller{
		Pool:    p,
		MaxSize: int64(fileMaxSize),
		FileKey: fileKey,
	}
	r.Post("/recognize", controller.Upload)
	log.Fatal().Err(http.ListenAndServe(":"+port, r))
}

package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"memery-recognizer/api"
	"memery-recognizer/recognizer/gosseract"
	"net/http"
	"os"
)

func main() {
	authToken := os.Getenv("AUTH_TOKEN")
	port := os.Getenv("PORT")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	controller := api.Controller{RecognizerFactory: gosseract.NewRecognizer, AuthToken: authToken}
	r.Post("/recognize", controller.Upload)
	log.Fatal(http.ListenAndServe(":" + port, r))
}

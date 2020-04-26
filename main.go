package main

import (
	"net/http"

	"./models"
	"./routes"
	"./utils"
)

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	router := routes.NewRouter()
	http.Handle("/", router)
	http.ListenAndServe("localhost:9999", nil)
}

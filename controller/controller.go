package controller

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
	regService "reg-api/services"
)
const port = 8081
func InitServer() {
	router := mux.NewRouter()
	router.HandleFunc("/api/clear/{type}", regService.ClearCache).Methods("GET")
	router.HandleFunc("/api/find/{cid}/{year}/{semester}", regService.FindCourseAll).Methods("GET")
	router.HandleFunc("/api/course/major/{id}", regService.CourseMajor).Methods("GET")
	router.HandleFunc("/api/major/list/{id}", regService.GetMajor).Methods("GET")
	router.HandleFunc("/api/import/{stdid}/{acadyear}/{semester}", regService.ImportFormReg).Methods("GET")
	router.HandleFunc("/api/{id}/{year}/{semester}", regService.ScrapingCourseDetail).Methods("GET")

	mcors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in producti
		Debug: true,
	})
	handler := mcors.Handler(router)
	fmt.Print("server port:", port)
	http.ListenAndServe(":8081", (handler))
}


package main

import (
	"fmt"
	"mahasiswa/controller/auth"
	"mahasiswa/controller/mahasiswa"
	"mahasiswa/controller/kelas"
	"mahasiswa/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()
	fmt.Println("Hello World")

	router := mux.NewRouter()

	router.HandleFunc("/regis", auth.Registration).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	//Router handler mahasiswa
	router.HandleFunc("/mahasiswa", mahasiswa.GetMahasiswa).Methods("GET")
	router.HandleFunc("/mahasiswa", auth.JWTAuth(mahasiswa.PostMahasiswa)).Methods("POST")
	router.HandleFunc("/mahasiswa/{id}", auth.JWTAuth(mahasiswa.PutMahasiswa)).Methods("PUT")
	router.HandleFunc("/mahasiswa/{id}", auth.JWTAuth(mahasiswa.DeleteMahasiswa)).Methods("DELETE")

	//Router handler Album
	router.HandleFunc("/kelas", kelas.GetKelas).Methods("GET")
	router.HandleFunc("/kelas", auth.JWTAuth(kelas.PostKelas)).Methods("POST")
	router.HandleFunc("/kelas/{id}", auth.JWTAuth(kelas.PutKelas)).Methods("PUT")
	router.HandleFunc("/kelas/{id}", auth.JWTAuth(kelas.DeleteKelas)).Methods("DELETE")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://127.0.0.1:5500"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
		Debug: true,
	})

	handler := c.Handler(router)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
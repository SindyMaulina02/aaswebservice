package main

import (
	"fmt"
	"log"
	"mahasiswa/controller/auth"
	"mahasiswa/controller/kelas"
	"mahasiswa/controller/mahasiswa"
	"mahasiswa/database"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// main.go


func main() {
    database.InitDB()
    fmt.Println("Hello World")

    router := mux.NewRouter()

    // Auth routes
    router.HandleFunc("/regis", auth.Registration).Methods("POST")
    router.HandleFunc("/login", auth.Login).Methods("POST")

    // Kelas routes
    router.HandleFunc("/kelas", kelas.GetKelas).Methods("GET")
    router.HandleFunc("/kelas/{id}", kelas.GetKelasByID).Methods("GET")
    router.HandleFunc("/kelas", auth.JWTAuth(kelas.PostKelas)).Methods("POST")
    router.HandleFunc("/kelas/{id}", auth.JWTAuth(kelas.PutKelas)).Methods("PUT")
    router.HandleFunc("/kelas/{id}", auth.JWTAuth(kelas.DeleteKelas)).Methods("DELETE")

    // Mahasiswa routes
    router.HandleFunc("/mahasiswa", mahasiswa.GetMahasiswa).Methods("GET")
    router.HandleFunc("/mahasiswa/{id}", mahasiswa.GetMahasiswaByID).Methods("GET")
    router.HandleFunc("/mahasiswa", auth.JWTAuth(mahasiswa.PostMahasiswa)).Methods("POST")
    router.HandleFunc("/mahasiswa/{id}", auth.JWTAuth(mahasiswa.PutMahasiswa)).Methods("PUT")
    router.HandleFunc("/mahasiswa/{id}", auth.JWTAuth(mahasiswa.DeleteMahasiswa)).Methods("DELETE")

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
        Debug:          true,
    })

    handler := c.Handler(router)

    fmt.Println("Server is running on http://localhost:8018")
    log.Fatal(http.ListenAndServe(":8018", handler))
}

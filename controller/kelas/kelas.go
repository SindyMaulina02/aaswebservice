package kelas

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"mahasiswa/database"
	"mahasiswa/model/kelas"

	"github.com/gorilla/mux"
)

// GetKelas handles GET requests to fetch all kelas
func GetKelas(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM kelas")
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var kelasList []kelas.Kelas
	for rows.Next() {
		var k kelas.Kelas
		if err := rows.Scan(&k.ID, &k.NamaKelas, &k.Semester); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		kelasList = append(kelasList, k)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error after rows iteration: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kelasList)
}

// PostKelas handles POST requests to add a new kelas
func PostKelas(w http.ResponseWriter, r *http.Request) {
	var pc kelas.Kelas
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO kelas (nama_kelas, semester) 
		VALUES (?, ?)`

	result, err := database.DB.Exec(query, pc.NamaKelas, pc.Semester)
	if err != nil {
		http.Error(w, "Failed to insert kelas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kelas added successfully",
		"id":      id,
	})
}

// PutKelas handles PUT requests to update an existing kelas by ID
func PutKelas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	var pc kelas.Kelas
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE kelas 
		SET nama_kelas=?, semester=?
		WHERE id=?`

	result, err := database.DB.Exec(query, pc.NamaKelas, pc.Semester, id)
	if err != nil {
		http.Error(w, "Failed to update Kelas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kelas updated successfully",
	})
}

// DeleteKelas handles DELETE requests to delete a kelas by ID
func DeleteKelas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		log.Println("ID not provided")
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %v", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Check if there are any associated mahasiswa records
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM mahasiswa WHERE kelas_id = ?", id).Scan(&count)
	if err != nil {
		log.Printf("Error checking associated mahasiswa: %v", err)
		http.Error(w, "Failed to check associated mahasiswa", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		// If there are associated mahasiswa records, return an error response
		http.Error(w, "Cannot delete kelas with associated mahasiswa records", http.StatusConflict)
		return
	}

	query := `DELETE FROM kelas WHERE id = ?`
	result, err := database.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting kelas: %v", err)
		http.Error(w, "Failed to delete kelas", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving affected rows: %v", err)
		http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		log.Println("No rows were deleted")
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kelas deleted successfully",
	})
}

// GetKelasByID handles GET requests to fetch a single kelas by its ID
func GetKelasByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var k kelas.Kelas
	query := "SELECT id, nama_kelas, semester FROM kelas WHERE id = ?"
	err = database.DB.QueryRow(query, id).Scan(&k.ID, &k.NamaKelas, &k.Semester)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Kelas not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(k)
}

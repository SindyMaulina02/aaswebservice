package kelas

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "mahasiswa/database"
    "mahasiswa/model/kelas"

    "github.com/gorilla/mux"
)

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

func PostKelas(w http.ResponseWriter, r *http.Request) {
	var pc kelas.Kelas
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new course
	query := `
		INSERT INTO kelas (nama_kelas, semester) 
		VALUES (?, ?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.NamaKelas, pc.Semester)
	if err != nil {
		http.Error(w, "Failed to insert kelas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kelas added successfully",
		"id":      id,
	})
}
func PutKelas(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
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

	// Decode JSON body
	var pc kelas.Kelas
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating the category admin
	query := `
		UPDATE kelas 
		SET nama_kelas=?, semester=?
		WHERE id=?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, pc.NamaKelas, pc.Semester, id)
	if err != nil {
		http.Error(w, "Failed to update Kelas: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were updated
	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kelas updated successfully",
	})
}

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

    query := `
        DELETE FROM kelas
        WHERE id = ?`

    result, err := database.DB.Exec(query, id)
    if err != nil {
        log.Printf("Error deleting from database: %v", err)
        http.Error(w, "Failed to delete kelas: "+err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error retrieving affected rows: %v", err)
        http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
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

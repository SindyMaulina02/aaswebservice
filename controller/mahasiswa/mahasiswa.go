package mahasiswa

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "strconv"

    "mahasiswa/database"
    "mahasiswa/model/mahasiswa"

    "github.com/gorilla/mux"
)

// GetMahasiswa handles GET requests to fetch all mahasiswa
func GetMahasiswa(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT * FROM mahasiswa")
    if err != nil {
        log.Printf("Error querying database: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var mahasiswas []mahasiswa.Mahasiswa
    for rows.Next() {
        var m mahasiswa.Mahasiswa
        if err := rows.Scan(&m.MahasiswaID, &m.Nama, &m.Nim, &m.Jurusan, &m.UserID, &m.KelasID); err != nil {
            log.Printf("Error scanning row: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        mahasiswas = append(mahasiswas, m)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error after rows iteration: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mahasiswas)
}

// PostMahasiswa handles POST requests to add a new mahasiswa
func PostMahasiswa(w http.ResponseWriter, r *http.Request) {
    var m mahasiswa.Mahasiswa
    if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    query := `
        INSERT INTO mahasiswa (nama, nim, jurusan, user_id, kelas_id) 
        VALUES (?, ?, ?, ?, ?)`

    res, err := database.DB.Exec(query, m.Nama, m.Nim, m.Jurusan, m.UserID, m.KelasID)
    if err != nil {
        log.Printf("Error inserting mahasiswa: %v", err)
        http.Error(w, "Failed to insert mahasiswa", http.StatusInternalServerError)
        return
    }

    id, err := res.LastInsertId()
    if err != nil {
        log.Printf("Error retrieving last insert ID: %v", err)
        http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Mahasiswa added successfully",
        "id":      id,
    })
}

// PutMahasiswa handles PUT requests to update an existing mahasiswa by ID
func PutMahasiswa(w http.ResponseWriter, r *http.Request) {
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

    var m mahasiswa.Mahasiswa
    if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    query := `
        UPDATE mahasiswa 
        SET nama=?, nim=?, jurusan=?, user_id=?, kelas_id=?
        WHERE id=?`

    result, err := database.DB.Exec(query, m.Nama, m.Nim, m.Jurusan, m.UserID, m.KelasID, id)
    if err != nil {
        log.Printf("Error updating mahasiswa: %v", err)
        http.Error(w, "Failed to update mahasiswa", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error retrieving affected rows: %v", err)
        http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "No rows were updated", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Mahasiswa updated successfully",
    })
}

// DeleteMahasiswa handles DELETE requests to delete a mahasiswa by ID
func DeleteMahasiswa(w http.ResponseWriter, r *http.Request) {
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

    query := `DELETE FROM mahasiswa WHERE id = ?`
    result, err := database.DB.Exec(query, id)
    if err != nil {
        log.Printf("Error deleting mahasiswa: %v", err)
        http.Error(w, "Failed to delete mahasiswa", http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error retrieving affected rows: %v", err)
        http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "No rows were deleted", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Mahasiswa deleted successfully",
    })
}

// GetMahasiswaByID handles GET requests to fetch a single mahasiswa by its ID
func GetMahasiswaByID(w http.ResponseWriter, r *http.Request) {
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

    var m mahasiswa.Mahasiswa
    query := "SELECT id, nama, nim, jurusan, user_id, kelas_id FROM mahasiswa WHERE id = ?"
    err = database.DB.QueryRow(query, id).Scan(&m.MahasiswaID, &m.Nama, &m.Nim, &m.Jurusan, &m.UserID, &m.KelasID)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Mahasiswa not found", http.StatusNotFound)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(m)
}

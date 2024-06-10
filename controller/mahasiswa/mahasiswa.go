package mahasiswa

import (
    "encoding/json"
    "net/http"
    "strconv"

    "mahasiswa/database"
    "mahasiswa/model/mahasiswa"

    "github.com/gorilla/mux"
)

func GetMahasiswa(w http.ResponseWriter, r *http.Request) {
    rows, err := database.DB.Query("SELECT * FROM mahasiswa")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var mahasiswas []mahasiswa.Mahasiswa
    for rows.Next() {
        var c mahasiswa.Mahasiswa
        if err := rows.Scan(&c.MahasiswaID, &c.Nama, &c.Nim, &c.Jurusan, &c.UserID, &c.KelasID); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        mahasiswas = append(mahasiswas, c)
    }

    if err := rows.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(mahasiswas)
}

func PostMahasiswa(w http.ResponseWriter, r *http.Request) {
    var pc mahasiswa.Mahasiswa
    if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    query := `
        INSERT INTO mahasiswa (nama, nim, jurusan, user_id, kelas_id) 
        VALUES (?, ?, ?, ?, ?)`

    res, err := database.DB.Exec(query, pc.Nama, pc.Nim, pc.Jurusan, pc.UserID, pc.KelasID)
    if err != nil {
        http.Error(w, "Failed to insert mahasiswa: "+err.Error(), http.StatusInternalServerError)
        return
    }

    id, err := res.LastInsertId()
    if err != nil {
        http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Mahasiswa added successfully",
        "id":      id,
    })
}

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

    var pc mahasiswa.Mahasiswa
    if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    query := `
        UPDATE mahasiswa
        SET nama=?, nim=?, jurusan=?, user_id=?, kelas_id=?
        WHERE id=?`

    result, err := database.DB.Exec(query, pc.Nama, pc.Nim, pc.Jurusan, pc.UserID, pc.KelasID, id)
    if err != nil {
        http.Error(w, "Failed to update mahasiswa: "+err.Error(), http.StatusInternalServerError)
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
        "message": "Mahasiswa updated successfully",
    })
}

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

    query := `
        DELETE FROM mahasiswa
        WHERE id = ?`

    result, err := database.DB.Exec(query, id)
    if err != nil {
        http.Error(w, "Failed to delete mahasiswa: "+err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
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

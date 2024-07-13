package mahasiswa

type Mahasiswa struct {
    MahasiswaID int    `json:"id"`
    Nama        string `json:"nama"`
    Nim         string `json:"nim"`
    Jurusan     string `json:"jurusan"`
    UserID      int    `json:"user_id"`
    KelasID     int    `json:"kelas_id"`
}
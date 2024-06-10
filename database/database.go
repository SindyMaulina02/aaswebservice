package database

import (
    "database/sql"
    "log"
    

    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
    var err error
    connStr := "root:@tcp(localhost:3306)/dbwebservice" // Ganti "root" dengan nama pengguna MySQL Anda jika diperlukan
    DB, err = sql.Open("mysql", connStr)
    if err != nil {
        log.Fatal(err)
    }

    // Pastikan koneksi ke database sudah berhasil
    err = DB.Ping()
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Database connected")
}

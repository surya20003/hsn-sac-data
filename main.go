
package main

import (
    "database/sql"
    "encoding/csv"
    "fmt"
    "log"
    "net/http"
    
    "strconv"

    "github.com/labstack/echo/v4"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    e := echo.New()

    // Connect to the database
    dsn := "root:suryakk07@tcp(127.0.0.1:3306)/testdb"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }
    defer db.Close()
    fmt.Println("✅ Database ready.")

    // Create the users table
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INT PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100)
    )`)
    if err != nil {
        log.Fatalf("❌ Table creation failed: %v", err)
    }
    fmt.Println("✅ Table ready.")

    // Define upload endpoint
    e.POST("/upload", func(c echo.Context) error {
        file, err := c.FormFile("file")
        if err != nil {
            return err
        }
        src, err := file.Open()
        if err != nil {
            return err
        }
        defer src.Close()

        reader := csv.NewReader(src)
        records, err := reader.ReadAll()
        if err != nil {
            return err
        }

        for i, row := range records {
            if i == 0 {
                continue // skip header
            }
            id, _ := strconv.Atoi(row[0])
            name := row[1]
            email := row[2]

            _, err := db.Exec("INSERT IGNORE INTO users (id, name, email) VALUES (?, ?, ?)", id, name, email)
            if err != nil {
                log.Printf("Insert failed: %v", err)
            }
        }
        return c.JSON(http.StatusOK, map[string]string{"message": "CSV uploaded successfully"})
    })

    fmt.Println("🚀 Server running at http://localhost:8080")
    e.Logger.Fatal(e.Start(":8080"))
}

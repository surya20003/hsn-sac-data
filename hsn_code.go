package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xuri/excelize/v2"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "suryakk07"
	dbname   = "hsn-db"
)

func main() {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	e := echo.New()

	//  CORS Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8081", "http://127.0.0.1:5500"},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	//  Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "Welcome to the HSN & SAC Master Data API Service!",
		})
	})

	e.POST("/upload_excel", func(c echo.Context) error {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "File is required"})
		}

		src, err := fileHeader.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to open file"})
		}
		defer src.Close()

		tempFilePath := "/tmp/" + fileHeader.Filename
		dst, err := os.Create(tempFilePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create temp file"})
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to save file"})
		}

		err = insertExcelDataToDB(context.Background(), conn, tempFilePath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{"message": "Excel data inserted successfully"})
	})

	e.GET("/hsn/:code", func(c echo.Context) error {
		hsncode := c.Param("code")
		var hsnCode, desc, rate string

		err := conn.QueryRow(context.Background(), "SELECT hsn_code, description, gst_rate FROM hsn_mstr WHERE hsn_code = $1", hsncode).Scan(&hsnCode, &desc, &rate)

		if err != nil {
			log.Printf("Error querying HSN code: %v", err) // Log the actual error
			return c.JSON(http.StatusNotFound, echo.Map{"error": "HSN code not found"})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"hsn_code":    hsnCode,
			"description": desc,
			"gst_rate":    rate,
		})
	})

	e.GET("/sac/:code", func(c echo.Context) error {
		code := c.Param("code")
		var sacCode, desc, hsn string

		err := conn.QueryRow(context.Background(), "SELECT sac_code, description, hsn_code FROM sac_mstr WHERE sac_code = $1", code).Scan(&sacCode, &desc, &hsn)
		if err != nil {
			log.Printf("Error querying SAC code: %v", err)
			return c.JSON(http.StatusNotFound, echo.Map{"error": "SAC code not found"})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"sac_code":    sacCode,
			"description": desc,
			"hsn_code":    hsn,
		})
	})

	//  Start the server
	e.Logger.Fatal(e.Start(":8081"))
}

func insertExcelDataToDB(ctx context.Context, conn *pgx.Conn, filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %v", err)
	}

	// Insert HSN data
	hsnRows, err := f.GetRows("hsn_mstr")
	if err == nil {
		for i, row := range hsnRows {
			if i == 0 || len(row) < 3 {
				continue
			}
			_, err := conn.Exec(ctx,
				"INSERT INTO hsn_mstr (hsn_code, description, gst_rate) VALUES ($1, $2, $3) ON CONFLICT (hsn_code) DO NOTHING",
				row[0], row[1], row[2],
			)
			if err != nil {
				log.Printf("Failed to insert HSN row: %v", err)
			}
		}
	} else {
		log.Println("HSN sheet not found.")
	}

	// Insert SAC data
	sacRows, err := f.GetRows("sac_mstr")
	if err == nil {
		for i, row := range sacRows {
			if i == 0 || len(row) < 3 {
				continue
			}

			// Check if sac_code already exists
			var exists bool
			err := conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM sac_mstr WHERE sac_code = $1)", row[0]).Scan(&exists)
			if err != nil {
				log.Printf("Error checking existence for SAC code %s: %v", row[0], err)
				continue
			}

			if !exists {
				_, err := conn.Exec(ctx,
					"INSERT INTO sac_mstr (sac_code, description, hsn_code) VALUES ($1, $2, $3)",
					row[0], row[1], row[2],
				)
				if err != nil {
					log.Printf("Failed to insert SAC row: %v", err)
				}
			} else {
				log.Printf("Skipped duplicate SAC code: %s", row[0])
			}
		}
	} else {
		log.Println("SAC sheet not found.")
	}

	return nil
}










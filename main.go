// package main

// import (
// 	"database/sql"
// 	"image/png"
// 	"bytes"
// 	"log"
// 	"net/http"
// 	"html/template"
// 	"strconv" // For integer to string conversion
	
// 	_ "github.com/mattn/go-sqlite3" // SQLite driver
// 	"github.com/boombuler/barcode"
// 	"github.com/boombuler/barcode/qr"
// )

// // Global database variable
// var db *sql.DB


// type Page struct {
// 	Title string
// }

// func initDB() {
// 	var err error
// 	db, err = sql.Open("sqlite3", "qrdata.db")
// 	if err != nil {
// 		log.Fatal("Database connection error: ", err)
// 	}

// 	// Create the table if it doesn't exist
// 	createTable := `
// 	CREATE TABLE IF NOT EXISTS qrcodes (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		data TEXT NOT NULL,
// 		qr_code BLOB NOT NULL,
// 		scan_count INTEGER DEFAULT 0
// 	);`
// 	_, err = db.Exec(createTable)
// 	if err != nil {
// 		log.Fatal("Error creating table: ", err)
// 	}

// 	log.Println("Database initialized successfully.")
// }

// func HomePage(w http.ResponseWriter, r *http.Request) {
// 	p := Page{Title: "QR Code Generator"}

// 	tmpl, err := template.ParseFiles("generator.html")
// 	if err != nil {
// 		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	tmpl.Execute(w, p)
// }


// func CodePage(w http.ResponseWriter, r *http.Request) {
// 	dataString := r.FormValue("dataString")
// 	qrCode, err := qr.Encode(dataString, qr.L, qr.Auto)
// 	if err != nil {
// 		http.Error(w, "QR Code generation error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	qrCode, err = barcode.Scale(qrCode, 512, 512)
// 	if err != nil {
// 		http.Error(w, "QR Code scaling error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	var qrCodeBuffer bytes.Buffer
// 	err = png.Encode(&qrCodeBuffer, qrCode)
// 	if err != nil {
// 		http.Error(w, "QR Code encoding error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Save the QR code to the database
// 	_, err = db.Exec("INSERT INTO qrcodes (data, qr_code) VALUES (?, ?)", dataString, qrCodeBuffer.Bytes())
// 	if err != nil {
// 		http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "image/png")
// 	w.Write(qrCodeBuffer.Bytes())
// }

// func RetrievePage(w http.ResponseWriter, r *http.Request) {
// 	rows, err := db.Query("SELECT id, data FROM qrcodes")
// 	if err != nil {
// 		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	// Display the list of stored QR codes
// 	var htmlContent string = "<h1>Stored QR Codes</h1><ul>"
// 	for rows.Next() {
// 		var id int
// 		var data string
// 		rows.Scan(&id, &data)
// 		htmlContent += "<li><a href='/view/" + strconv.Itoa(id) + "'>" + data + "</a></li>"

// 	}
// 	htmlContent += "</ul>"

// 	w.Write([]byte(htmlContent))
// }

// func ViewQRCode(w http.ResponseWriter, r *http.Request) {
// 	id := r.URL.Path[len("/view/"):]
// 	// Validate if the ID is empty
//     if id == "" {
//         http.Error(w, "ID is required", http.StatusBadRequest)
//         return
//     }
// 	row := db.QueryRow("SELECT qr_code, scan_count FROM qrcodes WHERE id = ?", id)

// 	var qrCodeData []byte
// 	var scanCount int
// 	err := row.Scan(&qrCodeData, &scanCount)
// 	if err != nil {
// 		http.Error(w, "QR Code not found", http.StatusNotFound)
// 		return
// 	}

// 	// Update the scan count
// 	_, err = db.Exec("UPDATE qrcodes SET scan_count = scan_count + 1 WHERE id = ?", id)
// 	if err != nil {
// 		http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "image/png")
// 	w.Write(qrCodeData)
// }



// func main() {

// 	initDB()
// 	defer db.Close()

// 	http.HandleFunc("/", HomePage)
// 	http.HandleFunc("/generator/", CodePage)
// 	http.HandleFunc("/retrieve/", RetrievePage)
// 	http.HandleFunc("/view/", ViewQRCode)

// 	log.Println("Server running at http://localhost:8080")
// 	http.ListenAndServe(":8080",nil)
// }


package main

import (
	"database/sql"
	"image/png"
	"bytes"
	"log"
	"net/http"
	"html/template"
	"strconv" // For integer to string conversion

	_ "github.com/mattn/go-sqlite3" // SQLite driver for database operations
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

// Global database variable
var db *sql.DB

// Page struct for passing dynamic data to HTML templates
type Page struct {
	Title string
}

// Initialize the SQLite database
func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "qrdata.db") // Connect to the SQLite database
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	// Create the table if it doesn't exist
	createTable := `
	CREATE TABLE IF NOT EXISTS qrcodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		data TEXT NOT NULL,
		qr_code BLOB NOT NULL,
		scan_count INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createTable) // Execute the table creation query
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	log.Println("Database initialized successfully.")
}

// Handler for the homepage
func HomePage(w http.ResponseWriter, r *http.Request) {
	p := Page{Title: "QR Code Generator"}

	// Parse and render the HTML template
	tmpl, err := template.ParseFiles("generator.html")
	if err != nil {
		http.Error(w, "Template parsing error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, p) // Render the template with dynamic data
}

// Handler for generating QR codes
func CodePage(w http.ResponseWriter, r *http.Request) {
	dataString := r.FormValue("dataString") // Get user input from the form
	qrCode, err := qr.Encode(dataString, qr.L, qr.Auto) // Generate the QR code
	if err != nil {
		http.Error(w, "QR Code generation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Scale the QR code to a standard size
	qrCode, err = barcode.Scale(qrCode, 512, 512)
	if err != nil {
		http.Error(w, "QR Code scaling error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var qrCodeBuffer bytes.Buffer
	err = png.Encode(&qrCodeBuffer, qrCode) // Encode the QR code to PNG format
	if err != nil {
		http.Error(w, "QR Code encoding error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the QR code to the database
	_, err = db.Exec("INSERT INTO qrcodes (data, qr_code) VALUES (?, ?)", dataString, qrCodeBuffer.Bytes())
	if err != nil {
		http.Error(w, "Database insert error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the QR code image as the response
	w.Header().Set("Content-Type", "image/png")
	w.Write(qrCodeBuffer.Bytes())
}

// Handler for retrieving stored QR codes
func RetrievePage(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, data FROM qrcodes") // Query all QR codes
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Generate an HTML list of stored QR codes
	var htmlContent string = "<h1>Stored QR Codes</h1><ul>"
	for rows.Next() {
		var id int
		var data string
		rows.Scan(&id, &data)
		htmlContent += "<li><a href='/view/" + strconv.Itoa(id) + "'>" + data + "</a></li>" // Link to view individual QR codes
	}
	htmlContent += "</ul>"

	w.Write([]byte(htmlContent))
}

// Handler for viewing individual QR codes
func ViewQRCode(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/view/"):] // Extract the QR code ID from the URL

	// Validate if the ID is empty
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	row := db.QueryRow("SELECT qr_code, scan_count FROM qrcodes WHERE id = ?", id) // Query the specific QR code

	var qrCodeData []byte
	var scanCount int
	err := row.Scan(&qrCodeData, &scanCount)
	if err != nil {
		http.Error(w, "QR Code not found", http.StatusNotFound)
		return
	}

	// Update the scan count
	_, err = db.Exec("UPDATE qrcodes SET scan_count = scan_count + 1 WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Database update error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the QR code image as the response
	w.Header().Set("Content-Type", "image/png")
	w.Write(qrCodeData)
}

func main() {
	initDB()           // Initialize the database
	defer db.Close()    // Ensure database connection is closed

	// Define HTTP routes
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/generator/", CodePage)
	http.HandleFunc("/retrieve/", RetrievePage)
	http.HandleFunc("/view/", ViewQRCode)

	// Start the HTTP server
	log.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

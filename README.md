## QR Code Generator Project


## Project Overview:

This project is a QR Code Generator web application built using the Go programming language. It allows users to generate QR codes from input text, store them in a SQLite database, and retrieve or view stored QR codes.


Project Setup and How to Run:

Prerequisites

Go installed on your machine.
GCC compiler (required for SQLite support). Install via Mingw-w64 for Windows users.


Installation Steps:

Clone the repository or copy the project files.
git clone <repository-url>
cd <project-directory>

Install the necessary Go dependencies.
go mod tidy

Run the project using the following command:
go run main.go

Open your browser and navigate to:
http://localhost:8080


Directory Structure:

main.go: Contains the backend logic for the application.
generator.html: Frontend template for QR code generation.
qrdata.db: SQLite database file to store QR codes.




## Code Implementation Breakdown

1. **Single Entry Point (main.go)**  
   All routes and logic are implemented within a single Go file `main.go`, ensuring a compact and maintainable project structure.

2. **Database Initialization (initDB())**  
   - Creates a `qrcodes` table if it does not exist.
   - Fields: `id` (primary key), `data` (text), `qr_code` (binary data), and `scan_count`.

3. **QR Code Generation (`/generator/` route, handled by CodePage)**  
   - Generates a QR code for the input text using the Boombuler Barcode package.
   - Encodes the QR code as PNG and stores it in the database.

4. **QR Code Retrieval (`/retrieve/` route, handled by RetrievePage)**  
   - Queries the database and lists all stored QR codes with their associated data.

5. **QR Code Viewing (`/view/{id}` route, handled by ViewQRCode)**  
   - Fetches and displays the QR code for a specific ID.
   - Increments the scan count for each view.

6. **Home Page (`/` route, handled by HomePage)**  
   - Displays the user interface for inputting text and generating QR codes via the `generator.html` template.

7. **Templating (`generator.html`)**  
   - Provides a simple and clean UI for users to input text and generate QR codes.

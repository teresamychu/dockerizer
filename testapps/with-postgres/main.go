package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	db, _ := sql.Open("postgres", "postgres://localhost:5432/mydb")
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello from postgres app!")
	})
	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", nil)
}

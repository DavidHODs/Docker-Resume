package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var db *sql.DB

type dataStruct struct {
	Name    string
	Comment string
}

// GET renders the html page, POST inserts form record into the database
func renderPage(res http.ResponseWriter, req *http.Request) {
	page, err := template.ParseFiles("index.tmpl")
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	rows, err := db.Query(`SELECT name, comment FROM record ORDER BY created_at DESC;`)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	defer rows.Close()

	data := dataStruct{}
	results := []dataStruct{}

	for rows.Next() {
		var name string
		var comment string
		if err := rows.Scan(&name, &comment); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
		data.Name = name
		data.Comment = comment
		results = append(results, data)
	}

	// renders and passes data to the template file
	err = page.Execute(res, results)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func insertRecord(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	name := req.Form.Get("name")
	message := req.Form.Get("message")

	stmt := `insert into record(name, comment) values($1, $2)`
	_, err := db.Exec(stmt, name, message)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func main() {
	var port int = 8000

	router := mux.NewRouter()

	router.HandleFunc("/", renderPage).Methods("GET")
	router.HandleFunc("/insertrecord", insertRecord).Methods("POST")
	// loads static files to be rendered
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	database()

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("127.0.0.1:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Printf("listening on port %d\n", port)

	log.Fatal(srv.ListenAndServe())
}

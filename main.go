package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var HOST string
var PORT string
var USER string
var PASSWORD string
var DB string
var ADDRESS string
var SSL string

var db *sql.DB
var Router *mux.Router
var find *sql.Stmt

type Request struct {
	//Ip       string
	User     string
	Password string
	Path     string
	Action   string
}

func main() {
	flag.StringVar(&HOST, "host", "localhost", "Postgres host address.")
	flag.StringVar(&PORT, "port", "5432", "Postgres port.")
	flag.StringVar(&USER, "user", "postgres", "Postgres user.")
	flag.StringVar(&PASSWORD, "password", "postgres", "Postgres password.")
	flag.StringVar(&DB, "dbname", "rtmp-auth", "Postgres database name.")
	flag.StringVar(&ADDRESS, "address", ":1555", "HTTP server address.")
	flag.StringVar(&SSL, "ssl", "disable", "Postgres SSL mode.")
	flag.Parse()

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", HOST, PORT, USER, PASSWORD, DB, SSL))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	find, err = db.Prepare(
		`SELECT
CASE WHEN EXISTS
 (
   SELECT user FROM users
   WHERE users.user=$1
   AND users.pass=$2
   AND users.path=$3
 )
 THEN 1
 ELSE 0
END`)

	if err != nil {
		panic(err)
	}

	Router = mux.NewRouter()
	Router.HandleFunc("/", auth).Methods("POST")
	http.ListenAndServe(ADDRESS, Router)
}

func auth(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var req Request
	json.Unmarshal(body, &req)

	if req.Action == "read" {
		w.WriteHeader(200)
		return
	}

	var found bool
	find.QueryRow(req.User, req.Password, req.Path).Scan(&found)
	if found {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(401)
	}
}

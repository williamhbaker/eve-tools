package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

type application struct {
	clientID *sqlite.ClientIDModel
}

func main() {
	var newClientID string
	flag.StringVar(&newClientID, "set-id", "", "ID string to save as the client ID - passing this value will reset it in the database")
	flag.Parse()

	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	app := application{
		clientID: &sqlite.ClientIDModel{DB: db},
	}

	if newClientID != "" {
		app.clientID.RegisterID(newClientID)
		fmt.Printf("Set the new client id")
		fmt.Println(app.clientID.GetID()) // remove this for production - not secure
	}

}

func getNewToken() {
	c := make(chan struct{})

	go func() {
		http.HandleFunc("/esi", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
			c <- struct{}{}
		})

		http.ListenAndServe(":4949", nil)
	}()

	<-c
}

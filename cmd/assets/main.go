package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

const callbackURL = "http://localhost:4949/esi"
const listenURL = ":4949"
const scopes = "esi-assets.read_assets.v1"

type application struct {
	clientID     *sqlite.ClientIDModel
	clientSecret *sqlite.ClientSecretModel
	authToken    *sqlite.AuthTokenModel
}

func main() {
	var newClientID string
	var newClientSecret string
	var addCharacter bool

	flag.StringVar(&newClientID, "id", "", "ID string to save as the client ID - passing this value will reset it in the database")
	flag.StringVar(&newClientSecret, "secret", "", "String value for the client secret - passing this value will reset it in the database")
	flag.BoolVar(&addCharacter, "add-char", false, "Set true if you want to register a new character with the application")
	flag.Parse()

	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	app := application{
		clientID:     &sqlite.ClientIDModel{DB: db},
		clientSecret: &sqlite.ClientSecretModel{DB: db},
		authToken:    &sqlite.AuthTokenModel{DB: db},
	}

	if newClientID != "" {
		app.clientID.RegisterID(newClientID)
		fmt.Println("New client id set")
	}

	if newClientSecret != "" {
		app.clientSecret.RegisterSecret(newClientSecret)
		fmt.Println("New client secret set")
	}

	if addCharacter {
		fmt.Printf("Login URL is: %q\n", loginURL(callbackURL, app.clientID.GetID(), scopes))

		gotToken := lib.GetNewToken(listenURL, app.clientID.GetID(), app.clientSecret.GetSecret())
		token := models.AuthToken{
			AccessToken:  gotToken.AccessToken,
			ExpiresIn:    gotToken.ExpiresIn,
			RefreshToken: gotToken.RefreshToken,
			Issued:       gotToken.Issued,
		}

		app.authToken.RegisterToken(token)
		fmt.Println(app.authToken.GetToken())
	}
}

package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/atotto/clipboard"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wbaker85/eve-tools/pkg/lib"
	"github.com/wbaker85/eve-tools/pkg/models"
	"github.com/wbaker85/eve-tools/pkg/models/sqlite"
)

const callbackURL = "http://localhost:4949/esi"
const listenURL = ":4949"
const scopes = "esi-assets.read_assets.v1 esi-markets.read_character_orders.v1"

const forgeRegionID = 10000002
const jitaStationID = 60003760
const perimiterTTTStationID = 1028858195912

type application struct {
	api             *lib.Esi
	orders          *sqlite.OrderModel
	clientID        *sqlite.ClientIDModel
	clientSecret    *sqlite.ClientSecretModel
	authToken       *sqlite.AuthTokenModel
	characterOrders *sqlite.CharacterOrderModel
	characterAssets *sqlite.CharacterAssetModel
	charID          int
}

func main() {
	var newClientID string
	var newClientSecret string
	var addCharacter bool
	var smallSellThreshold float64
	var uaString string

	flag.StringVar(&newClientID, "id", "", "ID string to save as the client ID - passing this value will reset it in the database")
	flag.StringVar(&newClientSecret, "secret", "", "String value for the client secret - passing this value will reset it in the database")
	flag.BoolVar(&addCharacter, "add-char", false, "Set true if you want to register a new character with the application")
	flag.StringVar(&uaString, "ua", "", "The string to use as the user agent for ESI API calls - usually an email address. Provide this to update prices before doing the rest.")
	flag.Float64Var(&smallSellThreshold, "small-sell", 5000000, "Threshold for what determines if a sell order is small or not. Default is 5,000,000.")
	flag.Parse()

	db, _ := sql.Open("sqlite3", "./data.db")
	defer db.Close()

	api := lib.Esi{
		Client:          http.DefaultClient,
		UserAgentString: uaString,
	}

	app := application{
		api:             &api,
		orders:          &sqlite.OrderModel{DB: db},
		clientID:        &sqlite.ClientIDModel{DB: db},
		clientSecret:    &sqlite.ClientSecretModel{DB: db},
		authToken:       &sqlite.AuthTokenModel{DB: db},
		characterOrders: &sqlite.CharacterOrderModel{DB: db},
		characterAssets: &sqlite.CharacterAssetModel{DB: db},
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
		fmt.Println(string(app.authorizedRequest(charIDURL, "GET", false)))
	}

	if newClientID == "" && newClientSecret == "" && !addCharacter {
		if uaString != "" {
			app.updateOrdersByRegion(forgeRegionID, jitaStationID, perimiterTTTStationID)
		}

		var charData map[string]interface{}
		d := app.authorizedRequest(charIDURL, "GET", false)
		json.Unmarshal(d, &charData)

		app.charID = int(charData["CharacterID"].(float64))
		fmt.Printf("Got character ID: %d\n", app.charID)

		fmt.Println("Getting Character orders...")
		app.populateCharacterOrders()
		fmt.Println("Getting Character assets...")
		app.populateCharacterAssets()

		hangarAssets := app.characterAssets.GetGrouped()
		escrowAssets := app.characterOrders.SellingInventory()

		rules := parseRules("./trade_rules.csv")
		prices := app.orders.BuyPriceTable(jitaStationID, perimiterTTTStationID)

		pricedOut := tooExpensive(prices, rules)
		tooMuch := tooMuchInventory(hangarAssets, escrowAssets, rules)
		allBuys := app.characterOrders.Orders(true)
		shouldBuy := shouldBeBuying(rules, pricedOut, tooMuch)
		allSells := app.characterOrders.Orders(false)
		shouldSell := shouldBeSelling(hangarAssets, rules)

		var finalString strings.Builder

		finalString.WriteString(printCategory("1 - bad buys", sliceDiff(allBuys, shouldBuy)))
		finalString.WriteString(printCategory("2 - should be buying", shouldBuy))
		finalString.WriteString(printCategory("3 - should be buying but am not", sliceDiff(shouldBuy, allBuys)))
		finalString.WriteString(printCategory("4 - should be selling but am not", sliceDiff(shouldSell, allSells)))
		finalString.WriteString(printCategory("5 - should be selling", shouldSell))
		finalString.WriteString(printCategory("6 - all sells", allSells))
		finalString.WriteString(printCategory("7 - all buys", allBuys))
		finalString.WriteString(printCategory("8 - small sell orders", app.characterOrders.SmallSells(smallSellThreshold)))

		fmt.Println(finalString.String())
		clipboard.WriteAll(finalString.String())
	}
}

func (app *application) updateOrdersByRegion(regionID, sellStationID, buyStationID int) {
	orders := app.api.AllOrders(regionID, -1)

	sellStationPrices := lib.AggregateOrders(orders, sellStationID)
	buyStationPrices := lib.AggregateOrders(orders, buyStationID)

	app.api.AddNames(sellStationPrices, 1000)
	app.api.AddNames(buyStationPrices, 1000)

	app.orders.LoadData(sellStationID, sellStationPrices)
	app.orders.LoadData(buyStationID, buyStationPrices)
}

func printCategory(catName string, items []string) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("+ %v\n", catName))

	for _, i := range items {
		b.WriteString(fmt.Sprintf("- %v\n", i))
	}

	return b.String()

	// fmt.Printf("+ %v\n", catName)

	// for _, i := range items {
	// 	fmt.Printf("- %v\n", i)
	// }
}

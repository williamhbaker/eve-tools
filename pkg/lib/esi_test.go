package lib

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type TestClient struct{}

func (t TestClient) Do(r *http.Request) (*http.Response, error) {
	agent := r.Header.Get("User-Agent")
	url := r.URL.String()

	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(agent + " - " + url)),
	}, nil
}

func TestGet(t *testing.T) {
	c := TestClient{}

	s := "user@addr.com"
	u := "https://www.whatever.com"

	e := esi{client: c, userAgentString: s}

	res, _, _ := e.get(u)
	got := string(res)
	want := s + " - " + u

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestAggregateOrders(t *testing.T) {
	testData := []map[string]interface{}{
		{
			"type_id":      1234.0,
			"is_buy_order": "true",
			"price":        100.0,
		},
		{
			"type_id":      1234.0,
			"is_buy_order": "true",
			"price":        101.0,
		},
		{
			"type_id":      1234.0,
			"is_buy_order": "true",
			"price":        99.0,
		},
		{
			"type_id":      1234.0,
			"is_buy_order": "false",
			"price":        120.0,
		},
		{
			"type_id":      1234.0,
			"is_buy_order": "false",
			"price":        118.0,
		},
		{
			"type_id":      321.0,
			"is_buy_order": "false",
			"price":        400.0,
		},
		{
			"type_id":      321.0,
			"is_buy_order": "true",
			"price":        300.0,
		},
	}

	got := aggregateOrders(testData)
	want := map[int]*item{
		1234: {
			id:        1234,
			sellPrice: 118,
			buyPrice:  101,
		},
		321: {
			id:        321,
			sellPrice: 400,
			buyPrice:  300,
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v want %#v", got, want)
	}
}

// func TestAllOrders(t *testing.T) {
// 	c := http.DefaultClient

// 	s := "wbaker@gmail.com"
// 	u := "https://esi.evetech.net/v1/markets/10000002/orders?page=789"

// 	e := esi{client: c, userAgentString: s}

// 	_, status, _ := e.get(u)

// 	e.AllOrders(10000002)

// 	t.Errorf(strconv.Itoa(status))
// }

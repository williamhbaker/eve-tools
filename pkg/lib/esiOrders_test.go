package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

type TestClient struct {
	DoFunc func(r *http.Request) (*http.Response, error)
}

func (t *TestClient) Do(r *http.Request) (*http.Response, error) {
	return t.DoFunc(r)
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

func TestAllOrders(t *testing.T) {
	c := &TestClient{DoFunc: func(r *http.Request) (*http.Response, error) {
		re := regexp.MustCompile(`\d+$`)
		page := re.FindString(r.URL.String())
		pageNum, _ := strconv.Atoi(page)

		if pageNum <= 2 {
			return &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(page)),
			}, nil
		}

		return &http.Response{
			Status: "404",
			Body:   ioutil.NopCloser(strings.NewReader("ERROR ERROR ERROR")),
		}, nil
	},
	}

	e := esi{client: c}

	res, _, _ := e.get("https://esi.evetech.net/v1/markets/1234/orders?page=22")

	fmt.Println("res", string(res))

	t.Errorf(string(res))
}

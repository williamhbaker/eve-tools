package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

var testData = []map[string]interface{}{
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

func TestAggregateOrders(t *testing.T) {
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
	stationID := 1234
	midPoint := len(testData) / 2
	urlSpy := []string{}

	c := &TestClient{DoFunc: func(r *http.Request) (*http.Response, error) {
		urlSpy = append(urlSpy, r.URL.String())

		re := regexp.MustCompile(`\d+$`)
		page := re.FindString(r.URL.String())
		pageNum, _ := strconv.Atoi(page)

		var status int
		var body io.ReadCloser

		switch pageNum {
		case 1:
			status = 200
			jsonData, _ := json.Marshal(testData[:midPoint])
			body = ioutil.NopCloser(bytes.NewReader(jsonData))
		case 2:
			status = 200
			jsonData, _ := json.Marshal(testData[midPoint:])
			body = ioutil.NopCloser(bytes.NewReader(jsonData))
		default:
			status = 404
			body = ioutil.NopCloser(strings.NewReader(page))
		}

		return &http.Response{
			StatusCode: status,
			Body:       body,
		}, nil
	},
	}

	e := esi{client: c}

	data := e.AllOrders(stationID)

	wantUrls := []string{}

	for idx := 1; idx <= 3; idx++ {
		wantUrls = append(wantUrls, fmt.Sprintf(ordersFragment, stationID, idx))
	}

	if !reflect.DeepEqual(urlSpy, wantUrls) {
		t.Errorf("\nurls wrong\ngot %v\nwant %v", urlSpy, wantUrls)
	}

	if !reflect.DeepEqual(data, testData) {
		t.Errorf("\nresponse data wrong\ngot %#v\nwant %#v", data, testData)
	}
}

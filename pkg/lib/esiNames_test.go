package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/wbaker85/eve-tools/pkg/models"
)

var (
	resJSON = `[
		{
				"category": "inventory_type",
				"id": 52275,
				"name": "Jita Protest YC113"
		},
		{
				"category": "inventory_type",
				"id": 30486,
				"name": "Sisters Combat Scanner Probe"
		},
		{
				"category": "inventory_type",
				"id": 3253,
				"name": "Inherent Implants 'Squire' Energy Grid Upgrades EU-702"
		}
	]`
	wantNamesMap = map[int]string{
		3253:  "Inherent Implants 'Squire' Energy Grid Upgrades EU-702",
		52275: "Jita Protest YC113",
		30486: "Sisters Combat Scanner Probe",
	}
)

func TestAddNames(t *testing.T) {
	testFunc := func(r *http.Request) (*http.Response, error) {
		data, _ := ioutil.ReadAll(r.Body)
		var body []int
		json.Unmarshal(data, &body)

		res := make([]struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}, len(body))

		for idx := 0; idx < len(body); idx++ {
			res[idx].ID = body[idx]
			res[idx].Name = wantNamesMap[body[idx]]
		}

		jsonRes, _ := json.Marshal(res)

		return &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader(jsonRes)),
		}, nil
	}

	c := newTestClient(testFunc)

	e := Esi{
		Client:          c,
		UserAgentString: "somebody@whatever.com",
	}

	got := map[int]*models.OrderItem{
		3253:  {},
		52275: {},
		30486: {},
	}

	want := map[int]*models.OrderItem{
		3253:  {Name: "Inherent Implants 'Squire' Energy Grid Upgrades EU-702"},
		52275: {Name: "Jita Protest YC113"},
		30486: {Name: "Sisters Combat Scanner Probe"},
	}

	e.AddNames(got, 2)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot %#v\nwant %#v", got, want)
	}
}

func TestItemNameList(t *testing.T) {
	var spyBody string

	testFunc := func(r *http.Request) (*http.Response, error) {
		buf := bytes.NewBuffer([]byte{})
		buf.ReadFrom(r.Body)
		spyBody = buf.String()

		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(resJSON)),
		}, nil
	}

	c := newTestClient(testFunc)

	e := Esi{
		Client:          c,
		UserAgentString: "somebody@whatever.com",
	}

	listSlice := []int{3253, 52275, 30486}
	listString := "[3253,52275,30486]"

	got := e.itemNameList(listSlice)

	if spyBody != listString {
		t.Errorf("\ngot %q\nwant %q", spyBody, listString)
	}

	if !reflect.DeepEqual(got, wantNamesMap) {
		t.Errorf("\ngot %v\nwant %v", got, wantNamesMap)
	}

}

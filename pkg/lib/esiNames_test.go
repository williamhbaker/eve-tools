package lib

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

const resJSON = `[
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

func TestIemNameList(t *testing.T) {
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
	want := map[int]string{
		3253:  "Inherent Implants 'Squire' Energy Grid Upgrades EU-702",
		52275: "Jita Protest YC113",
		30486: "Sisters Combat Scanner Probe",
	}

	got := e.itemNameList(listSlice)

	if spyBody != listString {
		t.Errorf("\ngot %q\nwant %q", spyBody, listString)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot %v\nwant %v", got, want)
	}

}

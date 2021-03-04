package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"testing"

	"github.com/wbaker85/eve-tools/pkg/models"
)

func TestVolumeForItems(t *testing.T) {
	testRegionID := 1234
	testItems := []int{1, 2, 3}

	testResData1 := []struct {
		OrderCount int `json:"order_count"`
		Volume     int `json:"volume"`
	}{
		{
			OrderCount: 10,
			Volume:     20,
		},
		{
			OrderCount: 10,
			Volume:     20,
		},
		{
			OrderCount: 10,
			Volume:     20,
		},
	}

	testResData2 := []struct {
		OrderCount int `json:"order_count"`
		Volume     int `json:"volume"`
	}{
		{
			OrderCount: 10,
			Volume:     20,
		},
		{
			OrderCount: 20,
			Volume:     40,
		},
		{
			OrderCount: 30,
			Volume:     60,
		},
	}

	testResData3 := []struct {
		OrderCount int `json:"order_count"`
		Volume     int `json:"volume"`
	}{
		{
			OrderCount: 20,
			Volume:     20,
		},
		{
			OrderCount: 30,
			Volume:     30,
		},
		{
			OrderCount: 40,
			Volume:     40,
		},
	}

	client := newTestClient(func(r *http.Request) (*http.Response, error) {
		re := regexp.MustCompile(`\d+$`)
		item := re.FindString(r.URL.String())
		itemID, _ := strconv.Atoi(item)

		var body io.ReadCloser

		switch itemID {
		case 1:
			jsonData, _ := json.Marshal(testResData1)
			body = ioutil.NopCloser(bytes.NewReader(jsonData))
		case 2:
			jsonData, _ := json.Marshal(testResData2)
			body = ioutil.NopCloser(bytes.NewReader(jsonData))
		default:
			jsonData, _ := json.Marshal(testResData3)
			body = ioutil.NopCloser(bytes.NewReader(jsonData))
		}

		return &http.Response{
			Body: body,
		}, nil
	})

	api := Esi{
		Client: client,
	}

	got := api.VolumeForItems(testRegionID, testItems)
	want := []models.ItemHistoryData{
		{
			RegionID:  1234,
			ItemID:    1,
			NumDays:   3,
			OrdersAvg: 10,
			VolumeAvg: 20,
		},
		{
			RegionID:  1234,
			ItemID:    2,
			NumDays:   3,
			OrdersAvg: 20,
			VolumeAvg: 40,
		},
		{
			RegionID:  1234,
			ItemID:    3,
			NumDays:   3,
			OrdersAvg: 30,
			VolumeAvg: 30,
		},
	}

	sort.Slice(got, func(i, j int) bool {
		return got[i].ItemID < got[j].ItemID
	})

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot: \n%#v, \nwant: \n%#v", got, want)
	}

}

func TestAvgForPeriod(t *testing.T) {
	testData := []itemDailyVolume{
		{
			OrderCount: 10,
			Volume:     10,
		},
		{
			OrderCount: 20,
			Volume:     40,
		},
		{
			OrderCount: 10,
			Volume:     10,
		},
		{
			OrderCount: 30,
			Volume:     60,
		},
		{
			OrderCount: 5,
			Volume:     5,
		},
	}

	tests := []struct {
		name   string
		data   []itemDailyVolume
		length int
		want   models.ItemHistoryData
	}{
		{
			"empty list with period > list length",
			testData[len(testData):],
			10,
			models.ItemHistoryData{NumDays: 0, OrdersAvg: 0, VolumeAvg: 0},
		},
		{
			"non-empty list with period > list length",
			testData,
			10,
			models.ItemHistoryData{NumDays: 5, OrdersAvg: 15, VolumeAvg: 25},
		},
		{
			"period < list length",
			testData,
			3,
			models.ItemHistoryData{NumDays: 3, OrdersAvg: 15, VolumeAvg: 25},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := avgForPeriod(tt.data, tt.length)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v want %#v", got, tt.want)
			}
		})
	}

}

func TestTruncateLastN(t *testing.T) {
	tests := []struct {
		name string
		vals []itemDailyVolume
		n    int
		want []itemDailyVolume
	}{
		{
			"short list, larger N",
			[]itemDailyVolume{{Volume: 1}, {Volume: 2}, {Volume: 3}},
			10,
			[]itemDailyVolume{{Volume: 1}, {Volume: 2}, {Volume: 3}},
		},
		{
			"short list, smaller N",
			[]itemDailyVolume{{Volume: 1}, {Volume: 2}, {Volume: 3}},
			2,
			[]itemDailyVolume{{Volume: 2}, {Volume: 3}},
		},
		{
			"short list, even smaller N",
			[]itemDailyVolume{{Volume: 1}, {Volume: 2}, {Volume: 3}},
			1,
			[]itemDailyVolume{{Volume: 3}},
		},
		{
			"short list, smallest possible N",
			[]itemDailyVolume{{Volume: 1}, {Volume: 2}, {Volume: 3}},
			0,
			[]itemDailyVolume{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateLastN(tt.vals, tt.n)
			assertSlices(t, got, tt.want)
		})
	}

}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name    string
		findVal float64
		vals    []float64
		want    int
	}{
		{"Item is at the beginning", 1.0, []float64{1.0, 2.0, 3.0}, 0},
		{"Item is at the end", 3.0, []float64{1.0, 2.0, 3.0}, 2},
		{"Item not found", 4.0, []float64{1.0, 2.0, 3.0}, -1},
	}

	for _, tt := range tests {
		assertInts(t, indexOf(tt.findVal, tt.vals), tt.want)
	}
}

func TestRemovedByIndexes(t *testing.T) {
	testHist := []itemDailyVolume{
		{Volume: 100}, {Volume: 200}, {Volume: 300}, {Volume: 400}, {Volume: 500},
	}

	tests := []struct {
		name string
		hist []itemDailyVolume
		i    []int
		want []itemDailyVolume
	}{
		{
			"Removes a list of indexes",
			testHist,
			[]int{1, 4},
			[]itemDailyVolume{{Volume: 100}, {Volume: 300}, {Volume: 400}},
		},
		{
			"Empty index list",
			testHist,
			[]int{},
			testHist,
		},
		{
			"Remove all the indexes",
			testHist,
			[]int{0, 1, 2, 3, 4},
			[]itemDailyVolume{},
		},
		{
			"Remove the first one",
			testHist,
			[]int{0},
			testHist[1:],
		},
		{
			"Remove the last one",
			testHist,
			[]int{4},
			testHist[:len(testHist)-1],
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeByIndexes(tt.hist, tt.i)
			assertSlices(t, got, tt.want)
		})
	}
}

func TestFindOutliers(t *testing.T) {
	tests := []struct {
		name string
		hist []itemDailyVolume
		want []int
	}{
		{
			"No outliers",
			[]itemDailyVolume{{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300}},
			[]int{},
		},
		{
			"Two outliers - same value",
			[]itemDailyVolume{
				{Volume: 30000}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 30000},
			},
			[]int{0, 15},
		},
		{
			"One outlier at the start",
			[]itemDailyVolume{
				{Volume: 30000}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300},
			},
			[]int{0},
		},
		{
			"One outlier at the end",
			[]itemDailyVolume{
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 30000},
			},
			[]int{14},
		},
		{
			"One outlier in the middle somewhere",
			[]itemDailyVolume{
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 30000}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300},
			},
			[]int{5},
		},
		{
			"Two outliers with different values",
			[]itemDailyVolume{
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 30000}, {Volume: 300}, {Volume: 300},
				{Volume: 300}, {Volume: 300}, {Volume: 300}, {Volume: 300},
				{Volume: 30000}, {Volume: 300}, {Volume: 300},
			},
			[]int{5, 12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findOutliers(tt.hist)
			assertSlices(t, got, tt.want)
		})
	}
}

func TestYearlyMinMax(t *testing.T) {
	testData := []itemDailyVolume{
		{Highest: 100.5},
		{Highest: 90.1},
		{Highest: 110.7},
		{Highest: 76.3},
		{Highest: 100.5},
	}

	wantMax := 110.7
	wantMin := 76.3

	gotMax, gotMin := yearlyMinMax(testData)

	if gotMax != wantMax {
		t.Errorf("got %.2f want %.2f", gotMax, wantMax)
	}

	if gotMin != wantMin {
		t.Errorf("got %.2f want %.2f", gotMin, wantMin)
	}
}

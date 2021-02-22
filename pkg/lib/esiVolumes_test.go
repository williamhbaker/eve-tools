package lib

import "testing"

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

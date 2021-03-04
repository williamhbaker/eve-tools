package lib

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/montanaflynn/stats"
	"github.com/wbaker85/eve-tools/pkg/models"
)

const volumesFragment = `https://esi.evetech.net/v1/markets/%d/history?type_id=%d`

type itemDailyVolume struct {
	Date       string  `json:"date"`
	Highest    float64 `json:"highest"`
	Lowest     float64 `json:"lowest"`
	OrderCount int     `json:"order_count"`
	Volume     int     `json:"volume"`
}

// VolumeForItems gets the volume information for many items, and return a slice
// containing the results.
func (e *Esi) VolumeForItems(regionID int, itemIDs []int) []models.ItemHistoryData {
	output := []models.ItemHistoryData{}

	var count int

	for _, itemID := range itemIDs {
		count++
		fmt.Printf("Getting volumes for item %d of %d...\n", count, len(itemIDs))
		output = append(output, e.VolumeForItem(regionID, itemID))
	}

	return output
}

// VolumeForItem gets the volume information for a single item
func (e *Esi) VolumeForItem(regionID, itemID int) models.ItemHistoryData {
	u := fmt.Sprintf(volumesFragment, regionID, itemID)

	bytes, _, _ := e.get(u)

	var data []itemDailyVolume

	json.Unmarshal(bytes, &data)

	maxSellPrice, minSellPrice := yearlyMinMax(data, true)
	maxBuyPrice, minBuyPrice := yearlyMinMax(data, false)

	data = truncateLastN(data, 30)
	outliers := findOutliers(data)
	cleaned := removeByIndexes(data, outliers)

	averages := avgForPeriod(cleaned, 7)
	averages.RegionID = regionID
	averages.ItemID = itemID
	averages.YearMaxSell = maxSellPrice
	averages.YearMinSell = minSellPrice
	averages.YearMaxBuy = maxBuyPrice
	averages.YearMinBuy = minBuyPrice

	return averages
}

func yearlyMinMax(data []itemDailyVolume, findSellPrice bool) (float64, float64) {
	offset := int(math.Min(365, float64(len(data))))

	foundMin := math.MaxFloat64
	foundMax := -math.MaxFloat64

	for idx := len(data) - offset; idx < len(data); idx++ {
		item := data[idx]
		var val float64

		if findSellPrice {
			val = item.Highest
		} else {
			val = item.Lowest
		}

		foundMin = math.Min(val, foundMin)
		foundMax = math.Max(val, foundMax)
	}

	return foundMax, foundMin
}

func truncateLastN(data []itemDailyVolume, num int) []itemDailyVolume {
	var n int

	if len(data) < num {
		n = len(data)
	} else {
		n = num
	}

	return data[len(data)-n:]
}

func avgForPeriod(data []itemDailyVolume, length int) models.ItemHistoryData {
	var n int

	totalOrders := 0
	totalVolume := 0

	if len(data) < length {
		n = len(data)
	} else {
		n = length
	}

	if n == 0 {
		return models.ItemHistoryData{}
	}

	for idx := len(data) - n; idx < len(data); idx++ {
		totalOrders += data[idx].OrderCount
		totalVolume += data[idx].Volume
	}

	return models.ItemHistoryData{
		NumDays:   n,
		OrdersAvg: totalOrders / n,
		VolumeAvg: totalVolume / n,
	}
}

func findOutliers(data []itemDailyVolume) []int {
	volumes := make([]int, len(data))

	for idx, val := range data {
		volumes[idx] = val.Volume
	}

	d := stats.LoadRawData(volumes)
	outliers, _ := d.QuartileOutliers()
	mild := outliers.Mild
	extreme := outliers.Extreme
	combined := append(mild, extreme...)

	outlierIdx := []int{}

	for _, val := range combined {
		found := indexOf(val, d)
		d[found] = -2

		if found > -1 {
			outlierIdx = append(outlierIdx, found)
		}
	}

	return outlierIdx
}

func removeByIndexes(h []itemDailyVolume, i []int) []itemDailyVolume {
	removeMap := make(map[int]bool)

	for _, idx := range i {
		removeMap[idx] = true
	}

	output := []itemDailyVolume{}

	for idx, val := range h {
		if !removeMap[idx] {
			output = append(output, val)
		}
	}

	return output
}

func indexOf(val float64, s []float64) int {
	for idx, sVal := range s {
		if val == sVal {
			return idx
		}
	}

	return -1
}

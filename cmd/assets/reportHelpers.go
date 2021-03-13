package main

import "github.com/wbaker85/eve-tools/pkg/models"

func tooExpensive(prices map[string]float64, rules []itemRule) []string {
	output := []string{}

	for _, val := range rules {
		thisName := val.ItemName
		max := val.BuyTargetPrice

		if prices[thisName] >= max {
			output = append(output, thisName)
		}
	}

	return output
}

func tooMuchInventory(hangar, escrow []models.CharacterAsset, rules []itemRule) []string {
	output := []string{}
	counts := combinedAssetCount(hangar, escrow)

	for _, v := range rules {
		thisName := v.ItemName
		max := v.MaxInventory

		if counts[thisName] >= max {
			output = append(output, thisName)
		}
	}

	return output
}

func combinedAssetCount(l1, l2 []models.CharacterAsset) map[string]int {
	output := make(map[string]int)

	for _, v := range l1 {
		output[v.Name] += v.Quantity
	}

	for _, v := range l2 {
		output[v.Name] += v.Quantity
	}

	return output
}

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

func shouldBeBuying(rules []itemRule, pricedOut, tooMuch []string) []string {
	output := []string{}
	m := stringSliceToMap(pricedOut, tooMuch)

	for _, val := range rules {
		if !m[val.ItemName] {
			output = append(output, val.ItemName)
		}
	}

	return output
}

func buyingButShouldNotBe(am, should []string) []string {
	output := []string{}

	sMap := make(map[string]struct{})
	for _, val := range should {
		sMap[val] = struct{}{}
	}

	for _, val := range am {
		_, ok := sMap[val]
		if !ok {
			output = append(output, val)
		}
	}

	return output
}

func stringSliceToMap(s1, s2 []string) map[string]bool {
	output := make(map[string]bool)

	for _, val := range s1 {
		output[val] = true
	}

	for _, val := range s2 {
		output[val] = true
	}

	return output
}

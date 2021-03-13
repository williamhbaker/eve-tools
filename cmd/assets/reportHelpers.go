package main

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

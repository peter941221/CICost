package pricing

func ChargedMinutes(totalBillable, freeTier, alreadyUsed float64) float64 {
	remainingFree := freeTier - alreadyUsed
	if remainingFree < 0 {
		remainingFree = 0
	}
	charged := totalBillable - remainingFree
	if charged < 0 {
		return 0
	}
	return charged
}


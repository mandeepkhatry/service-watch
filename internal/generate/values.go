package generate

import (
	"math/rand"
	"time"
)

func GenerateInteger(properties map[string]interface{}) int {

	rand.Seed(time.Now().UnixNano())

	var min, max int

	if _, present := properties["maximum"]; !present {
		max = 100
	} else {
		max = properties["maximum"].(int)
	}

	if _, present := properties["minimum"]; !present {
		min = 0
	} else {
		min = properties["minimum"].(int)
	}

	return rand.Intn(max-min+1) + min
}

func GenerateFloat(properties map[string]interface{}) float64 {

	rand.Seed(time.Now().UnixNano())

	var min, max float64

	if _, present := properties["maximum"]; !present {
		max = 100.0
	} else {
		max = properties["maximum"].(float64)
	}

	if _, present := properties["minimum"]; !present {
		min = 0.0
	} else {
		min = properties["minimum"].(float64)
	}

	return rand.Float64()*(max-min) + min
}

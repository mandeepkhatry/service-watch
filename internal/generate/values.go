package generate

import (
	"fmt"
	"math/rand"
	"service-watch/internal/def"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/lucasjones/reggen"
)

func GenerateBoolean(properties map[string]interface{}) bool {
	return true
}

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

	if _, present := properties["enum"]; present {
		enum := properties["enum"].([]int)
		if len(enum) > 0 {
			return enum[0]
		}
		return 0
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

	if _, present := properties["enum"]; present {
		enum := properties["enum"].([]float64)
		if len(enum) > 0 {
			return enum[0]
		}
		return 0.0
	}

	return rand.Float64()*(max-min) + min
}

func GenerateString(properties map[string]interface{}) string {
	rand.Seed(time.Now().UnixNano())

	var minLength, maxLength int

	if _, present := properties["maxLength"]; !present {
		maxLength = 100.0
	} else {
		maxLength = properties["maxLength"].(int)
	}

	if _, present := properties["minLength"]; !present {
		minLength = 0.0
	} else {
		minLength = properties["minLength"].(int)
	}

	if _, present := properties["enum"]; present {
		enum := properties["enum"].([]string)
		if len(enum) > 0 {
			return enum[0]
		}
		return "test"
	}

	stringLength := rand.Intn(maxLength-minLength+1) + minLength

	b := make([]byte, stringLength)
	for i := range b {
		b[i] = def.CharSet[rand.Intn(len(def.CharSet))]
	}
	return string(b)
}

func GenerateEmail() string {
	var properties = map[string]interface{}{
		"maxLength": 10,
		"minLength": 5,
	}
	return GenerateString(properties) + "@xyz.com"
}

func GenerateNumberArray(properties map[string]interface{}) []float64 {

	var maxItems int

	if _, present := properties["maxItems"]; !present {
		maxItems = 10
	} else {
		maxItems = properties["maxItems"].(int)
	}

	// if _, present := properties["minItems"]; !present {
	// 	minItems = 0
	// } else {
	// 	minItems = properties["minItems"].(int)
	// }

	numberArray := make([]float64, 0)

	var arrayProperties = map[string]interface{}{
		"maximum": 10.0,
		"minimum": 5.0,
	}

	for i := 0; i < maxItems; i++ {
		numberArray = append(numberArray, GenerateFloat(arrayProperties))
	}

	return numberArray

}

func GenerateStringArray(properties map[string]interface{}) []string {

	var maxItems int

	if _, present := properties["maxItems"]; !present {
		maxItems = 10
	} else {
		maxItems = properties["maxItems"].(int)
	}

	// if _, present := properties["minItems"]; !present {
	// 	minItems = 0
	// } else {
	// 	minItems = properties["minItems"].(int)
	// }

	stringArray := make([]string, 0)

	var arrayProperties = map[string]interface{}{
		"maxLength": 10.0,
		"minLength": 5.0,
	}

	for i := 0; i < maxItems; i++ {
		stringArray = append(stringArray, GenerateString(arrayProperties))
	}

	return stringArray

}

func GenerateStringType(stringType string) string {
	if stringType == "ipv4" {
		return randomdata.IpV4Address()
	} else if stringType == "ipv6" {
		return randomdata.IpV6Address()
	} else if stringType == "date-time" {
		return "2018-11-13T20:20:39+00:00"
	} else if stringType == "time" {
		return "20:20:39+00:00"
	} else if stringType == "date" {
		return "2018-11-13"
	}
	return "unknown field"
}

func GenerateRegex(regex string) string {
	str, _ := reggen.Generate(regex, 1)
	return str
}

func GenerateArray(properties map[string]interface{}) interface{} {

	if _, typePresent := properties["items"]; !typePresent {
		var maxItems int

		if _, present := properties["maxItems"]; !present {
			maxItems = 10
		} else {
			maxItems = properties["maxItems"].(int)
		}

		// if _, present := properties["minItems"]; !present {
		// 	minItems = 0
		// } else {
		// 	minItems = properties["minItems"].(int)
		// }

		numberArray := make([]float64, 0)

		var arrayProperties = map[string]interface{}{
			"maximum": 10.0,
			"minimum": 5.0,
		}

		for i := 0; i < maxItems; i++ {
			numberArray = append(numberArray, GenerateFloat(arrayProperties))
		}

		return numberArray
	}

	arrayItemType := fmt.Sprintf("%T", properties["items"])

	if arrayItemType == "[]map[string]interface {}" {

		resultingArray := make([]interface{}, 0)

		for _, eachTypeProperties := range properties["items"].([]map[string]interface{}) {

			eachArrayType := eachTypeProperties["type"].(string)

			if eachArrayType == "number" {
				resultingArray = append(resultingArray, GenerateFloat(eachTypeProperties))

			} else if eachArrayType == "string" {
				resultingArray = append(resultingArray, GenerateString(eachTypeProperties))
			}

		}

		//Mix Type Array

	} else if arrayItemType == "map[string]interface {}" {

		//Single Type Array

		itemType := properties["items"].(map[string]interface{})["type"].(string)

		if itemType == "number" {

			return GenerateNumberArray(properties)

		} else if itemType == "string" {

			return GenerateNumberArray(properties)

		}

	}

	return []string{}

}

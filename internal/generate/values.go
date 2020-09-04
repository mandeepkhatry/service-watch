package generate

import (
	"fmt"
	"math/rand"
	"service-watch/internal/def"
	"strconv"
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

	if _, present := properties["exclusiveMaximum"]; present {
		max = properties["exclusiveMaximum"].(int) - 1
	}

	if _, present := properties["exclusiveMinimum"]; present {
		min = properties["exclusiveMinimum"].(int) + 1
	}

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

	if _, present := properties["multipleOf"]; present {
		multipleOf := properties["multipleOf"].(int)
		i := 0
		for {
			eachMultiple := multipleOf * i
			if eachMultiple >= min && eachMultiple <= max {
				return multipleOf * i
			} else if eachMultiple > max {
				return 0
			}

		}
	}

	return rand.Intn(max-min+1) + min
}

func GenerateFloat(properties map[string]interface{}) float64 {

	rand.Seed(time.Now().UnixNano())

	var min, max float64

	if _, present := properties["exclusiveMaximum"]; present {
		max = properties["exclusiveMaximum"].(float64) - 0.5
	}

	if _, present := properties["exclusiveMinimum"]; present {
		min = properties["exclusiveMinimum"].(float64) + 0.5
	}

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

	if _, present := properties["multipleOf"]; present {
		multipleOf := properties["multipleOf"].(float64)
		i := 0.0
		for {
			eachMultiple := multipleOf * i
			if eachMultiple >= min && eachMultiple <= max {
				return multipleOf * i
			} else if eachMultiple > max {
				return 0.0
			}

		}
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

	if _, present := properties["format"]; present {
		return GenerateStringFormat(properties["format"].(string))
	}

	if _, present := properties["pattern"]; present {
		return GenerateRegex(properties["pattern"].(string))
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

func GenerateStringFormat(stringType string) string {
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
	} else if stringType == "email" {
		return GenerateEmail()
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

			resultingArray = append(resultingArray, FieldToGenerator[eachArrayType])

		}

		//Mix Type Array

	} else if arrayItemType == "map[string]interface {}" {

		//Single Type Array

		itemType := properties["items"].(map[string]interface{})["type"].(string)

		return FieldToGenerator[itemType]

	}

	return []string{}

}

func GenerateObject(properties map[string]interface{}) map[string]interface{} {

	if _, present := properties["properties"]; !present {
		if _, minPropertiesPresent := properties["minProperties"]; minPropertiesPresent {
			minPropertiesPresent := properties["minProperties"].(int)

			randomObject := make(map[string]interface{})

			for i := 0; i < minPropertiesPresent; i++ {
				randomObject["test"+strconv.Itoa(i)] = i
			}

			return randomObject
		}

		return map[string]interface{}{
			"k1": "v1",
			"k2": "v2",
		}
	}

	if _, present := properties["propertyNames"]; present {
		pattern := properties["propertyNames"].(string)
		return map[string]interface{}{
			GenerateRegex(pattern): "value",
		}
	}

	generatedObject := make(map[string]interface{}, 0)

	if _, present := properties["patternProperties"]; present {

		patternProperties := properties["patternProperties"].(map[string]interface{})

		for field, fieldProperties := range patternProperties {
			fieldType := fieldProperties.(map[string]interface{})["type"].(string)
			patternFieldName := GenerateRegex(field)
			generatedObject[patternFieldName] = FieldToGenerator[fieldType]
		}

	}

	if _, present := properties["properties"]; present {
		objectProperties := properties["properties"].(map[string]interface{})

		for field, fieldProperties := range objectProperties {
			fieldType := fieldProperties.(map[string]interface{})["type"].(string)
			generatedObject[field] = FieldToGenerator[fieldType]

		}

	}
	return generatedObject

}

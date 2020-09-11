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

	if _, present := properties["enum"]; present {
		if len(properties["enum"].([]int)) != 0 {
			return properties["enum"].([]int)[0]
		}
	}

	if _, present := properties["constant"]; present {
		return int(properties["constant"].(float64))

	}

	rand.Seed(time.Now().UnixNano())

	var min, max int

	if _, present := properties["exclusiveMaximum"]; present {
		max = int(properties["exclusiveMaximum"].(float64)) - 1
	}

	if _, present := properties["exclusiveMinimum"]; present {
		min = int(properties["exclusiveMinimum"].(float64)) + 1
	}

	if _, present := properties["maximum"]; !present {
		max = def.DummyDataRange["integer"].(map[string]int)["maximum"]
	} else {
		max = int(properties["maximum"].(float64))
	}

	if _, present := properties["minimum"]; !present {
		min = def.DummyDataRange["integer"].(map[string]int)["minimum"]
	} else {
		min = int(properties["minimum"].(float64))
	}

	if _, present := properties["enum"]; present {
		enum := properties["enum"].([]int)
		if len(enum) > 0 {
			return enum[0]
		}
		return 0
	}

	if _, present := properties["multipleOf"]; present {
		multipleOf := int(properties["multipleOf"].(float64))
		i := 0
		for {
			eachMultiple := multipleOf * i
			if eachMultiple >= min && eachMultiple <= max {
				return multipleOf * i
			} else if eachMultiple > max {
				return 0
			}
			i++

		}
	}

	return rand.Intn(max-min+1) + min
}

func GenerateFloat(properties map[string]interface{}) float64 {

	if _, present := properties["enum"]; present {
		if len(properties["enum"].([]float64)) != 0 {
			return properties["enum"].([]float64)[0]
		}

	}

	if _, present := properties["constant"]; present {
		return properties["constant"].(float64)

	}

	rand.Seed(time.Now().UnixNano())

	var min, max float64

	if _, present := properties["exclusiveMaximum"]; present {
		max = properties["exclusiveMaximum"].(float64) - 0.5
	}

	if _, present := properties["exclusiveMinimum"]; present {
		min = properties["exclusiveMinimum"].(float64) + 0.5
	}

	if _, present := properties["maximum"]; !present {
		max = def.DummyDataRange["number"].(map[string]float64)["maximum"]
	} else {
		max = properties["maximum"].(float64)
	}

	if _, present := properties["minimum"]; !present {
		min = def.DummyDataRange["number"].(map[string]float64)["minimum"]
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

	if _, present := properties["enum"]; present {
		if len(properties["enum"].([]string)) != 0 {
			return properties["enum"].([]string)[0]
		}

	}

	if _, present := properties["constant"]; present {
		return properties["constant"].(string)

	}

	rand.Seed(time.Now().UnixNano())

	var minLength, maxLength int

	if _, present := properties["maxLength"]; !present {
		maxLength = def.DummyDataRange["string"].(map[string]int)["maxLength"]
	} else {
		maxLength = int(properties["maxLength"].(float64))
	}

	if _, present := properties["minLength"]; !present {
		minLength = def.DummyDataRange["string"].(map[string]int)["minLength"]
	} else {
		minLength = int(properties["minLength"].(float64))
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
			maxItems = def.DummyDataRange["array"].(map[string]int)["maximum"]
		} else {
			maxItems = int(properties["maxItems"].(float64))
		}

		// if _, present := properties["minItems"]; !present {
		// 	minItems = 0
		// } else {
		// 	minItems = properties["minItems"].(int)
		// }

		numberArray := make([]float64, 0)

		var arrayProperties = map[string]interface{}{
			"maximum": def.DummyDataRange["number"].(map[string]int)["maximum"],
			"minimum": def.DummyDataRange["number"].(map[string]int)["minimum"],
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

			if eachArrayType == "array" {
				resultingArray = append(resultingArray, GenerateArray(eachTypeProperties))
			} else if eachArrayType == "object" {
				resultingArray = append(resultingArray, GenerateObject(eachTypeProperties))
			} else {
				resultingArray = append(resultingArray, FieldToGenerator[eachArrayType](eachTypeProperties))
			}

		}
		return resultingArray

		//Mix Type Array

	} else if arrayItemType == "map[string]interface {}" {

		//Single Type Arrays

		itemType := properties["items"].(map[string]interface{})["type"].(string)

		if itemType == "array" {
			return GenerateArray(properties["items"].(map[string]interface{}))
		} else if itemType == "object" {
			return []interface{}{GenerateObject(properties["items"].(map[string]interface{}))}
		}

	}

	return []string{}

}

func GenerateObject(properties map[string]interface{}) map[string]interface{} {

	if _, present := properties["propertyNames"]; present {
		pattern := properties["propertyNames"].(map[string]interface{})["pattern"].(string)
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

			if fieldType == "array" {
				generatedObject[patternFieldName] = GenerateArray(fieldProperties.(map[string]interface{}))
			} else if fieldType == "object" {
				generatedObject[patternFieldName] = GenerateObject(fieldProperties.(map[string]interface{}))
			} else {

				generatedObject[patternFieldName] = FieldToGenerator[fieldType](fieldProperties.(map[string]interface{}))
			}

		}

	}

	if _, present := properties["properties"]; present {
		objectProperties := properties["properties"].(map[string]interface{})

		for field, fieldProperties := range objectProperties {
			fieldType := fieldProperties.(map[string]interface{})["type"].(string)

			if fieldType == "array" {
				generatedObject[field] = GenerateArray(fieldProperties.(map[string]interface{}))

			} else if fieldType == "object" {
				generatedObject[field] = GenerateObject(fieldProperties.(map[string]interface{}))
			} else {
				generatedObject[field] = FieldToGenerator[fieldType](fieldProperties.(map[string]interface{}))

			}

		}
	}

	if len(generatedObject) == 0 {
		if _, minPropertiesPresent := properties["minProperties"]; minPropertiesPresent {
			minPropertiesPresent := int(properties["minProperties"].(float64))

			for i := 0; i < minPropertiesPresent; i++ {
				generatedObject["test"+strconv.Itoa(i)] = i
			}

			return generatedObject
		}

		generatedObject["k1"] = "v1"
		generatedObject["k2"] = "v2"
	}

	return generatedObject

}

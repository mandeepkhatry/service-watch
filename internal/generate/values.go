package generate

import (
	"fmt"
	"math/rand"
	"service-watch/internal/def"
	"strconv"
	"time"

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

func GenerateStringFormat(stringType string) string {
	return def.StringFormat[stringType]

}

func GenerateRegex(regex string) string {
	str, _ := reggen.Generate(regex, 1)
	return str
}

func GenerateArray(properties map[string]interface{}) interface{} {

	minItems := def.DummyDataRange["array"].(map[string]int)["minItems"]

	if min, present := properties["minItems"]; present {
		minItems = int(min.(float64))
	}

	arrayItemType := fmt.Sprintf("%T", properties["items"])

	if arrayItemType == "[]interface {}" {

		resultingArray := make([]interface{}, 0)

		for _, eachTypeProperties := range properties["items"].([]interface{}) {

			eachProperties := eachTypeProperties

			eachArrayType := eachProperties.(map[string]interface{})["type"].(string)

			if eachArrayType == "array" {
				resultingArray = append(resultingArray, GenerateArray(eachProperties.(map[string]interface{})))
			} else if eachArrayType == "object" {
				resultingArray = append(resultingArray, GenerateObject(eachProperties.(map[string]interface{})))
			} else {
				resultingArray = append(resultingArray, FieldToGenerator[eachArrayType](eachProperties.(map[string]interface{})))
			}

		}
		return resultingArray

		//Mix Type Array

	} else if arrayItemType == "map[string]interface {}" {
		resultingArray := make([]interface{}, 0)

		prop := properties["items"]

		//Single Type Arrays

		itemType := prop.(map[string]interface{})["type"].(string)

		if itemType == "array" {
			for i := 0; i < minItems; i++ {
				resultingArray = append(resultingArray, GenerateArray(prop.(map[string]interface{})))
			}

		} else if itemType == "object" {
			for i := 0; i < minItems; i++ {
				resultingArray = append(resultingArray, GenerateObject(prop.(map[string]interface{})))
			}
		} else {
			for i := 0; i < minItems; i++ {
				resultingArray = append(resultingArray, FieldToGenerator[itemType](prop.(map[string]interface{})))
			}
		}

		return resultingArray
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

	if _, present := properties["allOf"]; present {
		for _, eachAllOfProp := range properties["allOf"].([]interface{}) {
			for k, v := range GenerateObject(eachAllOfProp.(map[string]interface{})) {
				generatedObject[k] = v
			}
		}
	}

	if _, present := properties["oneOf"]; present {
		for _, eachAllOfProp := range properties["allOf"].([]interface{}) {
			generatedObject = GenerateObject(eachAllOfProp.(map[string]interface{}))
			break

		}
	}

	if _, present := properties["anyOf"]; present {
		for _, eachAllOfProp := range properties["anyOf"].([]interface{}) {
			generatedObject = GenerateObject(eachAllOfProp.(map[string]interface{}))
			break
		}
	}

	//Generate object for then properties
	if _, present := properties["then"]; present {
		generatedObject = GenerateObject(properties["then"].(map[string]interface{})["properties"].(map[string]interface{}))
	}

	// if _, present := properties["required"]; present {
	// 	object := make(map[string]interface{})
	// 	for _, eachRequiredField := range properties["required"].([]interface{}) {
	// 		object[eachRequiredField.(string)] = generatedObject[eachRequiredField.(string)]
	// 	}
	// 	return object
	// }

	return generatedObject

}

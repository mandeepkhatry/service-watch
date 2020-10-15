package def

//DummyDataRange represents default range for fields
var DummyDataRange = map[string]interface{}{
	"integer": map[string]int{
		"minimum": 5,
		"maximum": 10,
	},
	"number": map[string]float64{
		"minimum": 5.0,
		"maximum": 10.0,
	},
	"string": map[string]int{
		"minLength": 5,
		"maxLength": 10,
	},
	"array": map[string]int{
		"minItems": 2,
		"maxItems": 5,
	},
}

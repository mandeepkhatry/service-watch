package def

//DummyDataRange represents default range for fields
var DummyDataRange = map[string]interface{}{
	"integer": map[string]int{
		"minimum": 0,
		"maximum": 10,
	},
	"number": map[string]float64{
		"minimum": 0.0,
		"maximum": 10.0,
	},
	"string": map[string]int{
		"minLength": 0,
		"maxLength": 10,
	},
	"array": map[string]int{
		"minimum": 0,
		"maximum": 10,
	},
}

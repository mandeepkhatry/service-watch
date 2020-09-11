package main

import (
	"fmt"
	"service-watch/watch"
)

var configPath = "config/watch.json"

func main() {
	fmt.Println("Testing")

	serviceWatch, err := watch.NewServiceWatcher(configPath)

	if err != nil {
		panic(err)
	}

	err = serviceWatch.Watch()

	if err != nil {
		panic(err)
	}

	// fmt.Println(generate.GenerateFloat(map[string]interface{}{"type": "number"}))

	// var prop = map[string]interface{}{
	// 	"type": "array",
	// 	"items": []map[string]interface{}{
	// 		map[string]interface{}{
	// 			"type": "number",
	// 		},
	// 		map[string]interface{}{
	// 			"type": "number",
	// 		},
	// 		map[string]interface{}{
	// 			"type": "string",
	// 		},
	// 		map[string]interface{}{
	// 			"type": "array",
	// 			"items": []map[string]interface{}{
	// 				map[string]interface{}{
	// 					"type": "number",
	// 				},
	// 				map[string]interface{}{
	// 					"type": "number",
	// 				},
	// 				map[string]interface{}{
	// 					"type": "string",
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// var prop = map[string]interface{}{
	// 	"type": "object",
	// 	// "patternProperties": map[string]interface{}{
	// 	// 	"^S_": map[string]interface{}{"type": "string"},
	// 	// 	"^I_": map[string]interface{}{"type": "integer", "multipleOf": 10, "minimum": 500},
	// 	// 	"^K_": map[string]interface{}{
	// 	// 		"type": "string",
	// 	// 		"enum": []string{"Street", "Avenue", "Boulevard"},
	// 	// 	},
	// 	// },
	// 	// "propertyNames": map[string]interface{}{
	// 	// 	"pattern": "^[A-Za-z_][A-Za-z0-9_]*$",
	// 	// },
	// 	"properties": map[string]interface{}{
	// 		"num": map[string]interface{}{
	// 			"type": "number",
	// 		},
	// 		"street_num": map[string]interface{}{
	// 			"type": "number",
	// 		},
	// 		"name": map[string]interface{}{
	// 			"type": "string",
	// 		},
	// 		"array": map[string]interface{}{
	// 			"type": "array",
	// 			"items": []map[string]interface{}{
	// 				map[string]interface{}{
	// 					"type": "number",
	// 				},
	// 				map[string]interface{}{
	// 					"type": "number",
	// 				},
	// 				map[string]interface{}{
	// 					"type": "string",
	// 				},
	// 			},
	// 		},
	// 		"inside-object": map[string]interface{}{
	// 			"type": "object",
	// 			"properties": map[string]interface{}{
	// 				"mynumber": map[string]interface{}{
	// 					"type": "number",
	// 				},
	// 				"name": map[string]interface{}{
	// 					"type": "string",
	// 				},
	// 				"inside-inside-object": map[string]interface{}{
	// 					"type": "object",
	// 					"properties": map[string]interface{}{
	// 						"insidemynumber": map[string]interface{}{
	// 							"type": "number",
	// 						},
	// 						"insidename": map[string]interface{}{
	// 							"type": "string",
	// 						},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// x := generate.GenerateDummyData(prop)

	// fmt.Println(x)

}

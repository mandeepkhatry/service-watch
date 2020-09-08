package main

import (
	"fmt"
	"log"
	"service-watch/internal/generate"
	"service-watch/watch"
)

func main() {
	fmt.Println("Testing")

	watcher := watch.ServiceWatcher{}
	watcher.Init("static/gateway-config/test.toml")
	err := watcher.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(watcher.AppConfig.Api["mock"])

	fmt.Println(generate.GenerateFloat(map[string]interface{}{"type": "number"}))

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

	var prop = map[string]interface{}{
		"type": "object",
		// "patternProperties": map[string]interface{}{
		// 	"^S_": map[string]interface{}{"type": "string"},
		// 	"^I_": map[string]interface{}{"type": "integer", "multipleOf": 10, "minimum": 500},
		// 	"^K_": map[string]interface{}{
		// 		"type": "string",
		// 		"enum": []string{"Street", "Avenue", "Boulevard"},
		// 	},
		// },
		// "propertyNames": map[string]interface{}{
		// 	"pattern": "^[A-Za-z_][A-Za-z0-9_]*$",
		// },
		"properties": map[string]interface{}{
			"num": map[string]interface{}{
				"type": "number",
			},
			"street_num": map[string]interface{}{
				"type": "number",
			},
			"name": map[string]interface{}{
				"type": "string",
			},
			"array": map[string]interface{}{
				"type": "array",
				"items": []map[string]interface{}{
					map[string]interface{}{
						"type": "number",
					},
					map[string]interface{}{
						"type": "number",
					},
					map[string]interface{}{
						"type": "string",
					},
				},
			},
		},
	}

	x := generate.GenerateObject(prop)

	fmt.Println(x)

}

package def

var SchemaBasedMethods = map[string]bool{
	"POST":  true,
	"PUT":   true,
	"PATCH": true,
}

var RearMostMethod = map[string]bool{
	"DELETE": true,
}

var AccpetedStatus = map[string]map[int]bool{
	"POST": map[int]bool{
		201: true,
	},
	"GET": map[int]bool{
		200: true,
	},
	"PUT": map[int]bool{
		200: true,
		204: true,
	},
	"PATCH": map[int]bool{
		200: true,
		204: true,
	},
	"DELETE": map[int]bool{
		200: true,
		204: true,
	},
}

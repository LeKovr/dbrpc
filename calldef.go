package main

// -----------------------------------------------------------------------------

// CallDef holds function call definition attributes
type CallDef struct {
	Cache uint32                 `json:"cache"` // uint32 : 0 to 4294967295
	Name  *string                `json:"nsp"`
	Proc  *string                `json:"proc"`
	Args  map[string]interface{} `json:"arg"`
	Age   int                    `json:"age"`
	Lang  *string                `json:"lang"`
	TZ    *string                `json:"tz"`
}

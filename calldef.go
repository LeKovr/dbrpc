package main

// -----------------------------------------------------------------------------

// CallDef holds function call definition attributes
type CallDef struct {
	Name *string                `json:"nsp"`
	Proc *string                `json:"proc"`
	Args map[string]interface{} `json:"arg"`
	Age  int                    `json:"age"`
	Lang *string                `json:"lang"`
	TZ   *string                `json:"tz"`
}

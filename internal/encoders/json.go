package encoders

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
)

// GoQRCodeLogJSONSerializer is our custom version of echo's JSON serializing struct
type GoQRCodeLogJSONSerializer struct {
	defaultSerializer echo.DefaultJSONSerializer
}

func NewGoQRCodeLogJSONSerializer() *GoQRCodeLogJSONSerializer {
	s := GoQRCodeLogJSONSerializer{
		defaultSerializer: echo.DefaultJSONSerializer{},
	}
	return &s
}

// Serialize converts an interface into a json and writes it to the response.
// You can optionally use the indent parameter to produce pretty JSONs.
// Our version sets HTML chars escaping to false, so passwords are returned correctly when creating users
func (d GoQRCodeLogJSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	enc.SetEscapeHTML(false)
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

// Deserialize reads a JSON from a request body and converts it into an interface.
func (d GoQRCodeLogJSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	return d.defaultSerializer.Deserialize(c, i)
}

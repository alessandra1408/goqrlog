package echoutil

import (
	"reflect"
	"strings"

	"github.com/alessandra1408/goqrlog/internal/encoders"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.JSONSerializer = encoders.NewGoQRCodeLogJSONSerializer()

	e.Use(middleware.Recover())

	e.IPExtractor = echo.ExtractIPFromRealIPHeader()

	e.HideBanner = true

	return e
}

func SetupValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tag, ok := fld.Tag.Lookup("form")
		if !ok {
			tag = fld.Tag.Get("json")
		}

		tagSplit := strings.SplitN(tag, ",", 2)

		if len(tagSplit) > 0 {
			name := tagSplit[0]

			if name == "-" {
				return ""
			}

			return name
		}

		return ""
	})

	return v
}

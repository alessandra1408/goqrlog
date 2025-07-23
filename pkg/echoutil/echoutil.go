package echoutil

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/alessandra1408/goqrlog/internal/encoders"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.JSONSerializer = encoders.NewGoQRLogJSONSerializer()

	e.Use(middleware.Recover())

	e.Use(CORSMiddleware())

	e.IPExtractor = echo.ExtractIPFromRealIPHeader()

	e.HideBanner = true

	return e
}

// Adicione esta função no mesmo arquivo (ou em um arquivo middleware.go)
func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "https://qr-scanner-dom-jaime.surge.sh")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Response().Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

			if c.Request().Method == "OPTIONS" {
				return c.NoContent(http.StatusOK)
			}

			return next(c)
		}
	}
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

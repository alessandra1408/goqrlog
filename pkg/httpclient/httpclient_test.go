package httpclient

import (
	"testing"
	"time"

	"github.com/alessandra1408/goqrlog/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClient(t *testing.T) {
	cases := []struct {
		name         string
		expctTimeout time.Duration
		oTelEnabled  bool
	}{
		{
			name:         "should return a http client with the correct timeout",
			expctTimeout: time.Second * 15,
			oTelEnabled:  false,
		},
		{
			name:         "should return a http client with OpenTelemetry enabled",
			expctTimeout: time.Second,
			oTelEnabled:  true,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			mockCfg := config.Config{
				App: &config.App{HTTPTimeout: cs.expctTimeout},
			}
			client := NewHTTPClient(&mockCfg)

			if assert.NotNil(t, client) {
				assert.Equal(t, cs.expctTimeout, client.Timeout)
			}
		})
	}
}

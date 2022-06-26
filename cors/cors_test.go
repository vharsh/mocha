package cors

import (
	"github.com/stretchr/testify/assert"
	"github.com/vitorsalgado/mocha/internal/header"
	"github.com/vitorsalgado/mocha/internal/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCORS(t *testing.T) {
	msg := "hello world"
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	}

	ts := httptest.NewServer(
		middleware.Compose(CORS(*Configure().
			AllowMethods("GET", "POST").
			AllowedHeaders("x-allow-this", "x-allow-that").
			ExposeHeaders("x-expose-this").
			AllowOrigin("*").
			AllowCredentials(true).
			Build())).Root(http.HandlerFunc(handler)))
	defer ts.Close()

	t.Run("should allow request", func(t *testing.T) {
		// check preflight request
		req, _ := http.NewRequest(http.MethodOptions, ts.URL, nil)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Equal(t, "*", res.Header.Get(header.AccessControlAllowOrigin))
		assert.Equal(t, "x-expose-this", res.Header.Get(header.AccessControlExposeHeaders))
		assert.Equal(t, "true", res.Header.Get(header.AccessControlAllowCredentials))
		assert.Equal(t, "GET,POST", res.Header.Get(header.AccessControlAllowMethods))
		assert.Equal(t, "x-allow-this,x-allow-that", res.Header.Get(header.AccessControlAllowHeaders))

		// check the actual request
		res, err = http.Get(ts.URL)
		if err != nil {
			log.Fatal(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.True(t, strings.Contains(string(body), msg))
		assert.Equal(t, "*", res.Header.Get(header.AccessControlAllowOrigin))
		assert.Equal(t, "x-expose-this", res.Header.Get(header.AccessControlExposeHeaders))
		assert.Equal(t, "true", res.Header.Get(header.AccessControlAllowCredentials))
		assert.Equal(t, "text/plain", res.Header.Get("content-type"))
	})
}
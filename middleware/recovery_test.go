package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	g := gin.Default()
	g.Use(Recovery())
	g.GET("/ping", func(c *gin.Context) {
		panic("error")
	})

	ts := httptest.NewServer(g)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/ping")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal("Did not expect http.Get to fail")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("%s\n", `{"message":"Internal Server Error"}`), string(body))

}

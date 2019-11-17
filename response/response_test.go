package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rianekacahya/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestError(t *testing.T) {

	g := gin.Default()
	g.GET("/ping", func(c *gin.Context) {
		err := errors.New(errors.BADREQUEST, errors.Message("Error"))
		Error(c, err)
	})

	ts := httptest.NewServer(g)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/ping")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal("Did not expect http.Get to fail")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("%s\n", `{"message":"Error"}`), string(body))
}

func TestRender(t *testing.T) {

	g := gin.Default()
	g.GET("/pong", func(c *gin.Context) {
		Render(c, http.StatusOK, nil)
	})

	ts := httptest.NewServer(g)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/pong")
	defer resp.Body.Close()
	if err != nil {
		t.Fatal("Did not expect http.Get to fail")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("%s\n", `{"message":"success"}`), string(body))
}

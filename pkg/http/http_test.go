package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/spezifisch/silphtelescope/pkg/pogo"
	"github.com/stretchr/testify/assert"
)

var (
	responseOK   = "OK\n"
	responseFAIL = "FAIL\n"
)

func readTestFile(f string) string {
	content, err := ioutil.ReadFile("../../test/data/" + f)
	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}

func testMadWebhookRequest(data string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/webhook/mad", strings.NewReader(data))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	return c, rec
}

func TestSetup(t *testing.T) {
	Init()

	Bind = "invalid stuff"
	Run()
	Stop()
}

func TestRoot(t *testing.T) {
	data := `[]`
	c, rec := testMadWebhookRequest(data)

	if assert.NoError(t, hello(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestMadWebhookEmpty(t *testing.T) {
	data := `[]`
	c, rec := testMadWebhookRequest(data)

	if assert.NoError(t, madWebhook(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, responseOK, rec.Body.String())
	}
}

func TestMadWebhookInvalid(t *testing.T) {
	data := `foo`
	c, rec := testMadWebhookRequest(data)

	if assert.NoError(t, madWebhook(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, responseFAIL, rec.Body.String())
	}
}

func TestMadWebhookAllTypes(t *testing.T) {
	data := readTestFile("mad-webhook-all-types.json")
	c, rec := testMadWebhookRequest(data)

	GymUpdates = make(chan pogo.Gym, 50)
	RaidUpdates = make(chan pogo.Raid, 50)
	SpawnUpdates = make(chan pogo.Spawn, 200)

	if assert.NoError(t, madWebhook(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, responseOK, rec.Body.String())
	}
}

func TestMadWebhookUnhandledType(t *testing.T) {
	data := readTestFile("mad-webhook-unhandled-type.json")
	c, rec := testMadWebhookRequest(data)

	GymUpdates = make(chan pogo.Gym, 50)
	RaidUpdates = make(chan pogo.Raid, 50)
	SpawnUpdates = make(chan pogo.Spawn, 200)

	if assert.NoError(t, madWebhook(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, responseOK, rec.Body.String())
	}
}

func TestMadWebhookBrokenMessages(t *testing.T) {
	data := readTestFile("mad-webhook-broken.json")
	c, rec := testMadWebhookRequest(data)

	GymUpdates = make(chan pogo.Gym, 50)
	RaidUpdates = make(chan pogo.Raid, 50)
	SpawnUpdates = make(chan pogo.Spawn, 200)

	if assert.NoError(t, madWebhook(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, responseOK, rec.Body.String())
	}
}

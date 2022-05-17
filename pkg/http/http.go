package http

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

var (
	// Bind specifies "<host>:<port>" to listen on
	Bind string
	// GymUpdates receiver
	GymUpdates chan pogo.Gym
	// SpawnUpdates receiver
	SpawnUpdates chan pogo.Spawn
	// RaidUpdates receiver
	RaidUpdates chan pogo.Raid

	e *echo.Echo
)

// Init sets up the routes
func Init() {
	// Echo instance
	e = echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.POST("/webhook/mad", madWebhook)
}

// Run starts the httpd
func Run() {
	// Start server
	err := e.Start(Bind)
	if err != nil {
		log.WithError(err).Error("http server quit:")
	}
}

// Stop stops the http server immediately
func Stop() {
	log.Info("stopping http server")
	e.Close()
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "This is SilphTelescope.")
}

func madWebhook(c echo.Context) error {
	var envelopes []Envelope
	err := json.NewDecoder(c.Request().Body).Decode(&envelopes)
	if err != nil {
		log.Warnln("can't decode request to:", c.Request().URL, "error:", err)
		return c.String(http.StatusBadRequest, "FAIL\n")
	}

	for _, msg := range envelopes {
		var dst interface{}
		switch msg.Type {
		case "gym":
			dst = new(GymMessage)
		case "pokemon":
			dst = new(PokemonMessage)
		case "raid":
			dst = new(RaidMessage)
		case "pokestop":
			continue
		case "weather":
			continue
		default:
			log.Debugln("unhandled type", msg.Type)
			continue
		}

		err := json.Unmarshal(msg.Message, dst)
		if err != nil {
			log.Warn("can't decode message as type ", msg.Type, ": ", err)
			continue
		}

		switch msg.Type {
		case "gym":
			sendGymUpdate(dst.(*GymMessage))
		case "pokemon":
			sendSpawnUpdate(dst.(*PokemonMessage))
		case "raid":
			sendRaidUpdate(dst.(*RaidMessage))
		}
	}

	return c.String(http.StatusOK, "OK\n")
}

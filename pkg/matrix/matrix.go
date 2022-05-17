package matrix

import (
	"time"

	"github.com/matrix-org/gomatrix"
	log "github.com/sirupsen/logrus"
	"github.com/spezifisch/silphtelescope/pkg/roomservice"
)

// Matrix contains authentication info and our homeserver
type Matrix struct {
	// Homeserver URL in format https://example.com
	Homeserver string
	// UserID for homeserver in format @user:home.server
	UserID string
	// AccessToken for homeserver retrieved from login
	AccessToken string

	cli    *gomatrix.Client
	poster *roomservice.Poster

	// set to true when we're should shut down
	stopping bool
}

// New Matrix API
func New(homeserver, userID, accessToken string) (m *Matrix) {
	m = &Matrix{
		Homeserver:  homeserver,
		UserID:      userID,
		AccessToken: accessToken,
		stopping:    false,
	}
	m.Init()
	return
}

// Init needs to be called before using other methods. Doesn't need to be called when using New()
func (m *Matrix) Init() {
	// Custom interfaces must be set prior to calling functions on the client.
	var err error
	if m.cli, err = gomatrix.NewClient(m.Homeserver, m.UserID, m.AccessToken); err != nil {
		log.Panic("can't connect to homeserver")
		return
	}

	// we apparently need an own syncer to add the callbacks
	customSyncer := gomatrix.NewDefaultSyncer(m.UserID, m.cli.Store)
	m.cli.Syncer = customSyncer

	// add message callbacks
	customSyncer.OnEventType("m.room.message", m.onRoomMessage)
	customSyncer.OnEventType("m.room.member", m.onRoomMember)

	log.Print("Using Homeserver: ", m.Homeserver, ", UserID: ", m.UserID)
}

// SetPoster connects the matrix callbacks to a poster that gets configured by room commands
func (m *Matrix) SetPoster(p *roomservice.Poster) {
	m.poster = p
}

// Tick runs a single sync for use in a main loop
func (m *Matrix) Tick() {
	if err := m.cli.Sync(); err != nil {
		log.Errorln("Sync() returned", err)
	}
}

// Run matrix sync main loop
func (m *Matrix) Run() {
	for {
		m.Tick()
		if m.stopping {
			break
		}

		// Wait a period of time before trying to sync again.
		time.Sleep(2 * time.Second)
	}
}

// Stop stops the matrix client
func (m *Matrix) Stop() {
	log.Info("stopping matrix client")
	m.stopping = true
	m.cli.StopSync()
}

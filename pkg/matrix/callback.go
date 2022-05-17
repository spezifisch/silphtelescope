/*
 * this file is gplv3
 *
 * using Apache 2.0 licensed code:
 * - using parts of onRoomMemberEvent https://git.logansnow.com/logan/jeeves-riot/raw/commit/804def8087e6a8838464327d7f631b2c7783cbb7/src/github.com/matrix-org/go-neb/clients/clients.go
 *   in onRoomMember
 */

package matrix

import (
	"github.com/matrix-org/gomatrix"
	log "github.com/sirupsen/logrus"

	"github.com/spezifisch/silphtelescope/pkg/roomservice"
)

func (m *Matrix) contextFromEvent(e *gomatrix.Event) (r roomservice.Context) {
	r.Cli = m.cli
	r.Chatter = &Chatter{
		cli: m.cli,
	}

	r.Sender = e.Sender
	r.Type = e.Type
	r.Timestamp = e.Timestamp
	r.ID = e.ID
	r.RoomID = e.RoomID
	return
}

// onRoomMessage is called for m.room.message
func (m *Matrix) onRoomMessage(e *gomatrix.Event) {
	if m.isFromMe(e) {
		return
	}

	log.Debugln("onRoomMessage", e.Timestamp, e.RoomID, e.Sender)

	if body, ok := e.Body(); ok {
		if mtype, ok := e.MessageType(); ok {
			log.Debugln("- type:", mtype)
		}
		log.Debugln("- content:", body)

		ctx := m.contextFromEvent(e)
		m.poster.ParseMessage(body, ctx)
	}
}

// onRoomMessage is called for m.room.member
func (m *Matrix) onRoomMember(e *gomatrix.Event) {
	if m.isFromMe(e) {
		return
	}

	log.Debugln("onRoomMember", e.Timestamp, e.RoomID, e.Sender)

	if m.isInvite(e) {
		// TODO check if it's my admin user
		log.Infoln("got invite to", e.RoomID, "from", e.Sender)

		if _, err := m.cli.JoinRoom(e.RoomID, "", nil); err != nil {
			log.Errorln("Failed to join room", err)
		} else {
			log.Debugln("Joined room")
		}
	}
}

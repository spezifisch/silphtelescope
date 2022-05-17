/*
 * this file is gplv3
 *
 * using Apache 2.0 licensed code:
 * - using parts of onRoomMemberEvent https://git.logansnow.com/logan/jeeves-riot/raw/commit/804def8087e6a8838464327d7f631b2c7783cbb7/src/github.com/matrix-org/go-neb/clients/clients.go
 *   in isInvite
 */

package matrix

import "github.com/matrix-org/gomatrix"

func (m *Matrix) isForMe(e *gomatrix.Event) bool {
	return e.StateKey != nil && *e.StateKey == m.cli.UserID
}

func (m *Matrix) isFromMe(e *gomatrix.Event) bool {
	return e.Sender == m.cli.UserID
}

func (m *Matrix) isInvite(e *gomatrix.Event) bool {
	if !m.isForMe(e) {
		return false // not our member event
	}

	mem := e.Content["membership"]
	membership, ok := mem.(string)
	if !ok {
		return false
	}
	return membership == "invite"
}

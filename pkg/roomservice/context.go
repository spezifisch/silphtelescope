/*
 * this file is gplv3
 *
 * using Apache 2.0 licensed code:
 * - using parts of Event https://github.com/matrix-org/gomatrix/blob/7dd5e2a05bcda194c84dbe6a38c024ae787a568e/events.go#L9
 */

package roomservice

import "github.com/matrix-org/gomatrix"

// Context is the room state description
type Context struct {
	Cli     *gomatrix.Client // Client pointer
	Chatter Chatter
	Poster  *Poster

	Sender    string `json:"sender"`           // The user ID of the sender of the event
	Type      string `json:"type"`             // The event type
	Timestamp int64  `json:"origin_server_ts"` // The unix timestamp when this message was sent by the origin server
	ID        string `json:"event_id"`         // The unique ID of this event
	RoomID    string `json:"room_id"`          // The room the event was sent to. May be nil (e.g. for presence)
}

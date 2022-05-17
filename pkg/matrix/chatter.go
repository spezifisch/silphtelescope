package matrix

import (
	"github.com/matrix-org/gomatrix"
)

// Chatter decouples gomatrix.Client from our code
type Chatter struct {
	cli *gomatrix.Client
}

// SendText wraps gomatrix.SendText
func (c *Chatter) SendText(roomID, text string) {
	c.cli.SendText(roomID, text)
}

// SendFormattedText wraps gomatrix.SendFormattedText
func (c *Chatter) SendFormattedText(roomID, text, formattedText string) {
	c.cli.SendFormattedText(roomID, text, formattedText)
}

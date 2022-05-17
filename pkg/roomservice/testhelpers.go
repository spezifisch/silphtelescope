package roomservice

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// testChatter mocks up the matrix output for roomservice
type testChatter struct {
	LastRoomID, LastText, LastFormattedText string
	MessageReceived                         chan bool
	Counter                                 int
	LastExpectedMessage                     int
}

func (c *testChatter) SendText(roomID, text string) {
	c.LastRoomID = roomID
	c.LastText = text
	c.LastFormattedText = ""
	c.Counter++
	log.Debugln(roomID, text)
	c.MessageReceived <- true
}

func (c *testChatter) SendFormattedText(roomID, text, formattedText string) {
	c.LastRoomID = roomID
	c.LastText = text
	c.LastFormattedText = formattedText
	c.Counter++
	log.Debugln(roomID, "text:", text)
	c.MessageReceived <- true
}

func (c *testChatter) ExpectMessage(t *testing.T) {
	<-c.MessageReceived

	assert.Equal(t, c.LastExpectedMessage+1, c.Counter, "message counter didn't increase by one between messages")

	c.LastExpectedMessage = c.Counter
}

func (c *testChatter) ExpectNoMessage(t *testing.T) {
	assert.Equal(t, c.LastExpectedMessage, c.Counter, "a message was sent that shouldn't have been sent")
}

func (c *testChatter) PrintLastMessage() {
	log.Println("<Bot>", c.LastText)
}

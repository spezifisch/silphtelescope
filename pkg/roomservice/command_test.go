package roomservice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	c := &testChatter{
		// we need to buffer one message because we're running
		// the sender in the same thread as the receiver
		MessageReceived: make(chan bool, 1),
	}
	roomID := "!bar@example.com"
	ctx := Context{
		Chatter: c,
		RoomID:  roomID,
	}
	p := ctx.Poster

	var handled bool
	handled, _ = p.ParseMessage("help", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "commands:")
	assert.Equal(t, roomID, c.LastRoomID)

	handled, _ = p.ParseMessage("sdf sdf gdfg dsfg", ctx)
	c.ExpectNoMessage(t)
	assert.Equal(t, false, handled)

	handled, _ = p.ParseMessage("admin", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: admin")
	assert.Equal(t, roomID, c.LastRoomID)

	handled, _ = p.ParseMessage("status", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, c.LastText, "not ready")
	assert.Equal(t, roomID, c.LastRoomID)

	handled, _ = p.ParseMessage("fort", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "geodex deactivated")
	assert.Equal(t, roomID, c.LastRoomID)

	handled, _ = p.ParseMessage("filter", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: filter")
	assert.Equal(t, roomID, c.LastRoomID)

	handled, _ = p.ParseMessage("spawn", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: spawn")
	assert.Equal(t, roomID, c.LastRoomID)

	handled, _ = p.ParseMessage("raid", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: raid")
	assert.Equal(t, roomID, c.LastRoomID)

	// test things that need the poster object
	ctx.Poster = NewPoster(c, nil)
	p = ctx.Poster

	handled, _ = p.ParseMessage("status", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Bot uptime")
	assert.Contains(t, c.LastText, "Last MAD data: never")
	assert.Equal(t, roomID, c.LastRoomID)

	p.startTime = time.Now()
	p.lastDataTime = time.Now()
	handled, _ = p.ParseMessage("status", ctx)
	c.ExpectMessage(t)
	c.PrintLastMessage()
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Bot uptime")
	assert.Contains(t, c.LastText, "Last MAD data:")
	assert.NotContains(t, c.LastText, "Last MAD data: never")
	assert.Equal(t, roomID, c.LastRoomID)
}

func TestParseSpawn(t *testing.T) {
	c := &testChatter{
		// we need to buffer one message because we're running
		// the sender in the same thread as the receiver
		MessageReceived: make(chan bool, 1),
	}
	p := NewPoster(c, nil)
	roomID := "!bar@example.com"
	ctx := Context{
		Chatter: c,
		RoomID:  roomID,
		Poster:  p,
	}

	var handled bool
	handled, _ = p.ParseMessage("spawn foo", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: spawn")

	handled, _ = p.ParseMessage("spawn help", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: spawn")

	handled, _ = p.ParseMessage("spawn rm", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)

	handled, _ = p.ParseMessage("spawn add", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: spawn")

	handled, _ = p.ParseMessage("spawn rm", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: spawn")

	handled, _ = p.ParseMessage("spawn add 0 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "failed: roomconfig doesn't exist", c.LastText)

	handled, _ = p.ParseMessage("spawn add -1 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "failed: roomconfig doesn't exist", c.LastText)

	handled, _ = p.ParseMessage("spawn add a 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "invalid parameter", c.LastText)

	handled, _ = p.ParseMessage("spawn add 0 a", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "invalid parameter", c.LastText)

	handled, _ = p.ParseMessage("spawn rm 0 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "failed: roomconfig doesn't exist", c.LastText)

	handled, _ = p.ParseMessage("spawn rm -1 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "failed: roomconfig doesn't exist", c.LastText)

	handled, _ = p.ParseMessage("spawn rm y 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "invalid parameter", c.LastText)

	handled, _ = p.ParseMessage("spawn rm 23.23 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "invalid parameter", c.LastText)

	// and now valid stuff
	handled, _ = p.ParseMessage("filter add spawn 0 0 0", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)

	handled, _ = p.ParseMessage("spawn add 0 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "added to filter", c.LastText)
	assert.Equal(t, 1, len(p.roomConfigs[roomID].Filter))
	assert.Equal(t, true, p.roomConfigs[roomID].Filter[0].ListWanted)
	assert.Equal(t, false, p.roomConfigs[roomID].Filter[0].ListRaids)
	assert.Equal(t, 1, len(p.roomConfigs[roomID].Filter[0].PokemonIDs))
	assert.Equal(t, 1, p.roomConfigs[roomID].Filter[0].PokemonIDs[0])

	handled, _ = p.ParseMessage("spawn add 0 23", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "added to filter", c.LastText)
	assert.Equal(t, 1, len(p.roomConfigs[roomID].Filter))
	assert.Equal(t, 2, len(p.roomConfigs[roomID].Filter[0].PokemonIDs))
	assert.Equal(t, 1, p.roomConfigs[roomID].Filter[0].PokemonIDs[0])
	assert.Equal(t, 23, p.roomConfigs[roomID].Filter[0].PokemonIDs[1])

	handled, _ = p.ParseMessage("spawn rm 0 1", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "removed from filter", c.LastText)
	assert.Equal(t, 1, len(p.roomConfigs[roomID].Filter))
	assert.Equal(t, 1, len(p.roomConfigs[roomID].Filter[0].PokemonIDs))
	assert.Equal(t, 23, p.roomConfigs[roomID].Filter[0].PokemonIDs[0])
}

func TestParseRaid(t *testing.T) {
	c := &testChatter{
		// we need to buffer one message because we're running
		// the sender in the same thread as the receiver
		MessageReceived: make(chan bool, 1),
	}
	p := NewPoster(c, nil)
	roomID := "!bar@example.com"
	ctx := Context{
		Chatter: c,
		RoomID:  roomID,
		Poster:  p,
	}

	var handled bool
	handled, _ = p.ParseMessage("raid foo", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: raid")
}

func TestAdmin(t *testing.T) {
	c := &testChatter{
		// we need to buffer one message because we're running
		// the sender in the same thread as the receiver
		MessageReceived: make(chan bool, 1),
	}
	p := NewPoster(c, nil)
	p.saveStateAndQuit = make(chan bool, 1)
	roomID := "!bar@example.com"
	ctx := Context{
		Chatter: c,
		RoomID:  roomID,
		Poster:  p,
	}

	handled, _ := p.ParseMessage("admin help", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: admin")

	handled, _ = p.ParseMessage("admin foo", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "Usage: admin")

	// roomstate doesn't exist
	handled, _ = p.ParseMessage("admin roomstate", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "room state doesn't exist", c.LastText)

	// roomconfig doesn't exist
	handled, _ = p.ParseMessage("admin roomconfig", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Equal(t, "no roomconfig found", c.LastText)

	// add roomstate
	p.getOrCreateRoomState(roomID)

	// roomstate exists
	handled, _ = p.ParseMessage("admin roomstate", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "RoomState")
	assert.Contains(t, c.LastFormattedText, "<code>")

	// add roomconfig
	rc := getTestRoomConfig(roomID)
	p.UpdateRoomConfig(rc)

	// roomconfig exists
	handled, _ = p.ParseMessage("admin roomconfig", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	assert.Contains(t, c.LastText, "RoomConfig")
	assert.Contains(t, c.LastFormattedText, "<code>")

	handled, _ = p.ParseMessage("admin shutdown", ctx)
	c.ExpectNoMessage(t)
	assert.Equal(t, true, handled)
	<-p.saveStateAndQuit

	// break the room
	borkedRoom := "!bork@woof"
	ctx.RoomID = borkedRoom

	// roomstate is nil
	handled, _ = p.ParseMessage("admin roomstate", ctx)
	c.ExpectMessage(t)
	assert.Equal(t, true, handled)
	c.PrintLastMessage()
}

func TestCommandList(t *testing.T) {
	generateCommandList()
	assert.Contains(t, commandList, "commands:")

	commands = []Command{}
	generateCommandList()
	assert.Equal(t, "there are no commands", commandList)
}

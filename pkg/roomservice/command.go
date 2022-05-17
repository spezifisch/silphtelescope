package roomservice

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spezifisch/silphtelescope/pkg/geodex"
	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

const (
	roomConfigFilterLimit = 100 // maximum filters a room can have
)

var (
	commands = []Command{
		{"help", helpCallback, true},
		{"admin", adminCallback, false},
		{"status", statusCallback, true},
		{"mon", monCallback, true},
		{"fort", fortCallback, true},
		{"filter", filterCallback, true},
		{"spawn", spawnCallback, true},
		{"raid", raidCallback, true},
	}
	commandList string
)

// CommandCallback contains the line triggering it and room context
type CommandCallback func(args []string, context Context) (handled bool, err error)

// Command is triggered by a line starting with the keyword
type Command struct {
	// Command without / or ! prefix
	Command  string
	Callback CommandCallback
	// accessible by everyone? false=admin only
	Public bool
}

func generateCommandList() {
	if len(commands) == 0 {
		commandList = "there are no commands"
		return
	}

	commandList = "commands:"
	for _, cmd := range commands {
		commandList = fmt.Sprintf("%s %s", commandList, cmd.Command)
	}
}

// ParseMessage processes a single line and calls the callback if it contains a command
func (p *Poster) ParseMessage(msg string, context Context) (handled bool, err error) {
	if commandList == "" {
		generateCommandList()
	}

	msg = strings.TrimSpace(msg)
	msgParts := strings.Split(msg, " ")

	for _, cmd := range commands {
		if msgParts[0] == cmd.Command {
			log.Println("Command called:", msg)
			if context.Poster == nil {
				context.Poster = p
			}
			handled, err = cmd.Callback(msgParts, context)
			return
		}
	}

	log.Println("Command not handled:", msgParts[0])
	handled = false
	return
}

func simpleResponse(context Context, text string) {
	context.Chatter.SendText(context.RoomID, text)
}

func jsonResponse(context Context, title, roomID string, jsonData []byte) {
	text := fmt.Sprintf("%s for `%s`:\n```json\n%s\n```", title, roomID, jsonData)
	formattedText := fmt.Sprintf("<p>%s for <code>%s</code>:</p>\n<pre><code class=\"language-json\">%s\n</code></pre>\n", title, roomID, jsonData)

	context.Chatter.SendFormattedText(context.RoomID, text, formattedText)
}

func adminPostRoomState(context Context) (err error) {
	rs, ok := context.Poster.roomStates[context.RoomID]
	if !ok {
		simpleResponse(context, "room state doesn't exist")
		return
	}

	// Marshal can't fail with our RoomState type,
	// see: https://stackoverflow.com/questions/33903552/what-input-will-cause-golangs-json-marshal-to-return-an-error
	rsText, _ := json.MarshalIndent(rs, "", "    ")
	jsonResponse(context, "RoomState", context.RoomID, rsText)
	return
}

func adminClearRoomState(context Context) (err error) {
	rs, ok := context.Poster.roomStates[context.RoomID]
	if !ok {
		simpleResponse(context, "room state doesn't exist")
		return
	}

	rs.clear(true, true)
	return
}

func adminPostRoomConfig(context Context) (err error) {
	if rc, ok := context.Poster.GetRoomConfig(context.RoomID); ok {
		jsonText, _ := json.MarshalIndent(rc, "", "  ")
		jsonResponse(context, "RoomConfig", context.RoomID, jsonText)
	} else {
		simpleResponse(context, "no roomconfig found")
	}
	return
}

func adminCallback(args []string, context Context) (handled bool, err error) {
	handled = true

	subCmd := "help"
	if len(args) >= 2 {
		subCmd = args[1]
	}

	switch subCmd {
	case "roomconfig":
		err = adminPostRoomConfig(context)
	case "roomstate":
		err = adminPostRoomState(context)
	case "roomstate_clear":
		err = adminClearRoomState(context)
	case "shutdown":
		context.Poster.saveStateAndQuit <- true
	case "help":
		fallthrough
	default:
		simpleResponse(context, "Usage: admin [roomconfig|roomstate[_clear]|shutdown]")
	}

	return
}

func statusCallback(args []string, context Context) (handled bool, err error) {
	handled = true
	if context.Poster == nil {
		simpleResponse(context, "not ready")
		return
	}

	emptyTime := time.Time{}
	now := time.Now()
	uptime := now.Sub(context.Poster.startTime)

	lastData := now.Sub(context.Poster.lastDataTime)
	lastDataStr := "never"
	if context.Poster.lastDataTime != emptyTime {
		lastDataStr = fmt.Sprintf("%s ago", lastData)
	}

	text := fmt.Sprintf("Bot uptime: %s\nLast MAD data: %s", uptime, lastDataStr)

	simpleResponse(context, text)
	return
}

func helpCallback(args []string, context Context) (handled bool, err error) {
	simpleResponse(context, commandList)
	handled = true
	return
}

func monCallback(args []string, context Context) (handled bool, err error) {
	handled = true

	if context.Poster == nil || context.Poster.Pokedex == nil {
		simpleResponse(context, "pokedex deactivated by admin")
		return
	}

	arg := NewArgParser(args)
	if arg.Count() != 2 {
		simpleResponse(context, "Usage: mon <id[,id2[,id3...]]|name>\nLookup name or id in Pokedex")
		return
	}

	// try to parse as single id or as comma-separated id array
	var idArray []int
	gotIDArray := false
	if id, err := arg.AsInt(1); err == nil {
		idArray = []int{id}
		gotIDArray = true
	} else if ids, err := arg.AsIntArray(1); err == nil {
		idArray = ids
		gotIDArray = true
	}

	if gotIDArray {
		text := ""
		for _, id := range idArray {
			nameEN, nameDE, err := context.Poster.Pokedex.GetNamesByID(id)
			if err == nil {
				text = fmt.Sprintf("%s\n#%d English: %s, German: %s", text, id, nameEN, nameDE)
			} else {
				text = fmt.Sprintf("%s\n#%d not found", text, id)
			}
		}
		text = strings.TrimSpace(text)
		simpleResponse(context, text)
		return
	}

	// try to parse as single name
	if name, err := arg.AsString(1); err == nil {
		var text string
		id, nameEN, nameDE, err := context.Poster.Pokedex.GetIDByName(name)
		if err == nil {
			text = fmt.Sprintf("#%d English: %s German: %s", id, nameEN, nameDE)
		} else {
			text = "Pokemon not found"
		}
		simpleResponse(context, text)
	}

	return
}

func filterCallback(args []string, context Context) (handled bool, err error) {
	handled = true
	arg := NewArgParser(args)

	subCmd := "help"
	if arg.Count() >= 2 {
		subCmd, _ = arg.AsString(1)
	}

	switch subCmd {
	case "add":
		if arg.Count() != 6 {
			simpleResponse(context, "Usage: filter add <raid|spawn> <lat> <lon> <radius_m>\nAdd a new filter that matches raids or spawns around the given location.")
			return
		}

		// parse args
		typ, err2 := arg.AsString(2)
		listRaids := typ == "raid"
		listWanted := typ == "spawn"
		validTyp := listRaids || listWanted

		area, err3 := arg.AsLocationRadius(3, 4, 5)
		if err2 != nil || err3 != nil || !validTyp {
			simpleResponse(context, "invalid parameter")
			return
		}

		// check limits
		if rc, ok := context.Poster.GetRoomConfig(context.RoomID); ok {
			// RoomConfig exists
			if len(rc.Filter) >= roomConfigFilterLimit {
				simpleResponse(context, "you've reached the allowed limit of filters a room can have")
				return
			}
		}

		// add roomconfig
		rc := &RoomConfig{
			RoomID:     context.RoomID,
			FormatText: true,
			Filter: []PokemonFilter{
				{
					ListRaids:  listRaids,
					ListWanted: listWanted,
					PokemonIDs: []int{},
					Area:       area,
				},
			},
		}
		context.Poster.UpdateRoomConfig(rc)
		simpleResponse(context, "added filter to roomconfig")
	case "rm":
		if arg.Count() != 3 {
			simpleResponse(context, "Usage: spawn rm <filter_id>\nRemove filter from RoomConfig.")
			return
		}

		filterID, err2 := arg.AsInt(2)
		if err2 != nil {
			simpleResponse(context, "invalid parameter")
			return
		}

		change := &RoomConfigChange{
			Operation:   RoomConfigOperationRemoveFilter,
			FilterIndex: filterID,
		}
		newValues := &RoomConfig{}
		err := context.Poster.ChangeRoomConfig(context.RoomID, change, newValues)
		if err == nil {
			simpleResponse(context, "removed from filter")
		} else {
			text := fmt.Sprintf("failed: %s", err.Error())
			simpleResponse(context, text)
		}
	case "area":
		if arg.Count() != 6 {
			simpleResponse(context, "Usage: filter area <filter_id> <lat> <lon> <radius_m>\nChange area for given filter id.")
			return
		}

		filterID, err2 := arg.AsInt(2)
		area, err3 := arg.AsLocationRadius(3, 4, 5)
		if err2 != nil || err3 != nil {
			simpleResponse(context, "invalid parameters")
			return
		}

		change := &RoomConfigChange{
			Operation:    RoomConfigOperationUpdateFilter,
			FilterIndex:  filterID,
			FilterChange: FilterChangeArea,
		}
		newValues := &RoomConfig{
			Filter: []PokemonFilter{
				{
					Area: area,
				},
			},
		}
		err := context.Poster.ChangeRoomConfig(context.RoomID, change, newValues)
		if err == nil {
			simpleResponse(context, "filter area updated")
		} else {
			text := fmt.Sprintf("failed: %s", err.Error())
			simpleResponse(context, text)
		}
	case "drop":
		simpleResponse(context, "this removes ALL filters from this room IRREVOCABLY! type \"filter dropreally\" if you really intend to do this.")
	case "dropreally":
		deleted := context.Poster.DeleteFilters(context.RoomID)
		if deleted {
			simpleResponse(context, "removed all filters in roomconfig")
		} else {
			simpleResponse(context, "no roomconfig found")
		}
	case "help":
		fallthrough
	default:
		simpleResponse(context, "Usage: filter [add|rm|area|drop]")
	}
	return
}

func changeFilterMon(verb changeFilterMonType, args []string, context Context) (handled bool, err error) {
	handled = true
	arg := NewArgParser(args)

	subCmd := "help"
	if arg.Count() >= 2 {
		subCmd, _ = arg.AsString(1)
	}

	switch subCmd {
	case "add":
		if arg.Count() != 4 {
			text := fmt.Sprintf("Usage: %s add <filter_id> <pkmn_id[,id2[,id3...]]>\nAppend Pokemon ID(s) to filter.", verb)
			simpleResponse(context, text)
			return
		}

		filterID, err2 := arg.AsInt(2)
		monIDs, err3 := arg.AsIntArray(3)
		if err2 != nil || err3 != nil {
			simpleResponse(context, "invalid parameter")
			return
		}

		change := &RoomConfigChange{
			Operation:    RoomConfigOperationUpdateFilter,
			FilterIndex:  filterID,
			FilterChange: FilterChangeAddPokemon,
		}
		newValues := &RoomConfig{
			Filter: []PokemonFilter{
				{
					PokemonIDs: monIDs,
				},
			},
		}
		err := context.Poster.ChangeRoomConfig(context.RoomID, change, newValues)
		if err == nil {
			simpleResponse(context, "added to filter")
		} else {
			text := fmt.Sprintf("failed: %s", err.Error())
			simpleResponse(context, text)
		}
	case "rm":
		if arg.Count() != 4 {
			text := fmt.Sprintf("Usage: %s rm <filter_id> <pkmn_id[,id2[,id3...]]>\nRemove Pokemon ID(s) from filter.", verb)
			simpleResponse(context, text)
			return
		}

		filterID, err2 := arg.AsInt(2)
		monIDs, err3 := arg.AsIntArray(3)
		if err2 != nil || err3 != nil {
			simpleResponse(context, "invalid parameter")
			return
		}

		change := &RoomConfigChange{
			Operation:    RoomConfigOperationUpdateFilter,
			FilterIndex:  filterID,
			FilterChange: FilterChangeRemovePokemon,
		}
		newValues := &RoomConfig{
			Filter: []PokemonFilter{
				{
					PokemonIDs: monIDs,
				},
			},
		}
		err := context.Poster.ChangeRoomConfig(context.RoomID, change, newValues)
		if err == nil {
			simpleResponse(context, "removed from filter")
		} else {
			text := fmt.Sprintf("failed: %s", err.Error())
			simpleResponse(context, text)
		}
	case "help":
		fallthrough
	default:
		text := fmt.Sprintf("Usage: %s [add|rm]", verb)
		simpleResponse(context, text)
	}

	return
}

type changeFilterMonType string

const (
	changeFilterMonRaid  changeFilterMonType = "raid"
	changeFilterMonSpawn                     = "spawn"
)

func raidCallback(args []string, context Context) (handled bool, err error) {
	return changeFilterMon(changeFilterMonRaid, args, context)
}

func spawnCallback(args []string, context Context) (handled bool, err error) {
	return changeFilterMon(changeFilterMonSpawn, args, context)
}

func fortCallback(args []string, context Context) (handled bool, err error) {
	handled = true

	if context.Poster == nil || context.Poster.GeoDex == nil {
		simpleResponse(context, "geodex deactivated by admin")
		return
	}

	subCmd := "help"
	if len(args) >= 2 {
		subCmd = args[1]
	}

	switch subCmd {
	case "near":
		if len(args[2:]) != 2 && len(args[2:]) != 3 {
			simpleResponse(context, "Usage: fort near <lat> <lon> [radius_m=500]")
			return
		}

		lat, err1 := strconv.ParseFloat(args[2], 64)
		lon, err2 := strconv.ParseFloat(args[3], 64)
		var err3 error = nil
		radiusM := 500.0
		if len(args[2:]) == 3 {
			radiusM, err3 = strconv.ParseFloat(args[4], 64)
		}
		if err1 != nil || err2 != nil || err3 != nil {
			simpleResponse(context, "invalid float")
			return
		}

		center := pogo.Location{
			Latitude:  lat,
			Longitude: lon,
		}
		fort, err := context.Poster.GeoDex.Tile.GetNearestFort(center, radiusM)
		if err != nil {
			text := fmt.Sprintf("no fort found near (%f,%f) in %f m radius",
				center.Latitude, center.Longitude, radiusM)
			simpleResponse(context, text)
		} else {
			text := fmt.Sprintf("t38: %s", fort.ToString())

			// lookup name from geodex
			nFort, err := context.Poster.GeoDex.Disk.GetFort(*fort.GUID)
			if err == nil {
				text = fmt.Sprintf("%s\ndisk: %s", text, nFort.ToString())
			} else {
				text = fmt.Sprintf("%s\ndisk: not found", text)
				log.WithError(err).Warnf("guid %s from t38 is not on disk", *fort.GUID)
			}

			text = fmt.Sprintf("%s\ndistance from fort to (%f,%f): %fm, bearing %f°",
				text, center.Latitude, center.Longitude,
				fort.Location().DistanceTo(&center),
				fort.Location().BearingTo(&center))

			simpleResponse(context, text)
		}
	case "info":
		if len(args[2:]) != 1 {
			simpleResponse(context, "Usage: fort info <GUID>")
			return
		}

		guid := args[2]
		if !geodex.IsValidGUID(guid) {
			simpleResponse(context, "GUID contains invalid characters")
			return
		}

		fort, err := context.Poster.GeoDex.Disk.GetFort(guid)
		if err != nil {
			simpleResponse(context, "disk: no fort found with that GUID")
		} else {
			text := fmt.Sprintf("disk: %s", fort.ToString())
			simpleResponse(context, text)
		}
	case "help":
		fallthrough
	default:
		simpleResponse(context, "Usage: fort [near|info]")
	}

	return
}

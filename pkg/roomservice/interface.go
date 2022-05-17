package roomservice

// Chatter posts messages into a chatroom
type Chatter interface {
	SendText(roomID, text string)
	SendFormattedText(roomID, text, formattedText string)
}

// Persister persists states and configs between runs
type Persister interface {
	SaveRoomConfig(roomID string, roomConfig *RoomConfig)
	SaveRoomConfigs(roomConfigs map[string]*RoomConfig)
	ReadRoomConfigs(roomConfigs map[string]*RoomConfig)

	SaveRoomState(roomID string, roomState *RoomState)
	SaveRoomStates(roomStates map[string]*RoomState)
	ReadRoomStates(roomStates map[string]*RoomState)
}

// MainController stops main when Stop is called
type MainController interface {
	Stop()
}

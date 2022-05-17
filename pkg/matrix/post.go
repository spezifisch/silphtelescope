package matrix

// SendText sends a message to a room
func (m *Matrix) SendText(roomID, text string) {
	m.cli.SendText(roomID, text)
}

// SendFormattedText sends a message to a room
func (m *Matrix) SendFormattedText(roomID, text, formattedText string) {
	m.cli.SendFormattedText(roomID, text, formattedText)
}

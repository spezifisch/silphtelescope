package geodex

// BOQCell are the values from the outer dict that BookOfQuests returns
type BOQCell struct {
	Stops []*BOQStop `json:"stops"`
}

// BOQStop has name and location for POIs
type BOQStop struct {
	Name      string      `json:"name"`
	IsPortal  bool        `json:"portal"`
	IsGym     bool        `json:"gym"`
	IsStop    bool        `json:"stop"`
	Timestamp int64       `json:"ts"`
	S2Level20 string      `json:"s2l20"`
	Location  BOQGeometry `json:"loc"`
}

// BOQGeometry always is a Point with Lat/Lon here
type BOQGeometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

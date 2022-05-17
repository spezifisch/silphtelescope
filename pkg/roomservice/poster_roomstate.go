package roomservice

import log "github.com/sirupsen/logrus"

func (p *Poster) saveRoomStates() {
	if p.db == nil {
		return
	}

	log.Infof("saving RoomState for %d rooms", len(p.roomStates))
	p.db.SaveRoomStates(p.roomStates)
}

func (p *Poster) readRoomStates() {
	if p.db == nil {
		return
	}

	p.db.ReadRoomStates(p.roomStates)
	log.Infof("read RoomState for %d rooms", len(p.roomStates))
}

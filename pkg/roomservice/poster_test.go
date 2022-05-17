package roomservice

import (
	"testing"
	"time"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

func TestRoomConfigChanges(t *testing.T) {
	c := &testChatter{}
	p := NewPoster(c, nil)

	testRoom := "!foo@example.com"

	// add one config
	rcA := getTestRoomConfig(testRoom)
	rcA.AcceptCommands = false
	rcA.FormatText = true
	created := p.UpdateRoomConfig(rcA)

	assert.Equal(t, true, created)
	assert.Equal(t, rcA, p.roomConfigs[testRoom])
	assert.Equal(t, 1, len(p.roomConfigs[testRoom].Filter))

	// rc should now exist
	gotRC, ok := p.GetRoomConfig(testRoom)
	assert.Equal(t, true, ok)
	assert.Contains(t, gotRC.ToString(), testRoom)

	// add another filter
	rcB := getTestRoomConfig(testRoom)
	rcB.AcceptCommands = true
	rcB.FormatText = false
	rcB.Filter[0].ListRaids = true
	rcB.Filter[0].ListWanted = false
	rcB.Filter[0].PokemonIDs = []int{150, 151}
	rcB.Filter[0].Area.RadiusM = 5000
	created = p.UpdateRoomConfig(rcB)

	assert.Equal(t, false, created)
	assert.Equal(t, true, p.roomConfigs[testRoom].AcceptCommands)
	assert.Equal(t, false, p.roomConfigs[testRoom].FormatText)
	assert.Equal(t, 2, len(p.roomConfigs[testRoom].Filter))
}

func TestRoomStateChanges(t *testing.T) {
	rs := NewRoomState()

	encA := "one"
	assert.Equal(t, false, rs.spawnIsPosted(encA))

	// mark spawn as posted
	spawnA := getTestSpawn()
	spawnA.EncounterID = encA
	rs.postedSpawn(&spawnA, true)

	assert.Equal(t, true, rs.spawnIsPosted(encA))
	assert.Equal(t, false, rs.spawnIsPosted("blubb"))

	// not yet expired
	now := spawnA.EndTime - 600
	rs.removeExpired(now)
	assert.Equal(t, true, rs.spawnIsPosted(encA))

	now = spawnA.EndTime - 1
	rs.removeExpired(now)
	assert.Equal(t, true, rs.spawnIsPosted(encA))

	// we see the EndTime as "valid until including t"
	now = spawnA.EndTime
	rs.removeExpired(now)
	assert.Equal(t, true, rs.spawnIsPosted(encA))

	// expire
	now = spawnA.EndTime + 1
	rs.removeExpired(now)
	assert.Equal(t, false, rs.spawnIsPosted(encA))

	// a later expiration time
	rs.postedSpawn(&spawnA, true) // add again
	assert.Equal(t, true, rs.spawnIsPosted(encA))

	now = spawnA.EndTime + 600
	rs.removeExpired(now)
	assert.Equal(t, false, rs.spawnIsPosted(encA))
}

func TestPosterTicker(t *testing.T) {
	c := &testChatter{
		MessageReceived: make(chan bool, 1),
	}
	p := NewPoster(c, nil)

	// make the check period so short that it triggers instantly
	p.ExpiryCheckPeriod = 60 * time.Millisecond
	p.GymUpdates = make(chan pogo.Gym)
	p.RaidUpdates = make(chan pogo.Raid)
	p.SpawnUpdates = make(chan pogo.Spawn, 1)
	p.Quit = make(chan bool)

	// give the expiry check something to do
	testRoom := "!foo@example.com"
	rcA := getTestRoomConfig(testRoom)
	p.UpdateRoomConfig(rcA)
	s := getTestSpawn()
	p.SpawnUpdates <- s

	go func() {
		p.Run()
	}()

	time.Sleep(100 * time.Millisecond)
	p.Quit <- true
}

// startPoster created
func startPoster() (*Poster, chan bool, *testChatter) {
	done := make(chan bool)
	c := &testChatter{
		MessageReceived: make(chan bool, 10),
	}

	p := NewPoster(c, nil)
	// make the check period so long that it doesn't trigger
	p.ExpiryCheckPeriod = 24 * time.Hour
	p.GymUpdates = make(chan pogo.Gym)
	p.RaidUpdates = make(chan pogo.Raid)
	p.SpawnUpdates = make(chan pogo.Spawn)
	p.Quit = make(chan bool)

	// p.Run blocks, so wrap it in a goroutine
	go func() {
		p.Run()
		done <- true
	}()

	return p, done, c
}

// startPoster but with resumable RoomState/Config
func startPosterResumable(doResume bool, mockDB Persister, mockMainControl MainController) (*Poster, chan bool, *testChatter) {
	done := make(chan bool)
	c := &testChatter{
		MessageReceived: make(chan bool, 10),
	}

	p := NewPoster(c, mockDB)
	p.ResumeStateOnStartup = doResume
	p.MainControl = mockMainControl
	// make the check period so long that it doesn't trigger
	p.ExpiryCheckPeriod = 24 * time.Hour
	p.GymUpdates = make(chan pogo.Gym)
	p.RaidUpdates = make(chan pogo.Raid)
	p.SpawnUpdates = make(chan pogo.Spawn)
	p.Quit = make(chan bool)

	// p.Run blocks, so wrap it in a goroutine
	go func() {
		p.Run()
		done <- true
	}()

	return p, done, c
}

func startPosterPokedex(t *testing.T, pokedexFile string) (*Poster, chan bool, *testChatter) {
	done := make(chan bool)
	c := &testChatter{
		MessageReceived: make(chan bool, 10),
	}

	p := NewPoster(c, nil)
	p.ExpiryCheckPeriod = 24 * time.Hour
	p.GymUpdates = make(chan pogo.Gym)
	p.RaidUpdates = make(chan pogo.Raid)
	p.SpawnUpdates = make(chan pogo.Spawn)
	p.Quit = make(chan bool)

	var err error
	p.Pokedex, err = pogo.NewPokedex(pokedexFile)
	assert.NoError(t, err, "failed reading pokedex")

	// p.Run blocks, so wrap it in a goroutine
	go func() {
		p.Run()
		done <- true
	}()

	return p, done, c
}

func TestPosterAllTypesEmpty(t *testing.T) {
	p, done, _ := startPoster()

	// each line blocks until p.Run processes it
	p.GymUpdates <- pogo.Gym{}
	p.RaidUpdates <- pogo.Raid{}
	p.SpawnUpdates <- pogo.Spawn{}
	// no need to read from MessageReceived as it has space for 10 msgs
	p.Quit <- true

	// wait until it quits
	<-done
}

func getTestSpawn() pogo.Spawn {
	var endTime int64 = 1613800682 // 6:58:02
	//now := endTime - 25*60         // 25m before

	return pogo.Spawn{
		EncounterID:        "abc",
		VerifiedSpawnpoint: true,
		Pokemon: pogo.Pokemon{
			ID:   16, // Pidgey
			Name: "",
		},
		TimestampRange: pogo.TimestampRange{
			StartTime: 0,
			EndTime:   endTime,
		},
		Location: pogo.Location{
			// some horse stable in Cairo
			Latitude:  30.05113,
			Longitude: 31.21918,
		},
	}
}

func getTestEgg() pogo.Raid {
	var endTime int64 = 1613800682 // 6:58:02
	startTime := endTime - 45*60

	return pogo.Raid{
		Hash:    "abc",
		Level:   3,
		Pokemon: nil,
		TimestampRange: pogo.TimestampRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		GymID: "bepis",
		Location: pogo.Location{
			Latitude:  30.05113,
			Longitude: 31.21918,
		},
	}
}

func getTestRaid() pogo.Raid {
	var endTime int64 = 1613800682 // 6:58:02
	startTime := endTime - 45*60

	return pogo.Raid{
		Hash:  "abcfefsdf",
		Level: 5,
		Pokemon: &pogo.Pokemon{
			ID: 150,
		},
		TimestampRange: pogo.TimestampRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		GymID: "conke",
		Location: pogo.Location{
			Latitude:  30.05113,
			Longitude: 31.21918,
		},
	}
}

func getTestPoint2KMAway() pogo.Location {
	return pogo.Location{
		Latitude:  30.0495,
		Longitude: 31.2592,
	}
}

func getTestRoomConfig(testRoom string) *RoomConfig {
	return &RoomConfig{
		RoomID:         testRoom,
		AcceptCommands: false,
		FormatText:     false,
		Filter: []PokemonFilter{
			{
				ListRaids:  false,
				ListWanted: true,
				PokemonIDs: []int{2, 4, 8, 16, 25, 32},
				Area: pogo.LocationRadius{
					Location: pogo.Location{
						Latitude:  30.04896,
						Longitude: 31.22366,
					},
					RadiusM: 1000, // the spawn is at ~300m distance
				},
			},
		},
	}
}

func TestPosterSpawns(t *testing.T) {
	p, done, c := startPoster()

	// set roomservice filter
	testRoom := "!foo@example.com"

	rc := getTestRoomConfig(testRoom)
	p.UpdateRoomConfig(rc)

	// send spawn
	s := getTestSpawn()
	p.SpawnUpdates <- s
	// wait until spawn message is posted
	c.ExpectMessage(t)
	c.PrintLastMessage()

	assert.Equal(t, testRoom, c.LastRoomID)
	//assert.Equal(t, "a spawn", c.LastText)
	//assert.Equal(t, "", c.LastFormattedText)

	// send same spawn again, should be ignored
	p.SpawnUpdates <- s
	c.ExpectNoMessage(t)

	// send non-interesting spawn
	s.Pokemon.ID = 5
	s.EncounterID = "x"
	p.SpawnUpdates <- s
	c.ExpectNoMessage(t)
	s.EncounterID = "y"
	p.SpawnUpdates <- s
	c.ExpectNoMessage(t)
	s.EncounterID = "z"
	p.SpawnUpdates <- s
	c.ExpectNoMessage(t)

	// with formatted text
	rc.FormatText = true
	p.UpdateRoomConfig(rc)

	// send spawn
	s = getTestSpawn()
	s.EncounterID = "sdhfgksjdhf"
	p.SpawnUpdates <- s
	// wait until spawn message is posted
	c.ExpectMessage(t)

	assert.Equal(t, testRoom, c.LastRoomID)
	//assert.Equal(t, "a spawn", c.LastText)
	//assert.Equal(t, "", c.LastFormattedText)

	// spawn out of range
	s = getTestSpawn()
	s.EncounterID = "spawnOORtest"
	s.Location = getTestPoint2KMAway()
	p.SpawnUpdates <- s
	c.ExpectNoMessage(t)

	// wait
	p.Quit <- true
	<-done
}

func TestPosterResumeWithoutDB(t *testing.T) {
	mockMainControl := getMockMainControl()

	doResume := true
	p, done, _ := startPosterResumable(doResume, nil, mockMainControl)

	p.saveStateAndQuit <- true
	<-done

	assert.Equal(t, true, mockMainControl.stopCalled)
}

func TestPosterQuitAndResume(t *testing.T) {
	mockDB := getMockDB()
	mockMainControl := getMockMainControl()

	doResume := false
	p, done, c := startPosterResumable(doResume, mockDB, mockMainControl)

	// set roomservice filter
	testRoom := "!asdfadf@example.com"

	rc := getTestRoomConfig(testRoom)
	p.UpdateRoomConfig(rc)

	assert.Equal(t, 1, len(p.roomConfigs))
	assert.Equal(t, 0, len(p.roomStates))

	// send spawn
	s := getTestSpawn()
	p.SpawnUpdates <- s
	// wait until spawn message is posted
	c.ExpectMessage(t)

	assert.Equal(t, 1, len(p.roomConfigs))
	assert.Equal(t, 1, len(p.roomStates))

	// quit gracefully, saving states
	p.saveStateAndQuit <- true
	<-done

	assert.Equal(t, true, mockMainControl.stopCalled)

	// make sure it's saved
	assert.Equal(t, 1, len(mockDB.savedRoomConfigs))
	assert.Equal(t, 1, len(mockDB.savedRoomStates))

	// resume
	mockMainControl.stopCalled = false
	doResume = true
	p, done, c = startPosterResumable(doResume, mockDB, mockMainControl)

	// Now we need to wait until Poster.Run() initialized the config.
	// We can use the spawn channel without cache as it block until Run's main loop reads it.
	p.SpawnUpdates <- pogo.Spawn{}

	assert.Equal(t, mockDB, p.db)
	assert.Equal(t, true, p.ResumeStateOnStartup)

	assert.Equal(t, 1, len(p.roomConfigs))
	assert.Equal(t, 1, len(p.roomStates))

	// quit it harshly
	p.Quit <- true
	<-done

	assert.Equal(t, false, mockMainControl.stopCalled)
}

type mockMainController struct {
	stopCalled bool
}

func (m *mockMainController) Stop() {
	m.stopCalled = true
}

func getMockMainControl() *mockMainController {
	return &mockMainController{
		stopCalled: false,
	}
}

type mockPersister struct {
	savedRoomConfigs map[string]*RoomConfig
	savedRoomStates  map[string]*RoomState
}

func (p *mockPersister) SaveRoomConfig(roomID string, roomConfig *RoomConfig) {
	p.savedRoomConfigs[roomID] = roomConfig
}
func (p *mockPersister) SaveRoomConfigs(roomConfigs map[string]*RoomConfig) {
	copier.Copy(&p.savedRoomConfigs, roomConfigs)
}
func (p *mockPersister) ReadRoomConfigs(roomConfigs map[string]*RoomConfig) {
	copier.Copy(&roomConfigs, p.savedRoomConfigs)
}

func (p *mockPersister) SaveRoomState(roomID string, roomState *RoomState) {
	p.savedRoomStates[roomID] = roomState
}
func (p *mockPersister) SaveRoomStates(roomStates map[string]*RoomState) {
	copier.Copy(&p.savedRoomStates, roomStates)
}
func (p *mockPersister) ReadRoomStates(roomStates map[string]*RoomState) {
	copier.Copy(&roomStates, p.savedRoomStates)
}

func getMockDB() *mockPersister {
	return &mockPersister{
		savedRoomConfigs: map[string]*RoomConfig{},
		savedRoomStates:  map[string]*RoomState{},
	}
}

var testPokedexFile = "../../data/pokedex.json"

func TestPosterSpawnsWithPokedex(t *testing.T) {
	p, done, c := startPosterPokedex(t, testPokedexFile)

	// set roomservice filter
	testRoom := "!foo@example.com"

	rc := getTestRoomConfig(testRoom)
	p.UpdateRoomConfig(rc)

	// send spawn
	s := getTestSpawn()
	p.SpawnUpdates <- s
	// wait until spawn message is posted
	c.ExpectMessage(t)
	assert.Equal(t, testRoom, c.LastRoomID)
	assert.Contains(t, c.LastText, "Pidgey")
	assert.Contains(t, c.LastText, "Taubsi")

	// coverage: send a Pikachu whose name is the same for DE/EN
	s.Pokemon.ID = 25
	s.EncounterID = "zwei"
	p.SpawnUpdates <- s
	c.ExpectMessage(t)
	assert.Equal(t, testRoom, c.LastRoomID)
	assert.Contains(t, c.LastText, "Pikachu")

	// ...again with formatted text
	rc.FormatText = true
	p.UpdateRoomConfig(rc)

	s = getTestSpawn()
	s.EncounterID = "drei"
	p.SpawnUpdates <- s
	c.ExpectMessage(t)
	assert.Equal(t, testRoom, c.LastRoomID)
	assert.Contains(t, c.LastFormattedText, "Pidgey")
	assert.Contains(t, c.LastFormattedText, "Taubsi")

	s.Pokemon.ID = 25
	s.EncounterID = "vier"
	p.SpawnUpdates <- s
	c.ExpectMessage(t)
	assert.Equal(t, testRoom, c.LastRoomID)
	assert.Contains(t, c.LastFormattedText, "Pikachu")

	// wait
	p.Quit <- true
	<-done
}

func TestPosterRaids(t *testing.T) {
	p, done, c := startPoster()

	// set roomservice filter
	testRoom := "!foo@example.com"
	formattedRoom := "!formatted@example.com"

	// add raid config
	rc := getTestRoomConfig(testRoom)
	rc.Filter[0].ListWanted = false
	rc.Filter[0].ListRaids = true
	rc.Filter[0].PokemonIDs = []int{150, 151}
	p.UpdateRoomConfig(rc)

	// add raid config for room with formatted text
	rc = getTestRoomConfig(formattedRoom)
	rc.FormatText = true
	rc.Filter[0].ListWanted = false
	rc.Filter[0].ListRaids = true
	rc.Filter[0].PokemonIDs = []int{1}
	p.UpdateRoomConfig(rc)

	// add spawn config for coverage
	rc = getTestRoomConfig(testRoom)
	rc.Filter[0].PokemonIDs = []int{152}
	p.UpdateRoomConfig(rc)

	// send spawn
	s := getTestSpawn()
	p.SpawnUpdates <- s
	c.ExpectNoMessage(t)

	// send egg
	r := getTestEgg()
	p.RaidUpdates <- r
	c.ExpectNoMessage(t)

	// send raid
	r = getTestRaid()
	postedRaidEndTime := r.EndTime
	postedRaidHash := "abc1"
	r.Hash = postedRaidHash
	p.RaidUpdates <- r
	c.ExpectMessage(t)
	assert.Equal(t, testRoom, c.LastRoomID)

	r = getTestRaid()
	r.Hash = "abc2"
	p.RaidUpdates <- r
	c.ExpectMessage(t)
	c.PrintLastMessage()
	assert.Equal(t, testRoom, c.LastRoomID)

	// send same hash again
	r.Hash = "abc2"
	p.RaidUpdates <- r
	c.ExpectNoMessage(t)

	r.Pokemon.ID = 25
	r.Hash = "abc3"
	r.Level = 3
	p.RaidUpdates <- r
	c.ExpectNoMessage(t)

	// out of area
	r = getTestRaid()
	r.Hash = "abc4"
	r.Location = getTestPoint2KMAway()
	p.RaidUpdates <- r
	c.ExpectNoMessage(t)

	// we should have 2 raids now
	rs := p.getOrCreateRoomState(testRoom)
	assert.Equal(t, 2, len(rs.Raids))
	assert.Equal(t, 0, len(rs.Spawns))

	// test expiry...
	// not yet expired
	now := postedRaidEndTime - 1
	rs.removeExpired(now)
	assert.Equal(t, true, rs.raidIsPosted(postedRaidHash))

	// not yet expired
	now = postedRaidEndTime
	rs.removeExpired(now)
	assert.Equal(t, true, rs.raidIsPosted(postedRaidHash))

	// expired
	now = postedRaidEndTime + 1
	rs.removeExpired(now)
	assert.Equal(t, false, rs.raidIsPosted(postedRaidHash))

	// post it again
	// this should work now as we removed the hash above
	r = getTestRaid()
	r.Hash = postedRaidHash
	p.RaidUpdates <- r
	c.ExpectMessage(t)
	assert.Equal(t, 1, len(rs.Raids))
	assert.Equal(t, 0, len(rs.Spawns))

	// clear all states
	rs.clear(true, true)
	assert.Equal(t, 0, len(rs.Raids))
	assert.Equal(t, 0, len(rs.Spawns))

	// send raid to formatted room
	r = getTestRaid()
	r.Hash = "form1"
	r.Pokemon.ID = 1
	p.RaidUpdates <- r
	c.ExpectMessage(t)
	assert.Equal(t, formattedRoom, c.LastRoomID)

	// clear filters
	assert.NotNil(t, p.roomConfigs[formattedRoom].Filter)
	deleted := p.DeleteFilters(formattedRoom)
	assert.Equal(t, true, deleted)
	assert.Nil(t, p.roomConfigs[formattedRoom].Filter)

	// wait
	p.Quit <- true
	<-done
}

func TestPosterRaidsDex(t *testing.T) {
	p, done, c := startPosterPokedex(t, testPokedexFile)

	// set roomservice filter
	testRoom := "!foo@example.com"

	rc := getTestRoomConfig(testRoom)
	rc.Filter[0].ListWanted = false
	rc.Filter[0].ListRaids = true
	rc.Filter[0].PokemonIDs = []int{150, 151}
	p.UpdateRoomConfig(rc)

	// send raid
	r := getTestRaid()
	p.RaidUpdates <- r
	c.ExpectMessage(t)
	assert.Contains(t, c.LastText, "Mewtwo")
	assert.Contains(t, c.LastText, "conke")
	c.PrintLastMessage()

	r.Pokemon.ID = 151
	r.Hash = "itsamew"
	r.Level = 6
	p.RaidUpdates <- r
	c.ExpectMessage(t)
	assert.Contains(t, c.LastText, "Mew ")
	assert.Contains(t, c.LastText, "conke")
	c.PrintLastMessage()

	// wait
	p.Quit <- true
	<-done
}

func TestPoster_ChangeRoomConfig(t *testing.T) {
	testRoom := "!foo@example.com"

	type args struct {
		roomID    string
		rcChange  *RoomConfigChange
		newValues *RoomConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "bad room",
			args: args{
				roomID: "",
			},
			wantErr: true,
		},
		{
			name: "non-observed room",
			args: args{
				roomID: "!other@example.com",
			},
			wantErr: true,
		},
		{
			name: "rcChange is nil",
			args: args{
				roomID:   testRoom,
				rcChange: nil,
			},
			wantErr: true,
		},
		{
			name: "newValues is nil",
			args: args{
				roomID: testRoom,
				rcChange: &RoomConfigChange{
					ChangeAcceptCommands: true,
				},
				newValues: nil,
			},
			wantErr: true,
		},
		{
			name: "append filter with nil filter",
			args: args{
				roomID: testRoom,
				rcChange: &RoomConfigChange{
					Operation: RoomConfigOperationAppendFilter,
				},
				newValues: &RoomConfig{
					Filter: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "append filter with empty filter",
			args: args{
				roomID: testRoom,
				rcChange: &RoomConfigChange{
					Operation: RoomConfigOperationAppendFilter,
				},
				newValues: &RoomConfig{
					Filter: []PokemonFilter{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &testChatter{}
			p := NewPoster(c, nil)
			rc := getTestRoomConfig(testRoom)
			p.UpdateRoomConfig(rc)
			if err := p.ChangeRoomConfig(tt.args.roomID, tt.args.rcChange, tt.args.newValues); (err != nil) != tt.wantErr {
				t.Errorf("Poster.ChangeRoomConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

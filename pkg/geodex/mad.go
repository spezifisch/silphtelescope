package geodex

import (
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // mysql support for goqu
	"github.com/doug-martin/goqu/v9/exec"
	log "github.com/sirupsen/logrus"
)

// GetPokestopName gets the pokestop's name from MAD's DB
func (sdb *SQLDB) GetPokestopName(GUID string) (name string, err error) {
	found, err := sdb.db.From("pokestop").
		Select("name").
		Where(goqu.Ex{"pokestop_id": GUID}).
		Limit(1).
		ScanVal(&name)

	if !found {
		name = "nil"
	}
	return
}

// GetGymName gets the gym's name from MAD's DB
func (sdb *SQLDB) GetGymName(GUID string) (name string, err error) {
	found, err := sdb.db.From("gymdetails").
		Select("name").
		Where(goqu.Ex{"gym_id": GUID}).
		Limit(1).
		ScanVal(&name)

	if !found || name == "unknown" {
		name = "nil"
	}
	return
}

// MADPokestopScanner reads pokestops from MAD
type MADPokestopScanner struct {
	sdb     *SQLDB
	scanner exec.Scanner
}

// NewMADPokestopScanner sets up the scanner
func (sdb *SQLDB) NewMADPokestopScanner() (m *MADPokestopScanner, err error) {
	m = &MADPokestopScanner{sdb: sdb}
	m.scanner, err = sdb.db.From("pokestop").
		Select("pokestop_id", "latitude", "longitude", "name").
		Where(goqu.Ex{"enabled": "1"}).
		Order(goqu.I("pokestop_id").Asc()).
		Executor().
		Scanner()
	return
}

// Next prepares to read the next row
func (m *MADPokestopScanner) Next() bool {
	return m.scanner.Next()
}

// ScanPokestop returns the next row as a Pokestop
func (m *MADPokestopScanner) ScanPokestop() (f Pokestop, err error) {
	f = Pokestop{}
	err = m.scanner.ScanStruct(&f)
	if err == nil {
		// MAD sets name="unknown" when its.. unknown, that's not our way
		if f.Name != nil && *f.Name == "unknown" {
			f.Name = nil
		}
	}
	return
}

// Close closes the scanner
func (m *MADPokestopScanner) Close() {
	if m.scanner != nil {
		if m.scanner.Err() != nil {
			log.WithError(m.scanner.Err()).Error("pokestop scanner failed")
			return
		}
		m.scanner.Close()
	}
}

// MADGymScanner reads gyms from MAD
type MADGymScanner struct {
	sdb     *SQLDB
	scanner exec.Scanner
}

// NewMADGymScanner sets up the scanner
func (sdb *SQLDB) NewMADGymScanner() (m *MADGymScanner, err error) {
	m = &MADGymScanner{sdb: sdb}
	m.scanner, err = sdb.db.From("gym").
		Select("gym.gym_id", "latitude", "longitude", "gymdetails.name").
		Join(
			goqu.T("gymdetails"),
			goqu.On(goqu.Ex{"gym.gym_id": goqu.I("gymdetails.gym_id")}),
		).
		Where(goqu.Ex{"enabled": "1"}).
		Order(goqu.I("gym_id").Asc()).
		Executor().
		Scanner()
	return
}

// Next prepares to read the next row
func (m *MADGymScanner) Next() bool {
	return m.scanner.Next()
}

// ScanGym returns the next row as a Pokestop
func (m *MADGymScanner) ScanGym() (f Gym, err error) {
	f = Gym{}
	err = m.scanner.ScanStruct(&f)
	if err == nil {
		// MAD sets name="unknown" when its.. unknown, that's not our way
		if f.Name != nil && *f.Name == "unknown" {
			f.Name = nil
		}
	}
	return
}

// Close closes the scanner
func (m *MADGymScanner) Close() {
	if m.scanner != nil {
		if m.scanner.Err() != nil {
			log.WithError(m.scanner.Err()).Error("gym scanner failed")
			return
		}
		m.scanner.Close()
	}
}

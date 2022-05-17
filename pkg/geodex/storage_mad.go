package geodex

import (
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql" // mysql support for goqu
	_ "github.com/go-sql-driver/mysql"
)

// SQLDB connects to MAD's MariaDB to get pokestops and gyms
type SQLDB struct {
	db     *goqu.Database
	realDB *sql.DB
}

// NewSQLDB returns a usable MAD db object
func NewSQLDB(hostname, database, username, password string) (db *SQLDB, err error) {
	db = &SQLDB{}

	// connect to mysql
	sourceName := fmt.Sprintf("%s:%s@(%s)/%s", username, password, hostname, database)
	db.realDB, err = sql.Open("mysql", sourceName)
	if err != nil {
		return
	}

	// setup goqu with mysql dialect
	db.db = goqu.New("mysql", db.realDB)
	return
}

// Close closes the database connection
func (sdb *SQLDB) Close() {
	if sdb.realDB != nil {
		sdb.realDB.Close()
		sdb.realDB = nil
	}
}

// GetVersion gets the mysql server version
func (sdb *SQLDB) GetVersion() (version string, err error) {
	_, err = sdb.db.Select(goqu.L("VERSION()")).ScanVal(&version)
	return
}

// Pokestop matches MAD's pokestop table
type Pokestop struct {
	GUID      *string `db:"pokestop_id"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
	Name      *string `db:"name"`
}

// Gym matches MAD's gym table joined with gymdetails
type Gym struct {
	GUID      *string `db:"gym_id"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
	Name      *string `db:"name"`
}

// ToString returns human-readable pokestop info
func (p *Pokestop) ToString() string {
	empty := "nil"
	guid := &empty
	if p.GUID != nil {
		guid = p.GUID
	}
	name := &empty
	if p.Name != nil {
		name = p.Name
	}
	return fmt.Sprintf("GUID=%s Lat=%f Lon=%f Name=%s", *guid, p.Latitude, p.Longitude, *name)
}

// ToFort returns a Fort for TDB
func (p *Pokestop) ToFort() *Fort {
	return &Fort{
		GUID:      p.GUID,
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
		Name:      p.Name,
		Type:      FortTypeStop,
	}
}

// ToString returns human-readable gym info
func (p *Gym) ToString() string {
	empty := "nil"
	guid := &empty
	if p.GUID != nil {
		guid = p.GUID
	}
	name := &empty
	if p.Name != nil {
		name = p.Name
	}
	return fmt.Sprintf("GUID=%s Lat=%f Lon=%f Name=%s", *guid, p.Latitude, p.Longitude, *name)
}

// ToFort returns a Fort for TDB
func (p *Gym) ToFort() *Fort {
	return &Fort{
		GUID:      p.GUID,
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
		Name:      p.Name,
		Type:      FortTypeGym,
	}
}

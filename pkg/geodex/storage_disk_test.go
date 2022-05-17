package geodex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiskDB(t *testing.T) {
	var err error

	basePath := "test-data-2342"
	db := NewDiskDB(&basePath)
	defer db.Drop()

	_, err = db.GetFort("nonexistent_guid")
	assert.Error(t, err)

	badFort := Fort{}
	err = db.SaveFort(&badFort)
	assert.Error(t, err)

	goodGUID := "8d07e423eb2898c9f853d7b9aec08905.16"
	goodName := "Good Gym"
	goodFort := Fort{
		GUID:      &goodGUID,
		Latitude:  52.503355,
		Longitude: 13.435746,
		Name:      &goodName,
		Type:      FortTypeGym,
	}
	err = db.SaveFort(&goodFort)
	assert.NoError(t, err)

	retFort, err := db.GetFort(goodGUID)
	assert.NoError(t, err)
	assert.Equal(t, &goodFort, retFort)
	assert.NotEqual(t, nil, retFort.GUID)
	assert.NotEqual(t, nil, retFort.Name)
	assert.Equal(t, goodGUID, *retFort.GUID)
	assert.Equal(t, goodName, *retFort.Name)
}

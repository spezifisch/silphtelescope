package geodex

import (
	"reflect"
	"testing"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
)

func TestFort_ToString(t *testing.T) {
	type fields struct {
		GUID      *string
		Latitude  float64
		Longitude float64
		Name      *string
		Type      FortType
	}
	guid := "guido"
	name := "nameo"
	lat := 12.34
	lon := 56.78

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "valid with name",
			fields: fields{
				GUID:      &guid,
				Latitude:  lat,
				Longitude: lon,
				Name:      &name,
				Type:      FortTypeGym,
			},
			want: "GUID=guido Type=Gym (12.340000,56.780000) Name: nameo",
		},
		{
			name: "valid without name",
			fields: fields{
				GUID:      &guid,
				Latitude:  lat,
				Longitude: lon,
				Name:      nil,
				Type:      FortTypeStop,
			},
			want: "GUID=guido Type=Stop (12.340000,56.780000) Name: nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fort{
				GUID:      tt.fields.GUID,
				Latitude:  tt.fields.Latitude,
				Longitude: tt.fields.Longitude,
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
			}
			if got := f.ToString(); got != tt.want {
				t.Errorf("Fort.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFort_GetName(t *testing.T) {
	guid := "guido"
	name := "nameo"
	lat := 12.34
	lon := 56.78

	type fields struct {
		GUID      *string
		Latitude  float64
		Longitude float64
		Name      *string
		Type      FortType
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "valid with name",
			fields: fields{
				GUID:      &guid,
				Latitude:  lat,
				Longitude: lon,
				Name:      &name,
				Type:      FortTypeGym,
			},
			want: name,
		},
		{
			name: "valid without name",
			fields: fields{
				GUID:      &guid,
				Latitude:  lat,
				Longitude: lon,
				Name:      nil,
				Type:      FortTypeStop,
			},
			want: "Stop:guido",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fort{
				GUID:      tt.fields.GUID,
				Latitude:  tt.fields.Latitude,
				Longitude: tt.fields.Longitude,
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
			}
			if got := f.GetName(); got != tt.want {
				t.Errorf("Fort.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFort_Location(t *testing.T) {
	guid := "guido"
	name := "nameo"
	lat := 12.34
	lon := 56.78

	type fields struct {
		GUID      *string
		Latitude  float64
		Longitude float64
		Name      *string
		Type      FortType
	}
	tests := []struct {
		name   string
		fields fields
		want   *pogo.Location
	}{
		{
			name: "valid location",
			fields: fields{
				GUID:      &guid,
				Latitude:  lat,
				Longitude: lon,
				Name:      &name,
				Type:      FortTypeStop,
			},
			want: &pogo.Location{
				Latitude:  lat,
				Longitude: lon,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fort{
				GUID:      tt.fields.GUID,
				Latitude:  tt.fields.Latitude,
				Longitude: tt.fields.Longitude,
				Name:      tt.fields.Name,
				Type:      tt.fields.Type,
			}
			if got := f.Location(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fort.Location() = %v, want %v", got, tt.want)
			}
		})
	}
}

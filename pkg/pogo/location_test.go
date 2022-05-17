package pogo

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocation_DistanceTo(t *testing.T) {
	type fields struct {
		Latitude  float64
		Longitude float64
	}
	type args struct {
		p2 *Location
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "same point",
			fields: fields{
				Latitude:  12.345,
				Longitude: -54.321,
			},
			args: args{
				&Location{12.345, -54.321},
			},
			want: 0,
		},
		{
			name: "10m distance",
			fields: fields{
				Latitude:  12.345,
				Longitude: -54.321,
			},
			args: args{
				&Location{12.345064, -54.320935},
			},
			want: 10,
		},
		{
			name: "10km distance",
			fields: fields{
				Latitude:  12.345,
				Longitude: -54.321,
			},
			args: args{
				&Location{12.281462, -54.255982},
			},
			// wanted 10km distance from starting point
			// used https://rechneronline.de/geo-coordinates/#bearing to calculate args point
			// the difference is probably because we don't consider earth's radius dependence on latitude
			want: 9990,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Location{
				Latitude:  tt.fields.Latitude,
				Longitude: tt.fields.Longitude,
			}
			// round result
			if got := l.DistanceTo(tt.args.p2); math.Round(got) != tt.want {
				t.Errorf("Location.DistanceTo() = %v rounded to %v, want %v", got, math.Round(got), tt.want)
			}
		})
	}
}

func TestLocation_BearingTo(t *testing.T) {
	type fields struct {
		Latitude  float64
		Longitude float64
	}
	type args struct {
		p2 *Location
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "same point",
			fields: fields{
				Latitude:  12.345,
				Longitude: -54.321,
			},
			args: args{
				&Location{12.345, -54.321},
			},
			want: 0,
		},
		{
			name: "45° bearing",
			fields: fields{
				Latitude:  12.345,
				Longitude: -54.321,
			},
			args: args{
				&Location{12.345064, -54.320935},
			},
			want: 45,
		},
		{
			name: "135° bearing",
			fields: fields{
				Latitude:  12.345,
				Longitude: -54.321,
			},
			args: args{
				&Location{12.281462, -54.255982},
			},
			want: 135,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Location{
				Latitude:  tt.fields.Latitude,
				Longitude: tt.fields.Longitude,
			}
			// round result
			if got := l.BearingTo(tt.args.p2); math.Round(got) != tt.want {
				t.Errorf("Location.BearingTo() = %v rounded to %v, want %v", got, math.Round(got), tt.want)
			}
		})
	}
}

func TestLocationRadius(t *testing.T) {
	lat := 12.345
	lon := -54.321
	nearLoc := Location{12.345064, -54.320935}
	farLoc := Location{12.281462, -54.255982}

	lr := NewLocationRadius(lat, lon, 11)
	assert.Equal(t, lat, lr.Latitude)
	assert.Equal(t, lon, lr.Longitude)
	assert.Equal(t, 11.0, lr.RadiusM)

	res := lr.Contains(&nearLoc)
	assert.Equal(t, true, res)

	res = lr.Contains(&farLoc)
	assert.Equal(t, false, res)

	lr = NewLocationRadius(lat, lon, 9)
	assert.Equal(t, lat, lr.Latitude)
	assert.Equal(t, lon, lr.Longitude)
	assert.Equal(t, 9.0, lr.RadiusM)

	res = lr.Contains(&nearLoc)
	assert.Equal(t, false, res)

	res = lr.Contains(&farLoc)
	assert.Equal(t, false, res)
}

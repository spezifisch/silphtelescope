package roomservice

import (
	"reflect"
	"testing"

	"github.com/spezifisch/silphtelescope/pkg/pogo"
	"github.com/stretchr/testify/assert"
)

func TestArgParser(t *testing.T) {
	args := []string{"foo", "bar", "1", "3.14", "-1"}
	ap := NewArgParser(args)
	assert.Equal(t, args, ap.args)
	assert.Equal(t, 5, ap.Count())
	assert.Equal(t, true, ap.isValidIndex(0))
	assert.Equal(t, true, ap.isValidIndex(1))
	assert.Equal(t, true, ap.isValidIndex(2))
	assert.Equal(t, true, ap.isValidIndex(3))
	assert.Equal(t, true, ap.isValidIndex(4))
	assert.Equal(t, false, ap.isValidIndex(5))
	assert.Equal(t, false, ap.isValidIndex(-1))

	var err error
	// AsInt
	valInt, err := ap.AsInt(2)
	assert.NoError(t, err)
	assert.Equal(t, 1, valInt)
	valInt, err = ap.AsInt(4)
	assert.NoError(t, err)
	assert.Equal(t, -1, valInt)
	_, err = ap.AsInt(3)
	assert.Error(t, err)
	_, err = ap.AsInt(1)
	assert.Error(t, err)
	_, err = ap.AsInt(99)
	assert.Error(t, err)

	// AsFloat
	valFloat, err := ap.AsFloat(2)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, valFloat)
	valFloat, err = ap.AsFloat(3)
	assert.NoError(t, err)
	assert.Equal(t, 3.14, valFloat)
	valFloat, err = ap.AsFloat(4)
	assert.NoError(t, err)
	assert.Equal(t, -1.0, valFloat)
	_, err = ap.AsFloat(0)
	assert.Error(t, err)
	_, err = ap.AsFloat(1)
	assert.Error(t, err)
	_, err = ap.AsFloat(99)
	assert.Error(t, err)

	// AsString
	valString, err := ap.AsString(0)
	assert.NoError(t, err)
	assert.Equal(t, "foo", valString)
	valString, err = ap.AsString(2)
	assert.NoError(t, err)
	assert.Equal(t, "1", valString)
	_, err = ap.AsString(23)
	assert.Error(t, err)
}

func TestArgParser_AsIntArray(t *testing.T) {
	type fields struct {
		args []string
	}
	type args struct {
		index int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantArr []int
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				args: []string{},
			},
			args:    args{0},
			wantErr: true,
		},
		{
			name: "bad int",
			fields: fields{
				args: []string{"hi", "bye"},
			},
			args:    args{0},
			wantErr: true,
		},
		{
			name: "single int",
			fields: fields{
				args: []string{"23"},
			},
			args:    args{0},
			wantArr: []int{23},
		},
		{
			name: "multiple ints",
			fields: fields{
				args: []string{"hi", "23,5,42", "bye"},
			},
			args:    args{1},
			wantArr: []int{23, 5, 42},
		},
		{
			name: "multiple ints with error",
			fields: fields{
				args: []string{"hi", "23,x,42", "bye"},
			},
			args:    args{1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArgParser{
				args: tt.fields.args,
			}
			gotArr, err := a.AsIntArray(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArgParser.AsIntArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotArr, tt.wantArr) {
				t.Errorf("ArgParser.AsIntArray() = %v, want %v", gotArr, tt.wantArr)
			}
		})
	}
}

func TestArgParser_AsLocation(t *testing.T) {
	type fields struct {
		args []string
	}
	type args struct {
		latIdx int
		lonIdx int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantL   pogo.Location
		wantErr bool
	}{
		{
			name: "valid location",
			fields: fields{
				args: []string{"1.2", "3.4", "bye"},
			},
			args:  args{0, 1},
			wantL: pogo.Location{Latitude: 1.2, Longitude: 3.4},
		},
		{
			name: "bad latitude",
			fields: fields{
				args: []string{"x1.2", "3.4", "bye"},
			},
			args:    args{0, 1},
			wantErr: true,
		},
		{
			name: "bad longitude",
			fields: fields{
				args: []string{"1.2", "3.y4", "bye"},
			},
			args:    args{0, 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArgParser{
				args: tt.fields.args,
			}
			gotL, err := a.AsLocation(tt.args.latIdx, tt.args.lonIdx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArgParser.AsLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotL, tt.wantL) {
				t.Errorf("ArgParser.AsLocation() = %v, want %v", gotL, tt.wantL)
			}
		})
	}
}

func TestArgParser_AsLocationRadius(t *testing.T) {
	type fields struct {
		args []string
	}
	type args struct {
		latIdx    int
		lonIdx    int
		radiusIdx int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantLr  pogo.LocationRadius
		wantErr bool
	}{
		{
			name: "valid location+radius",
			fields: fields{
				args: []string{"1.2", "3.4", "100"},
			},
			args: args{0, 1, 2},
			wantLr: pogo.LocationRadius{
				Location: pogo.Location{
					Latitude:  1.2,
					Longitude: 3.4,
				},
				RadiusM: 100,
			},
		},
		{
			name: "bad latitude",
			fields: fields{
				args: []string{"foo", "3.4", "100"},
			},
			args:    args{0, 1, 2},
			wantErr: true,
		},
		{
			name: "bad longitude",
			fields: fields{
				args: []string{"1.2", "if", "100"},
			},
			args:    args{0, 1, 2},
			wantErr: true,
		},
		{
			name: "bad radius",
			fields: fields{
				args: []string{"1.2", "100", "bye"},
			},
			args:    args{0, 1, 2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &ArgParser{
				args: tt.fields.args,
			}
			gotLr, err := a.AsLocationRadius(tt.args.latIdx, tt.args.lonIdx, tt.args.radiusIdx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ArgParser.AsLocationRadius() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotLr, tt.wantLr) {
				t.Errorf("ArgParser.AsLocationRadius() = %v, want %v", gotLr, tt.wantLr)
			}
		})
	}
}

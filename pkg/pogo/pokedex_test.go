package pogo

import (
	"testing"
)

func TestPokedex_GetNamesByID(t *testing.T) {
	type fields struct {
		fileName string
		entries  []*PokedexEntry
	}
	testFields := fields{
		fileName: "foo",
		entries: []*PokedexEntry{
			{
				ID:     1,
				NameEN: "en1",
				NameDE: "de1",
			},
			{
				ID:     2,
				NameEN: "en2",
				NameDE: "de2",
			},
		},
	}

	type args struct {
		id int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantNameEN string
		wantNameDE string
		wantErr    bool
	}{
		{
			name:       "valid lookup 1",
			fields:     testFields,
			args:       args{1},
			wantNameEN: "en1",
			wantNameDE: "de1",
			wantErr:    false,
		},
		{
			name:       "valid lookup 2",
			fields:     testFields,
			args:       args{2},
			wantNameEN: "en2",
			wantNameDE: "de2",
			wantErr:    false,
		},
		{
			name:    "not found lookup",
			fields:  testFields,
			args:    args{3},
			wantErr: true,
		},
		{
			name:    "invalid id 0",
			fields:  testFields,
			args:    args{0},
			wantErr: true,
		},
		{
			name:    "invalid id 0",
			fields:  testFields,
			args:    args{-1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pokedex{
				fileName: tt.fields.fileName,
				entries:  tt.fields.entries,
			}
			gotNameEN, gotNameDE, err := p.GetNamesByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pokedex.GetNamesByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotNameEN != tt.wantNameEN {
				t.Errorf("Pokedex.GetNamesByID() gotNameEN = %v, want %v", gotNameEN, tt.wantNameEN)
			}
			if gotNameDE != tt.wantNameDE {
				t.Errorf("Pokedex.GetNamesByID() gotNameDE = %v, want %v", gotNameDE, tt.wantNameDE)
			}
		})
	}
}

func TestPokedex_GetIDByName(t *testing.T) {
	type fields struct {
		fileName string
		entries  []*PokedexEntry
	}
	testFields := fields{
		fileName: "foo",
		entries: []*PokedexEntry{
			{
				ID:     1,
				NameEN: "en1",
				NameDE: "de1",
			},
			{
				ID:     2,
				NameEN: "en2",
				NameDE: "de2",
			},
		},
	}

	type args struct {
		wantedName string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantID     int
		wantNameEN string
		wantNameDE string
		wantErr    bool
	}{
		{
			name:       "found 1 english",
			fields:     testFields,
			args:       args{"en1"},
			wantID:     1,
			wantNameEN: "en1",
			wantNameDE: "de1",
			wantErr:    false,
		},
		{
			name:       "found 1 english case",
			fields:     testFields,
			args:       args{"EN1"},
			wantID:     1,
			wantNameEN: "en1",
			wantNameDE: "de1",
			wantErr:    false,
		},
		{
			name:       "found 1 german",
			fields:     testFields,
			args:       args{"de1"},
			wantID:     1,
			wantNameEN: "en1",
			wantNameDE: "de1",
			wantErr:    false,
		},
		{
			name:       "found 1 german case",
			fields:     testFields,
			args:       args{"dE1"},
			wantID:     1,
			wantNameEN: "en1",
			wantNameDE: "de1",
			wantErr:    false,
		},
		{
			name:       "found 2 english",
			fields:     testFields,
			args:       args{"en2"},
			wantID:     2,
			wantNameEN: "en2",
			wantNameDE: "de2",
			wantErr:    false,
		},
		{
			name:    "not found 1",
			fields:  testFields,
			args:    args{"en2x"},
			wantErr: true,
		},
		{
			name:    "not found 2",
			fields:  testFields,
			args:    args{""},
			wantErr: true,
		},
		{
			name:    "not found 3",
			fields:  testFields,
			args:    args{"Ä"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pokedex{
				fileName: tt.fields.fileName,
				entries:  tt.fields.entries,
			}
			gotID, gotNameEN, gotNameDE, err := p.GetIDByName(tt.args.wantedName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pokedex.GetIDByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotID != tt.wantID {
				t.Errorf("Pokedex.GetIDByName() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotNameEN != tt.wantNameEN {
				t.Errorf("Pokedex.GetIDByName() gotNameEN = %v, want %v", gotNameEN, tt.wantNameEN)
			}
			if gotNameDE != tt.wantNameDE {
				t.Errorf("Pokedex.GetIDByName() gotNameDE = %v, want %v", gotNameDE, tt.wantNameDE)
			}
		})
	}
}

package geodex

import "testing"

func TestIsValidGUID(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid guid 16",
			args: args{"39665b7dd75b4c10b1f4c796676c2cf4.16"},
			want: true,
		},
		{
			name: "valid guid 11",
			args: args{"464d1951148045aaaffe14ce2e956397.11"},
			want: true,
		},
		{
			name: "valid guid 12",
			args: args{"db0a966397fa4e289523d17eacdaaf17.12"},
			want: true,
		},
		{
			name: "invalid guid empty",
			args: args{""},
			want: false,
		},
		{
			name: "invalid guid bad chars 1",
			args: args{"xyz"},
			want: false,
		},
		{
			name: "invalid guid bad chars 2",
			args: args{"000:12"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidGUID(tt.args.val); got != tt.want {
				t.Errorf("IsValidGUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

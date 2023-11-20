package types

import (
	"reflect"
	"testing"
	"time"
)

func TestParseJSONTime(t *testing.T) {
	type args struct {
		date string
	}
	tests := []struct {
		name    string
		args    args
		want    JSONTime
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				date: "2018-01-01T00:00:00",
			},
			want: JSONTime{
				Time: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseJSONTime(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJSONTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseJSONTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}

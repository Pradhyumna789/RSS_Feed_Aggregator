package config

import (
	"reflect"
	"testing"
)

func TestConfig_SetUser(t *testing.T) {
	type fields struct {
		DbURL           string
		CurrentUserName string
	}
	type args struct {
		current_user_name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{
				DbURL:           tt.fields.DbURL,
				CurrentUserName: tt.fields.CurrentUserName,
			}
			config.SetUser(tt.args.current_user_name)
		})
	}
}

func TestRead(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Read(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

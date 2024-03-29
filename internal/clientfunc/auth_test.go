package clientfunc

import (
	"reflect"
	"testing"

	"github.com/gambruh/simplevault/internal/auth"
)

func Test_getUserDataFromFile(t *testing.T) {
	// creation of a temp file

	// filling the file with some values

	tests := []struct {
		name    string
		want    auth.LoginData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUserDataFromFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserDataFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUserDataFromFile() = %v, want %v", got, tt.want)
			}
		})
	}

	// deleting temp file

}

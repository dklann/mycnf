package mycnf

import (
	"os"
	"reflect"
	"testing"
)

func TestReadMyCnf(t *testing.T) {
	var configurationFile string
	configurationFile = os.Getenv("HOME") + "/.my.cnf"
	var myConfigFile = &configurationFile
	type want map[string]string

	type args struct {
		configFile *string
		profile    string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			"should not work 1",
			args{
				myConfigFile,
				"client",
			},
			want{
				"dbhost": "service",
				"dbpass": "YHYgXEE",
				"dbname": "",
				"dbuser": "",
				"dbport": "3306",
			},
			false,
		},
		{
			"should not work 2",
			args{
				myConfigFile,
				"mysql",
			},
			want{
				"dbhost": "localhost",
				"dbpass": "6fd49tieoaniPNDkwAlNJEaPkwib",
				"dbname": "",
				"dbuser": "root",
				"dbport": "3306",
			},
			false,
		},
		{
			"should work 1",
			args{
				myConfigFile,
				"client",
			},
			want{
				"dbhost": "localhost",
				"dbuser": "dklann",
				"dbpass": "1q2w3e4r",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadMyCnf(tt.args.configFile, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadMyCnf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadMyCnf() = %v, want %v", got, tt.want)
			}
		})
	}
}

package crossref

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		apiBase string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "A",
			args:    args{apiBase: "http://api.crossref.org/v1"},
			want:    []string{"http://api.crossref.org/v1/", "http://api.crossref.org/v1/works"},
			wantErr: false,
		},
		{
			name:    "B",
			args:    args{apiBase: "http://api.crossref.org/v1/"},
			want:    []string{"http://api.crossref.org/v1/", "http://api.crossref.org/v1/works"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.apiBase)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.APIBase().String(), tt.want[0]) {
				t.Errorf("New() = %v, want %v", got, tt.want[0])
			}
			if !reflect.DeepEqual(got.WorksPath().String(), tt.want[1]) {
				t.Errorf("New() = %v, want %v", got, tt.want[1])
			}
		})
	}
}

func TestGetWork(t *testing.T) {
	c, err := New("http://api.crossref.org/v1")
	if err != nil {
		t.Fatal(err)
	}

	id := DOI("10.3390/act7010007")
	work, err := GetWork(c, id)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("work: %+v", work)
	t.Fail()
}

func Test_decodeWork(t *testing.T) {
	f, err := os.Open("testdata/work.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	work, err := decodeWork(f)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", work)
}

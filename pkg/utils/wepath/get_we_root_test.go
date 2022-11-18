package wepath

import (
	"testing"
)

func TestGetWorkEnvironmentRoot(t *testing.T) {
	type args struct {
		from string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should find a work environment",
			args: args{
				from: "/Users/sven.liebig/workspace/repositories/isbj/commons/ansible-paas",
			},
			want: "/Users/sven.liebig/workspace/.work-environment",
		},
		{
			name: "should not find a work environment",
			args: args{
				from: "/Users/sven.liebig",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetWorkEnvironmentRoot(tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWorkEnvironmentRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetWorkEnvironmentRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

// BenchmarkGitClone-8   	   93578	     12059 ns/op	    2240 B/op	      21 allocs/op
func BenchmarkGetWorkEnvironmentRoot(b *testing.B) {
	for n := 0; n < b.N; n++ {
		GetWorkEnvironmentRoot("/Users/sven.liebig/workspace/repositories/isbj/commons/ansible-paas")
	}
}

func MapWith4EntriesGettingThe3rd() string {
	m := map[string]string{"a": "asdf", "b": "asdf", "c": "asdf", "d": "asdf"}

	v, ok := m["c"]

	if ok {
		return v
	}

	return ""
}

func ArrayWith4EntriesGettingThe3rd() string {
	m := []string{"a", "b", "c", "d"}

	for _, k := range m {
		if k == "c" {
			return k
		}
	}

	return ""
}

func BenchmarkGetWorkReadingFromMapWith4EntriesGettingThe3rd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		MapWith4EntriesGettingThe3rd()
	}
}

func BenchmarkGetWorkReadingFromArrayWith4EntriesGettingThe3rd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ArrayWith4EntriesGettingThe3rd()
	}
}

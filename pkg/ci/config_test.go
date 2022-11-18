package ci

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	type args struct {
		configPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *CIConfig
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadConfig(tt.args.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func MapWith4EntriesGettingThe3rdComplex() *CI {
	m := map[string]*CI{"a": {Identifier: "a"}, "b": {Identifier: "b"}, "c": {Identifier: "c"}, "d": {Identifier: "d"}}

	v, ok := m["c"]

	if ok {
		return v
	}

	return nil
}

func ArrayWith4EntriesGettingThe3rdComplex() *CI {
	m := []*CI{{Identifier: "a"}, {Identifier: "b"}, {Identifier: "c"}, {Identifier: "d"}}

	for _, k := range m {
		if k.Identifier == "c" {
			return k
		}
	}

	return nil
}

func BenchmarkGetWorkReadingFromMapWith4EntriesGettingThe3rdComplex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		MapWith4EntriesGettingThe3rdComplex()
	}
}

func BenchmarkGetWorkReadingFromArrayWith4EntriesGettingThe3rdComplex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		ArrayWith4EntriesGettingThe3rdComplex()
	}
}

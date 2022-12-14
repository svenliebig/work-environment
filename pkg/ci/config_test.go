package ci

import (
	"fmt"
	"testing"
)

type C interface {
	Hello() string
}

type D interface {
	C
}

type A struct {
	Hello string
}

func (a *A) Get() string {
	return a.Hello
}

type B struct {
}

func (b *B) Hello() string {
	return ""
}

var (
	d D = &B{}
)

func TestThings(t *testing.T) {
	t.Run("inheritance things", func(t *testing.T) {
		fmt.Println(d)
		t.Fatal()
	})
}

type CI struct {
	Identifier string
}

// @comm @alex

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

// @comm @alex, why is that
func BenchmarkArrayCreationCap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = make([]string, 0, 100)
	}
}

func BenchmarkArrayCreationSize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = make([]string, 100)
	}
}

package array

import "testing"

func TestRemove(t *testing.T) {
	t.Run("should remove index 0 (1 element)", func(t *testing.T) {
		arr, err := Remove([]string{"h"}, "h")

		if err != nil {
			t.Fatal(err)
		}

		if len(arr) != 0 {
			t.Fatalf("expected array to have lenght 0 but it was %d.", len(arr))
		}
	})

	t.Run("should remove index 0 (2 elements)", func(t *testing.T) {
		arr, err := Remove([]string{"h", "f"}, "h")

		if err != nil {
			t.Fatal(err)
		}

		if len(arr) != 1 {
			t.Fatalf("expected array to have lenght 1 but it was %d -> %s.", len(arr), arr)
		}

		if arr[0] != "f" {
			t.Fatalf("expected element 0 in arr to be 'f', but it was %s.", arr[0])
		}
	})

	t.Run("should remove index 1 (2 elements)", func(t *testing.T) {
		arr, err := Remove([]string{"h", "f"}, "f")

		if err != nil {
			t.Fatal(err)
		}

		if len(arr) != 1 {
			t.Fatalf("expected array to have lenght 1 but it was %d -> %s.", len(arr), arr)
		}

		if arr[0] != "h" {
			t.Fatalf("expected element 0 in arr to be 'h', but it was %s.", arr[0])
		}
	})

	t.Run("should remove index 2 (3 elements)", func(t *testing.T) {
		arr, err := Remove([]string{"h", "f", "s"}, "f")

		if err != nil {
			t.Fatal(err)
		}

		if len(arr) != 2 {
			t.Fatalf("expected array to have lenght 2 but it was %d -> %s.", len(arr), arr)
		}

		if arr[0] != "h" {
			t.Fatalf("expected element 0 in arr to be 'h', but it was %s.", arr[0])
		}

		if arr[1] != "s" {
			t.Fatalf("expected element 0 in arr to be 's', but it was %s.", arr[0])
		}
	})
}

// BenchmarkRemoveFromThree-8   	37061571	        32.21 ns/op	      32 B/op	       1 allocs/op
// BenchmarkRemoveFromThree-8   	37251589	        32.13 ns/op	      32 B/op	       1 allocs/op
func BenchmarkRemoveFromThree(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Remove([]string{"h", "f", "s"}, "f")
	}
}

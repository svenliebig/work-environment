package array

func Remove(arr []string, s string) ([]string, error) {
	for i, e := range arr {
		if e == s {
			// TODO copy last element to position i and return the base slice
			c := make([]string, 0, len(arr)-1)
			p := arr[:i]
			c = append(c, p...)
			t := arr[i+1:]
			c = append(c, t...)
			return c, nil
		}
	}
	return arr, nil
}

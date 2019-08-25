package util

func FirstIndexOf(inputs []string, key string) int {
	for index, val := range inputs {
		if key == val {
			return index
		}
	}
	return -1
}

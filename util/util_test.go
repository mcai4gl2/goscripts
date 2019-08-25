package util

import "testing"

func TestCanFind(t *testing.T) {
	result := FirstIndexOf([]string{"1", "2", "3"}, "3")
	if result != 2 {
		t.Errorf("Searching 3 in [1, 2, 3] shall give 2 instead of %d", result)
	}
}

func TestCannotFind(t *testing.T) {
	result := FirstIndexOf([]string{"1", "2", "3"}, "4")
	if result != -1 {
		t.Errorf("Searching 4 in [1, 2,3] shall give -1 as 4 is not in instead of %d",
			result)
	}
}

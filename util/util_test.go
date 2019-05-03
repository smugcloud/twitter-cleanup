package util

import (
	"testing"
	"time"
)

var getDateTests = []int{
	12, 3,
}

func TestDate(t *testing.T) {
	time := time.Date(2019, time.October, 3, 0, 0, 0, 0, time.UTC)
	expected := []string{"201810030000", "201907030000"}
	for k, v := range getDateTests {
		r := GetToDate(v, time)

		if r != expected[k] {

			t.Errorf("Expected toDate not returned: %v - %v\n", r, expected[k])
		}
	}
}

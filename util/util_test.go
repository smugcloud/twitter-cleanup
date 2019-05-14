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

func TestURLParse(t *testing.T) {
	scenarios := []struct {
		given    string
		expected string
	}{
		{"localhost", "https://localhost/"},
		{"http://127.0.0.1:53941", "http://127.0.0.1:53941/"},
		{"https://api.twitter.com/", "https://api.twitter.com/"},
	}

	for _, i := range scenarios {
		p := URLParse(i.given)

		if i.expected != p.String() {
			t.Errorf("Expected URL does not match returned: %v, %v\n", i.expected, p)
		}
	}
}

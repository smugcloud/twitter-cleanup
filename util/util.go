package util

import (
	"strconv"
	"time"
)

func GetToDate(months int, now time.Time) string {
	// Gross time manipulation
	var smo, sd string
	var mo int
	ad := now.AddDate(0, -months, 0)

	mo = int(ad.Month())

	if mo < 10 {
		smo = "0" + strconv.Itoa(mo)

	} else {
		smo = strconv.Itoa(mo)
	}

	d := ad.Day()
	if d < 10 {
		sd = "0" + strconv.Itoa(d)

	} else {
		sd = strconv.Itoa(d)
	}

	return strconv.Itoa(ad.Year()) + smo + sd
}

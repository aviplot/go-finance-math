package financial

import (
	"math"
)

func GetRmz(fZins, fZzr, fBw, fZw float64, nF int64) float64 {
	var fRmz float64
	if fZins == 0.0 {
		fRmz = (fBw + fZw) / fZzr
	} else {
		fTerm := math.Pow(1.0+fZins, fZzr)
		if nF > 0 {
			fRmz = (fZw*fZins/(fTerm-1.0) + fBw*fZins/(1.0-1.0/fTerm)) / (1.0 + fZins)
		} else {
			fRmz = fZw*fZins/(fTerm-1.0) + fBw*fZins/(1.0-1.0/fTerm)
		}
	}
	return -fRmz
}

func GetZw(fZins, fZzr, fRmz, fBw float64, nF int64) float64 {
	var fZw float64
	if fZins == 0.0 {
		fZw = fBw + fRmz*fZzr
	} else {
		fTerm := math.Pow(1.0+fZins, fZzr)
		if nF > 0 {
			fZw = fBw*fTerm + fRmz*(1.0+fZins)*(fTerm-1.0)/fZins
		} else {
			fZw = fBw*fTerm + fRmz*(fTerm-1.0)/fZins
		}
	}

	return -fZw
}

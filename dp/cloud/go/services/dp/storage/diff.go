package storage

import (
	"math"
)

const (
	radiusM            = 6371e3
	distanceThresholdM = 10
)

func ShouldEnodebdUpdateInstallationParams(prev *DBCbsd, next *DBCbsd) bool {
	// TODO this should probably moved out from storage
	return paramsChanges(next, prev) || coordinatesChanged(next, prev)
}

func paramsChanges(prev *DBCbsd, next *DBCbsd) bool {
	return prev.HeightM != next.HeightM ||
		prev.HeightType != next.HeightType ||
		prev.IndoorDeployment != next.IndoorDeployment ||
		prev.AntennaGainDbi != next.AntennaGainDbi
}

func coordinatesChanged(prev *DBCbsd, next *DBCbsd) bool {
	return coordinatesAreEmpty(prev, next) ||
		distanceM(
			prev.LatitudeDeg.Float64, prev.LongitudeDeg.Float64,
			next.LatitudeDeg.Float64, next.LongitudeDeg.Float64,
		) > distanceThresholdM
}

func coordinatesAreEmpty(prev *DBCbsd, next *DBCbsd) bool {
	return !prev.LatitudeDeg.Valid ||
		!prev.LongitudeDeg.Valid ||
		!next.LatitudeDeg.Valid ||
		!next.LongitudeDeg.Valid
}

func distanceM(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	lat1 = lat1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180
	dLon := lon2 - lon1
	dLat := lat2 - lat1
	sinLat := math.Sin(dLat / 2)
	sinLon := math.Sin(dLon / 2)
	res := sinLat*sinLat + math.Cos(lat1)*math.Cos(lat2)*sinLon*sinLon
	return 2 * math.Asin(math.Sqrt(res)) * radiusM
}

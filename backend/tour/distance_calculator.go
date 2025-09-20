package main

import (
	"math"
	"sort"
)

// DistanceCalculator provides static methods for calculating distances
type DistanceCalculator struct{}

func (DistanceCalculator) HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Differences
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Haversine formula
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (dc DistanceCalculator) CalculateTourDistance(keyPoints []KeyPoint) float64 {
	if len(keyPoints) < 2 {
		return 0
	}

	sortedKeyPoints := make([]KeyPoint, len(keyPoints))
	copy(sortedKeyPoints, keyPoints)
	sort.Slice(sortedKeyPoints, func(i, j int) bool {
		return sortedKeyPoints[i].Order < sortedKeyPoints[j].Order
	})

	totalDistance := 0.0
	for i := 0; i < len(sortedKeyPoints)-1; i++ {
		current := sortedKeyPoints[i]
		next := sortedKeyPoints[i+1]

		distance := dc.HaversineDistance(
			current.Latitude, current.Longitude,
			next.Latitude, next.Longitude,
		)
		totalDistance += distance
	}

	return totalDistance
}

var Calculator = DistanceCalculator{}

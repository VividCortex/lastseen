// Package lastseen implements a last-seen sketch.
package lastseen

import (
	"fmt"
	"math"
	"time"
)

// zeroTime is the default value for time.Time.
var zeroTime time.Time

const (
	multiplier = 10
	numBuckets = 7
)

// Sketch is a last-seen sketch.
type Sketch [numBuckets][]time.Time

// NewSketch returns a sketch that can hold N elements,
// offering ~1% false positives for lookups. N should be
// at most 45805.
func NewSketch(N int) Sketch {
	sk := Sketch{}
	if primeFloor := math.Ceil(float64(N*multiplier) / numBuckets); primeFloor > 0 &&
		primeFloor < math.MaxUint16 {
		for i, prime := range GetPrimesFrom(uint16(primeFloor), numBuckets) {
			sk[i] = make([]time.Time, prime)
		}
	}
	if len(sk[numBuckets-1]) == 0 { // could not get numBuckets primes
		panic("lastseen: unsupported size for NewSketch")
	}
	return sk
}

// Update updates the timestamps in the sketch to now.
// It returns true if the key was seen before.
func (sk Sketch) Update(key uint64, now time.Time) bool {
	// Assume we have seen it.
	present := true

	for i := 0; i < numBuckets; i++ {
		index := int(key % uint64(len(sk[i])))
		if sk[i][index].IsZero() {
			// not seen yet
			present = false
		}
		sk[i][index] = now
	}

	return present
}

// Lookup returns the latest possible time that the
// key could have been updated in the sketch. The zero
// time is returned if the key has never been updated.
func (sk Sketch) Lookup(key uint64) time.Time {
	var lastSeen time.Time

	for i := 0; i < numBuckets; i++ {
		index := int(key % uint64(len(sk[i])))
		if t := sk[i][index]; t.IsZero() {
			return zeroTime
		} else if t.Before(lastSeen) || lastSeen.IsZero() {
			lastSeen = t
		}
	}

	return lastSeen
}

// LookupAndUpdate returns the current timestamp for key and updates it.
func (sk Sketch) LookupAndUpdate(key uint64, now time.Time) time.Time {
	var lastSeen time.Time = time.Unix(1<<62-1, 0)

	for i := 0; i < numBuckets; i++ {
		index := int(key % uint64(len(sk[i])))
		if t := sk[i][index]; t.Before(lastSeen) {
			lastSeen = t
		}
		sk[i][index] = now
	}

	return lastSeen
}

// String returns the string representation of the sketch.
func (sk Sketch) String() string {
	str := ""

	for i := 0; i < numBuckets; i++ {
		str += fmt.Sprintf("Modulus %d: %v\n", len(sk[i]), timeSlice(sk[i]))
	}

	return str
}

// Capacity returns the maximum number of distinct elements
// that can be held by the sketch with ~1% false positives
func (sk Sketch) Capacity() int {
	capacity := 0
	for i := 0; i < numBuckets; i++ {
		capacity += len(sk[i])
	}

	return capacity / multiplier
}

// CountDistinct returns the minimum number of distinct
// elements seen since the given time.
func (sk Sketch) CountDistinct(since time.Time) int {
	uniqueTimestamps := map[time.Time]int{}

	for i := 0; i < numBuckets; i++ {

		// Get the counts for each bucket
		bucketCounts := map[time.Time]int{}

		for _, t := range sk[i] {
			if t.After(since) {
				bucketCounts[t]++
			}
		}

		// uniqueTimestamps[t] = max(uniqueTimestamps[t], bucketCounts[t])
		for t, count := range bucketCounts {
			if count > uniqueTimestamps[t] {
				uniqueTimestamps[t] = count
			}
		}
	}

	total := 0
	for _, count := range uniqueTimestamps {
		total += count
	}

	return total
}

// timeSlice is a slice of time.Time.
type timeSlice []time.Time

// String returns a string representation of []time.Time.
func (ts timeSlice) String() string {
	str := "["

	for _, t := range ts {
		if !t.IsZero() {
			str += fmt.Sprintf(" %v", t)
		}
	}

	return str + " ]"
}

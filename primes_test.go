package lastseen

import (
	"crypto/md5"
	"fmt"
	"testing"
)

func sliceEqual(s1, s2 []uint16) bool {
	return fmt.Sprint(s1) == fmt.Sprint(s2)
}

func TestPrimes(t *testing.T) {
	// first primes
	if !sliceEqual(GetPrimesFrom(0, 1), []uint16{2}) {
		t.Error("fail", GetPrimesFrom(0, 1))
	}
	for i := uint16(0); i <= 2; i++ {
		if !sliceEqual(GetPrimesFrom(i, 3), []uint16{2, 3, 5}) {
			t.Error("fail")
		}
	}
	if !sliceEqual(GetPrimesFrom(5, 3), []uint16{5, 7, 11}) {
		t.Error("fail")
	}

	// qword boundary
	if !sliceEqual(GetPrimesFrom(0, 40), []uint16{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29,
		31, 37, 41, 43, 47, 53, 59, 61, 67, 71,
		73, 79, 83, 89, 97, 101, 103, 107, 109, 113,
		127, 131, 137, 139, 149, 151, 157, 163, 167, 173}) {
		t.Error("fail", GetPrimesFrom(0, 40))
	}

	// last primes
	for i := uint16(65480); i <= 65497; i++ {
		if !sliceEqual(GetPrimesFrom(i, 3), []uint16{65497, 65519, 65521}) {
			t.Error("fail")
		}
	}
	for i := uint16(65498); i <= 65519; i++ {
		if !sliceEqual(GetPrimesFrom(i, 3), []uint16{65519, 65521}) {
			t.Error("fail")
		}
	}
	for i := uint16(65520); i <= 65521; i++ {
		if !sliceEqual(GetPrimesFrom(i, 3), []uint16{65521}) {
			t.Error("fail")
		}
	}
	for i := 65522; i <= 65535; i++ {
		if !sliceEqual(GetPrimesFrom(uint16(i), 3), []uint16{}) {
			t.Error("fail")
		}
	}

	// all
	all1 := GetPrimesFrom(0, 99999)
	if len(all1) != 6542 {
		t.Error("fail", len(all1), "!=", 6542)
	}
	all2 := []uint16{}
	for i := 0; i <= 65535; i++ {
		if prime := GetPrimesFrom(uint16(i), 1); len(prime) == 1 {
			all2 = append(all2, prime[0])
			i = int(prime[0])
		} else if len(prime) != 0 {
			t.Error("fail")
		}
	}
	if !sliceEqual(all1, all2) {
		t.Error("fail")
	}

	// verify content
	hash := md5.Sum([]byte(fmt.Sprint(all1)))
	if fmt.Sprint(hash) != fmt.Sprint([]byte{
		106, 142, 35, 126, 140, 58, 220, 125,
		38, 27, 30, 135, 58, 131, 44, 111}) {
		t.Error("fail")
	}
}

func BenchmarkPrimes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for loops := 0; loops < b.N; loops++ {
		for i := uint16(0); i < 65535; i++ {
			_ = GetPrimesFrom(i, 7)
		}
	}
}

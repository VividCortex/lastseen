package lastseen

import (
	"testing"
	"time"
)

func verifyKnown(t *testing.T, sk Sketch, key uint64, known bool) {
	if lup := sk.Lookup(key); known == lup.IsZero() {
		t.Errorf("key %d expected w/ known=%t", key, known)
	}
}

func updateAndVerify_noCollisions(t *testing.T, sk Sketch, key uint64, isNew bool, ts time.Time) {
	if ls := sk.Lookup(key); isNew && !ls.IsZero() {
		t.Errorf("key %d expected to get Lookup = %s, got %s", key, zeroTime, ls)
	} else if !isNew && ls.IsZero() {
		t.Errorf("key %d expected to get Lookup != %s", key, zeroTime)
	}

	present := sk.Update(key, ts)
	if present == isNew {
		t.Errorf("key %d expected to get present = %t, got %t", key, !isNew, present)
	}

	latest := sk.Lookup(key)
	if latest != ts {
		t.Errorf("key %d expected to get Lookup = %s, got %s", key, ts, latest)
	}
}

func lookupUpdateAndVerify_noCollisions(t *testing.T, sk Sketch, key uint64, isNew bool, ts time.Time) {
	ls := sk.LookupAndUpdate(key, ts)
	present := !ls.IsZero()

	if isNew && present {
		t.Errorf("key %d expected to get Lookup = %s, got %s", key, zeroTime, ls)
	} else if !isNew && !present {
		t.Errorf("key %d expected to get Lookup != %s", key, zeroTime)
	}

	if present == isNew {
		t.Errorf("key %d expected to get present = %t, got %t", key, !isNew, present)
	}

	if latest := sk.Lookup(key); latest != ts {
		t.Errorf("key %d expected to get Lookup = %s, got %s", key, ts, latest)
	}
}

func verifyCountDistinctSince(t *testing.T, sk Sketch, ts time.Time, expectedDistinct int) {
	if distinct := sk.CountDistinct(ts); distinct != expectedDistinct {
		t.Errorf("expected %d distinct entries, got %d", expectedDistinct, distinct)
	}
}

func TestSketch(t *testing.T) {
	sk := NewSketch(1000)

	if c := sk.Capacity(); c != 1011 {
		t.Errorf("expected Capacity %d, got %d", 1011, c)
	}

	for _, updateVerifyFunc := range []func(*testing.T, Sketch, uint64, bool, time.Time){
		updateAndVerify_noCollisions, lookupUpdateAndVerify_noCollisions} {

		sk = NewSketch(2500)

		if c := sk.Capacity(); c != 2521 {
			t.Errorf("expected Capacity %d, got %d", 2521, c)
		}

		now := time.Now()

		verifyCountDistinctSince(t, sk, now.Add(-time.Second), 0)

		updateVerifyFunc(t, sk, 1, true, now)
		updateVerifyFunc(t, sk, 3000, true, now.Add(time.Second))
		updateVerifyFunc(t, sk, 4000, true, now.Add(2*time.Second))
		verifyCountDistinctSince(t, sk, now.Add(-time.Second), 3)

		verifyKnown(t, sk, 1, true)
		verifyKnown(t, sk, 3000, true)
		verifyKnown(t, sk, 4000, true)

		updateVerifyFunc(t, sk, 997, true, now.Add(3*time.Second))
		verifyCountDistinctSince(t, sk, now.Add(-time.Second), 4)

		updateVerifyFunc(t, sk, 0, true, now.Add(4*time.Second))
		verifyCountDistinctSince(t, sk, now.Add(-time.Second), 5)

		verifyCountDistinctSince(t, sk, now.Add(time.Second*2), 2) // 997, 0

		verifyKnown(t, sk, 0, true)
		verifyKnown(t, sk, 997, true)

		// key 3000 already exists, update
		updateVerifyFunc(t, sk, 3000, false, now.Add(5*time.Second))
		verifyCountDistinctSince(t, sk, now.Add(-time.Second), 5)

		verifyCountDistinctSince(t, sk, now.Add(time.Second*2), 3) // now key 3000 is also counted

		verifyKnown(t, sk, 999, false)
		verifyKnown(t, sk, 3000, true)
	}
}

func TestSmallSketch(t *testing.T) {
	sk := NewSketch(1)

	if len(sk) < 1 {
		t.Errorf("expected sketch size of at least 1, got %v", len(sk))
	}
}

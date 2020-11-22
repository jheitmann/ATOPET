package main

import (
	"testing"
)

func TestSplitSecret(t *testing.T) {
	secret := uint64(12)
	splits := SplitIntoShares(3, secret, T)
	res := uint64(0)
	for _, elem := range splits {
		res += elem
	}
	if res != secret {
		t.Errorf("Secret reconstruction failed ! Got %d and wanted %d !", res, secret)
	}
}

func TestCalculsSplitSecret(t *testing.T) {
	s1 := uint64(18)
	s2 := uint64(7)
	s3 := uint64(42)

	expectedResult := uint64(s1 + s2 + s3)
	splitsS1 := SplitIntoShares(3, s1, T)
	splitsS2 := SplitIntoShares(3, s2, T)
	splitsS3 := SplitIntoShares(3, s3, T)
	res := []uint64{uint64(0), uint64(0), uint64(0)}
	for i := 0; i < 3; i++ {
		res[i] = splitsS1[i] + splitsS2[i] + splitsS3[i]
	}
	actualResult := uint64(0)
	for i := 0; i < 3; i++ {
		actualResult += res[i]
	}
	if actualResult != expectedResult {
		t.Errorf("Result of secure computation failed ! Got %d and wanted %d !", actualResult, expectedResult)
	}
}

func TestBeaverGeneration(t *testing.T) {
	a, b, c := beaverCreate(T)
	if a*b != c {
		t.Errorf("a*b != c !")
	}
}

func TestBeaverLogic(t *testing.T) {
	a, b, c := beaverCreate(T)
	y := RandomUInt64(T)
	x := RandomUInt64(T)
	res := c + x*(y-b) + y*(x-a) - (x-a)*(y-b)
	if res != x*y {
		t.Errorf("you idiot")
	}

}

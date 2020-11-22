package main

import (
	"math/big"

	"github.com/ldsec/lattigo/ring"
)

// AddModT computes addition mod T
func AddModT(a, b, T uint64) uint64 {
	res := (a + b) % T
	return res
}

// SubModT computes subtraction mod T
func SubModT(a, b, T uint64) uint64 {
	res := (a + T - b) % T
	return res
}

// MulModT computes multiplication mod T
func MulModT(a, b, T uint64) uint64 {
	res := (a * b) % T
	return res
}

// RandomUInt64 generates a random uint64, bounded by some value T
func RandomUInt64(T uint64) uint64 {
	return ring.RandInt(new(big.Int).SetUint64(T)).Uint64()
}

// SplitIntoShares splits a secret into a number shares with a crytographically secure generator
func SplitIntoShares(numberOfPeers int, secret uint64, T uint64) []uint64 {
	arrayResult := make([]uint64, numberOfPeers)
	indexAt0 := secret
	for i := 1; i < numberOfPeers; i++ {
		generated := RandomUInt64(T)
		arrayResult[i] = generated
		indexAt0 = SubModT(indexAt0, generated, T)
	}
	arrayResult[0] = indexAt0
	return arrayResult
}

func beaverCreate(T uint64) (uint64, uint64, uint64) {
	a := RandomUInt64(T)
	b := RandomUInt64(T)
	c := MulModT(a, b, T)
	return a, b, c
}

// BeaverToShares returns a beaver triplet decomposed into additive shares
func BeaverToShares(nbPeers int, T uint64) ([]uint64, []uint64, []uint64) {
	a, b, c := beaverCreate(T)
	sharesA := SplitIntoShares(nbPeers, a, T)
	sharesB := SplitIntoShares(nbPeers, b, T)
	sharesC := SplitIntoShares(nbPeers, c, T)
	return sharesA, sharesB, sharesC
}

// NewRandomVec generates a random vector of N numbers modulo T
func NewRandomVec(n, T uint64) []uint64 {
	res := []uint64{}
	for i := uint64(0); i < n; i++ {
		res = append(res, RandomUInt64(T))
	}
	return res
}

// AddVec adds two vectors component wise modulo T
func AddVec(a, b []uint64, T uint64) []uint64 {
	res := make([]uint64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = AddModT(a[i], b[i], T)
	}
	return res
}

// SubVec subtracts two vectors component wise modulo T
func SubVec(a, b []uint64, T uint64) []uint64 {
	res := make([]uint64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = SubModT(a[i], b[i], T)
	}
	return res
}

// MulVec multiplies two vectors component wise modulo T
func MulVec(a, b []uint64, T uint64) []uint64 {
	res := make([]uint64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = MulModT(a[i], b[i], T)
	}
	return res
}

// NegVec yields the inverse of all components mod T
func NegVec(a []uint64, T uint64) []uint64 {
	res := make([]uint64, len(a))
	for i := 0; i < len(a); i++ {
		res[i] = (T - a[i])
	}
	return res
}

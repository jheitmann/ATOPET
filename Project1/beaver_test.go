package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ldsec/lattigo/bfv"
)

// TestBeaverExchange is a test to see if two peers can create the beaver triplets
func TestBeaverExchange(t *testing.T) {
	params := bfv.DefaultParams[bfv.PN13QP218]
	N := uint64(3)
	P := make([]*LocalParty, N, N)
	protocol := make([]*BeaverProtocol, N, N)
	testedCircuit := Circuit3

	var err error
	wg := new(sync.WaitGroup)
	for i := range testedCircuit.Peers {
		P[i], err = NewLocalParty(i, testedCircuit.Peers)
		P[i].WaitGroup = wg
		check(err)

		protocol[i] = P[i].NewBeaverProtocol()
	}

	network := GetTestingTCPNetwork(P)
	fmt.Println("parties connected")

	for i, Pi := range protocol {
		Pi.BindNetwork(network[i])
	}

	for _, p := range protocol {
		p.Add(1)
		go p.Run()
	}
	wg.Wait()

	for _, p := range protocol {
		fmt.Println(p, "completed")
	}

	fmt.Println("test completed")
	fmt.Println("Verifying the outputs")

	// Getting the results of the protocol
	a1s := protocol[0].Output[0]
	b1s := protocol[0].Output[1]
	c1s := protocol[0].Output[2]
	a2s := protocol[1].Output[0]
	b2s := protocol[1].Output[1]
	c2s := protocol[1].Output[2]
	a3s := protocol[2].Output[0]
	b3s := protocol[2].Output[1]
	c3s := protocol[2].Output[2]

	// Computing the relation to verify
	as_added := AddVec(a1s, a2s, params.T)
	as_added = AddVec(as_added, a3s, params.T)
	bs_added := AddVec(b1s, b2s, params.T)
	bs_added = AddVec(bs_added, b3s, params.T)
	cs_added := AddVec(c1s, c2s, params.T)
	cs_added = AddVec(cs_added, c3s, params.T)
	// Adding the shares of C
	cs_reconstructed := MulVec(as_added, bs_added, params.T)

	// Verifying that the result is satisfactory
	for i := 0; i < 10; i++ {
		if cs_reconstructed[i] != cs_added[i] {
			t.Errorf("Addition failed ! Got %d and wanted %d !", cs_reconstructed[i], cs_added[i])
		}
	}
	// No bug if reached here
}

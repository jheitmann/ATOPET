package main

import (
	"sync"
	"testing"
)

var Mul1 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:7660",
		1: "localhost:7661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Mult{
			In1: 0,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 6,
		},
	},
	ExpOutput: 115,
}

var Mul2 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:8660",
		1: "localhost:8661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Mult{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Mult{
			In1: 2,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 6,
		},
	},
	ExpOutput: 115,
}

var Mul3 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:9660",
		1: "localhost:9661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Mult{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Mult{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Mult{
			In1: 3,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 6,
		},
	},
	ExpOutput: 115,
}

var Mul4 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Mult{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Mult{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Mult{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Mult{
			In1: 4,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 6,
		},
	},
	ExpOutput: 115,
}

var Mul5 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Mult{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Mult{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Mult{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Mult{
			In1: 4,
			In2: 1,
			Out: 5,
		},
		&Mult{
			In1: 5,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}

var Mul6 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Mult{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Mult{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Mult{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Mult{
			In1: 4,
			In2: 1,
			Out: 5,
		},
		&Mult{
			In1: 5,
			In2: 1,
			Out: 6,
		},
		&Mult{
			In1: 6,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}

var Mul15 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Mult{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Mult{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Mult{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Mult{
			In1: 4,
			In2: 1,
			Out: 5,
		},
		&Mult{
			In1: 5,
			In2: 1,
			Out: 6,
		},
		&Mult{
			In1: 6,
			In2: 1,
			Out: 7,
		},
		&Mult{
			In1: 7,
			In2: 1,
			Out: 8,
		},
		&Mult{
			In1: 8,
			In2: 1,
			Out: 9,
		},
		&Mult{
			In1: 9,
			In2: 1,
			Out: 10,
		},
		&Mult{
			In1: 10,
			In2: 1,
			Out: 11,
		},
		&Mult{
			In1: 11,
			In2: 1,
			Out: 12,
		},
		&Mult{
			In1: 12,
			In2: 1,
			Out: 13,
		},
		&Mult{
			In1: 13,
			In2: 1,
			Out: 14,
		},
		&Mult{
			In1: 14,
			In2: 1,
			Out: 15,
		},
		&Mult{
			In1: 15,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}

var P5 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
		2: "localhost:5662",
		3: "localhost:5663",
		4: "localhost:5664",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
		2: {2: 5},
		3: {3: 7},
		4: {3: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Input{
			Party: 3,
			Out:   3,
		},
		&Input{
			Party: 4,
			Out:   4,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 5,
		},
		&Add{
			In1: 2,
			In2: 3,
			Out: 6,
		},
		&Add{
			In1: 5,
			In2: 6,
			Out: 10,
		},
		&Add{
			In1: 10,
			In2: 4,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}
var P4 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
		2: "localhost:5662",
		3: "localhost:5663",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
		2: {2: 5},
		3: {3: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Input{
			Party: 3,
			Out:   3,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 5,
		},
		&Add{
			In1: 2,
			In2: 3,
			Out: 6,
		},
		&Add{
			In1: 5,
			In2: 6,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}
var P3 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
		2: "localhost:5662",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
		2: {2: 5},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 5,
		},
		&Add{
			In1: 2,
			In2: 5,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}
var P2 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Input{
			Party: 2,
			Out:   2,
		},
		&Input{
			Party: 3,
			Out:   3,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}

var Add1 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:7660",
		1: "localhost:7661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 6,
		},
	},
	ExpOutput: 115,
}

var Add2 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:8660",
		1: "localhost:8661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Add{
			In1: 2,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 6,
		},
	},
	ExpOutput: 19,
}

var Add3 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:9660",
		1: "localhost:9661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Add{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Add{
			In1: 3,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 6,
		},
	},
	ExpOutput: 115,
}

var Add4 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Add{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Add{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Add{
			In1: 4,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 6,
		},
	},
	ExpOutput: 115,
}

var Add5 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Add{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Add{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Add{
			In1: 4,
			In2: 1,
			Out: 5,
		},
		&Add{
			In1: 5,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}

var Add6 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Add{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Add{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Add{
			In1: 4,
			In2: 1,
			Out: 5,
		},
		&Add{
			In1: 5,
			In2: 1,
			Out: 6,
		},
		&Add{
			In1: 6,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}

var Add15 = TestCircuit{
	// f(a,b,c) = (a + b + c) * K
	Peers: map[PartyID]string{
		0: "localhost:5660",
		1: "localhost:5661",
	},
	Inputs: map[PartyID]map[GateID]uint64{
		0: {0: 5},
		1: {1: 7},
	},
	Circuit: []Operation{
		&Input{
			Party: 0,
			Out:   0,
		},
		&Input{
			Party: 1,
			Out:   1,
		},
		&Add{
			In1: 0,
			In2: 1,
			Out: 2,
		},
		&Add{
			In1: 2,
			In2: 1,
			Out: 3,
		},
		&Add{
			In1: 3,
			In2: 1,
			Out: 4,
		},
		&Add{
			In1: 4,
			In2: 1,
			Out: 5,
		},
		&Add{
			In1: 5,
			In2: 1,
			Out: 6,
		},
		&Add{
			In1: 6,
			In2: 1,
			Out: 7,
		},
		&Add{
			In1: 7,
			In2: 1,
			Out: 8,
		},
		&Add{
			In1: 8,
			In2: 1,
			Out: 9,
		},
		&Add{
			In1: 9,
			In2: 1,
			Out: 10,
		},
		&Add{
			In1: 10,
			In2: 1,
			Out: 11,
		},
		&Add{
			In1: 11,
			In2: 1,
			Out: 12,
		},
		&Add{
			In1: 12,
			In2: 1,
			Out: 13,
		},
		&Add{
			In1: 13,
			In2: 1,
			Out: 14,
		},
		&Add{
			In1: 14,
			In2: 1,
			Out: 15,
		},
		&Add{
			In1: 15,
			In2: 1,
			Out: 20,
		},
		&Reveal{
			In:  20,
			Out: 75,
		},
	},
	ExpOutput: 115,
}

func BenchmarkMulCircuit1(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(35), Mul1)
	}
}
func BenchmarkMulCircuit2(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Mul2)
	}
}
func BenchmarkMulCircuit3(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Mul3)
	}
}
func BenchmarkMulCircuit4(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Mul4)
	}
}
func BenchmarkMulCircuit5(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Mul5)
	}
}
func BenchmarkMulCircuit6(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Mul6)
	}
}
func BenchmarkMulCircuit15(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Mul15)
	}
}

func BenchmarkAddCircuit1(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(12), Add1)
	}
}
func BenchmarkAddCircuit2(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Add2)
	}
}
func BenchmarkAddCircuit3(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Add3)
	}
}
func BenchmarkAddCircuit4(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Add4)
	}
}
func BenchmarkAddCircuit5(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Add5)
	}
}
func BenchmarkAddCircuit6(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Add6)
	}
}
func BenchmarkAddCircuit15(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Add15)
	}
}
func BenchmarkPeer2(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), P2)
	}
}
func BenchmarkPeer3(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), Circuit1)
	}
}

func BenchmarkPeer4(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), P4)
	}
}

func BenchmarkPeer5(t *testing.B) {
	for n := 0; n < t.N; n++ {
		CircuitBench(uint64(666), P5)
	}
}

func CircuitBench(expectedOutput uint64, testedCircuit TestCircuit) {
	N := uint64(len(testedCircuit.Peers))
	P := make([]*LocalParty, N, N)
	protocol := make([]*Protocol, N, N)

	inputs, circuit, tripletWires := ParseCircuit(testedCircuit)
	var err error
	wg := new(sync.WaitGroup)
	for i := range testedCircuit.Peers {
		P[i], err = NewLocalParty(i, testedCircuit.Peers)
		P[i].WaitGroup = wg
		check(err)

		protocol[i] = P[i].NewProtocol(inputs[i])
	}

	network := GetTestingTCPNetwork(P)

	for i, Pi := range protocol {
		Pi.BindNetwork(network[i])
	}

	for _, p := range protocol {
		p.Add(1)
		go p.Run(circuit, tripletWires)
	}
	wg.Wait()
}

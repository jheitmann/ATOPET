package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestAddition1(t *testing.T) {
	// Create a simple circuit to evaluate

	// Need to create a map of inputs
	inputs := make(map[PartyID]uint64)
	outputWireToGate := make(map[WireID]Operation)

	// Test circuit to be implemented
	inputs[PartyID(0)] = uint64(2)
	inputs[PartyID(1)] = uint64(3)

	input0 := &Input{Party: 0, Out: 0}
	input1 := &Input{Party: 1, Out: 1}
	add1 := &Add{In1: 0, In2: 1, Out: 3}
	outputWireToGate[0] = input0
	outputWireToGate[1] = input1
	res := add1.Eval(&inputs, &outputWireToGate, &Protocol{}, nil)
	if res != 5 {
		t.Errorf("Addition failed ! Got %d and wanted %d !", res, 5)
	}
}
func TestAddition2(t *testing.T) {
	// Create a simple circuit to evaluate

	// Need to create a map of inputs
	inputs := make(map[PartyID]uint64)
	outputWireToGate := make(map[WireID]Operation)

	// Test circuit to be implemented
	inputs[PartyID(0)] = uint64(2)
	inputs[PartyID(1)] = uint64(3)
	inputs[PartyID(2)] = uint64(7)
	inputs[PartyID(3)] = uint64(9)

	input0 := &Input{Party: 0, Out: 0}
	input1 := &Input{Party: 1, Out: 1}
	input2 := &Input{Party: 2, Out: 2}
	input3 := &Input{Party: 3, Out: 3}
	add1 := &Add{In1: 0, In2: 1, Out: 4}
	add2 := &Add{In1: 2, In2: 3, Out: 5}
	add3 := &Add{In1: 4, In2: 5, Out: 6}
	outputWireToGate[0] = input0
	outputWireToGate[1] = input1
	outputWireToGate[2] = input2
	outputWireToGate[3] = input3
	outputWireToGate[4] = add1
	outputWireToGate[5] = add2
	outputWireToGate[6] = add3
	res := add3.Eval(&inputs, &outputWireToGate, &Protocol{}, nil)
	if res != 21 {
		t.Errorf("Addition failed ! Got %d and wanted %d !", res, 21)
	}
}

func CircuitTest(t *testing.T, testedCircuit TestCircuit) func(t *testing.T) {
	return func(t *testing.T) {
		N := uint64(len(testedCircuit.Peers))
		P := make([]*LocalParty, N, N)
		protocol := make([]*Protocol, N, N)

		inputs, circuit, tripletWires := ParseCircuit(testedCircuit)
		fmt.Println(testedCircuit)
		fmt.Println(inputs)
		var err error
		wg := new(sync.WaitGroup)
		for i := range testedCircuit.Peers {
			P[i], err = NewLocalParty(i, testedCircuit.Peers)
			P[i].WaitGroup = wg
			check(err)

			protocol[i] = P[i].NewProtocol(inputs[i])
		}

		network := GetTestingTCPNetwork(P)
		fmt.Println("parties connected")

		for i, Pi := range protocol {
			Pi.BindNetwork(network[i])
		}

		for _, p := range protocol {
			p.Add(1)
			go p.Run(circuit, tripletWires)
		}
		wg.Wait()

		for _, p := range protocol {
			fmt.Println(p, "completed with output", p.Output)
			if p.Output != testedCircuit.ExpOutput {
				t.Errorf("Wrong output ! Expected %d and got %d !", testedCircuit.ExpOutput, p.Output)
			}
		}

		fmt.Println("test completed")
	}
}

func TestEval(t *testing.T) {
	t.Run("circuit1", CircuitTest(t, Circuit1))
	t.Run("circuit2", CircuitTest(t, Circuit2))
	t.Run("circuit3", CircuitTest(t, Circuit3))
	t.Run("circuit4", CircuitTest(t, Circuit4))
	t.Run("circuit5", CircuitTest(t, Circuit5))
	t.Run("circuit6", CircuitTest(t, Circuit6))
	t.Run("circuit7", CircuitTest(t, Circuit7))
	t.Run("circuit8", CircuitTest(t, Circuit8))
	t.Run("dating", CircuitTest(t, DateCircuit))
}

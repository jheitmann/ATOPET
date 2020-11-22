package main

import "fmt"

type Circuit struct {
	OutputWireToGate *map[WireID]Operation
	Out              *Reveal
}

// ParseCircuit parses a TestCircuit, so that it can be further processed by the MPC protocol
func ParseCircuit(testCircuit TestCircuit) (map[PartyID]uint64, Circuit, []WireID) {
	inputs := make(map[PartyID]uint64)
	tripletWires := []WireID{}

	// Read inputs
	for partyID, gateInput := range testCircuit.Inputs {
		for _, input := range gateInput {
			inputs[partyID] = input
		}
	}

	// Create structure for recursive evaluation
	outputWireToGate := make(map[WireID]Operation)
	out := &Reveal{}
	for _, gate := range testCircuit.Circuit {
		// Identify multiplication gates for beaver triplet generation
		potentialMult, isMult := gate.(*Mult)
		fmt.Println("isMult:", isMult)
		if isMult {
			tripletWires = append(tripletWires, potentialMult.Output())
		}

		// Identify reveal gate for circuit endpoint
		potentialReveal, isReveal := gate.(*Reveal)
		fmt.Println("isReveal:", isReveal)
		if isReveal {
			out = potentialReveal
		}

		wireID := gate.Output()
		outputWireToGate[wireID] = gate
	}

	circuit := Circuit{
		OutputWireToGate: &outputWireToGate,
		Out:              out,
	}

	return inputs, circuit, tripletWires
}

package main

import "testing"

func TestParseCircuit(t *testing.T) {
	var Circuit1 = TestCircuit{
		Peers: map[PartyID]string{
			0: "localhost:6660",
			1: "localhost:6661",
			2: "localhost:6662",
		},
		Inputs: map[PartyID]map[GateID]uint64{
			0: {0: 18},
			1: {1: 7},
			2: {2: 42},
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
				Out: 3,
			},
			&Reveal{
				In:  4,
				Out: 5,
			},
		},
		ExpOutput: 67,
	}

	_, circuit, _ := ParseCircuit(Circuit1)

	if circuit.Out.In != 4 {
		t.Errorf("Reveal parsing failed ! Got %d and wanted %d !", circuit.Out.In, 4)
	}
	if circuit.Out.Out != 5 {
		t.Errorf("Reveal parsing failed ! Got %d and wanted %d !", circuit.Out.In, 5)
	}
}

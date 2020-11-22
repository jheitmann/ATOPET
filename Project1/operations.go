package main

import (
	"fmt"
)

type WireID uint64

type GateID uint64

type Operation interface {
	Output() WireID
	Eval(*map[PartyID]uint64, *map[WireID]Operation, *Protocol, *map[WireID][]uint64) uint64
}

type Input struct {
	Party PartyID
	Out   WireID
}

type Add struct {
	In1 WireID
	In2 WireID
	Out WireID
}

type Sub struct {
	In1 WireID
	In2 WireID
	Out WireID
}

type MultCst struct {
	In       WireID
	CstValue uint64
	Out      WireID
}

type AddCst struct {
	In       WireID
	CstValue uint64
	Out      WireID
}

type Mult struct {
	In1 WireID
	In2 WireID
	Out WireID
}

type Reveal struct {
	In  WireID
	Out WireID
}

// ======= Operation functions =======

func (io Input) Output() WireID {
	return io.Out
}

// Eval for the Input operation will return the secret shares
func (io Input) Eval(inputs *map[PartyID]uint64, outputWireToGate *map[WireID]Operation, prot *Protocol, beavertriplets *map[WireID][]uint64) uint64 {
	return (*inputs)[io.Party]
}

func (ao Add) Output() WireID {
	return ao.Out
}

// Eval for the Add operation will compute the sum of two shared secrets
func (ao Add) Eval(inputs *map[PartyID]uint64, outputWireToGate *map[WireID]Operation, prot *Protocol, beavertriplets *map[WireID][]uint64) uint64 {
	leftside := (*outputWireToGate)[ao.In1].Eval(inputs, outputWireToGate, prot, beavertriplets)
	rightside := (*outputWireToGate)[ao.In2].Eval(inputs, outputWireToGate, prot, beavertriplets)
	return AddModT(leftside, rightside, T)
}

func (so Sub) Output() WireID {
	return so.Out
}

// Eval for the Sub operation will compute the difference between two shared secrets
func (so Sub) Eval(inputs *map[PartyID]uint64, outputWireToGate *map[WireID]Operation, prot *Protocol, beavertriplets *map[WireID][]uint64) uint64 {
	leftside := (*outputWireToGate)[so.In1].Eval(inputs, outputWireToGate, prot, beavertriplets)
	rightside := (*outputWireToGate)[so.In2].Eval(inputs, outputWireToGate, prot, beavertriplets)
	return SubModT(leftside, rightside, T)
}

func (mco MultCst) Output() WireID {
	return mco.Out
}

// Eval for the AddCst operation will multiply the shared secret by a constant
func (mco MultCst) Eval(inputs *map[PartyID]uint64, outputWireToGate *map[WireID]Operation, prot *Protocol, beavertriplets *map[WireID][]uint64) uint64 {
	return MulModT((*outputWireToGate)[mco.In].Eval(inputs, outputWireToGate, prot, beavertriplets), mco.CstValue, T)
}

func (aco AddCst) Output() WireID {
	return aco.Out
}

// Eval for the AddCst operation will add a constant to the shared secret
func (aco AddCst) Eval(inputs *map[PartyID]uint64, outputWireToGate *map[WireID]Operation, prot *Protocol, beavertriplets *map[WireID][]uint64) uint64 {
	if prot.LocalParty.Party.ID == 0 {
		return AddModT((*outputWireToGate)[aco.In].Eval(inputs, outputWireToGate, prot, beavertriplets), aco.CstValue, T)
	} else {
		return (*outputWireToGate)[aco.In].Eval(inputs, outputWireToGate, prot, beavertriplets)

	}
}

func (mo Mult) Output() WireID {
	return mo.Out
}

// Eval for the Mult operation will perform the necessary steps using beaver triplets to multiply two shared secrets
func (mo Mult) Eval(inputs *map[PartyID]uint64, outputWireToGate *map[WireID]Operation, prot *Protocol, beavertriplets *map[WireID][]uint64) uint64 {
	// Computing the inputs to this gate
	leftside := (*outputWireToGate)[mo.In1].Eval(inputs, outputWireToGate, prot, beavertriplets)
	rightside := (*outputWireToGate)[mo.In2].Eval(inputs, outputWireToGate, prot, beavertriplets)

	triplet := (*beavertriplets)[mo.Out]
	a := triplet[0]
	b := triplet[1]
	c := triplet[2]

	// Computing x_i-a and y_i-b
	leftShare := SubModT(leftside, a, T)   // x_i - a_i
	rightShare := SubModT(rightside, b, T) // y_i - b_i

	// Broadcasting x_i - a_i
	for _, peer := range prot.Peers {
		if peer.ID != prot.ID {
			peer.Chan <- Message{prot.ID, leftShare}
			peer.Chan <- Message{prot.ID, rightShare}
		}
	}

	// Initializing x - a and y - b
	leftSecretMinusA := leftShare   // x - a
	rightSecretMinusb := rightShare // y - b

	// Receiving the other shares
	newMessages := prot.SmartMessageReceiver(2)
	for peer := range prot.Peers {
		if peer != prot.ID {
			leftSecretMinusA = AddModT(leftSecretMinusA, newMessages[peer][0].Value, T)
			rightSecretMinusb = AddModT(rightSecretMinusb, newMessages[peer][1].Value, T)
		}
	}

	res := AddModT(AddModT(c, MulModT(leftside, rightSecretMinusb, T), T), MulModT(rightside, leftSecretMinusA, T), T)
	if prot.ID == 0 {
		return SubModT(res, MulModT(leftSecretMinusA, rightSecretMinusb, T), T)
	}
	return res
}

func (ro Reveal) Output() WireID {
	return ro.Out
}

// Eval for the Reveal operation will compute the result of the entire function
func (ro Reveal) Eval(inputs *map[PartyID]uint64, outputWireToGate *map[WireID]Operation, prot *Protocol, beavertriplets *map[WireID][]uint64) uint64 {
	// Evaluating the function
	result := (*outputWireToGate)[ro.In].Eval(inputs, outputWireToGate, prot, beavertriplets)
	// Sending the result to others
	for _, peer := range prot.Peers {
		if peer.ID != prot.ID {
			peer.Chan <- Message{prot.ID, result}
		}
	}

	// Receiving the shares from the others
	newMessages := prot.SmartMessageReceiver(1)
	for peer := range prot.Peers {
		if peer != prot.ID {
			m := newMessages[peer][0]
			fmt.Println(prot, "received message from", m.Party, ":", m.Value)
			result = AddModT(result, m.Value, T)
		}
	}
	close(prot.Chan)
	return result
}

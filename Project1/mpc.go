package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Message struct {
	Party PartyID
	Value uint64
}

type Protocol struct {
	*LocalParty
	Chan  chan Message
	Peers map[PartyID]*Remote

	PendingMessages []Message

	Input  uint64
	Output uint64
}

type Remote struct {
	*RemoteParty
	Chan chan Message
}

func (lp *LocalParty) NewProtocol(input uint64) *Protocol {
	cep := new(Protocol)
	cep.LocalParty = lp
	cep.Chan = make(chan Message, 32)
	cep.Peers = make(map[PartyID]*Remote, len(lp.Peers))
	for i, rp := range lp.Peers {
		cep.Peers[i] = &Remote{
			RemoteParty: rp,
			Chan:        make(chan Message, 32),
		}
	}

	cep.Input = input
	cep.Output = input
	return cep
}

func (cep *Protocol) BindNetwork(nw *TCPNetworkStruct) {
	for partyID, conn := range nw.Conns {

		if partyID == cep.ID {
			continue
		}

		rp := cep.Peers[partyID]

		// Receiving loop from remote
		go func(conn net.Conn, rp *Remote) {
			for {
				var id, val uint64
				var err error
				err = binary.Read(conn, binary.BigEndian, &id)
				check(err)
				err = binary.Read(conn, binary.BigEndian, &val)
				check(err)
				msg := Message{
					Party: PartyID(id),
					Value: val,
				}
				cep.Chan <- msg
			}
		}(conn, rp)

		// Sending loop of remote
		go func(conn net.Conn, rp *Remote) {
			var m Message
			var open = true
			for open {
				m, open = <-rp.Chan
				check(binary.Write(conn, binary.BigEndian, m.Party))
				check(binary.Write(conn, binary.BigEndian, m.Value))
			}
		}(conn, rp)
	}
}

// Run executes the protocol
func (cep *Protocol) Run(circuit Circuit, beaverTripletWires []WireID) {
	multiplicationBeaverTriplets := make(map[WireID][]uint64)
	if len(beaverTripletWires) > 0 {
		a, b, c := cep.generateNBeaverTriplets(len(beaverTripletWires))
		for index, wireID := range beaverTripletWires {
			multiplicationBeaverTriplets[wireID] = []uint64{a[index], b[index], c[index]}
		}
	}

	// Splitting the secret into shares
	shares := SplitIntoShares(len(cep.Peers), cep.Input, T)

	inputs := make(map[PartyID]uint64)

	for index, peer := range cep.Peers {
		if peer.ID != cep.ID {
			fmt.Println(cep, " has sent share :", shares[index])
			peer.Chan <- Message{cep.ID, shares[index]} // Sending the shares independently
		} else {
			inputs[cep.ID] = shares[index]
		}
	}

	newMessages := cep.SmartMessageReceiver(1)
	for peer := range cep.Peers {
		if cep.ID != peer {
			m := newMessages[peer][0]
			fmt.Println(cep, "received message from", m.Party, ":", m.Value)

			inputs[m.Party] = m.Value
		}
	}

	cep.Output = circuit.Out.Eval(&inputs, circuit.OutputWireToGate, cep, &multiplicationBeaverTriplets)

	if cep.WaitGroup != nil {
		cep.WaitGroup.Done()
	}
}

// SmartMessageReceiver deals with concurrency issues by maintaining N-1 queues, N being the number of users, and adding messages to them
func (cep *Protocol) SmartMessageReceiver(nbMessageAwaitedPerPeer int) map[PartyID][]Message {
	// Setting up counting of incoming packets
	numberOfPacketsReceived := make(map[PartyID]int)
	processedMsg := 0

	resultingMessages := make(map[PartyID][]Message)
	for key := range cep.Peers {
		if key != cep.ID {
			numberOfPacketsReceived[key] = 0
			resultingMessages[key] = make([]Message, nbMessageAwaitedPerPeer)
		}
	}

	// Treating messages previously inserted in the queue
	newqueue := make([]Message, 0)
	for _, receivedMSG := range cep.PendingMessages {
		if numberOfPacketsReceived[receivedMSG.Party] < nbMessageAwaitedPerPeer {
			resultingMessages[receivedMSG.Party][numberOfPacketsReceived[receivedMSG.Party]] = receivedMSG
			numberOfPacketsReceived[receivedMSG.Party] = numberOfPacketsReceived[receivedMSG.Party] + 1
			processedMsg++ // Counting the number of messages taken from the queue
		} else {
			fmt.Println("Warning, message arrived before this routine had time to proceed !")
			newqueue = append(newqueue, receivedMSG)
		}
	}
	cep.PendingMessages = newqueue // Saving the new queue
	fmt.Println("Waiting for : ", nbMessageAwaitedPerPeer*(len(cep.Peers)-1)-processedMsg, " messages !")
	// Waiting for new messages
	for i := 0; i < nbMessageAwaitedPerPeer*(len(cep.Peers)-1)-processedMsg; i++ {
		receivedMSG := <-cep.Chan
		if numberOfPacketsReceived[receivedMSG.Party] < nbMessageAwaitedPerPeer {
			resultingMessages[receivedMSG.Party][numberOfPacketsReceived[receivedMSG.Party]] = receivedMSG
			numberOfPacketsReceived[receivedMSG.Party] = numberOfPacketsReceived[receivedMSG.Party] + 1
		} else {
			fmt.Println("Warning, message arrived before this routine had time to proceed !")
			cep.PendingMessages = append(cep.PendingMessages, receivedMSG)
			i = i - 1
		}
	}
	fmt.Println(cep, "finished receiving!")
	return resultingMessages
}

// Initiates a new beaverprotocol to generate the triplets necessary to evaluate the circuit
func (cep *Protocol) generateBeaverTriplets() ([]uint64, []uint64, []uint64) {

	beaverPeers := make(map[PartyID]string)
	for partyID, peer := range cep.Peers {
		address := peer.Addr
		beaverPeers[partyID] = address[:len(address)-4] + "8" + address[len(address)-3:]
	}

	lp, err1 := NewLocalParty(cep.LocalParty.Party.ID, beaverPeers)
	check(err1)
	network, err2 := NewTCPNetwork(lp)
	check(err2)
	err3 := network.Connect(lp)
	check(err3)

	b := lp.NewBeaverProtocol()
	b.BindNetwork(network)

	<-time.After(time.Second) // Leave time for others to connect

	b.Run()
	return b.Output[0], b.Output[1], b.Output[2]
}

// Generates beaver triplets by running the beaver protocol multiple times if necessary
func (cep *Protocol) generateNBeaverTriplets(numberRequired int) ([]uint64, []uint64, []uint64) {
	a, b, c := cep.generateBeaverTriplets()
	if len(a) < numberRequired {
		for len(a) < numberRequired {
			na, nb, nc := cep.generateBeaverTriplets()
			a = append(a, na...)
			b = append(b, nb...)
			c = append(c, nc...)
		}
	}
	return a, b, c
}

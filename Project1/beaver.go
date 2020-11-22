package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/ldsec/lattigo/bfv"
	"github.com/ldsec/lattigo/ring"
)

var T uint64 = bfv.DefaultParams[bfv.PN13QP218].T

type BeaverMessage struct {
	Party PartyID
	Ct    *bfv.Ciphertext
}

type BeaverProtocol struct {
	*LocalParty
	Chan  chan BeaverMessage
	Peers map[PartyID]*BeaverRemote

	PendingMessages []BeaverMessage

	Output [][]uint64
}

type BeaverRemote struct {
	*RemoteParty
	Chan chan BeaverMessage
}

func (lp *LocalParty) NewBeaverProtocol() *BeaverProtocol {
	cep := new(BeaverProtocol)
	cep.LocalParty = lp
	cep.Chan = make(chan BeaverMessage, 32)
	cep.Peers = make(map[PartyID]*BeaverRemote, len(lp.Peers))
	for i, rp := range lp.Peers {
		cep.Peers[i] = &BeaverRemote{
			RemoteParty: rp,
			Chan:        make(chan BeaverMessage, 32),
		}
	}

	cep.Output = make([][]uint64, 3)
	return cep

}

func (cep *BeaverProtocol) BindNetwork(nw *TCPNetworkStruct) {
	for partyID, conn := range nw.Conns {

		if partyID == cep.ID {
			continue
		}

		rp := cep.Peers[partyID]

		// Receiving loop from remote
		go func(conn net.Conn) {
			for {
				var id uint64
				var dataLen int32
				var err error
				err = binary.Read(conn, binary.BigEndian, &id)
				check(err)
				err = binary.Read(conn, binary.BigEndian, &dataLen)
				check(err)
				data := make([]byte, dataLen)
				_, err = io.ReadFull(conn, data)
				check(err)
				ciphertext := new(bfv.Ciphertext)
				// Obtain the ciphertext from an array of bytes
				err = ciphertext.UnmarshalBinary(data)
				check(err)

				msg := BeaverMessage{
					Party: PartyID(id),
					Ct:    ciphertext,
				}

				cep.Chan <- msg
			}
		}(conn)

		// Sending loop of remote
		go func(conn net.Conn, rp *BeaverRemote) {
			var m BeaverMessage
			var open = true
			for open {
				m, open = <-rp.Chan
				id := m.Party
				ct := m.Ct
				// The lattigo API provides a method to serialize the ciphertext into an array of bytes
				marshalledCiphertext, err := ct.MarshalBinary()
				check(err)
				dataLen := int32(len(marshalledCiphertext))

				check(binary.Write(conn, binary.BigEndian, id))
				check(binary.Write(conn, binary.BigEndian, dataLen))
				_, err = conn.Write(marshalledCiphertext)
				check(err)
			}
		}(conn, rp)
	}
}

// Run executes the Beaver's triple generation protocol
func (cep *BeaverProtocol) Run() {
	// BFV parameters (128 bit security)
	params := bfv.DefaultParams[bfv.PN13QP218]

	N := uint64(1 << params.LogN)

	encoder := bfv.NewEncoder(params)

	// Keygen
	kgen := bfv.NewKeyGenerator(params)

	sk, _ := kgen.GenKeyPair()

	decryptor := bfv.NewDecryptor(params, sk)

	encryptorSk := bfv.NewEncryptorFromSk(params, sk)

	evaluator := bfv.NewEvaluator(params)

	// Protocol logic
	aShares := NewRandomVec(N, T)
	bShares := NewRandomVec(N, T)
	cShares := MulVec(aShares, bShares, T)

	// Encryption

	fmt.Println("============================================")
	fmt.Println("Homomorphic computations on batched integers")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Printf("Parameters : N=%d, T=%d, Q = %d bits, sigma = %f \n",
		1<<params.LogN, T, params.LogQP(), params.Sigma)
	fmt.Println()

	aPlaintext := bfv.NewPlaintext(params)
	encoder.EncodeUint(aShares, aPlaintext)

	bPlaintext := bfv.NewPlaintext(params)
	encoder.EncodeUint(bShares, bPlaintext)

	aCiphertext := encryptorSk.EncryptNew(aPlaintext)

	for _, peer := range cep.Peers {
		if peer.ID != cep.ID {
			peer.Chan <- BeaverMessage{cep.ID, aCiphertext} // Sending the shares independently
		}
	}

	newMessages := cep.SmartMessageReceiver(1)
	for _, peer := range cep.Peers {
		if peer.ID != cep.ID {
			m := newMessages[peer.ID][0]
			dj := m.Ct
			rShares := NewRandomVec(N, T)
			cShares = SubVec(cShares, rShares, T)
			rPlaintext := bfv.NewPlaintext(params)
			encoder.EncodeUint(rShares, rPlaintext)

			// Generate error
			dij := evaluator.AddNew(evaluator.MulNew(dj, bPlaintext), rPlaintext)

			context, err := ring.NewContextWithParams(N, params.Qi)
			check(err)
			context.Add(dij.Value()[0], context.SampleGaussianNew(3.2, 19), dij.Value()[0])
			context.Add(dij.Value()[1], context.SampleGaussianNew(3.2, 19), dij.Value()[1])

			peer.Chan <- BeaverMessage{cep.ID, dij} // Sending the shares independently
		}
	}

	newMessages = cep.SmartMessageReceiver(1)
	cPrime := bfv.NewCiphertext(params, 1)
	for peer := range cep.Peers {
		if cep.ID != peer {
			m := newMessages[peer][0]
			dji := m.Ct
			// Add homomorphically
			evaluator.Add(cPrime, dji, cPrime)
		}
	}

	cShares = AddVec(cShares, encoder.DecodeUint(decryptor.DecryptNew(cPrime)), T)
	cep.Output[0] = aShares
	cep.Output[1] = bShares
	cep.Output[2] = cShares

	if cep.WaitGroup != nil {
		cep.WaitGroup.Done()
	}
}

// SmartMessageReceiver deals with concurrency issues by maintaining N-1 queues, N being the number of users, and adding messages to them
func (cep *BeaverProtocol) SmartMessageReceiver(nbMessageAwaitedPerPeer int) map[PartyID][]BeaverMessage {
	// Setting up counting of incoming packets
	numberOfPacketsReceived := make(map[PartyID]int)
	processedMsg := 0

	resultingMessages := make(map[PartyID][]BeaverMessage)
	for key := range cep.Peers {
		if key != cep.ID {
			numberOfPacketsReceived[key] = 0
			resultingMessages[key] = make([]BeaverMessage, nbMessageAwaitedPerPeer)
		}
	}

	// Treating messages previously inserted in the queue
	newqueue := make([]BeaverMessage, 0)
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
	fmt.Println("Finished receiving ! ", resultingMessages)
	return resultingMessages
}

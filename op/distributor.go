package op

// import (
// 	//"fmt"

// 	gomidi "gitlab.com/gomidi/midi"
// 	midi "github.com/plewto/pigiron/midi"
// )


// type Distributor struct {
// 	Operator
// }

// func makeDistributor(name string) *Distributor {
// 	op := new(Distributor)
// 	initOperator(&op.Operator, "Distributor", name. midi.MultiChannel)
// 	return op
// }

// func (op *Distributor) Send(msg gomidi.Message) {
// 	if op.MIDIEnabled() {
// 		raw := msg.Raw()
// 		if isChannelMessage(raw[0]) {
// 			// TODO
// 		} else {
// 			op.distribute(msg)
// 		}
// 	}
// }

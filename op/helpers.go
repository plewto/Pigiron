package op


const channelMask byte = 0xF0
const statusMask byte = 0x0F


func isChannelMessage(status byte) bool {
	var sig byte = 0xf0
	s := status & channelMask
	return s & sig != sig
}

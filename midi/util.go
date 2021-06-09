package midi

import (
	//"fmt"
)

const (
	CHANNEL_MASK (int64) = 0xF0
	STATUS_MASK (int64) = 0x0F
)

var mnemonics = make(map[byte]string)

func init() {
	mnemonics[0x80] = "NOTE-OFF"   
	mnemonics[0x90] = "NOTE-ON"    
	mnemonics[0xA0] = "TOUCH"      
	mnemonics[0xB0] = "CTRL"       
	mnemonics[0xC0] = "PROG"       
	mnemonics[0xD0] = "PRESS"      
	mnemonics[0xE0] = "PITCH"      
	mnemonics[0xF0] = "SYSEX"      
	mnemonics[0xF1] = "TIMECODE"   
	mnemonics[0xF2] = "SONG-POS"        
	mnemonics[0xF3] = "SONG-SEL"    
	mnemonics[0xF6] = "TUNE"       
	mnemonics[0xF8] = "CLOCK"      
	mnemonics[0xFA] = "START"      
	mnemonics[0xFB] = "CONTINUE"   
	mnemonics[0xFC] = "STOP"       
	mnemonics[0xFE] = "ASENS"      
	mnemonics[0xFF] = "RESET"  
}
	

func IsChannelStatus(s int64) bool {
	return s & CHANNEL_MASK != 0xF0
}

func IsSystemStatus(s int64) bool {
	return s >= 0xF0
}

func StatusChannelCommand(s int64) int64 {
	return s & CHANNEL_MASK
}

func StatusChannelIndex(s int64) int64 {
	return s & STATUS_MASK
}

func SetStatusChannelIndex(s int64, ci int64) int64 {
	return StatusChannelCommand(s) | ci
}


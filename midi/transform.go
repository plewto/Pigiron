package midi

import (
	"fmt"
	gomidi "gitlab.com/gomidi/midi/v2"
)


/*
** Transform interface defines general byte transformation fn(n) -> n'
**
*/
type Transform interface {
	TransformRange() (floor byte, ceiling byte)
	Reset()
	Value(index byte) (byte, error)
	SetValue(index byte, value byte) error
	Dump() string
	Plot() string
}

/*
** DataTable implements Transform using 128-slot lookup table.
**
*/ 
type DataTable struct {
	table [128]byte
}

// NewDataTable returns pointer to DataTable
// 
func NewDataTable() *DataTable {
	dt := DataTable{[128]byte{}}
	dt.Reset()
	return &dt
}

// TransformRange() returns table index range 
//
func (dt *DataTable) TransformRange()(floor byte, ceiling byte) {
	return 0, byte(len(dt.table))
}

// Reset() sets table to identity f(x) -> x
//
func (dt *DataTable) Reset() {
	f, c := dt.TransformRange()
	for i := f; i < c; i++ {
		dt.table[i] = byte(i)
	}
}

func (dt *DataTable) validate(n byte, param string) error {
	var err error
	f, c := dt.TransformRange()
	if n < f || c <= n {
		msg := "DataTable %s out of bounds: %d"
		err = fmt.Errorf(msg, param, n)
	}
	return err
}

// Value() returns indexed table value.
// Returns non-nil error if index is out of range.
//
func (dt *DataTable) Value(index byte) (byte, error) {
	err := dt.validate(index, "index")
	if err != nil {
		return 0, err
	}
	return dt.table[index], err
}


// SetValue() sets indexed table value.
// Returns non-nil error if either index or value is out of range.
//
func (dt *DataTable) SetValue(index byte, value byte) error {
	err := dt.validate(index, "index")
	if err != nil {
		return err
	}
	err = dt.validate(value, "value")
	if err != nil {
		return err
	}
	dt.table[index] = value
	return err
}

func (dt *DataTable) Dump() string {
	acc := "\tDataTable:"
	width := byte(8)
	f, c := dt.TransformRange()
	for i := f; i < c; i++ {
		if i % width == 0 {
			acc += fmt.Sprintf("\n\t[%02X] ", i)
		}
		acc += fmt.Sprintf("%02X ", dt.table[i])
	}
	acc += "\n"
	return acc
}

func (dt *DataTable) Plot() string {
	scale := 0.50
	acc := ""
	f, c := dt.TransformRange()
	for i := f; i < c; i++ {
		v, _ := dt.Value(i)
		acc += fmt.Sprintf("[%02x] %02x ", i, v)
		for q := 0; q < int(float64(v) * scale); q++ {
			acc += "-"
		}
		acc += "\n"
	}
	return acc
}

/*
** TransformBank is a ProgramBank of DataTable.
** Implements:
**     ProgramBank
**     ChannelSelector (for single channel)
**     Transform (for currently selected program)
**
*/
type TransformBank struct {
	programs []Transform
	current byte
	channelIndex MIDIChannelNibble
}

// NewTransformBank creates new TransformBank with n slots.
//
func NewTransformBank(n int) *TransformBank {
	p := make([]Transform, n)
	bank := TransformBank{p, 0, 0}
	for i := 0; i < n; i++ {
		p[i] = NewDataTable()
	}
	return &bank
}

// currentTransform returns the currently selected Transform.
//
func (bank *TransformBank) currentTransform() Transform {
	return bank.programs[bank.current]
}

// TransformRange() returns TransformRange of current transform.
//
func (bank *TransformBank) TransformRange() (floor byte, ceiling byte) {
	tr := bank.currentTransform()
	return tr.TransformRange()
}

// Reset() sets current transform to identity.
//
func (bank *TransformBank) Reset() {
	bank.currentTransform().Reset()
}

// Value() returns the indexed value from current transform.
// Returns non-nil error if index is out of bounds.
//
func (bank *TransformBank) Value(index byte) (byte, error) {
	return bank.currentTransform().Value(index)
}

// SetValue() sets indexed position of current transform table.
// Returns non-nil error if either index or value are out of bounds.
//
func (bank *TransformBank) SetValue(index byte, value byte) error {
	return bank.currentTransform().SetValue(index, value)
}


func (bank *TransformBank) ChannelMode() ChannelMode {
	return SingleChannel
}

// EnableChannel selects MIDI channel.
// EnableChannel is required by the ChannelSelector interface,
// The SelectChannel method is more convenient.
//
func (bank *TransformBank) EnableChannel(c MIDIChannel, _ bool) error {
	var err error
	ci := MIDIChannelNibble(c - 1)
	if ci < 0 || ci > 15 {
		msg := "Illegal MIDI channel: %d"
		err = fmt.Errorf(msg, byte(c))
	}
	bank.channelIndex = ci
	return err
}
		
// SelectChannel selects indicated MIDI channel.
//
func (bank *TransformBank) SelectChannel(c MIDIChannel) error {
	return bank.EnableChannel(c, true)
}

// SelectedChannelIndexes returns slice of current MIDI channel.
//
func (bank *TransformBank) SelectedChannelIndexes() []MIDIChannelNibble {
	rs := make([]MIDIChannelNibble, 1)
	rs[0] = bank.channelIndex
	return rs
}

// ChannelIndexSelected returns true if argument is equal to current MIDI channel.
//
func (bank *TransformBank) ChannelIndexSelected(ci MIDIChannelNibble) bool {
	return ci == bank.channelIndex
}

// DeselectAllChannels() is required by ChannelSelector interface, does nothing.
//
func (bank *TransformBank) DeselectAllChannels(){}

// SelectAllChannels() is required by ChannelSelector interface, does nothing.
//
func (bank *TransformBank) SelectAllChannels(){}

// ProgramTransformRange() returns valid program-number range.
// Program numbers outside this range are ignored.
//
func (bank *TransformBank) ProgramRange() (floor byte, ceiling byte) {
	floor, ceiling = byte(0), byte(len(bank.programs))
	return
}

// CurrentProgram() returns the current program-number.
//
func (bank *TransformBank) CurrentProgram() byte {
	return bank.current
}

// ChangeProgram() selects a new current-program.
// Out of range program-numbers are ignored.
//
func (bank *TransformBank) ChangeProgram(msg gomidi.Message) {
	st := StatusByte(msg.Data[0])
	ci := MIDIChannelNibble(st & 0x0F)
	if st == PROGRAM && ci == bank.channelIndex {
		n := byte(msg.Data[1])
		if 0 <= n && n < byte(len(bank.programs)) {
			bank.current = n
		}
	}
}


func (bank *TransformBank) String() string {
	s := "TransformBank(%d), channel: %d, current program: %d"
	return fmt.Sprintf(s, len(bank.programs), int(bank.channelIndex) + 1, bank.current)
}

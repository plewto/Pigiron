package op

import (
	"time"
	"fmt"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/smf"
	gomidi "gitlab.com/gomidi/midi/v2"
)

const (
	PLAYER_START_DELAY = 200 // msec
)

type PlayerState byte
const (
	READY PlayerState = iota
	STOP
	STOPPING
)

type MIDIPlayer struct {
	baseOperator
	midifile *smf.SMF
	noteQueue midi.NoteQueue
	state chan PlayerState
	eventIndex int
	tempo float64
	tempoScale float64
	tickDuration uint64 // usec
	currentTime uint64  // usec
	isPlaying bool
	enableMIDITransport bool
}

func newMIDIPlayer(name string) *MIDIPlayer {
	op := new(MIDIPlayer)
	initOperator(&op.baseOperator, "MIDIPlayer", name, midi.NoChannel)
	op.midifile = smf.NewSMF()
	op.noteQueue = *midi.MakeNoteQueue()
	initTransportHandlers(op)
	op.Reset()
	op.enableMIDITransport = true
	return op
}

func (op *MIDIPlayer) Reset() {
	op.Stop()
	op.noteQueue.Reset()
}
	


func (op *MIDIPlayer) MediaFilename() string {
	return op.midifile.Filename()
}

func (op *MIDIPlayer) LoadMedia(filename string) error {
	mf, err := smf.ReadSMF(filename)
	if err != nil {
		return err
	}
	op.midifile = mf
	return err
}

func (op *MIDIPlayer) Reload() error {
	var err error
	fname := op.MediaFilename()
	if fname == "" {
		errmsg := "Reload faild, No MIDI file specified"
		err = fmt.Errorf(errmsg)
		return err
	}
	err = op.LoadMedia(fname)
	return err
}

func (op *MIDIPlayer) resetControllers() {
	// ISSUE: TODO implement MIDIPlayer.resetControllers()
	fmt.Println("MIDIPlayer.resetControllers() not implemented.")
}


func (op *MIDIPlayer) killActiveNotes() {
	counter := 0
	for _, msg := range op.noteQueue.OffEvents() {
		op.distribute(msg)
		if counter % 16 == 0 {
			time.Sleep(2 * time.Millisecond)
		}
		counter++
	}
	op.noteQueue.Reset()
}


func (op *MIDIPlayer) Stop() {
	fmt.Printf("\nMIDIPlayer %s: STOP\n", op.Name())
	op.state <- STOPPING
	op.isPlaying = false
	time.Sleep(20 * time.Millisecond)
	op.killActiveNotes()
	op.resetControllers()
	op.state <- READY
}

func (op *MIDIPlayer) Continue() error {
	var err error
	if op.MediaFilename() == "" {
		errmsg := "No MIDI file loaded"
		err = fmt.Errorf(errmsg)
		return err
	}
	op.noteQueue.Reset()
	go op.playLoop()
	return err
}

func (op *MIDIPlayer) Play() error {
	var err error
	err = op.Reload()
	if err != nil {
		return err
	}
	op.currentTime = 0
	op.eventIndex = 0
	err = op.Continue()
	return err
}
	
func (op *MIDIPlayer) playLoop() error {
	var err error
	var track smf.Track
	if <- op.state != READY {
		errmsg := "MIDIPlayer %s not ready"
		err = fmt.Errorf(errmsg, op.Name())
		return err
	}
	time.Sleep(PLAYER_START_DELAY * time.Millisecond)
	track, err = op.midifile.Track(0)
	if err != nil {
		return err
	}
	events := track.Events()
	op.eventIndex = 0
	op.isPlaying = true
	for op.eventIndex < len(events) {
		event := events[op.eventIndex]
		delay := time.Duration(op.tickDuration * event.DeltaTime() * uint64(1e6))
		time.Sleep(delay * time.Microsecond)
		d := event.Message().Data
		if len(d) == 0 {
			continue
		}
		st := midi.StatusByte(d[0])
		msg := event.Message()
		switch {
		case midi.IsChannelStatus(st):
			op.distribute(msg)
			op.noteQueue.Update(msg)
		case midi.IsSystemStatus(st):
			op.distribute(msg)
		case midi.IsMetaStatus(st):
			var exitFlag bool
			exitFlag, err = op.handleMeta(msg)
			if exitFlag || err != nil {
				break
			}
		default:
			// ignore
		} 
		select {
		case s := <- op.state:
			if s == STOP || s == STOPPING {
				break
			}
		default:
			// ignore
		} 
		op.currentTime += uint64(delay)
		op.eventIndex++
	}
	op.state <- STOPPING
	op.isPlaying = false
	return err
}

func metaTempoMicroseconds(data []byte) (uint64, error) {
	// 0xff 0x51 0x03 v1 v2 v3
	var err error
	var acc uint64
	if len(data) < 6 || data[1] != 0x51 {
		errmsg := "Malformed meta tempo: %v"
		err = fmt.Errorf(errmsg, data)
		return 0, err
	}
	for i, shift := 3, 16; i < 6; i, shift = i+1, shift-8 {
		acc += uint64(data[i]) << shift
	}
	return acc, err
}

func tempoMicroToBPM(μsec uint64) float64 {
	if μsec == 0 {
		return 60.0
	}
	var k float64 = 60000000
	return k/float64(μsec)
}


func (op *MIDIPlayer) setTickDuration(division int, tempo float64) {
	division = division & 0x7FFF
	if tempo == 0 {
		tempo = 60.0
	}
	var qdur float64 = 60.0/tempo
	dur := qdur/float64(division)
	op.tickDuration = uint64(dur)
}

func (op *MIDIPlayer) handleMeta(msg gomidi.Message) (exitFlag bool, err error) {
	d := msg.Data
	if len(d) < 2 {
		errmsg := "Malformed meta message at event index %d"
		err = fmt.Errorf(errmsg, op.eventIndex)
		return
	}
	mtype := midi.MetaType(msg.Data[1])
	switch {
	case mtype == midi.META_TEMPO:
		var μsec uint64
		μsec, err = metaTempoMicroseconds(d)
		if err != nil {
			return
		}
		op.tempo = tempoMicroToBPM(μsec)
		op.setTickDuration(op.midifile.Division(), op.tempo)
		exitFlag = false
	case midi.IsMetaText(mtype):
		fmt.Printf("META TEXT t %d  %v\n", op.currentTime, d)
		exitFlag = false
	case mtype == midi.META_END_OF_TRACK:
		exitFlag = true
	default:
		// ignore
	}
	return exitFlag, err
}
	
func (op *MIDIPlayer) IsPlaying() bool {
	return op.isPlaying
}

func (op *MIDIPlayer) Duration() uint64 {
	if op.MediaFilename() == "" {
		return 0
	} else {
		return op.midifile.Duration()
	}
}

func (op *MIDIPlayer) Position() uint64 {
	return op.currentTime
}
	
func (op *MIDIPlayer) EnableMIDITransport(flag bool) {
	op.enableMIDITransport = flag
}

func (op *MIDIPlayer) MIDITransportEnabled() bool {
	return op.enableMIDITransport
}

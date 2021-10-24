package op

import (
	"time"
	"fmt"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/pigerr"
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
	PLAYING
)

func (st PlayerState) String() string {
	var s string
	switch st {
	case READY: s = "READY"
	case STOP: s = "STOP"
	case STOPPING: s = "STOPPING"
	case PLAYING: s = "PLAYING"
	default:
		s = "?"
	}
	return s
}
		



type MIDIPlayer struct {
	baseOperator
	midifile *smf.SMF
	noteQueue midi.NoteQueue
	state PlayerState
	eventIndex int
	tempo float64
	tempoScale float64
	tickDuration uint64 // μseconds
	currentTime float64 // seconds
	enableMIDITransport bool
}

func newMIDIPlayer(name string) *MIDIPlayer {
	op := new(MIDIPlayer)
	initOperator(&op.baseOperator, "MIDIPlayer", name, midi.NoChannel)
	op.midifile = smf.NewSMF()
	op.noteQueue = *midi.MakeNoteQueue()
	initTransportHandlers(op)
	op.enableMIDITransport = true
	go op.Reset()
	return op
}

func (op *MIDIPlayer) Reset() {
	op.killActiveNotes()
	op.resetControllers()
	op.state = READY
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
		errmsg := "Reload failed, No MIDI file specified"
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
	op.state = STOPPING
	time.Sleep(20 * time.Millisecond)
	op.killActiveNotes()
	op.resetControllers()
	op.state = READY
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
	op.currentTime = 0.0
	op.eventIndex = 0
	err = op.Continue()
	return err
}
	
func (op *MIDIPlayer) playLoop() error {
	var err error
	var track smf.Track
	if !(op.state == READY) {
		errmsg := "MIDIPlayer %s is not ready, try again in a few seconds."
		err = fmt.Errorf(errmsg, op.Name())
		return err
	}
	op.state = PLAYING
	time.Sleep(PLAYER_START_DELAY * time.Millisecond)
	fmt.Printf("\nMIDIPlayer %s: PLAYING\n", op.Name())
	track, err = op.midifile.Track(0)
	if err != nil {
		return err
	}
	events := track.Events()
	op.eventIndex = 0
	for op.eventIndex < len(events) {
		event := events[op.eventIndex]
		delay := time.Duration(op.tickDuration * event.DeltaTime())
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
		if op.state != PLAYING {
			break
		}
		op.currentTime += float64(delay * 1e6)
		op.eventIndex++
	}
	op.Stop()
	return err
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
	case midi.IsTempoChange(msg):
		op.tempo, err = midi.MetaTempoBPM(msg)
		if err != nil {
			errmsg := "Meta tempo message looks weird, using default 120 BPM"
			pigerr.Warning(errmsg, err.Error())
			op.tempo = 120.0
		}
		div := op.midifile.Division()
		tck := smf.TickDuration(div, op.tempo)
		op.tickDuration = uint64(tck * 1e6)
		exitFlag, err = false, nil
		return
		// ISSUE: TODO Handle text message
	// case midi.IsMetaText(msg):
	// 	fmt.Printf("META TEXT t %f  %v\n", op.currentTime, d)
	// 	exitFlag = false
	case mtype == midi.META_END_OF_TRACK:
		exitFlag = true
	default:
		// ignore
	}
	return exitFlag, err
}

func (op *MIDIPlayer) IsReady() bool {
	return op.state == READY
}


func (op *MIDIPlayer) IsPlaying() bool {
	return op.state == PLAYING
}
	
func (op *MIDIPlayer) Duration() float64 {
	if op.MediaFilename() == "" {
		return 0
	} else {
		return op.midifile.Duration()
	}
}

func (op *MIDIPlayer) Position() float64 {
	return op.currentTime
}
	
func (op *MIDIPlayer) EnableMIDITransport(flag bool) {
	op.enableMIDITransport = flag
}

func (op *MIDIPlayer) MIDITransportEnabled() bool {
	return op.enableMIDITransport
}

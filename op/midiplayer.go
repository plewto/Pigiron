package op

import (
	"fmt"
	"time"
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/pigerr"
	"github.com/plewto/pigiron/pigpath"
)


// Implements Operator and Transport interfaces for MIDI file playback.
//
type MIDIPlayer struct {
	baseOperator
	smf *midi.SMF
	noteQueue *midi.NoteQueue
	isPlaying bool
	eventIndex int
	tempo float64
	tickDuration float64
	currentTime int // msec
	delayStart int  // msec
	enableMIDITransport bool
	
}

func newMIDIPlayer(name string) *MIDIPlayer {
	op := new(MIDIPlayer)
	initOperator(&op.baseOperator, "MIDIPlayer", name, midi.NoChannel)
	op.noteQueue = midi.MakeNoteQueue()
	op.isPlaying = false
	op.eventIndex = 0
	op.tempo = 120.0
	op.tickDuration = midi.TickDuration(24, op.tempo)
	op.currentTime = 0.0  // milliseconds
	op.delayStart = 200   // milliseconds
	op.enableMIDITransport = true
	initTransportHandlers(op)
	return op
}

func (op *MIDIPlayer) Info() string {
	acc := op.commonInfo()
	acc += fmt.Sprintf("\tfilename      : \"%s\"\n", op.MediaFilename())
	acc += fmt.Sprintf("\tIsPlaying     : %v\n", op.IsPlaying())
	acc += fmt.Sprintf("\teventIndex    : %d\n", op.eventIndex)
	acc += fmt.Sprintf("\tduration      : %d msec\n", op.Duration())
	acc += fmt.Sprintf("\tposition      : %d msec\n", op.Position())
	acc += fmt.Sprintf("\tMIDItransport : %v\n", op.enableMIDITransport)
	return acc
}
	
func (op *MIDIPlayer) Reset() {
	op.Stop()
	op.eventIndex = 0
	op.tempo = 120.0
	if op.smf != nil {
		op.tickDuration = midi.TickDuration(op.smf.Division(), op.tempo)
	} else {
		op.tickDuration = midi.TickDuration(24, op.tempo)
	}
	op.currentTime = 0.0
}

func (op *MIDIPlayer) Panic() {
	op.Reset()
}

// op.Stop() see Transport interface.
//
func (op *MIDIPlayer) Stop() {
	fmt.Printf("\nMIDIPlayer %s Stop()\n", op.Name())
	op.isPlaying = false
	for ci := 0; ci < 16; ci++ {
		ch := midi.MIDIChannelNibble(ci)
		ev := midi.MakeControllerEvent(ch, 120, 127)
		op.distribute(ev.PortmidiEvent())
		time.Sleep(1 * time.Millisecond)
		ev = midi.MakeControllerEvent(ch, 120, 0)
		op.distribute(ev.PortmidiEvent())
	}
	op.killActiveNotes()
}

// killActiveNotes() transmits MIDI note-off messages for all unresolved note-on events.
//
func (op *MIDIPlayer) killActiveNotes() {
	fmt.Printf("MIDIPlayer %s killActiveNotes\n", op.Name())
	counter := 0
	for _, pme := range op.noteQueue.OffEvents() {
		op.distribute(*pme)
		if counter % 16 == 0 {
			time.Sleep(2 * time.Millisecond)
		}
		counter++
	}
	op.noteQueue.Reset()
}
	

// op.Play() see Transport interface.
//
func (op *MIDIPlayer) Play() error {
	fmt.Printf("MIDIPlayer %s Play()\n", op.Name())
	var err error
	if op.isPlaying {
		return err
	}
	op.Stop()
	if op.smf == nil {
		errmsg := "No MIDI file loaded"
		err = pigerr.New(errmsg)
		return err
	}
	op.eventIndex = 0
	op.tempo = 120.0
	op.tickDuration = midi.TickDuration(op.smf.Division(), op.tempo)
	op.currentTime = 0.0
	op.isPlaying = true
	go op.playloop()
	return err
}

// op.Continue() see Transport interface.
//
func (op *MIDIPlayer) Continue() error {
	fmt.Printf("MIDIPlayer %s Continue()", op.Name())
	var err error
	if op.isPlaying {
		return err
	}
	if op.smf == nil {
		errmsg := "NO MIDI file loaded"
		return pigerr.New(errmsg)
	}
	op.isPlaying = true
	go op.playloop()
	return err
}

// op.IsPlaying() see Transport interface.
//
func (op *MIDIPlayer) IsPlaying() bool {
	return op.isPlaying
}

func (op *MIDIPlayer) LoadMedia(filename string) error {
	var err error
	var smf *midi.SMF
	filename = pigpath.SubSpecialDirectories(filename)
	smf, err = midi.ReadSMF(filename)
	if err != nil {
		errmsg := "Can not open MIDI file %s"
		err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, filename))
		return err
	}
	op.Stop()
	op.smf = smf
	return err
}

// op.MediaFilename() see Transport interface.
//
func (op *MIDIPlayer) MediaFilename() string {
	if op.smf == nil {
		return ""
	}
	return op.smf.Filename()
}

// op.Duration() see Transport interface.
//
func (op *MIDIPlayer) Duration() int {
	if op.smf == nil {
		return 0.0
	}
	return int(1000 * op.smf.Duration())
}

// op.Position() see Transport interface.
//
func (op *MIDIPlayer) Position() int {
	return op.currentTime
}

func (op *MIDIPlayer) playloop() {
	time.Sleep(time.Duration(op.delayStart) * time.Millisecond)
	fmt.Printf("MIDIPlayer %s playing\n", op.Name())
	var err error
	var track *midi.SMFTrack
	track, err = op.smf.Track(0)
	if err != nil {
		return
	}
	events := track.Events()
	op.noteQueue.Reset()
	for op.isPlaying && op.eventIndex < len(events) {
		event := events[op.eventIndex]
		delay := time.Duration((op.tickDuration * float64(event.DeltaTime())) * 1000) // msec
		time.Sleep(delay * time.Millisecond)
		switch {
		case event.IsChannelEvent():
			pe := event.PortmidiEvent()
			op.distribute(pe)
			op.noteQueue.Update(&pe)
		case event.IsSystemEvent():
			pe := event.PortmidiEvent()
			op.distribute(pe)
		default: // meta
			mtype := byte(event.MetaType())
			switch {
			case mtype == byte(midi.META_TEMPO):
				bpm, _ := event.MetaTempoBPM()
				op.tempo = bpm
				fmt.Printf("tempo = %f\n", op.tempo)
			case midi.IsMetaTextType(mtype):
				bytes, err := event.MetaData()
				if err == nil {
					fs := "time %8d  %s : %s\n"
					fmt.Printf(fs, op.currentTime, midi.MetaType(mtype), string(bytes))
				}
			case mtype == byte(midi.META_END_OF_TRACK):
				op.isPlaying = false
			default:
				// ignore
			}
		}
		op.currentTime += int(delay)
		op.eventIndex++
	}
	op.isPlaying = false
	fmt.Printf("MIDIPlayer %s play-stop\n", op.Name())
	op.killActiveNotes()
	return
}

	
// op.EnableMIDITransport() see Transport interface.
//
func (op *MIDIPlayer) EnableMIDITransport(flag bool) {
	op.enableMIDITransport = flag
}
	
// op.MIDITransportEnabled() see Transport interface.
//
func (op *MIDIPlayer) MIDITransportEnabled() bool {
	return op.enableMIDITransport
}
			
func (op *MIDIPlayer) Send(event portmidi.Event) {
	op.distribute(event)
	if op.enableMIDITransport {
		switch int(event.Status) {
		case 0xFA:
			op.Play()
		case 0xFB:
			op.Continue()
		case 0xFC:
			op.Stop()
		default:
			// ignore
		}
	}
}

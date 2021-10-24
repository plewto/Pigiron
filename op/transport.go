package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
)


// Transport interface defines all media players.
//
// All types which implement Transport should call initTransportHandlers()
// during construction.
//
//   t.Name() The operators name.
//
//   t.Stop() halts playback.
//      osc command  /pig/op <name>, stop
//
//   t.Play() commence playback starting at beginning.
//       Returns non-nil error if no media has been loaded.
//       osc command /pig/op <name>, play
//
//   t.Continue() continues playback from current position.
//       Returns non-nil error if no media has been loaded.
//       osc command /pig/op <name>, continue
//
//   t.IsPlaying() returns true if playback is current in progress.
//       osc command /pig/op <name>, q-is-playing
//       osc returns bool
//
//   t.LoadMedia() loads named media file.
//       Returns non-nil error if file could not be loaded.
//       osc command /pig/op <name>, load, <filename>
//       osc returns error if file could not be loaded.
//
//  t.MediaFilename() returns filename for currently loaded media.
//       Returns empty-string if no media is loaded.
//       osc command /pig/op <name>, q-media-filename
//       osc returns filename
//
//  t.Duration() returns approximate media length in seconds.
//       osc command /pig/op <name>, q-duration
//       osc returns int time in milliseconds.
//
//  t.Position() returns current playback position in seconds.
//       osc command /pig/op <name>, q-position
//       osc returns int time in milliseconds.
//
//  t.EnableMIDITransport() enable/disable MIDI transport control.
//       If enabled the player will stop/start/continue on reception
//       of corresponding MIDI system commands.
//       osc command /pig/op <name>, enable-midi-transport, <bool>
//       osc returns ACK
//
//  t.MIDITransportEnabled() returns true if MIDItransport is enabled.
//      osc command /pig/op <name>, q-midi-transport-enabled
//      osc returns bool
//
type Transport interface {
	Name() string
	Stop()
	Play() error
	Continue() error
	IsPlaying() bool
	LoadMedia(filename string) error
	MediaFilename() string
	Duration() uint64
	Position() uint64
	EnableMIDITransport(flag bool)
	MIDITransportEnabled() bool
	addCommandHandler(command string, handler func(*goosc.Message)([]string, error))
}

// initTransportHandlers() adds transport OSC handlers.
// All type implementing Transport should call initTransportHandlers()
// during construction.
//
func initTransportHandlers(transport Transport) {

	formatResponse := func(command string, values ...string) []string {
		acc := make([]string, 0, len(values) + 1)
		acc = append(acc, fmt.Sprintf("subcommand = %s.%s", transport.Name(), command))
		for _, v := range values {
			acc = append(acc, v)
		}
		return acc
	}

	formatErrorResponse := func(command string) []string {
		acc := []string{fmt.Sprintf("subcommand = %s.%s", transport.Name(), command)}
		return acc
	}
			
	
	// /pig/op <name>, stop
	//
	remoteStop := func(msg *goosc.Message) ([]string, error) {
		var err error
		transport.Stop()
		return  formatResponse("stop"), err
	}

	// /pig/op <name>, play
	//
	remotePlay := func(msg *goosc.Message) ([]string, error) {
		err := transport.Play()
		if err != nil {
			rs := formatErrorResponse("play")
			return rs, err
		}
		rs := formatResponse("play", transport.MediaFilename())
		return rs, err
	}

	// /pig/op <name>, continue
	//
	remoteContinue := func(msg *goosc.Message) ([]string, error) {
		err := transport.Continue()
		if err != nil {
			return formatErrorResponse("continue"), err
		}
		return formatResponse("continue"), err
	}

	// /pig/op <name> load <filename>
	//
	remoteLoad := func(msg *goosc.Message) ([]string, error) {
		value, err := ExpectMsg("oss", msg)
		if err != nil {
			return formatErrorResponse("load"), err
		}
		filename := value[2].S
		err = transport.LoadMedia(filename)
		if err != nil {
			return formatErrorResponse("load"), err
		}
		return formatResponse("load", filename), err
	}

	// op <name>, enable-midi-transport, <bool>
	//
	remoteEnableMIDITransport := func(msg *goosc.Message) ([]string, error) {
		value, err := ExpectMsg("osb", msg)
			if err != nil {
				return formatErrorResponse("enable-midi-transport"), err
			}
		flag := value[2].B
		transport.EnableMIDITransport(flag)
		return formatResponse("enable-midi-transport", fmt.Sprintf("%v", flag)), err
	}

	// op <name> q-midi-transport-enabled  --> bool
	//
	remoteQueryMIDITransport := func(msg *goosc.Message) ([]string, error) {
		var err error
		flag := fmt.Sprintf("%v", transport.MIDITransportEnabled())
		return formatResponse("q-midi-transport-enabled", fmt.Sprintf("%v", flag)), err
	}

	// op <name> q-is-playing  --> bool
	//
	remoteQueryIsPlaying := func(msg *goosc.Message) ([]string, error) {
		var err error
		flag := fmt.Sprintf("%v", transport.IsPlaying())
		return formatResponse("q-is-playing", fmt.Sprintf("%v", flag)), err
	}

	// op <name> q-duration  --> time(sec)
	//
	remoteQueryDuration := func(msg *goosc.Message) ([]string, error) {
		var err error
		dur := fmt.Sprintf("%d", transport.Duration())
		return formatResponse("q-duration", dur), err
	}

	// op <name> q-position  --> time(sec)
	//
	remoteQueryPosition := func(msg *goosc.Message) ([]string, error) {
		var err error
		pos := fmt.Sprintf("%d", transport.Position())
		return formatResponse("q-position", pos), err
	}

	// op <name> q-media-filename  --> filename
	//
	remoteQueryMediaName := func(msg *goosc.Message) ([]string, error) {
		var err error
		name := transport.MediaFilename()
		return formatResponse("q-media-name", name), err
	}

	transport.addCommandHandler("stop", remoteStop)
	transport.addCommandHandler("play", remotePlay)
	transport.addCommandHandler("continue", remoteContinue)
	transport.addCommandHandler("load", remoteLoad)
	transport.addCommandHandler("enable-midi-transport", remoteEnableMIDITransport)
	transport.addCommandHandler("q-midi-transport-enabled", remoteQueryMIDITransport)
	transport.addCommandHandler("q-is-playing", remoteQueryIsPlaying)
	transport.addCommandHandler("q-duration", remoteQueryDuration)
	transport.addCommandHandler("q-position", remoteQueryPosition)
	transport.addCommandHandler("q-media-filename", remoteQueryMediaName)
}

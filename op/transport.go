package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
)


type Transport interface {
	Name() string
	Stop()
	Play() error
	Continue() error
	IsPlaying() bool
	LoadMedia(filename string) error
	MediaFilename() string
	Duration() int // msec
	Position() int // msec
	EnableMIDITransport(flag bool)
	MIDITransportEnabled() bool
	addCommandHandler(command string, handler func(*goosc.Message)([]string, error))
}


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

	// op <name> q-duration  --> time(msec)
	//
	remoteQueryDuration := func(msg *goosc.Message) ([]string, error) {
		var err error
		dur := fmt.Sprintf("%d", transport.Duration())
		return formatResponse("q-duration", dur), err
	}

	// op <name> q-position  --> time(msec)
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

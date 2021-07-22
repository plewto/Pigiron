package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
)


type Transport interface {
	Stop()
	Play() error
	Continue() error
	IsPlaying() bool
	LoadMedia(filename string) error
	MediaFilename() string
	Duration() float64
	Position() float64
	addCommandHandler(command string, handler func(*goosc.Message)([]string, error))
}


func initTransportHandlers(transport Transport) {

	// stop
	// play
	// continue
	// load
	// q-is-playing
	// q-duration
	// q-position

	remoteStop := func(msg *goosc.Message) ([]string, error) {
		var err error
		transport.Stop()
		return empty, err
	}

	remotePlay := func(msg *goosc.Message) ([]string, error) {
		var err error
		transport.Play()
		return empty, err
	}

	remoteContinue := func(msg *goosc.Message) ([]string, error) {
		var err error
		transport.Continue()
		return empty, err
	}

	// op <id>, load, <filename>
	remoteLoad := func(msg *goosc.Message) ([]string, error) {
		value, err := ExpectMsg("oss", msg)
		if err != nil {
			return empty, err
		}
		filename := value[2].S
		err = transport.LoadMedia(filename)
		return []string{filename}, err
	}

	remoteQueryIsPlaying := func(msg *goosc.Message) ([]string, error) {
		var err error
		flag := fmt.Sprintf("%v", transport.IsPlaying())
		return []string{flag}, err
	}

	remoteQueryDuration := func(msg *goosc.Message) ([]string, error) {
		var err error
		dur := fmt.Sprintf("%f", transport.Duration())
		return []string{dur}, err
	}
			
	remoteQueryPosition := func(msg *goosc.Message) ([]string, error) {
		var err error
		pos := fmt.Sprintf("%f", transport.Position())
		return []string{pos}, err
	}

	remoteQueryMediaName := func(msg *goosc.Message) ([]string, error) {
		var err error
		name := transport.MediaFilename()
		return []string{name}, err
	}
	
	transport.addCommandHandler("stop", remoteStop)
	transport.addCommandHandler("play", remotePlay)
	transport.addCommandHandler("continue", remoteContinue)
	transport.addCommandHandler("load", remoteLoad)
	transport.addCommandHandler("q-is-playing", remoteQueryIsPlaying)
	transport.addCommandHandler("q-duration", remoteQueryDuration)
	transport.addCommandHandler("q-position", remoteQueryPosition)
	transport.addCommandHandler("q-media-filename", remoteQueryMediaName)
}

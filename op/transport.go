package op

type Transport interface {
	Stop()
	Play()
	Continue()
	IsPlaying() bool
	LoadMedia(filename string) error
	MediaFilename() string
	// Goto(position time)
	// Length() time
	// Position() time
	// Tempo() float64
}



# Pigiron README

Pigiron is a fully configurable MIDI routing utility written in Go.  It
includes a MIDI file player and has a comprehensive OSC interface. 

The primary Pigiron object is called an **Operator**.  Each Operator has
zero or more MIDI inputs and zero or more MIDI outputs.   When an Operator
receives a MIDI message, it determines if the message should be forwarded
to it's outputs.  An Operator may also modify the message prior to
re-sending it.  

## Operator Types

The following Operators are currently available:


- ChannelFilter - filter events by MIDI channel.
- Distributor - transmit events over several MIDI channels.
- MIDIInput - wrapper for MIDI input device.
- MIDIOutput - wrapper for MIDI output device.
- MIDIPlayer - MIDI file player.
- Monitor - print incoming MIDI messages.
- Transposer - manipulate MIDI data bytes.
- Delay - Repeat notes with optional key mapping.


There are three distinct ways to interact with Pigiron.

1. Remotely via OSC messages.
2. Manually enter commands at terminal prompt.
3. Load a batch file of commands.


The command syntax for these modes are nearly identical.  The only real
difference is that the command 'foo' at the terminal prompt and in a batch
file is entered directly, while as an OSC message it is prefixed with the
application OSC id (by default /pig/).

	foo        command at terminal or in batch file.
	/pig/foo   as OSC message.
	

## Getting help

On startup Pigiron displays a command prompt (by default /pig: ).   For a
list of help topics enter

	/pig:  help topics
	
For a list of commands enter

	/pig: q-commands
	
For details on how OSC messages are handled, enter

	/pig: help OSC
	
In general commands which begin with 'q-' (for query) returns some
information.  

The resources directory contains the batch file 'example.osc'
which sets up a basic MIDI process.  It is heavily annotated.



## Dependencies
    go
	github.com/pelletier/go-toml
	github.com/rakyll/portmidi
    github.com/hypebeast/go-osc


## Installation

**Build Pigiron**

In a terminal cd into the pigiron directory and enter

    [pigiron]$ go build .


**Install**

Either copy the pigiron executable to a location included on $HOME/$PATH, typical
locations would be ~/.local/bin or ~.bin

Alternately you may install with the go command.

    [pigiron]$ go install 


**Configuration Directory**

- Linux   : ~/.config/pigiron/
- Windows : To be determined.
- OSX     : To be determined.

Make 'batch' directrory inside the configuration directory

     ~/.config/pigiron/
	     |
		 +-- config.toml
		 |
		 +-- batch/
		 |
		 +-- resources/
		       |
			   +-- help/
			   +-- testFiles/

   
The location of the configuration directory is printed as Pigiron starts.
You should either copy, or make a symbolic link, within the configuration
directory to the pigiron/resources directory.  The resources directly
contains all the help documents and MIDI files necessary for the test
suit.   Pigiron will function without these but testing will fail and help
will not be available.


## Command Line Options

Pigiron has the following command line options:

--config filename    # Use alternate configuration file.
--batch filename     # Load named batch file.

In general filenames within Pigiron may be prefixed with one of two special
characters.  

  ~/foo names the file foo relative to the user's home directory.
  !/foo names a file relative to the configuratin directory.
  
## GUI?

Pigiron is strictly a terminal based, however due to it's OSC
interface it should be realtivly easy to write a GUI client app.  


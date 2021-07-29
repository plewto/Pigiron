# Pigiron README

Pigiron is a fully configurable MIDI routing utility with integrated MIDI
file player and comprehensive OSC interface.  It provides a series of
*Operators* which may be freely linked to form a *MIDI process tree*.  Each
Operator may have any number of MIDI inputs and outputs.  Currently the
following Operator types are available:

- MIDIInput - wrapper for MIDI input device.
- MIDIOutput - wrapper for MIDI output device.
- ChannelFilter - filter events by MIDI channel.
- Distributor - transmit events over several MIDI channels.
- MIDIPlayer - MIDI file player.
- Monitor - print incoming MIDI messages.


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
	
In general commands which begin with 'q-' (q for query) returns some
information.  

The resources directory contains the batch file 'example.osc'
which sets up a basic MIDI process.   


## Dependencies
	github.com/pelletier/go-toml
	github.com/rakyll/portmidi
    github.com/hypebeast/go-osc


## Installation

**Configuration Directory**
- Linux   : ~/.config/pigiron/
- Windows : To be determined.
- OSX     : o be determined.
   
The configuration directory location is printed on program startup,
immediately under the Pigiron banner.  
  
  
1. Create symbolic link in the configuration directory to the pigiron 
   resources directory.
   
**On Linux**
	`# cd into the configuration directory
	 #
	 $ cd ~/.config/pigiron
	 
	 # Create symbolic link to resources directory.
	 #
	 $ ln -s <pigiron>/resources .
	 #
	 # where <pigiron> is the location of the main pigiron project 
	 # directory.  IE the directory containing main.go
	 
     # The resources directory contains an example configuration file 
	 #       ~/.config/pigiron/resources/config.toml
     # You may wish to either copy or link to this file in the top-level
	 # config directory.
	 #
	 $ cd ~/.config/pigiron
	 $ ln -s ./resources/config.toml .
	 
	 
   





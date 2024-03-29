# Pigiron README

(c) 2021 Steven Jones  

**Pigiron** is a fully configurable MIDI routing utility written in Go.  It
includes a MIDI file player and has a comprehensive OSC interface. 

The primary Pigiron object is called an **Operator**.  Each Operator has
zero or more MIDI inputs and zero or more MIDI outputs.   When an Operator
receives a MIDI message, it determines if the message should be forwarded
to its outputs.  An Operator may also modify the message prior to
re-sending it.  

Version 0.2.0 replaces the original portmidi backend with gomidi.   After
approximately 6 weeks of practical use, it appears to be more stable.  There
are still timing issues with the MIDIPlayer.  It is OK for auditioning MIDI
files, but not quite there for production use. 


## Operator Types

The following Operators are currently available:


- ChannelFilter - filter events by MIDI channel.
- SingleChannelFilter - More efficient channel filter fir single channel filtering.
- Distributor - transmit events over several MIDI channels.
- MIDIInput - wrapper for MIDI input device.
- MIDIOutput - wrapper for MIDI output device.
- MIDIPlayer - MIDI file player.
- Monitor - print incoming MIDI messages.
- Transformer - manipulate MIDI data bytes.


There are three distinct ways to interact with Pigiron.

1. Remotely via OSC messages.
2. Manually enter commands at terminal prompt.
3. Load a batch file of commands.


The command syntax for these modes are nearly identical.  The only
difference is that the command 'foo' at the terminal prompt or in a batch
file is entered directly, while as an OSC message it is prefixed with the
application OSC id (by default /pig/).

        foo        command at terminal or in batch file.
        /pig/foo   as OSC message.
        

## Getting help

On startup Pigiron displays a command prompt (by default /pig: ).   For a
list of help topics enter

**/pig:  help topics**
        
For a list of commands enter  

**/pig: q-commands**
        

For details on how OSC messages are handled, enter

**/pig: help OSC**
        
In general commands which begin with 'q-' (for query) returns some
information.  

The resources/batch directory contains the batch file 'example'
which illustrates several commands to set up a basic MIDI process.
It is heavily annotated.


## Dependencies

    Pigiron has been compiled and tested on Linux.   For compilation on
    Ubuntu 22.04 the following packages were required:

    - build-essential
    - librtmidi-dev
     
    Golang dependencies are:
    
    - go 1.16
    - github.com/pelletier/go-toml
    - github.com/rakyll/portmidi
    - github.com/hypebeast/go-osc

  

## Installation

**Build Pigiron**


In a terminal navigate into the Pigiron directory and enter

    [pigiron]$ go build .


**Install**

Copy the generated pigiron executable to a location included on $HOME/$PATH, typical
locations would be ~/.local/bin or ~.bin


**Configuration Directory**

- Linux   : ~/.config/pigiron/
- Windows : To be determined.
- OSX     : To be determined.

The structure within the .config/pigiron/ directory is:

     ~/.config/
         |
         +--pigiron/
              |
              +-- config.toml
              +-- log
              |
              +-- batch/
              |
              +-- resources/
                    |
                    +-- help/
                    +-- testFiles/

   

Within .config/pigiron/ create the batch directory and either simlink or
copy the resources directory from the Pigiron project.   Pigiron can
operate without the resources directory but help will not be available and
some unit test will fail.  



## Command Line Options

Pigiron has the following command line options:

    --config filename    # Use alternate configuration file.
    --batch filename     # Load named batch file.

In general filenames within Pigiron may be prefixed with one of two special
characters.  

    '~/foo' names the file foo relative to the user's home directory.
    '!/foo' names a file relative to the configuration directory.
  
## GUI?

Pigiron is strictly a terminal based, however due to it's OSC
interface it should be relatively easy to write a GUI client app.  


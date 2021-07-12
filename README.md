# Pigiron Readme


## Requirements
portmidi


## Installation

**Configuration Directory**
- Linux   : ~/.config/pigiron/
- Windows : To be determined.
- OSX     : o be determined.
   
The configuration directory location is printed on program startup,
immediately under the version number.  
  
  
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
	 
	 
   





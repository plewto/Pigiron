/*
** config package establishes global parameters on pigiron start.
**
** 1) Global parameters are first set to default values.
** 2) The command line is then checked to see if an alternate config file has
**    been specified.   
** 3) The configuration file, in toml format, is read to set global
**    values. 
**
** Config file structure and corresponding global parameters
**
**    [log]
**          enable = bool               GlobalParameters.EnableLogging
**          logfile = filename          GlobalParameters.Logfile
**
**    [batch]
**          directory                   Location for batch files.
**
**    Enables logging and specifies location of the log file.
**
**    [osc-server]
**          root = string               GlobalParameters.OSCServerRoot
**          host = ip-address           GlobalParameters.OSCServerHost
**          port = int                  GlobalParameters.OSCServerPort
**
**    [osc-client]
**          root = string               GlobalParameters.OSCClientRoot
**          host = ip-address           GlobalParameters.OSCClientHost
**          port = int                  GlobalParameters.OSCClientPort
**          file = filename             GlobalParameters.OSCClientFilename
**
**     Two sets of OSC values are defined:
**     1) The 'server' represents this instance of pigiron.
**     2) The 'client' sends OSC messages to pigiron.  For each received
**        message pigiron sends a response message back to the client.  The
**        response is optionally sent to a specified client file.
**
**    [tree]
**          max-depth = int             GlobalParameters.MaxTreeDepth
**
**     Sets the maximum depth of the MIDI process tree.  No path through
**     the tree, from root to final output, is allowed to pass through more
**     then max-depth operators.
**
**    [midi-input]
**          buffer-size = int           GlobalParameters.MIDIInputBufferSize
**          poll-interval = int (msec)  GlobalParameters.MIDIInputPollInterval
**
**    [midi-output]
**          buffer-size = int           GlobalParameters.MIDIOutputBufferSize         
**          latency = int               GlobalParameters.MIDIOutputLatency
**
**    buffer-size and latency are parameters passed to the portmidi library.
**    poll-interval is the time, in milliseconds, between polling the MIDI
**    inputs for incoming messages.  A value of 0 checks input as fast as
**    possible.
**
**    [color]
**          banner = color-name         GlobalParameters.BannerColor
**          text = color-name           GlobalParameters.TextColor
**          error = color-name          GlobalParameters.ErrorColor
**
**    The following named colors are defined:
**    red, green, yellow, blue, purple, cyan, gray and white.
**    Colors are not supported on Windows.
**
*/

package config

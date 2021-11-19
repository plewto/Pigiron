/*
** osc package defines communication between Pigiron and external clients via OSC.
**
** There are two general OSC components:
**    server  - the pigiron application.
**    client  - application sending OSC message to pigiron.
**

** The client is represented within Pigiron by the Responder interface.  In
** general for each OSC message received, a response message is sent back
** to the client.  There are two types of response:
**   1) 'ACK', for Acknowledgment, responses indicates the received OSC message
**             was processed without error.
**   2) 'ERROR'
**
** Both the ACK and Error responses include the original message address
** An ACK response may include requested data.
**
*/

package osc

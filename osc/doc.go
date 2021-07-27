// osc package defines communication between Pigiron and external programs via OSC.
//
// There are two general OSC components:
//    server  - the pigiron application.
//    client  - application sending OSC message to pigiron.
//
// The client is represented within Pigiron by the Responder interface.  In
// general for each OSC message received by pigiron, a response message is
// sent back to the client.   There are two types of response.  An ACK for
// (Acknowledgment) responses indicates the received OSC message was processed
// without error.   An Error response indicates the received OSC message
// was malformed or otherwise produced an error. 
//
// Both the ACK and Error responses include the original message
// address. An ACK response may include requested data and an Error
// response includes error messages.


package osc

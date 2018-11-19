// Package nbd implements the NBD network protocol.
//
// You can find a full description of the protocol at
// https://sourceforge.net/p/nbd/code/ci/master/tree/doc/proto.md
//
// This package implements both the client and the server side of the protocol,
// as well as (on Linux) utilities to hook up the kernel NBD client to use a
// server as a block device. The protocol is split into two phases: The
// handshake phase, which allows the client and server to negotiate their
// respective capabilities and what export to use. And the transmission phase,
// for actually reading/writing to the block device.
//
// The client side of the handshake is done with the Client type. Its methods
// can be used to list the exports a server provides and their respective
// capabilities. Its Go method enters transmission phase. The returned Export
// can then be passed to Configure (linux only) to hook it up to an NBD device
// (/dev/nbdX).
//
// The server side combines both handshake and transmission phase into the
// Serve or ListenAndServe functions. The user is expected to implement the
// Device interface to serve actual reads/writes. Under linux, the Loopback
// function serves as a convenient way to use a given Device as a block device.
package nbd

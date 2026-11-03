package main

import (
 "encoding/binary"
 "encoding/hex"
 "errors"
 "strings"
)

// CircuitIdentity holds values returned by a successful G1 XML-RPC login.
// These values are secrets and must never be sent to the browser or logged.
type CircuitIdentity struct {
 AgentID string
 SessionID string
 CircuitCode uint32
}

// uuidWire parses a canonical UUID into the 16-byte LLUDP representation.
func uuidWire(value string) ([16]byte,error) {
 var id [16]byte
 parts:=strings.Split(value,"-")
 if len(parts)!=5 || len(parts[0])!=8 || len(parts[1])!=4 || len(parts[2])!=4 || len(parts[3])!=4 || len(parts[4])!=12 {
  return id,errors.New("invalid UUID format")
 }
 raw,err:=hex.DecodeString(strings.Join(parts,""))
 if err!=nil || len(raw)!=16 {return id,errors.New("invalid UUID value")}
 copy(id[:],raw)
 return id,nil
}

// BuildUseCircuitCode constructs the OpenSim LLUDP low-frequency message 3.
// The caller must obtain the identity from a verified login response and
// choose the sequence number within a per-avatar reliable circuit.
func BuildUseCircuitCode(identity CircuitIdentity,sequence uint32) (LLUDPHeader,error) {
 agent,err:=uuidWire(identity.AgentID)
 if err!=nil {return LLUDPHeader{},err}
 session,err:=uuidWire(identity.SessionID)
 if err!=nil {return LLUDPHeader{},err}
 if identity.CircuitCode==0 {return LLUDPHeader{},errors.New("missing circuit code")}
 payload:=make([]byte,36)
 binary.LittleEndian.PutUint32(payload[:4],identity.CircuitCode)
 copy(payload[4:20],session[:])
 copy(payload[20:36],agent[:])
 return LLUDPHeader{
  Flags:packetReliable,
  Sequence:sequence,
  MessageID:3,
  MessageIDLength:4,
  Payload:payload,
 },nil
}

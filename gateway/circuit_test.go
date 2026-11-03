package main

import (
 "encoding/binary"
 "testing"
)

func TestUseCircuitCodePacket(t *testing.T) {
 identity:=CircuitIdentity{
  AgentID:"00112233-4455-6677-8899-aabbccddeeff",
  SessionID:"ffeeddcc-bbaa-9988-7766-554433221100",
  CircuitCode:0x12345678,
 }
 header,err:=BuildUseCircuitCode(identity,7)
 if err!=nil {t.Fatal(err)}
 packet,err:=EncodeLLUDP(header)
 if err!=nil {t.Fatal(err)}
 decoded,err:=DecodeLLUDP(packet)
 if err!=nil {t.Fatal(err)}
 if decoded.MessageID!=3 || decoded.MessageIDLength!=4 || decoded.Sequence!=7 || decoded.Flags&packetReliable==0 {
  t.Fatalf("unexpected circuit header: %+v",decoded)
 }
 if len(decoded.Payload)!=36 || binary.LittleEndian.Uint32(decoded.Payload[:4])!=identity.CircuitCode {
  t.Fatal("circuit code payload mismatch")
 }
 if decoded.Payload[4]!=0xff || decoded.Payload[20]!=0x00 || decoded.Payload[35]!=0xff {
  t.Fatal("UUID wire encoding mismatch")
 }
}
func TestUseCircuitCodeRejectsInvalidIdentity(t *testing.T) {
 _,err:=BuildUseCircuitCode(CircuitIdentity{AgentID:"invalid",SessionID:"ffeeddcc-bbaa-9988-7766-554433221100",CircuitCode:1},1)
 if err==nil {t.Fatal("invalid agent ID accepted")}
 _,err=BuildUseCircuitCode(CircuitIdentity{AgentID:"00112233-4455-6677-8899-aabbccddeeff",SessionID:"ffeeddcc-bbaa-9988-7766-554433221100"},1)
 if err==nil {t.Fatal("empty circuit code accepted")}
}

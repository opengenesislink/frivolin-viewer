package main

import "testing"

func TestEncodeRoundTrip(t *testing.T) {
 cases:=[]LLUDPHeader{
  {Flags:packetReliable,Sequence:12,MessageID:3,MessageIDLength:1,Payload:[]byte{4,5}},
  {Sequence:14,MessageID:23,MessageIDLength:2,Payload:[]byte{0,1}},
  {Sequence:15,MessageID:512,MessageIDLength:4,Payload:[]byte{8}},
  {Flags:packetAppendedACK,Sequence:16,MessageID:7,MessageIDLength:1,ACKs:[]uint32{10,20}},
 }
 for _,h:=range cases {
  data,err:=EncodeLLUDP(h)
  if err!=nil {t.Fatal(err)}
  got,err:=DecodeLLUDP(data)
  if err!=nil {t.Fatal(err)}
  if got.Sequence!=h.Sequence || got.MessageID!=h.MessageID || len(got.ACKs)!=len(h.ACKs) {
   t.Fatalf("roundtrip failed: %+v",got)
  }
 }
}
func TestEncodeRejectsOversize(t *testing.T) {
 _,err:=EncodeLLUDP(LLUDPHeader{MessageIDLength:1,Payload:make([]byte,1400)})
 if err==nil {t.Fatal("expected MTU limit error")}
}
func TestEncodeRejectsUnimplementedZeroCoding(t *testing.T) {
 _,err:=EncodeLLUDP(LLUDPHeader{Flags:packetZerocoded,MessageIDLength:1})
 if err==nil {t.Fatal("expected zerocode rejection")}
}

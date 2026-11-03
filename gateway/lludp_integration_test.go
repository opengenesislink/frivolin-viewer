package main

import (
 "net"
 "testing"
 "time"
)

// This test runs against a loopback UDP simulator fixture. It does not
// authenticate with or connect to the production grid.
func TestSimulatorDatagramLoopback(t *testing.T) {
 listener,err:=net.ListenUDP("udp",&net.UDPAddr{IP:net.IPv4(127,0,0,1),Port:0})
 if err!=nil {t.Fatal(err)}
 defer listener.Close()
 client,err:=DialSimulatorUDP(listener.LocalAddr().(*net.UDPAddr))
 if err!=nil {t.Fatal(err)}
 defer client.Close()
 sent:=LLUDPHeader{Flags:packetReliable,Sequence:19,MessageID:4,MessageIDLength:1,Payload:[]byte{1,2,3}}
 if err:=client.Send(sent);err!=nil {t.Fatal(err)}
 buf:=make([]byte,1500)
 if err:=listener.SetReadDeadline(time.Now().Add(2*time.Second));err!=nil {t.Fatal(err)}
 n,peer,err:=listener.ReadFromUDP(buf)
 if err!=nil {t.Fatal(err)}
 parsed,err:=DecodeLLUDP(buf[:n])
 if err!=nil {t.Fatal(err)}
 if parsed.Sequence!=19 || parsed.MessageID!=4 {t.Fatalf("bad outbound datagram: %+v",parsed)}
 reply,err:=EncodeLLUDP(LLUDPHeader{Sequence:20,MessageID:5,MessageIDLength:1,Payload:[]byte{8}})
 if err!=nil {t.Fatal(err)}
 if _,err:=listener.WriteToUDP(reply,peer);err!=nil {t.Fatal(err)}
 received,err:=client.Receive()
 if err!=nil {t.Fatal(err)}
 if received.Sequence!=20 || received.MessageID!=5 || len(received.Payload)!=1 || received.Payload[0]!=8 {
  t.Fatalf("bad inbound datagram: %+v",received)
 }
}

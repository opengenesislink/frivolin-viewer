package main

import "testing"

func TestDecodePacket(t *testing.T) {
 p:=[]byte{0x40,0,0,0,7,0,5,9,8}
 h,err:=DecodeLLUDP(p)
 if err!=nil {t.Fatal(err)}
 if h.Sequence!=7 || h.MessageID!=5 || len(h.Payload)!=2 {t.Fatalf("unexpected packet: %+v",h)}
}
func TestDecodeTruncated(t *testing.T) {
 for _,p:=range [][]byte{{},{0,0,0},{0,0,0,0,1,0}} {
  if _,err:=DecodeLLUDP(p);err==nil {t.Fatal("expected parse error")}
 }
}
func TestDecodeZeroRun(t *testing.T) {
 data,err:=decodeZeroCode([]byte{3,0,4,9},100)
 if err!=nil {t.Fatal(err)}
 if len(data)!=6 || data[0]!=3 || data[5]!=9 {t.Fatalf("unexpected data %v",data)}
}

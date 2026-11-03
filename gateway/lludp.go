package main

import (
 "encoding/binary"
 "errors"
 "fmt"
)

// LLUDP packet framing. This layer is intentionally transport-independent.
// It does not yet establish an authenticated circuit or decode region objects.
const (
 packetReliable byte = 0x40
 packetResent byte = 0x20
 packetZerocoded byte = 0x80
 packetAppendedACK byte = 0x10
)
type LLUDPHeader struct {
 Flags byte
 Sequence uint32
 Extra []byte
 MessageID uint32
 MessageIDLength int
 Payload []byte
 ACKs []uint32
}
func DecodeLLUDP(packet []byte) (LLUDPHeader,error) {
 var h LLUDPHeader
 if len(packet)<7 {return h,errors.New("packet too short")}
 h.Flags=packet[0]
 h.Sequence=binary.BigEndian.Uint32(packet[1:5])
 extraLength:=int(packet[5])
 if extraLength>len(packet)-6 {return h,errors.New("invalid extra header length")}
 h.Extra=append([]byte(nil),packet[6:6+extraLength]...)
 data:=packet[6+extraLength:]
 if h.Flags&packetAppendedACK!=0 {
  if len(data)<1 {return h,errors.New("missing ACK count")}
  n:=int(data[len(data)-1])
  if n> (len(data)-1)/4 {return h,errors.New("invalid ACK count")}
  ackStart:=len(data)-1-n*4
  for i:=ackStart;i<ackStart+n*4;i+=4 {h.ACKs=append(h.ACKs,binary.BigEndian.Uint32(data[i:i+4]))}
  data=data[:ackStart]
 }
 if h.Flags&packetZerocoded!=0 {
  expanded,err:=decodeZeroCode(data,1024*1024)
  if err!=nil {return h,err}
  data=expanded
 }
 if len(data)<1 {return h,errors.New("missing message identifier")}
 switch data[0] {
 case 0xff:
  if len(data)<2 {return h,errors.New("truncated medium message identifier")}
  if data[1]==0xff {
   if len(data)<4 {return h,errors.New("truncated low message identifier")}
   h.MessageID=uint32(binary.BigEndian.Uint16(data[2:4]))
   h.MessageIDLength=4
  } else {
   h.MessageID=uint32(data[1])
   h.MessageIDLength=2
  }
 default:
  h.MessageID=uint32(data[0])
  h.MessageIDLength=1
 }
 h.Payload=append([]byte(nil),data[h.MessageIDLength:]...)
 return h,nil
}
func decodeZeroCode(input []byte,limit int) ([]byte,error) {
 out:=make([]byte,0,len(input))
 for i:=0;i<len(input);i++ {
  if input[i]!=0 {
   if len(out)>=limit {return nil,errors.New("zerocode output exceeds limit")}
   out=append(out,input[i]);continue
  }
  i++
  if i>=len(input) {return nil,errors.New("truncated zerocode run")}
  count:=int(input[i])
  if count==0 {return nil,fmt.Errorf("invalid zero run")}
  if count>limit-len(out) {return nil,errors.New("zerocode output exceeds limit")}
  for j:=0;j<count;j++ {out=append(out,0)}
 }
 return out,nil
}

package main

import (
 "encoding/binary"
 "errors"
 "net"
 "time"
)

// EncodeLLUDP frames an outbound simulator message without credentials.
// MessageIDLength selects high (1), medium (2) or low (4) frequency.
func EncodeLLUDP(h LLUDPHeader) ([]byte,error) {
 if len(h.Extra)>255 {return nil,errors.New("extra header too long")}
 if h.Flags&packetZerocoded!=0 {return nil,errors.New("outbound zerocoding not implemented")}
 if len(h.ACKs)>255 {return nil,errors.New("too many ACKs")}
 if h.Flags&packetAppendedACK==0 && len(h.ACKs)>0 {return nil,errors.New("ACK flag missing")}
 id:=make([]byte,0,4)
 switch h.MessageIDLength {
 case 1:
  if h.MessageID>=255 {return nil,errors.New("invalid high message ID")}
  id=append(id,byte(h.MessageID))
 case 2:
  if h.MessageID>254 {return nil,errors.New("invalid medium message ID")}
  id=append(id,0xff,byte(h.MessageID))
 case 4:
  if h.MessageID>65535 {return nil,errors.New("invalid low message ID")}
  id=append(id,0xff,0xff,byte(h.MessageID>>8),byte(h.MessageID))
 default: return nil,errors.New("invalid message ID length")
 }
 length:=6+len(h.Extra)+len(id)+len(h.Payload)
 if h.Flags&packetAppendedACK!=0 {length+=len(h.ACKs)*4+1}
 if length>1400 {return nil,errors.New("outbound datagram exceeds MTU budget")}
 out:=make([]byte,6,length)
 out[0]=h.Flags
 binary.BigEndian.PutUint32(out[1:5],h.Sequence)
 out[5]=byte(len(h.Extra))
 out=append(out,h.Extra...)
 out=append(out,id...)
 out=append(out,h.Payload...)
 if h.Flags&packetAppendedACK!=0 {
  for _,ack:=range h.ACKs {
   var b [4]byte
   binary.BigEndian.PutUint32(b[:],ack)
   out=append(out,b[:]...)
  }
  out=append(out,byte(len(h.ACKs)))
 }
 return out,nil
}

// SimulatorDatagram is a bounded UDP socket for later authenticated circuit
// handling. It deliberately has no automatic connection or login side effects.
type SimulatorDatagram struct {conn *net.UDPConn}
func DialSimulatorUDP(addr *net.UDPAddr) (*SimulatorDatagram,error) {
 if addr==nil || addr.IP==nil || addr.Port<1 || addr.Port>65535 {return nil,errors.New("invalid simulator address")}
 c,err:=net.DialUDP("udp",nil,addr)
 if err!=nil {return nil,err}
 return &SimulatorDatagram{conn:c},nil
}
func (s *SimulatorDatagram) Send(h LLUDPHeader) error {
 data,err:=EncodeLLUDP(h)
 if err!=nil {return err}
 if err=s.conn.SetWriteDeadline(time.Now().Add(3*time.Second));err!=nil {return err}
 _,err=s.conn.Write(data)
 return err
}
func (s *SimulatorDatagram) Receive() (LLUDPHeader,error) {
 var h LLUDPHeader
 if err:=s.conn.SetReadDeadline(time.Now().Add(3*time.Second));err!=nil {return h,err}
 buf:=make([]byte,4096)
 n,err:=s.conn.Read(buf)
 if err!=nil {return h,err}
 return DecodeLLUDP(buf[:n])
}
func (s *SimulatorDatagram) Close() error {return s.conn.Close()}

package main

import (
 "encoding/xml"
 "strings"
 "testing"
)

func TestLoginXMLIsEscaped(t *testing.T) {
 payload:=makeLoginXML(loginRequest{First:"A<&",Last:"Viewer",Password:"<secret>&",Start:"last"})
 if strings.Contains(payload,"<secret>") || !strings.Contains(payload,"&lt;secret&gt;&amp;") {
  t.Fatal("XML escaping failed")
 }
 var v struct { XMLName xml.Name }
 if err:=xml.Unmarshal([]byte(payload),&v);err!=nil {t.Fatal(err)}
}
func TestLoginResponseParsing(t *testing.T) {
 data:=`<?xml version="1.0"?><methodResponse><params><param><value><struct><member><name>login</name><value><string>false</string></value></member><member><name>reason</name><value><string>key</string></value></member></struct></value></param></params></methodResponse>`
 var response xmlResponse
 if err:=xml.Unmarshal([]byte(data),&response);err!=nil {t.Fatal(err)}
 if len(response.Params)!=1 || valueMap(response.Params[0])["reason"]!="key" {t.Fatal("unexpected response")}
}

func TestLoginPresenceResponse(t *testing.T) {
 data:=`<?xml version="1.0"?><methodResponse><params><param><value><struct><member><name>login</name><value><string>false</string></value></member><member><name>reason</name><value><string>presence</string></value></member><member><name>message</name><value><string>Logins are currently restricted.</string></value></member></struct></value></param></params></methodResponse>`
 var response xmlResponse
 if err:=xml.Unmarshal([]byte(data),&response);err!=nil {t.Fatal(err)}
 fields:=valueMap(response.Params[0])
 if fields["login"]!="false" || fields["reason"]!="presence" {t.Fatalf("unexpected response: %v",fields)}
}
func TestLoginFaultResponse(t *testing.T) {
 data:=`<?xml version="1.0"?><methodResponse><fault><value><struct><member><name>faultCode</name><value><i4>-32603</i4></value></member></struct></value></fault></methodResponse>`
 var response xmlResponse
 if err:=xml.Unmarshal([]byte(data),&response);err!=nil {t.Fatal(err)}
 if response.Fault==nil || valueMap(*response.Fault)["faultCode"]!="-32603" {t.Fatal("XML-RPC fault was not decoded")}
}

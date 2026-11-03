package main

import (
 "bytes"
 "encoding/json"
 "encoding/xml"
 "io"
 "net/http"
 "net/url"
 "strings"
 "time"
)

type loginRequest struct {
 First string `json:"first"`
 Last string `json:"last"`
 Password string `json:"password"`
 Start string `json:"start"`
}
type xmlValue struct {
 String string `xml:"string"`
 Int string `xml:"i4"`
 Struct xmlStruct `xml:"struct"`
}
type xmlMember struct {
 Name string `xml:"name"`
 Value xmlValue `xml:"value"`
}
type xmlStruct struct { Members []xmlMember `xml:"member"` }
type xmlResponse struct {
 Fault *xmlValue `xml:"fault>value"`
 Params []xmlValue `xml:"params>param>value"`
}
func escapeXML(s string) string {
 var b bytes.Buffer
 _ = xml.EscapeText(&b, []byte(s))
 return b.String()
}
func member(k,v string) string {
 return "<member><name>"+k+"</name><value><string>"+escapeXML(v)+"</string></value></member>"
}
func makeLoginXML(p loginRequest) string {
 return `<?xml version="1.0"?><methodCall><methodName>login_to_simulator</methodName><params><param><value><struct>`+
 member("first",p.First)+member("last",p.Last)+member("passwd",p.Password)+
 member("start",p.Start)+member("channel","Frivolin Viewer")+member("version","0.2.0")+
 `</struct></value></param></params></methodCall>`
}
func valueMap(v xmlValue) map[string]string {
 m:=make(map[string]string)
 for _, item:=range v.Struct.Members {
  m[item.Name]=item.Value.String
  if m[item.Name]=="" {m[item.Name]=item.Value.Int}
 }
 return m
}
func handleG1Login(w http.ResponseWriter,r *http.Request, endpoint *url.URL) {
 if r.Header.Get("Content-Type")!="application/json" {
  jsonResponse(w,415,map[string]string{"code":"unsupported_media_type"});return
 }
 r.Body=http.MaxBytesReader(w,r.Body,4096)
 defer r.Body.Close()
 var p loginRequest
 dec:=json.NewDecoder(r.Body)
 dec.DisallowUnknownFields()
 if err:=dec.Decode(&p);err!=nil {
  jsonResponse(w,400,map[string]string{"code":"invalid_request"});return
 }
 if len(p.First)<1 || len(p.First)>64 || len(p.Last)<1 || len(p.Last)>64 || len(p.Password)<1 || len(p.Password)>256 {
  jsonResponse(w,400,map[string]string{"code":"invalid_credentials_format"});return
 }
 if p.Start=="" {p.Start="last"}
 if p.Start!="last" && p.Start!="home" {
  jsonResponse(w,400,map[string]string{"code":"invalid_start"});return
 }
 // Never forward credentials over plaintext HTTP.
 if endpoint.Scheme!="https" {
  jsonResponse(w,503,map[string]string{"code":"secure_upstream_required"});return
 }
 body:=strings.NewReader(makeLoginXML(p))
 req,err:=http.NewRequestWithContext(r.Context(),http.MethodPost,endpoint.String(),body)
 if err!=nil {jsonResponse(w,502,map[string]string{"code":"upstream_request_error"});return}
 req.Header.Set("Content-Type","text/xml")
 client:=&http.Client{Timeout:10*time.Second,CheckRedirect:func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}}
 resp,err:=client.Do(req)
 if err!=nil {jsonResponse(w,502,map[string]string{"code":"upstream_unavailable"});return}
 defer resp.Body.Close()
 if resp.StatusCode!=200 {jsonResponse(w,502,map[string]string{"code":"upstream_http_error"});return}
 raw,err:=io.ReadAll(io.LimitReader(resp.Body,1024*1024))
 if err!=nil {jsonResponse(w,502,map[string]string{"code":"upstream_read_error"});return}
 var parsed xmlResponse
 if err:=xml.Unmarshal(raw,&parsed);err!=nil {jsonResponse(w,502,map[string]string{"code":"upstream_xml_error"});return}
 if parsed.Fault!=nil {
  jsonResponse(w,502,map[string]string{"code":"upstream_xmlrpc_fault"});return
 }
 if len(parsed.Params)==0 {jsonResponse(w,502,map[string]string{"code":"upstream_invalid_response"});return}
 fields:=valueMap(parsed.Params[0])
 if fields["login"]!="true" {
  reason:=fields["reason"]
  if reason=="" {reason="unknown"}
  jsonResponse(w,401,map[string]string{"code":"login_rejected","reason":reason,"message":"Anmeldung vom Grid abgelehnt."});return
 }
 // Simulator circuit, CAPS and server-side session store are not implemented.
 // Do not expose seed caps, agent/session IDs or other upstream login payload.
 jsonResponse(w,501,map[string]string{"code":"session_transport_not_implemented","message":"Grid hat die Zugangsdaten akzeptiert; Browser-Sitzung und Simulatorverbindung sind noch nicht implementiert."})
}

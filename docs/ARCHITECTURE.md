# Frivolin Viewer – Architecture (Generation 1 first)

## Product boundaries
- Only Stadt NexVerse G1; G2 adapter will follow **no earlier than 2027**.
- Never expose third-party grid selection.
- The browser client is untrusted: no administrative OpenSim secrets or direct database access.
- The bundled 3D scene is a **local placeholder**. It does not reproduce Magdeburg or another real region.

## Existing G1 system integration — required investigation
1. Verify actual deployed login protocol and authentication method against the G1 server (do not assume an undocumented JSON API).
2. Determine OpenSim viewer login, authentication, UDP circuit, packet handling, CAPS/seed capabilities, EventQueue and region handshake requirements.
3. Implement a controlled server-side protocol adapter with session/circuit lifecycle and per-user isolation.
4. Resolve assets through the authenticated, permission-aware service boundary. Consider cached mesh/texture format conversion for browsers.
5. Add region scene replication, object transforms, appearance, animations, movement, chat, teleports and later inventory/editing.
6. Add NexVoice (LiveKit) using short-lived participant tokens issued server-side.

## Security requirements
TLS for all public endpoints, WSS for realtime; explicit allowed origins; no wildcard credentialed CORS.
HTTP-only Secure SameSite session cookies, CSRF protection on state-changing endpoints, strict rate limits, input validation and strict Content Security Policy.
Do not embed administrator tokens, OpenSim database credentials or internal services in the web client.
No sensitive response or authenticated asset in an offline service-worker cache.
Audit inventory permissions, asset access and economic transactions. No client-authoritative NV$ balances.

## PWA
The initial manifest has a temporary SVG branding icon. For broad installability, replace it with authorized 192x192 and 512x512 PNG icons (and optionally maskable icons). The PWA app shell may work offline, but the virtual world cannot.

## G1 LLUDP progress (2026-11-03)

`gateway/lludp.go` implements a bounded decoder for basic LLUDP packet framing: flags, sequence, extra headers, high/medium/low message identifiers, appended ACKs and zero-coded payload expansion. Unit tests cover initial parsing and malformed input. This is **not** an active simulator connection. There is no UseCircuitCode handshake, reliable retransmission window, message-template decoder, terrain patch decoding, object replication, avatar appearance, CAPS seed/EventQueue, or live world rendering. The Babylon scene remains a local demonstration.

A real login must keep agent/session identifiers and simulator circuit credentials isolated on the server and bind them to a secure browser session. The browser must not receive raw circuit codes or seed capabilities. The G1 login upstream is HTTP, as specified by the grid operator; public HTTP exposes credentials in transit, so this risk must be handled explicitly before production login.

## LLUDP outbound transport milestone

`gateway/lludp_transport.go` adds a bounded LLUDP outbound encoder and an opt-in UDP socket wrapper with read/write deadlines. The encoder supports high/medium/low message identifiers, extra headers and appended ACKs; outbound zero coding is rejected rather than silently emitting malformed packets. Tests cover roundtrips and bounds. **No UDP socket is opened by the running HTTP gateway yet.** The wrapper must be attached only after a verified avatar login and session-scoped simulator endpoint validation. UDP packet framing alone cannot show terrain, objects, or avatars.

Next integration requirements: parse login result fields with their XML-RPC types (UUID, circuit code, simulator IP/port, seed capability), establish authenticated circuit, process reliability/ACK/resent, consume region handshake and terrain/object messages, then implement the browser world-state stream and Babylon scene replication. Never use a client-supplied simulator address.

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

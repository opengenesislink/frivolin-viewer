# G1 Gateway – KeyHelp deployment (Phase 1)

The browser application is currently a local 3D prototype. **There is no authenticated OpenSim login, inventory, region streaming, or avatar session yet.** The login interface is a non-submitting preview.

## Build
```sh
cd web && npm install && npm run build
cd ../gateway && go test ./... && go build -o frivolin-gateway .
```

## Run gateway (localhost only)
```sh
export LISTEN_ADDR=127.0.0.1:8090
export G1_LOGIN_URL=http://hg.stadt-nexverse.de:80/
./frivolin-gateway
```
The G1 endpoint was supplied by the grid operator. Confirm that this path accepts OpenSim XML-RPC login requests before enabling authentication. The current preflight sends HEAD and reports transport reachability only, not successful login or grid compatibility. The gateway-to-G1 connection currently uses plaintext HTTP; prefer a private network or TLS when available.

## Production routing
Configure KeyHelp's managed web server to reverse proxy **only** `/api/` and `/healthz` to `http://127.0.0.1:8090`. Preserve the existing document root serving `web/dist`. Do not expose port 8090 publicly. Configure reverse proxy in the supported KeyHelp mechanism, not by modifying generated vhost files.

The browser defaults to same-origin `/api/v1/capabilities`; `VITE_GATEWAY_URL` may override the base URL at build time. The frontend must be rebuilt after changing it.

## Security
- No credentials are accepted or forwarded in this phase.
- Never expose a raw arbitrary upstream URL as a client-provided parameter.
- TLS is terminated at the HTTPS frontend; secure the backend hop.
- Add CSRF defenses, rate limiting, XML-RPC response validation, secure session handling and UDP/CAPS transport before enabling actual login.

## Verified G1 XML-RPC probe (2026-11-03)

`POST http://hg.stadt-nexverse.de:80/` with method `login_to_simulator` and an empty params array returned HTTP 200, an XML-RPC fault (`-32603`), and `System.Net.IPEndPoint` to `System.Collections.Hashtable` cast exception. This establishes XML-RPC method dispatch, **not** successful authentication. The malformed probe omitted the required XML-RPC struct parameter. Do not infer a broken grid from this fault. Gateway capability `gridConnected` currently means HTTP transport reachable, not logged in.

## Structured login probe (2026-11-03)

A correctly structured XML-RPC `login_to_simulator` request with a dummy username/password returned HTTP 200 and a valid method response: `login=false`, `reason=presence`, `message=Logins are currently restricted. Please try again later.` This confirms XML-RPC request parsing and login method execution, **not** a successful authenticated session. Check G1 login restrictions, presence/login services and server logs before enabling credential submission. The previous `IPEndPoint`/`Hashtable` fault was triggered by an empty parameter list.

## XML-RPC login adapter (phase 2)

The Go gateway now implements XML-RPC request generation and response parsing for `login_to_simulator`. **Login remains disabled by default**. To enable upstream credential forwarding, use HTTPS **or** a privately routed HTTP connection (e.g. WireGuard) with a literal private or loopback IP as `G1_LOGIN_URL`. For private HTTP set both `G1_PRIVATE_TRANSPORT=true` and `G1_ENABLE_LOGIN=true`. The public hostname `http://hg.stadt-nexverse.de:80/` does **not** qualify: do not forward credentials over the public Internet in plaintext. Never disable TLS certificate verification.

Even with HTTPS and login enabled, the adapter does **not** establish a simulator circuit or maintain an authenticated browser session. Successful upstream authentication returns `session_transport_not_implemented`; do not treat it as a playable login. Before any public release implement per-IP rate limits, session storage, CSRF protections, login attempt throttling, and simulator transport. Do not publicly expose this experimental endpoint.

Example only after configuring and verifying the private tunnel: `G1_LOGIN_URL=http://10.78.0.5:80/ G1_PRIVATE_TRANSPORT=true G1_ENABLE_LOGIN=true ./frivolin-gateway` (replace example IP with the actual privately reachable G1 address). The flags alone do not create or secure a tunnel. No live avatar session is available yet.

## Browser login UI (2026-11-03)

The web client now has a real form wired to `POST /api/v1/session/login`. It is disabled unless the gateway reports `loginSupported=true` and the browser is on a secure same-origin connection. Password input is cleared after a response. **This does not mean an avatar session exists.** The Go adapter still returns `session_transport_not_implemented` on successful upstream authentication; the real simulator region is not streamed. The Babylon scene is a placeholder, not Magdeburg. Do not advertise this build as a playable world viewer.

Required for a real world view: store successful login circuit/agent/session information server-side; establish and maintain simulator UDP circuit, CAPS seed/EventQueue; fetch and decode terrain, region objects, textures, mesh and avatar appearance; expose authenticated per-session world state to the browser. No OpenSim UDP-to-WebSocket protocol is implemented yet.

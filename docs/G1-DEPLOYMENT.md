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

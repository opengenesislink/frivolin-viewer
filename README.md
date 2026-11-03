# Frivolin Viewer

**Frivolin Viewer** is the browser-based, installable PWA viewer for **Stadt NexVerse** powered by **OpenGenesisLINK**.

- Target: https://viewer.stadt-nexverse.de
- Primary grid: Generation 1 (`hg.stadt-nexverse.de`, OpenSim-derived)
- Future: Generation 2 (not before 2027)
- No third-party grids or user-selectable grid list

## Status

**Early architectural prototype.** The browser currently displays a local 3D scene, not a live OpenSim region. A genuine account login, region replication, inventory, avatars and in-world interactions are **not implemented** yet. No claim of production readiness is made.

## Quick start

Requires Node.js 22+ and npm.

```bash
cd web
npm install
npm run dev
```

Open the Vite URL shown in the terminal.

```bash
npm run build
```

The output is in `web/dist`. Serve it over HTTPS for PWA installation and service workers. PWA installability varies across browsers.

## Repository layout

- `web/` — TypeScript + Babylon.js 3D client and PWA shell
- `gateway/` — Go HTTP/WebSocket gateway scaffold (no simulated OpenSim authentication)
- `docs/` — Integration plan, protocol boundaries and security

## Gateway

```bash
cd gateway
go run .
```

For development the client can use `VITE_GATEWAY_URL=http://localhost:8090` (in `web/.env.local`). The gateway exposes `/healthz` and `/api/v1/capabilities`, reporting `gridConnected:false`. It does **not** accept or store in-world login credentials.

## Architecture

A browser cannot speak the OpenSim viewer's UDP circuits directly. The long-term G1 adapter must implement/bridge the viewer login, circuit lifecycle, CAPS/EventQueue, region state and asset delivery through audited HTTPS/WebSocket APIs. Credentials and session tokens must not be put in URLs, client storage or logs.

## License

MIT; see [LICENSE](LICENSE). Third-party packages remain under their respective licenses. This repository does not bundle Firestorm/OpenSim code.

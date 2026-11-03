import './style.css';
import { ArcRotateCamera, Color3, Color4, DirectionalLight, Engine, HemisphericLight, MeshBuilder, Scene, StandardMaterial, Vector3 } from '@babylonjs/core';

const root = document.querySelector<HTMLDivElement>('#app');
if (!root) throw new Error('Missing app element');
root.innerHTML = `
  <main class="shell">
    <canvas id="renderCanvas" aria-label="Frivolin 3D-Demowelt"></canvas>
    <header class="top">
      <section class="panel"><div class="brand"><span>F</span>rivolin Viewer</div><div class="sub">Stadt NexVerse · Generation 1</div></section>
      <section class="panel right"><div class="pill">G1-Adapter · Vorschau</div><div>OpenGenesisLINK WebViewer</div><div id="gateway-status" class="status">Gateway wird geprüft …</div></section>
    </header>
    <aside class="login panel" aria-label="G1 Anmeldung"><h2>Stadt NexVerse</h2><div class="sub">Generation 1 · Login-Vorschau</div><label>Vorname<input autocomplete="off" disabled placeholder="Avatar-Vorname" /></label><label>Nachname<input autocomplete="off" disabled placeholder="Avatar-Nachname" /></label><label>Passwort<input type="password" disabled placeholder="Noch nicht verfügbar" /></label><button class="btn" disabled>Anmelden (in Entwicklung)</button><p class="foot">Die Anmeldung ist absichtlich deaktiviert. Es werden keine Zugangsdaten übertragen.</p></aside>\n    <footer class="bottom">
      <div class="panel hint">Lokale 3D-Demoszene · Maus/Touch: Kamera drehen · Mausrad: Zoom<div class="foot">G1-Anmeldung und Live-Regionen sind noch nicht verfügbar.</div></div>
      <button id="install" class="btn" hidden>App installieren</button>
    </footer>
  </main>
`;

const canvas = document.querySelector<HTMLCanvasElement>('#renderCanvas')!;
const engine = new Engine(canvas, true, { preserveDrawingBuffer: false, stencil: true });
const scene = new Scene(engine);
scene.clearColor = new Color4(0.05, 0.07, 0.14, 1);
const camera = new ArcRotateCamera('camera', -Math.PI / 3, Math.PI / 3.1, 32, new Vector3(0, 1, 0), scene);
camera.lowerRadiusLimit = 8;
camera.upperRadiusLimit = 70;
camera.attachControl(canvas, true);
new HemisphericLight('skyLight', new Vector3(0, 1, 0), scene).intensity = 0.85;
const sun = new DirectionalLight('sun', new Vector3(-1, -2, -1), scene);
sun.intensity = 0.65;
const ground = MeshBuilder.CreateGround('demoGround', { width: 36, height: 36 }, scene);
const grass = new StandardMaterial('grass', scene);
grass.diffuseColor = new Color3(0.16, 0.32, 0.34);
ground.material = grass;
for (let x = -12; x <= 12; x += 6) {
  for (let z = -12; z <= 12; z += 6) {
    const box = MeshBuilder.CreateBox(`placeholder-${x}-${z}`, { size: 2, height: 1 + (Math.abs(x * z) % 5) }, scene);
    box.position.set(x, box.scaling.y + 1, z);
    const mat = new StandardMaterial(`mat-${x}-${z}`, scene);
    mat.diffuseColor = new Color3(0.35 + (x + 12) / 100, 0.20 + (z + 12) / 110, 0.48);
    box.material = mat;
  }
}
const monument = MeshBuilder.CreateCylinder('FrivolinDemoMarker', { height: 5, diameter: 2.5, tessellation: 6 }, scene);
monument.position.y = 2.5;
const markerMat = new StandardMaterial('marker', scene);
markerMat.emissiveColor = new Color3(0.8, 0.22, 0.6);
monument.material = markerMat;
engine.runRenderLoop(() => scene.render());
window.addEventListener('resize', () => engine.resize());

const gatewayStatus = document.querySelector<HTMLDivElement>('#gateway-status')!;
const gateway = (import.meta.env.VITE_GATEWAY_URL as string | undefined) || window.location.origin;
if (!gateway) {
  gatewayStatus.textContent = 'Gateway nicht konfiguriert · Demo-Modus';
} else {
  const ctrl = new AbortController();
  const timeout = window.setTimeout(() => ctrl.abort(), 3000);
  fetch(new URL('/api/v1/capabilities', new URL(gateway, window.location.origin)), { signal: ctrl.signal })
    .then(async r => {
      if (!r.ok) throw new Error(String(r.status));
      return r.json() as Promise<{ gridConnected?: boolean; loginSupported?: boolean; loginEndpointConfigured?: boolean }>;
    })
    .then(info => { gatewayStatus.textContent = info.gridConnected ? (info.loginSupported ? 'G1-Endpunkt erreichbar · Login-Adapter bereit' : 'G1-Endpunkt erreichbar · Login noch gesperrt') : (info.loginEndpointConfigured ? 'Gateway online · G1-Endpunkt nicht erreichbar' : 'Gateway online · G1-Endpunkt nicht konfiguriert'); })
    .catch(() => { gatewayStatus.textContent = 'Gateway nicht erreichbar · Demo-Modus'; })
    .finally(() => window.clearTimeout(timeout));
}

type InstallPrompt = Event & { prompt: () => Promise<void>; userChoice: Promise<{ outcome: string }> };
let deferredPrompt: InstallPrompt | null = null;
const installButton = document.querySelector<HTMLButtonElement>('#install')!;
window.addEventListener('beforeinstallprompt', (event: Event) => {
  event.preventDefault();
  deferredPrompt = event as InstallPrompt;
  installButton.hidden = false;
});
installButton.addEventListener('click', async () => {
  if (!deferredPrompt) return;
  await deferredPrompt.prompt();
  await deferredPrompt.userChoice;
  deferredPrompt = null;
  installButton.hidden = true;
});

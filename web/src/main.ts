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
    <aside class="login panel" aria-label="G1 Anmeldung"><h2>Stadt NexVerse</h2><div class="sub">Generation 1 · Login</div><form id="login-form"><label>Vorname<input id="login-first" required maxlength="64" disabled autocomplete="given-name" /></label><label>Nachname<input id="login-last" required maxlength="64" disabled autocomplete="family-name" /></label><label>Passwort<input id="login-password" type="password" required maxlength="256" disabled autocomplete="current-password" /></label><button id="login-submit" class="btn" type="submit" disabled>Anmelden</button></form><p id="login-message" class="foot" role="status">Gateway wird geprüft. Login noch nicht verfügbar.</p></aside>
    <footer class="bottom">
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
const loginForm = document.querySelector<HTMLFormElement>('#login-form')!;
const loginMessage = document.querySelector<HTMLParagraphElement>('#login-message')!;
const loginButton = document.querySelector<HTMLButtonElement>('#login-submit')!;
const fields = ['login-first','login-last','login-password'].map(id => document.getElementById(id) as HTMLInputElement);
let loginReady = false;
loginForm.addEventListener('submit', async event => {
  event.preventDefault();
  if (!loginReady) return;
  loginButton.disabled = true;
  loginMessage.textContent = 'Anmeldung wird geprüft …';
  try {
    const response = await fetch(new URL('/api/v1/session/login', new URL(gateway, window.location.origin)), {
      method: 'POST', credentials: 'same-origin', cache: 'no-store',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({first: fields[0].value.trim(), last: fields[1].value.trim(), password: fields[2].value, start: 'last'})
    });
    fields[2].value = '';
    const data = await response.json() as { code?: string; reason?: string; message?: string };
    loginMessage.textContent = data.code === 'login_rejected'
      ? (data.reason === 'key' ? 'Avatarname oder Passwort ungültig.' : 'Anmeldung abgelehnt: ' + (data.reason || 'unbekannt'))
      : (data.message || 'Loginstatus: ' + response.status);
  } catch {
    fields[2].value = '';
    loginMessage.textContent = 'Login-Gateway nicht erreichbar.';
  } finally {
    loginButton.disabled = !loginReady;
  }
});
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
    .then(info => {
      loginReady = Boolean(info.loginSupported && window.isSecureContext && new URL(gateway, window.location.origin).origin === window.location.origin);
      fields.forEach(field => { field.disabled = !loginReady; });
      loginButton.disabled = !loginReady;
      loginMessage.textContent = loginReady ? 'Login-Adapter bereit; Simulatorverbindung noch nicht implementiert.' : 'Login nicht verfügbar; keine Zugangsdaten eingeben.';
      gatewayStatus.textContent = info.gridConnected ? (loginReady ? 'G1 erreichbar · Login-Adapter bereit' : 'G1 erreichbar · Login gesperrt') : 'G1 derzeit nicht erreichbar';
    })
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

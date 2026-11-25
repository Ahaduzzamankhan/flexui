/**
 * Minimal static file server used to preview the example sites.
 *
 * Node APIs only, so it runs under Node or Bun:
 *   bun scripts/serve.ts examples/basic 4173
 *   node scripts/serve.ts examples/basic 4173
 */

import { createReadStream } from "node:fs";
import { stat } from "node:fs/promises";
import { createServer, type IncomingMessage, type ServerResponse } from "node:http";
import { extname, join, normalize, resolve, sep } from "node:path";

const MIME: Readonly<Record<string, string>> = {
  ".html": "text/html; charset=utf-8",
  ".css": "text/css; charset=utf-8",
  ".js": "text/javascript; charset=utf-8",
  ".mjs": "text/javascript; charset=utf-8",
  ".json": "application/json; charset=utf-8",
  ".svg": "image/svg+xml",
  ".png": "image/png",
  ".jpg": "image/jpeg",
  ".jpeg": "image/jpeg",
  ".webp": "image/webp",
  ".ico": "image/x-icon",
  ".woff": "font/woff",
  ".woff2": "font/woff2",
  ".txt": "text/plain; charset=utf-8",
  ".map": "application/json; charset=utf-8",
};

const rootArgument = process.argv[2] ?? ".";
const portArgument = process.argv[3] ?? process.env["PORT"] ?? "4173";
const root = resolve(rootArgument);
const port = Number.parseInt(portArgument, 10);
const host = "127.0.0.1";

function send(response: ServerResponse, status: number, body: string, type = "text/plain; charset=utf-8"): void {
  response.writeHead(status, { "content-type": type, "cache-control": "no-store" });
  response.end(body);
}

async function handle(request: IncomingMessage, response: ServerResponse): Promise<void> {
  const url = new URL(request.url ?? "/", `http://${host}:${port}`);
  let pathname = decodeURIComponent(url.pathname);
  if (pathname.endsWith("/")) {
    pathname += "index.html";
  }
  const target = normalize(join(root, pathname));
  if (target !== root && !target.startsWith(root + sep)) {
    send(response, 403, "Forbidden");
    return;
  }
  let info;
  try {
    info = await stat(target);
  } catch {
    send(response, 404, `Not found: ${pathname}`);
    return;
  }
  if (info.isDirectory()) {
    response.writeHead(302, { location: `${pathname}/` });
    response.end();
    return;
  }
  const type = MIME[extname(target).toLowerCase()] ?? "application/octet-stream";
  response.writeHead(200, {
    "content-type": type,
    "content-length": String(info.size),
    "cache-control": "no-store",
  });
  createReadStream(target).pipe(response);
}

createServer((request, response) => {
  void handle(request, response).catch((error: unknown) => {
    send(response, 500, error instanceof Error ? error.message : "Internal error");
  });
}).listen(port, host, () => {
  console.log(`serving ${root} at http://${host}:${port}/ (pid ${process.pid})`);
});

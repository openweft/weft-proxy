<p align="center"><img src="https://raw.githubusercontent.com/openweft/brand/main/social/openweft.png" alt="openweft" width="720"></p>

# weft-proxy

`weft-proxy` is the Caddy build weft-agent supervises when an operator
opts into the reverse-proxy plane (`weft agent --proxy`).

It is a stock Caddy binary plus the modules weft-agent's
[`agent/proxy`](https://github.com/openweft/weft/tree/main/agent/proxy)
subsystem references at runtime — today, that's the etcd3 storage
adapter so certificates are shared across all hosts in a 3-DC cluster
instead of being re-minted by every node on first reload (which would
hammer Let's Encrypt's rate limit).

## Why a separate binary

[`agent/proxy`](https://github.com/openweft/weft/tree/main/agent/proxy)
deliberately doesn't import Caddy as a Go library — see
`agent/proxy/doc.go` for the reasoning (vendor weight, crash
isolation, the admin-API supervisor pattern Caddy was designed for).
weft-agent launches a Caddy subprocess and talks to it over the
admin socket; this repo produces the subprocess.

Stock Caddy from upstream would work too — `weft-proxy` exists so the
operator gets a single binary with everything weft expects to find
already compiled in, no `xcaddy build --with …` step required on
each host.

## What's compiled in

- `github.com/caddyserver/caddy/v2/modules/standard` — the full
  upstream module set (http server, reverse_proxy, ACME issuer, …).
- `github.com/gsmlg-dev/caddy-storage-etcd` — etcd3 storage adapter.

## Build

```sh
go build -trimpath -ldflags "-s -w" -o weft-proxy .
```

Cross-build for Linux:

```sh
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
  go build -trimpath -ldflags "-s -w" -o weft-proxy-linux-arm64 .
```

Or grab a release artifact from
[ghcr.io/openweft/weft-proxy](https://github.com/openweft/weft-proxy/pkgs/container/weft-proxy)
once CI has published one (`v0.1.0` and up).

## Usage from weft-agent

```sh
weft agent --proxy \
  --proxy-caddy-binary=/usr/local/bin/weft-proxy \
  --proxy-state-dir=/var/lib/weft-agent/proxy
```

When `WEFT_PROXY_STORAGE_ETCD_ENDPOINTS` is also set in the agent's
environment, weft-agent emits a Caddy config block that selects the
etcd3 storage adapter — that's why the module has to be present in
the binary at build time, not loaded at runtime (Caddy doesn't
support runtime module loading).

## CLI

`weft-proxy` exposes the same CLI as upstream `caddy` — `run`, `stop`,
`reload`, `fmt`, `validate`, `version`, all the rest. weft-agent uses
`run --config -` to spawn it; ad-hoc operations during development
work directly via the binary.

```sh
weft-proxy version
weft-proxy run --config /etc/caddy/Caddyfile
weft-proxy validate --config /etc/caddy/Caddyfile
```

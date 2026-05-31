// weft-proxy is a Caddy build with the modules weft-agent needs at
// runtime baked in.
//
// Stock Caddy ships only the modules in github.com/caddyserver/caddy/v2/
// modules/standard. weft-agent's reverse-proxy plane (see
// github.com/openweft/weft/agent/proxy) speaks to Caddy through its
// admin API; the `caddy reload`-shaped requests it issues reference
// modules that aren't standard:
//
//   - etcd3 storage (so multiple weft-agent hosts share certificates
//     and avoid hammering Let's Encrypt's rate limit on every reload).
//
// xcaddy would handle this at build time too — we're producing the
// equivalent output from Go modules so the openweft build pipeline
// stays Go-native (no extra Docker step, no xcaddy host requirement).
//
// Operator usage:
//
//	weft agent --proxy --proxy-caddy-binary=/usr/local/bin/weft-proxy
//
// The weft-agent's proxy.Supervisor launches this binary as a child
// process and pokes its admin socket to push route config; from
// weft-agent's side nothing else changes.
package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// Core Caddy modules (http server, reverse_proxy handler, ACME
	// issuer, file_server, …). Without this the binary boots empty.
	_ "github.com/caddyserver/caddy/v2/modules/standard"

	// etcd storage adapter from darkweak/storages — selects when the
	// Caddy JSON config's top-level `storage.module` is "etcd". weft-
	// agent emits that block when the operator sets
	// WEFT_PROXY_STORAGE_ETCD_ENDPOINTS.
	_ "github.com/darkweak/storages/etcd/caddy"
)

func main() { caddycmd.Main() }

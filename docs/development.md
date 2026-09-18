# Development

You need Go and Docker. Node runs through Docker, no need to install it.

```sh
make dev      # Go server on :7676, Vite on :5173 with hot reload
make test     # Go tests and svelte-check
make ui       # build the interface into web/dist
make build    # interface and binary, into bin/caern
make run      # build, then serve on :7676
make docker   # build the image
make up       # docker compose up, on CAERN_PORT (7676 by default)
```

Set `PORT` to use another local port, for example `make dev PORT=7080`.

## Environment variables

Useful when running Caern outside Docker.

- `CAERN_ADDR`: listening address, `:7676` by default
- `CAERN_CONFIG_DIR`: configuration folder, `config` by default and `/config` in the image
- `CAERN_PORT`: host port used by `compose.yaml`, `7676` by default
- `CAERN_API`: where `make dev` sends the Vite requests, `http://localhost:7676` by default
- `PUID` and `PGID`: owner of `/config` in the image, root by default

## Structure

- `cmd/caern`: entry point
- `docker-entrypoint.sh`: gives `/config` to `PUID:PGID` and starts Caern as that user
- `internal/config`: reads and watches the YAML files
- `internal/server`: HTTP server, pushes updates to the page with server-sent events
- `web`: Svelte frontend, built into `web/dist` and embedded in the binary

Themes live in `web/src/styles/themes`. The types in `web/src/lib/types.ts` must match the Go structs in `internal/config`.

`npm run build` recreates `web/dist/.gitkeep`; without it, `go:embed` fails on a fresh clone.

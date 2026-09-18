# Caern

Self-hosted homepage for your services, configured in YAML.

Write your services in a YAML file and Caern shows them as tiles, in sections if you want. Edit the file and the page updates itself, no restart. Tiles can also be moved and resized from the page, which writes the file back.

Go backend, Svelte interface embedded in the binary, one container.

## Installation

```sh
docker run -d \
  --name caern \
  -p 7676:7676 \
  -v /path/to/config:/config \
  -e PUID=1000 \
  -e PGID=1000 \
  zerathblackfin/caern:latest
```

Open `http://<server>:7676`, where `<server>` is the machine you started it on.

Or with Docker Compose:

```yaml
services:
  caern:
    image: zerathblackfin/caern:latest
    container_name: caern
    ports:
      - 7676:7676
    volumes:
      - /path/to/config:/config
    environment:
      PUID: "1000"
      PGID: "1000"
    restart: unless-stopped
```

The same image is on GHCR, as `ghcr.io/zerathblackfin/caern`, if you prefer it.

To build it yourself instead:

```sh
git clone https://github.com/ZerathBlackfin/Caern.git
cd Caern
docker compose up --build
```

## Configuration

On first start, Caern creates two example files in the config folder:
- `settings.yaml` for the page settings
- `services.yaml` for the tiles

Your own icons go in `/config/icons` and background images in `/config/images`.

Mount the whole folder rather than single files, otherwise changes made with some editors won't be picked up.

If a file contains an error, Caern keeps showing the last valid version and lists the problems at the top of the page, with the file and line of each one.

Caern has no authentication: anyone who can open the page can rearrange the tiles. Put it behind a reverse proxy if it is reachable from outside your network.

## User and group

Caern writes the files in `/config`, so they should belong to you and not to root. Set `PUID` and `PGID` to your own user:

```sh
id your-user
# uid=1000(you) gid=1000(you)
```

On Synology the pair is usually `1026` and `100`. Left unset, the files belong to root and you need `sudo` to edit them.

## Documentation

- [settings.yaml](docs/settings.md)
- [services.yaml](docs/services.md)
- [Icons](docs/icons.md)
- [Development](docs/development.md)

## Updating

```sh
docker compose pull && docker compose up -d
```

## License

[AGPL-3.0](LICENSE). Use it, change it, run it. If you publish a modified version, or run one that other people use over a network, the source has to stay open.

# services.yaml

Tiles appear in the order of the file. They can be grouped in sections.

```yaml
- name: Proxmox
  url: https://proxmox.home.lan:8006
  icon: proxmox
  description: Hypervisor
  size: 2x1

- section: Media
  items:
    - name: Jellyfin
      url: http://jellyfin.home.lan:8096
      icon: jellyfin
      size: 2x2
    - name: Sonarr
      url: http://sonarr.home.lan:8989
      icon: sonarr
```

## Tiles

- `name`: required
- `url`: required
- `icon`: see [Icons](icons.md)
- `description`: shown below the name
- `size`: width and height in grid cells, like `2x1`. Defaults to `1x1`, up to `12x12`.
- `newTab`: overrides the value from `settings.yaml`

## Sections

A section has a title in `section` and its tiles in `items`. Sections can't be nested.

Tiles can also sit outside sections, anywhere in the file.

## Editing from the page

The **Edit** button in the header turns the tiles into a board: drag one to move it, inside its section or into another, and pull the corner to resize it. **Save** writes the new order and sizes into `services.yaml`, **Cancel** drops them.

Only the order and the `size` of the tiles change. Everything else, comments included, is left as it is.

## Layout

Tiles are placed from left to right, then top to bottom. Smaller tiles fill the gaps left by larger ones.

On narrow screens there are fewer columns (a phone usually gets 2), and tiles wider than the screen are shrunk to fit. The order stays the same.

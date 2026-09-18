# settings.yaml

All settings are optional.

```yaml
title: Home
theme: dusk
columns: 6
newTab: true
background: "#1b2230"
```

| Key | Default | Description |
| --- | --- | --- |
| `title` | `Caern` | Shown in the header and in the browser tab |
| `theme` | `dusk` | Color theme, `dusk` is the only one for now |
| `columns` | `6` | Columns on a large screen, from 1 to 24 |
| `newTab` | `true` | Open links in a new tab, can be changed per tile |
| `background` | theme color | A color or an image, see below |

The page gets wider as `columns` grows, up to the width of the window. Tiles stay around 150px wide or more, so a window too narrow for all the columns gets fewer of them.

## background

Colors must be quoted, otherwise YAML reads them as a comment:

```yaml
background: "#1b2230"
background: "rgb(27 34 48)"
```

Hex values are accepted, as well as the `rgb()`, `rgba()`, `hsl()`, `hsla()`, `hwb()`, `lab()`, `lch()`, `oklab()` and `oklch()` functions.

Images can be a URL or a file placed in `/config/images`:

```yaml
background: /images/wallpaper.jpg
background: https://example.com/wallpaper.jpg
```

A dark overlay is added on top of images to keep the text readable.

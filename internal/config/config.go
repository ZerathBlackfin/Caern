package config

import (
	"cmp"
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v3"
)

const (
	settingsFile = "settings.yaml"
	servicesFile = "services.yaml"

	defaultTitle   = "Caern"
	defaultTheme   = "dusk"
	defaultColumns = 6
	maxColumns     = 24
	maxTileSpan    = 12

	dashboardIconsURL = "https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons"
)

var (
	files  = []string{settingsFile, servicesFile}
	themes = []string{defaultTheme}
)

type settings struct {
	Title      string `yaml:"title"`
	Theme      string `yaml:"theme"`
	Background string `yaml:"background"`
	Columns    *int   `yaml:"columns"`
	NewTab     *bool  `yaml:"newTab"`
}

type item struct {
	Name        string `yaml:"name"`
	URL         string `yaml:"url"`
	Icon        string `yaml:"icon"`
	Description string `yaml:"description"`
	Size        string `yaml:"size"`
	NewTab      *bool  `yaml:"newTab"`

	Section string `yaml:"section"`
	Items   []item `yaml:"items"`
}

type Page struct {
	Title           string  `json:"title"`
	Theme           string  `json:"theme"`
	BackgroundColor string  `json:"backgroundColor,omitempty"`
	BackgroundImage string  `json:"backgroundImage,omitempty"`
	Columns         int     `json:"columns"`
	Blocks          []Block `json:"blocks"`
}

type Block struct {
	ID      string `json:"id,omitempty"`
	Section string `json:"section,omitempty"`
	Tiles   []Tile `json:"tiles"`
}

type Tile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	NewTab      bool   `json:"newTab"`
}

func parse(settingsData, servicesData []byte) (*Page, error) {
	var problems Problems

	var s settings
	settingsRoot := decode(settingsFile, settingsData, &s, &problems)
	page := s.page(settingsRoot, &problems)

	var items []item
	servicesRoot := decode(servicesFile, servicesData, &items, &problems)
	page.addItems(items, servicesRoot, s.newTab(), &problems)

	if len(problems) > 0 {
		problems.sort()
		return nil, problems
	}
	return page, nil
}

func (s *settings) page(root *yaml.Node, problems *Problems) *Page {
	p := &Page{
		Title:   cmp.Or(s.Title, defaultTitle),
		Theme:   cmp.Or(s.Theme, defaultTheme),
		Columns: defaultColumns,
		Blocks:  []Block{},
	}
	if s.Columns != nil {
		p.Columns = *s.Columns
	}
	p.setBackground(s.Background, root, problems)
	if !slices.Contains(themes, p.Theme) {
		problems.add(settingsFile, value(root, "theme"), "", "unknown theme %q, only dusk exists", p.Theme)
	}
	if p.Columns < 1 || p.Columns > maxColumns {
		problems.add(settingsFile, value(root, "columns"), "", "columns must be between 1 and %d, got %d", maxColumns, p.Columns)
	}
	return p
}

func (s *settings) newTab() bool {
	return s.NewTab == nil || *s.NewTab
}

var cssColor = regexp.MustCompile(`^(#([[:xdigit:]]{3,4}|[[:xdigit:]]{6}|[[:xdigit:]]{8})|(rgba?|hsla?|hwb|lab|lch|oklab|oklch)\([^()]*\))$`)

func (p *Page) setBackground(bg string, root *yaml.Node, problems *Problems) {
	key, val := entry(root, "background")
	unquotedColor := ""
	if key != nil {
		unquotedColor = strings.TrimSpace(key.LineComment)
	}
	switch bg = strings.TrimSpace(bg); {
	case bg == "" && cssColor.MatchString(unquotedColor):
		problems.add(settingsFile, key, "", `put the color in quotes: background: "%s"`, unquotedColor)
	case bg == "":
	case isLink(bg):
		p.BackgroundImage = bg
	case cssColor.MatchString(bg):
		p.BackgroundColor = bg
	default:
		problems.add(settingsFile, val, "", `background %q must be a color like "#1b2230", or an image like /images/wall.jpg`, bg)
	}
}

func (p *Page) addItems(items []item, list *yaml.Node, newTab bool, problems *Problems) {
	for i, it := range items {
		node := child(list, i)
		id := strconv.Itoa(i)
		if it.Section == "" {
			if tile, ok := it.tile(id, node, newTab, problems); ok {
				p.appendLoose(tile)
			}
			continue
		}

		if it.Name != "" || it.URL != "" {
			problems.add(servicesFile, node, it.Section, "a section can't also have a name or url")
		}
		block := Block{ID: id, Section: it.Section, Tiles: []Tile{}}
		children := value(node, "items")
		for j, c := range it.Items {
			cnode := child(children, j)
			if c.Section != "" {
				problems.add(servicesFile, cnode, c.Section, "sections can't be inside sections")
				continue
			}
			if tile, ok := c.tile(id+"."+strconv.Itoa(j), cnode, newTab, problems); ok {
				block.Tiles = append(block.Tiles, tile)
			}
		}
		p.Blocks = append(p.Blocks, block)
	}
}

func (p *Page) appendLoose(t Tile) {
	if n := len(p.Blocks); n > 0 && p.Blocks[n-1].Section == "" {
		p.Blocks[n-1].Tiles = append(p.Blocks[n-1].Tiles, t)
		return
	}
	p.Blocks = append(p.Blocks, Block{Tiles: []Tile{t}})
}

func (it *item) tile(id string, node *yaml.Node, pageNewTab bool, problems *Problems) (Tile, bool) {
	before := len(*problems)
	report := func(format string, args ...any) {
		problems.add(servicesFile, node, it.Name, format, args...)
	}
	if len(it.Items) > 0 {
		report("only sections can have items")
	}
	if it.Name == "" {
		report("a tile needs a name")
	}
	if it.URL == "" {
		report("a tile needs a url")
	}
	w, h, err := parseSize(it.Size)
	if err != nil {
		report("%v", err)
	}
	icon, err := resolveIcon(it.Icon)
	if err != nil {
		report("%v", err)
	}
	if len(*problems) > before {
		return Tile{}, false
	}

	newTab := pageNewTab
	if it.NewTab != nil {
		newTab = *it.NewTab
	}
	return Tile{
		ID:          id,
		Name:        it.Name,
		URL:         it.URL,
		Icon:        icon,
		Description: it.Description,
		Width:       w,
		Height:      h,
		NewTab:      newTab,
	}, true
}

func parseSize(s string) (w, h int, err error) {
	if s == "" {
		return 1, 1, nil
	}
	ws, hs, ok := strings.Cut(strings.ToLower(s), "x")
	if ok {
		w, err = strconv.Atoi(ws)
		if err == nil {
			h, err = strconv.Atoi(hs)
		}
	}
	if !ok || err != nil || w < 1 || h < 1 || w > maxTileSpan || h > maxTileSpan {
		return 0, 0, fmt.Errorf("size must be WIDTHxHEIGHT, from 1x1 to %dx%d, got %q", maxTileSpan, maxTileSpan, s)
	}
	return w, h, nil
}

func resolveIcon(s string) (string, error) {
	switch {
	case s == "":
		return "", nil
	case isLink(s):
		return s, nil
	}
	switch ext := path.Ext(s); ext {
	case "":
		return fmt.Sprintf("%s/png/%s.png", dashboardIconsURL, s), nil
	case ".svg", ".png", ".webp":
		return fmt.Sprintf("%s/%s/%s", dashboardIconsURL, ext[1:], s), nil
	}
	return "", fmt.Errorf("icon %q is not a Dashboard Icons name, an /icons/ file or a URL", s)
}

func isLink(s string) bool {
	return strings.HasPrefix(s, "/") || strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

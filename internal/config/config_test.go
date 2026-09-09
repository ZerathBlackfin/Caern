package config

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func starter(t *testing.T, name string) []byte {
	t.Helper()
	data, err := defaults.ReadFile("defaults/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestParseDefaults(t *testing.T) {
	p, err := parse(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Caern" || p.Theme != "dusk" || p.Columns != 6 || len(p.Blocks) != 0 {
		t.Fatalf("unexpected defaults: %+v", p)
	}
}

func TestParseStarterFiles(t *testing.T) {
	p, err := parse(starter(t, settingsFile), starter(t, servicesFile))
	if err != nil {
		t.Fatal(err)
	}
	if p.Title != "Home" || len(p.Blocks) != 3 || p.Blocks[0].Section != "" || p.Blocks[1].Section != "Media" {
		t.Fatalf("unexpected page: %+v", p)
	}
}

func TestParseGroupsLooseTiles(t *testing.T) {
	p, err := parse(nil, []byte(`
- {name: A, url: http://a}
- {name: B, url: http://b}
- section: S
  items:
    - {name: C, url: http://c}
- {name: D, url: http://d}
`))
	if err != nil {
		t.Fatal(err)
	}
	got := []string{}
	for _, b := range p.Blocks {
		names := []string{}
		for _, tile := range b.Tiles {
			names = append(names, tile.Name)
		}
		got = append(got, b.Section+":"+strings.Join(names, ","))
	}
	if strings.Join(got, " ") != ":A,B S:C :D" {
		t.Fatalf("got %v", got)
	}
}

func TestParseTile(t *testing.T) {
	p, err := parse([]byte("newTab: false"), []byte(`
- name: Sonarr
  url: http://sonarr
  icon: sonarr.svg
  size: 2X3
  newTab: true
`))
	if err != nil {
		t.Fatal(err)
	}
	want := Tile{
		Name:   "Sonarr",
		URL:    "http://sonarr",
		Icon:   dashboardIconsURL + "/svg/sonarr.svg",
		Width:  2,
		Height: 3,
		NewTab: true,
	}
	if got := p.Blocks[0].Tiles[0]; got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestParseBackground(t *testing.T) {
	for in, want := range map[string][2]string{
		`background: "#1b2230"`:             {"#1b2230", ""},
		`background: "rgb(20 30 40 / 50%)"`: {"rgb(20 30 40 / 50%)", ""},
		`background: /images/wall.jpg`:      {"", "/images/wall.jpg"},
		`background: https://x.dev/a.png`:   {"", "https://x.dev/a.png"},
		``:                                  {"", ""},
	} {
		p, err := parse([]byte(in), nil)
		if err != nil {
			t.Fatalf("%s: %v", in, err)
		}
		if got := [2]string{p.BackgroundColor, p.BackgroundImage}; got != want {
			t.Errorf("%s: got %q, want %q", in, got, want)
		}
	}
}

func TestResolveIcon(t *testing.T) {
	for in, want := range map[string]string{
		"":                  "",
		"sonarr":            dashboardIconsURL + "/png/sonarr.png",
		"sonarr.webp":       dashboardIconsURL + "/webp/sonarr.webp",
		"/icons/nas.png":    "/icons/nas.png",
		"https://x.dev/a.b": "https://x.dev/a.b",
	} {
		got, err := resolveIcon(in)
		if err != nil || got != want {
			t.Errorf("resolveIcon(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	if _, err := resolveIcon("nas.jpg"); err == nil {
		t.Error("expected an error for an unknown icon format")
	}
}

func TestParseProblems(t *testing.T) {
	for name, tc := range map[string]struct{ settings, services, want string }{
		"syntax":          {"title: [", "", "settings.yaml:1: did not find expected node content"},
		"unknown key":     {"title: x\ntitel: y", "", `settings.yaml:2: unknown setting "titel"`},
		"wrong type":      {"columns: nope", "", `settings.yaml:1: expected a whole number, got "nope"`},
		"columns":         {"title: x\ncolumns: 30", "", "settings.yaml:2: columns must be between 1 and 24, got 30"},
		"bad background":  {"background: wall.jpg", "", `settings.yaml:1: background "wall.jpg" must be a color`},
		"unknown theme":   {"theme: dsuk", "", `settings.yaml:1: unknown theme "dsuk", only dusk exists`},
		"zero columns":    {"columns: 0", "", "settings.yaml:1: columns must be between 1 and 24, got 0"},
		"unquoted color":  {"title: x\nbackground: #1b2230", "", `settings.yaml:2: put the color in quotes: background: "#1b2230"`},
		"not a list":      {"", "foo: bar", "services.yaml:1: expected a list"},
		"missing name":    {"", "- url: http://a", "services.yaml:1: a tile needs a name"},
		"missing url":     {"", "- name: A\n- name: B\n  url: u", "services.yaml:1: A: a tile needs a url"},
		"bad size":        {"", "- {name: A, url: u, size: big}", `services.yaml:1: A: size must be WIDTHxHEIGHT, from 1x1 to 12x12, got "big"`},
		"too wide":        {"", "- {name: A, url: u, size: 13x1}", `got "13x1"`},
		"bad icon":        {"", "- {name: A, url: u, icon: a.jpg}", `services.yaml:1: A: icon "a.jpg" is not`},
		"section + tile":  {"", "- {section: S, name: A}", "services.yaml:1: S: a section can't also have a name or url"},
		"nested":          {"", "- section: S\n  items:\n    - section: T", "services.yaml:3: T: sections can't be inside sections"},
		"tile with items": {"", "- {name: A, url: u, items: [{name: B}]}", "services.yaml:1: A: only sections can have items"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := parse([]byte(tc.settings), []byte(tc.services))
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want a problem containing %q", err, tc.want)
			}
		})
	}
}

func TestParseReportsEveryProblem(t *testing.T) {
	_, err := parse([]byte("columns: 99\ntitel: x"), []byte(`- name: Radarr
  ulr: http://radarr
  size: big
- section: Media
  items:
    - url: http://jellyfin
`))
	want := `settings.yaml:1: columns must be between 1 and 24, got 99
settings.yaml:2: unknown setting "titel"
services.yaml:1: Radarr: a tile needs a url
services.yaml:1: Radarr: size must be WIDTHxHEIGHT, from 1x1 to 12x12, got "big"
services.yaml:2: unknown setting "ulr"
services.yaml:6: a tile needs a name`
	if err == nil || err.Error() != want {
		t.Fatalf("got:\n%v\nwant:\n%s", err, want)
	}
}

func TestStoreKeepsLastValidPage(t *testing.T) {
	dir := t.TempDir()
	s, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if st := s.Current(); len(st.Problems) != 0 || len(st.Page.Blocks) != 3 {
		t.Fatalf("starter files not loaded: %+v", st)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go s.Watch(ctx)
	changed, unsubscribe := s.Subscribe()
	defer unsubscribe()
	time.Sleep(50 * time.Millisecond)

	write := func(name, content string) State {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		select {
		case <-changed:
		case <-time.After(2 * time.Second):
			t.Fatal("no reload after writing the file")
		}
		return s.Current()
	}

	st := write(settingsFile, "title: Lab")
	if len(st.Problems) != 0 || st.Page.Title != "Lab" {
		t.Fatalf("got %+v", st)
	}
	st = write(servicesFile, "- name: A")
	if len(st.Problems) != 1 || st.Problems[0].Line != 1 || st.Page.Title != "Lab" || len(st.Page.Blocks) != 3 {
		t.Fatalf("broken file should keep the last page and report the problem: %+v", st)
	}
	st = write(servicesFile, "- {name: A, url: http://a}")
	if len(st.Problems) != 0 || len(st.Page.Blocks) != 1 {
		t.Fatalf("fixed file should clear the problems: %+v", st)
	}
}

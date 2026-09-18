package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const servicesWithComments = `# Options: docs/services.md in the Caern repository

- name: Proxmox
  url: https://proxmox
  icon: proxmox
  size: 2x1

# media stack
- section: Media
  items:
    - name: Jellyfin
      url: http://jellyfin
    - name: Sonarr # shows
      url: http://sonarr
`

func storeWith(t *testing.T, services string) (*Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, servicesFile)
	for name, content := range map[string]string{settingsFile: "title: Home\n", servicesFile: services} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	s, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	return s, path
}

func TestSaveLayoutKeepsComments(t *testing.T) {
	s, path := storeWith(t, servicesWithComments)
	err := s.SaveLayout(Layout{Blocks: []LayoutBlock{
		{Tiles: []LayoutTile{{ID: "1.1", Size: "1x1"}}},
		{ID: "1", Tiles: []LayoutTile{{ID: "1.0", Size: "2x2"}}},
		{Tiles: []LayoutTile{{ID: "0", Size: "2x1"}}},
	}})
	if err != nil {
		t.Fatal(err)
	}

	want := `# Options: docs/services.md in the Caern repository

- name: Sonarr # shows
  url: http://sonarr

# media stack
- section: Media
  items:
    - name: Jellyfin
      url: http://jellyfin
      size: 2x2

- name: Proxmox
  url: https://proxmox
  icon: proxmox
  size: 2x1
`
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestSaveLayoutReloadsThePage(t *testing.T) {
	s, _ := storeWith(t, servicesWithComments)
	if err := s.SaveLayout(Layout{Blocks: []LayoutBlock{
		{ID: "1", Tiles: []LayoutTile{{ID: "1.0"}, {ID: "1.1"}, {ID: "0", Size: "3x2"}}},
	}}); err != nil {
		t.Fatal(err)
	}
	st := s.Current()
	if len(st.Problems) != 0 || len(st.Page.Blocks) != 1 {
		t.Fatalf("got %+v", st)
	}
	block := st.Page.Blocks[0]
	if block.Section != "Media" || len(block.Tiles) != 3 {
		t.Fatalf("got %+v", block)
	}
	if tile := block.Tiles[2]; tile.Name != "Proxmox" || tile.Width != 3 || tile.Height != 2 {
		t.Fatalf("got %+v", tile)
	}
}

func TestSaveLayoutRefuses(t *testing.T) {
	for name, tc := range map[string]struct {
		layout Layout
		want   string
	}{
		"unknown tile": {Layout{Blocks: []LayoutBlock{{Tiles: []LayoutTile{{ID: "9"}}}}}, `unknown tile "9"`},
		"tile twice":   {Layout{Blocks: []LayoutBlock{{Tiles: []LayoutTile{{ID: "0"}, {ID: "0"}}}}}, `"0" is used twice`},
		"bad size":     {Layout{Blocks: []LayoutBlock{{Tiles: []LayoutTile{{ID: "0", Size: "huge"}}}}}, "size must be WIDTHxHEIGHT"},
	} {
		t.Run(name, func(t *testing.T) {
			s, path := storeWith(t, servicesWithComments)
			err := s.SaveLayout(tc.layout)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want an error containing %q", err, tc.want)
			}
			got, _ := os.ReadFile(path)
			if string(got) != servicesWithComments {
				t.Fatalf("the file should not change, got:\n%s", got)
			}
		})
	}
}

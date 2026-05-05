package core

import (
	"testing"
)

func TestTopoSort_NoDeps(t *testing.T) {
	mods := []Module{
		{Name: "a"},
		{Name: "b"},
		{Name: "c"},
	}

	g := NewDepGraph(mods)
	sorted, err := g.TopoSort()
	if err != nil {
		t.Fatalf("TopoSort error: %v", err)
	}
	if len(sorted) != 3 {
		t.Errorf("expected 3, got %d", len(sorted))
	}
}

func TestTopoSort_WithDeps(t *testing.T) {
	mods := []Module{
		{Name: "nvchad", DependsOn: []string{"neovim"}},
		{Name: "neovim"},
		{Name: "starship"},
	}

	g := NewDepGraph(mods)
	sorted, err := g.TopoSort()
	if err != nil {
		t.Fatalf("TopoSort error: %v", err)
	}

	// neovim must come before nvchad
	nvimIdx := -1
	nvchadIdx := -1
	for i, m := range sorted {
		if m.Name == "neovim" {
			nvimIdx = i
		}
		if m.Name == "nvchad" {
			nvchadIdx = i
		}
	}
	if nvimIdx >= nvchadIdx {
		t.Errorf("neovim (idx=%d) should come before nvchad (idx=%d)", nvimIdx, nvchadIdx)
	}
}

func TestTopoSort_Cycle(t *testing.T) {
	mods := []Module{
		{Name: "a", DependsOn: []string{"b"}},
		{Name: "b", DependsOn: []string{"c"}},
		{Name: "c", DependsOn: []string{"a"}},
	}

	g := NewDepGraph(mods)
	_, err := g.TopoSort()
	if err == nil {
		t.Error("expected cycle error, got nil")
	}
}

func TestTopoSort_ExternalDep(t *testing.T) {
	// Module depends on something not in the graph
	mods := []Module{
		{Name: "nvchad", DependsOn: []string{"external-thing"}},
	}

	g := NewDepGraph(mods)
	sorted, err := g.TopoSort()
	if err != nil {
		t.Fatalf("should not error on external dep: %v", err)
	}
	if len(sorted) != 1 {
		t.Errorf("expected 1, got %d", len(sorted))
	}
}

func TestTopoSortSubset(t *testing.T) {
	mods := []Module{
		{Name: "nvchad", DependsOn: []string{"neovim"}},
		{Name: "neovim", DependsOn: []string{"build-essential"}},
		{Name: "build-essential"},
		{Name: "starship"},
	}

	g := NewDepGraph(mods)

	// Request only nvchad, should pull in neovim + build-essential
	sorted, err := g.TopoSortSubset([]string{"nvchad"}, func(name string) bool {
		return false // nothing is installed
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(sorted) != 3 {
		t.Errorf("expected 3 (nvchad + deps), got %d", len(sorted))
	}

	// If neovim is already installed, should only need nvchad
	sorted2, err := g.TopoSortSubset([]string{"nvchad"}, func(name string) bool {
		return name == "neovim" || name == "build-essential"
	})
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(sorted2) != 1 {
		t.Errorf("expected 1 (only nvchad), got %d", len(sorted2))
	}
}

func TestGetDependencies(t *testing.T) {
	mods := []Module{
		{Name: "a", DependsOn: []string{"b", "c"}},
		{Name: "b"},
		{Name: "c"},
	}
	g := NewDepGraph(mods)
	deps := g.GetDependencies("a")
	if len(deps) != 2 {
		t.Errorf("expected 2 deps, got %d", len(deps))
	}
}

func TestGetDependees(t *testing.T) {
	mods := []Module{
		{Name: "a", DependsOn: []string{"b"}},
		{Name: "b"},
		{Name: "c", DependsOn: []string{"b"}},
	}
	g := NewDepGraph(mods)
	dependees := g.GetDependees("b")
	if len(dependees) != 2 {
		t.Errorf("expected 2 dependees, got %d", len(dependees))
	}
}

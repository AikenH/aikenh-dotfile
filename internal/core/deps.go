package core

import "fmt"

// DepGraph represents a directed acyclic graph of module dependencies
type DepGraph struct {
	nodes map[string]*DepNode
}

// DepNode is a node in the dependency graph
type DepNode struct {
	Name      string
	Module    Module
	DependsOn []string // names of modules this depends on
	Dependees []string // names of modules that depend on this
}

// NewDepGraph builds a dependency graph from a list of modules
func NewDepGraph(modules []Module) *DepGraph {
	g := &DepGraph{
		nodes: make(map[string]*DepNode),
	}

	// Add all modules as nodes
	for _, m := range modules {
		g.nodes[m.Name] = &DepNode{
			Name:      m.Name,
			Module:    m,
			DependsOn: m.DependsOn,
		}
	}

	// Build reverse edges (dependees)
	for _, node := range g.nodes {
		for _, dep := range node.DependsOn {
			if depNode, ok := g.nodes[dep]; ok {
				depNode.Dependees = append(depNode.Dependees, node.Name)
			}
		}
	}

	return g
}

// TopoSort returns modules in dependency order (dependencies first)
// Returns error if there's a cycle
func (g *DepGraph) TopoSort() ([]Module, error) {
	visited := make(map[string]bool)
	temp := make(map[string]bool) // for cycle detection
	var result []Module

	var visit func(name string) error
	visit = func(name string) error {
		if temp[name] {
			return fmt.Errorf("dependency cycle detected involving %q", name)
		}
		if visited[name] {
			return nil
		}

		temp[name] = true

		node, ok := g.nodes[name]
		if !ok {
			// External dependency not in our graph — skip silently
			visited[name] = true
			temp[name] = false
			return nil
		}

		for _, dep := range node.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}

		temp[name] = false
		visited[name] = true
		result = append(result, node.Module)
		return nil
	}

	// Visit all nodes
	for name := range g.nodes {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// TopoSortSubset returns only the given modules in dependency order,
// including any transitive dependencies not already installed
func (g *DepGraph) TopoSortSubset(names []string, isInstalled func(string) bool) ([]Module, error) {
	needed := make(map[string]bool)

	// Recursively collect all dependencies
	var collect func(name string)
	collect = func(name string) {
		if needed[name] || isInstalled(name) {
			return
		}
		needed[name] = true
		if node, ok := g.nodes[name]; ok {
			for _, dep := range node.DependsOn {
				collect(dep)
			}
		}
	}

	for _, name := range names {
		collect(name)
	}

	// Now topo sort only the needed set
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	var result []Module

	var visit func(name string) error
	visit = func(name string) error {
		if !needed[name] {
			return nil
		}
		if temp[name] {
			return fmt.Errorf("dependency cycle detected involving %q", name)
		}
		if visited[name] {
			return nil
		}

		temp[name] = true
		node, ok := g.nodes[name]
		if ok {
			for _, dep := range node.DependsOn {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}

		temp[name] = false
		visited[name] = true
		if ok {
			result = append(result, node.Module)
		}
		return nil
	}

	for name := range needed {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// GetDependencies returns direct dependencies of a module
func (g *DepGraph) GetDependencies(name string) []string {
	if node, ok := g.nodes[name]; ok {
		return node.DependsOn
	}
	return nil
}

// GetDependees returns modules that directly depend on the given module
func (g *DepGraph) GetDependees(name string) []string {
	if node, ok := g.nodes[name]; ok {
		return node.Dependees
	}
	return nil
}

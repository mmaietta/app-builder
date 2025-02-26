package node_modules

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path"
	"testing"

	"github.com/develar/app-builder/pkg/fs"
	. "github.com/onsi/gomega"
	"github.com/samber/lo"
)

type NodeTreeDepItem struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type NodeTreeItem struct {
	Dir  string            `json:"dir"`
	Deps []NodeTreeDepItem `json:"deps"`
}

func nodeDepTree(t *testing.T, dir string, results []string) {
	g := NewGomegaWithT(t)
	rootPath := fs.FindParentWithFile(Dirname(), "go.mod")
	cmd := exec.Command("go", "run", path.Join(rootPath, "main.go"), "node-dep-tree", "--dir", dir)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("err", err)
	}
	g.Expect(err).NotTo(HaveOccurred())
	var j []NodeTreeItem
	json.Unmarshal(output, &j)
	r := lo.FlatMap(j, func(it NodeTreeItem, i int) []string {
		return lo.Map(it.Deps, func(it NodeTreeDepItem, i int) string {
			return it.Name
		})
	})
	g.Expect(r).To(ConsistOf(results))
}

func TestNodeDepTreeCmd(t *testing.T) {
	nodeDepTree(t, path.Join(Dirname(), "npm-demo"), []string{
		"react", "js-tokens", "loose-envify",
	})

	nodeDepTree(t, path.Join(Dirname(), "pnpm-demo"), []string{
		"react", "js-tokens", "loose-envify",
	})

	nodeDepTree(t, path.Join(Dirname(), "yarn-demo"), []string{
		"ms", "foo", "ms",
	})
}

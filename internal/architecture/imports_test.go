package architecture

import (
	"os/exec"
	"strings"
	"testing"
)

// TestFeaturePackagesDoNotImportEachOther ensures features stay independent.
// Only the composition root (internal/app) and transport layer (internal/server)
// are allowed to wire features together.
func TestFeaturePackagesDoNotImportEachOther(t *testing.T) {
	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}|{{join .Imports \" \"}}", "./internal/feature/...")
	cmd.Dir = "../.."
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list failed: %v", err)
	}

	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		pkg := parts[0]
		imports := strings.Fields(parts[1])

		for _, imp := range imports {
			if strings.HasPrefix(imp, "github.com/victorive/coursepack/internal/feature/") {
				if imp == pkg {
					continue
				}
				t.Errorf("%s imports %s; features must not import each other", pkg, imp)
			}
		}
	}
}

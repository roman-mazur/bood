package bood

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrepareContext(t *testing.T) {
	ctx := PrepareContext()

	ctx.MockFileSystem(map[string][]byte{
		"Blueprints": nil,
	})

	_, errs := ctx.PrepareBuildActions(NewConfig())
	if len(errs) > 0 {
		t.Fatalf("Unexpected errors preparing build actions: %s", errs)
	}
	buffer := new(bytes.Buffer)
	if err := ctx.WriteBuildFile(buffer); err == nil {
		text := buffer.String()
		t.Logf("Generated file:\n%s", text)
		if !strings.Contains(text, "out/bin: ") {
			t.Errorf("Generated file has no out/bin build actions")
		}
		if !strings.Contains(text, "builddir = ") {
			t.Error("Generated file has no build dir definition")
		}
	} else {
		t.Errorf("Unexpected error writing build actios: %s", err)
	}
}

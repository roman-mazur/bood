package bood

import (
	"github.com/google/blueprint"
)

const (
	binOutPath = "bin"

	BlueprintFileName = "build.bood"
)

var (
	pctx = blueprint.NewPackageContext(PackagePath)

	createDir = pctx.StaticRule("createDir", blueprint.RuleParams{
		Command:     "mkdir -p $out",
		Description: "Ensure that directory $out exists",
	})
)

// PrepareContext creates a new blueprint.Context registering some common modules/singletons
// that produce build actions related to Config features.
func PrepareContext() *blueprint.Context {
	ctx := blueprint.NewContext()

	// Register a singleton that will take care of config.BinOutputPath() directory creation.
	ctx.RegisterSingletonType("boodInternalBinOut", binOutFactory)

	return ctx
}

// binOutputGenerator generates build actions that create $out/bin directory.
type binOutputGenerator struct{}

func (binOut binOutputGenerator) GenerateBuildActions(ctx blueprint.SingletonContext) {
	config := ExtractConfig(ctx)
	config.Debug.Printf("Adding build actions for %s", config.BinOutputPath())
	ctx.SetNinjaBuildDir(pctx, config.BaseOutputDir)
	ctx.Build(pctx, blueprint.BuildParams{
		Rule:     createDir,
		Outputs:  []string{config.BinOutputPath()},
		Optional: true,
	})
}

func binOutFactory() blueprint.Singleton {
	return binOutputGenerator{}
}

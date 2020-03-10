package bood

import (
	"github.com/google/blueprint"
)

const PackagePath = "github.com/roman-mazur/bood"

var (
	pctx = blueprint.NewPackageContext(PackagePath)
)

// PrepareContext creates a new blueprint.Context registering some common modules/singletons
// that produce build actions related to Config features.
func PrepareContext() *blueprint.Context {
	ctx := blueprint.NewContext()

	// Register a singleton that will take care of config.BinOutputPath() directory creation.
	ctx.RegisterSingletonType("boodInternalBinOut", outputDefFactory)

	return ctx
}

// outputDef adds instructions necessary to configure ninja output directory.
type outputDef struct{}

func (outDef outputDef) GenerateBuildActions(ctx blueprint.SingletonContext) {
	config := ExtractConfig(ctx)
	config.Debug.Printf("Configure output as %s", config.BaseOutputDir)
	ctx.SetNinjaBuildDir(pctx, config.BaseOutputDir)
}

func outputDefFactory() blueprint.Singleton {
	return outputDef{}
}

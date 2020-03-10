package bood

import (
	"github.com/google/blueprint"
	"log"
	"os"
	"path"
	"path/filepath"
)

const BlueprintFileName = "build.bood"

// GenerateBuildFile creates build.ninja file in the base output directory (as specified by the config).
// This function will call methods on the context to parse blueprint files, prepare build actions,
// and write the result file. All the modules and singletons need to be registered before calling this function.
// This function is supposed to be called from the main function. In case of unexpected errors,
// it will terminate the process printing error messages using the loggers provided in the config.
func GenerateBuildFile(config *Config, ctx *blueprint.Context) string {
	workingDir, err := os.Getwd()
	if err != nil {
		config.Info.Fatalf("Cannot obtain working directory: %s", err)
	}
	rootBuildFile := path.Join(workingDir, BlueprintFileName)

	if _, err := os.Stat(config.BaseOutputDir); os.IsNotExist(err) {
		_ = os.MkdirAll(config.BaseOutputDir, os.ModePerm)
	}

	// Configure the context to read bood files in the current dir.
	if err := initModulesListFile(config, ctx); err != nil {
		log.Fatalf("Failed to initialize bood: %s", err)
	}

	deps, errs := ctx.ParseBlueprintsFiles(rootBuildFile, config)
	checkFatalErrors(config, "Problems parsing blueprint files", errs)
	config.Debug.Printf("Parsed blueprint files: %s", deps)

	_, errs = ctx.PrepareBuildActions(config)
	checkFatalErrors(config, "Problems preparing build actions", errs)

	config.Debug.Println("Start writing ninja build file...")
	ninjaBuildPath := path.Join(config.BaseOutputDir, "build.ninja")
	ninjaFile, err := os.Create(ninjaBuildPath)
	if err != nil {
		config.Info.Fatalf("Cannot create new build.ninja file: %s", err)
	}
	defer ninjaFile.Close()

	if err := ctx.WriteBuildFile(ninjaFile); err != nil {
		config.Info.Fatalf("Cannot write to build.ninja file: %s", err)
	}
	config.Info.Printf("Ninja build file is generated at %s", ninjaBuildPath)
	return ninjaBuildPath
}

func checkFatalErrors(config *Config, message string, errs []error) {
	if len(errs) > 0 {
		config.Info.Fatalf("%s: %s", message, errs)
	}
}

func initModulesListFile(config *Config, ctx *blueprint.Context) error {
	modulesList := filepath.Join(config.BaseOutputDir, ".bood-modules-list")

	ml, err := os.Create(modulesList)
	if err != nil {
		return err
	}
	closed := false
	defer func() {
		if !closed {
			_ = ml.Close()
		}
	}()

	if _, err := ml.WriteString(BlueprintFileName); err != nil {
		return err
	}
	if err := ml.Close(); err != nil {
		return err
	}
	closed = true

	ctx.SetModuleListFile(modulesList)
	return nil
}

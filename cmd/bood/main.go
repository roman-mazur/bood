package main

import (
	"flag"
	"github.com/google/blueprint"
	"github.com/roman-mazur/bood"
	"github.com/roman-mazur/bood/gomodule"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var (
	dryRun  = flag.Bool("dry-run", false, "Generate ninja build file but don't start the build")
	verbose = flag.Bool("v", false, "Display debugging logs")
)

func NewContext() *blueprint.Context {
	ctx := bood.PrepareContext()

	// Configure the context to read bood files in the current dir.
	if err := initModulesListFile(ctx); err != nil {
		log.Fatalf("Failed to initialize bood: %s", err)
	}

	ctx.RegisterModuleType("go_binary", gomodule.BinFactory)
	return ctx
}

func initModulesListFile(ctx *blueprint.Context) error {
	modulesList := filepath.Join(os.TempDir(), "bood-modules-list")

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

	if _, err := ml.WriteString(bood.BlueprintFileName); err != nil {
		return err
	}
	if err := ml.Close(); err != nil {
		return err
	}
	closed = true

	ctx.SetModuleListFile(modulesList)
	return nil
}

func main() {
	flag.Parse()

	config := bood.NewConfig()
	if !*verbose {
		config.Debug = log.New(ioutil.Discard, "", 0)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		config.Info.Fatalf("Cannot obtain working directory: %s", err)
	}
	rootBuildFile := path.Join(workingDir, bood.BlueprintFileName)

	if _, err := os.Stat(config.BaseOutputDir); os.IsNotExist(err) {
		_ = os.MkdirAll(config.BaseOutputDir, os.ModePerm)
	}

	ctx := NewContext()

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

	if !*dryRun {
		config.Info.Println("Starting the build now")

		cmd := exec.Command("ninja", "-f", ninjaBuildPath)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			config.Info.Fatal("Error invoking ninja build. See logs above.")
		}
	}
}

func checkFatalErrors(config *bood.Config, message string, errs []error) {
	if len(errs) > 0 {
		config.Info.Fatalf("%s: %s", message, errs)
	}
}

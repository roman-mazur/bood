package bood

import (
	"log"
	"os"
	"path"
)

const PackagePath = "github.com/roman-mazur/bood"

// Config represents build system configuration.
type Config struct {
	// Base path to the output directory.
	BaseOutputDir string

	// Information-level logger to be used to trace build logic.
	Info *log.Logger
	// Debug-level logger to be used to trace build logic.
	Debug *log.Logger
}

// BinOutputPath returns a path to the output binaries directory.
func (cfg *Config) BinOutputPath() string {
	return path.Join(cfg.BaseOutputDir, binOutPath)
}

// NewConfig creates a new instance of bood Config with default values.
func NewConfig() *Config {
	return &Config{
		BaseOutputDir: "out",
		Info:          log.New(os.Stderr, "INFO ", log.LstdFlags),
		Debug:         log.New(os.Stderr, "DEBUG ", log.LstdFlags|log.Lshortfile),
	}
}

// ConfigContext is used to work with module or singleton context that give access to a config instance.
type ConfigContext interface {
	Config() interface{}
}

// ExtractConfig returns an instance of Config from the input module or singleton context.
func ExtractConfig(ctx ConfigContext) *Config {
	return ctx.Config().(*Config)
}

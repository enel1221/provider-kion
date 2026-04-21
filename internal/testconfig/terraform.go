package testconfig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const (
	envTerraformVersion         = "TERRAFORM_VERSION"
	envTerraformProviderSource  = "TERRAFORM_PROVIDER_SOURCE"
	envTerraformProviderVersion = "TERRAFORM_PROVIDER_VERSION"
)

type TerraformConfig struct {
	Version         string
	ProviderSource  string
	ProviderVersion string
}

var (
	loadTerraformConfigOnce  sync.Once
	loadedTerraformConfig    TerraformConfig
	loadedTerraformConfigErr error
)

func LoadTerraformConfig() (TerraformConfig, error) {
	loadTerraformConfigOnce.Do(func() {
		loadedTerraformConfig, loadedTerraformConfigErr = loadTerraformConfig()
	})
	return loadedTerraformConfig, loadedTerraformConfigErr
}

func loadTerraformConfig() (TerraformConfig, error) {
	cfg := TerraformConfig{
		Version:         os.Getenv(envTerraformVersion),
		ProviderSource:  os.Getenv(envTerraformProviderSource),
		ProviderVersion: os.Getenv(envTerraformProviderVersion),
	}
	if cfg.isComplete() {
		return cfg, nil
	}

	makeVars, err := loadTerraformVarsFromMakefile()
	if err != nil {
		return TerraformConfig{}, err
	}

	if cfg.Version == "" {
		cfg.Version = makeVars[envTerraformVersion]
	}
	if cfg.ProviderSource == "" {
		cfg.ProviderSource = makeVars[envTerraformProviderSource]
	}
	if cfg.ProviderVersion == "" {
		cfg.ProviderVersion = makeVars[envTerraformProviderVersion]
	}

	if !cfg.isComplete() {
		return TerraformConfig{}, fmt.Errorf("terraform test config is incomplete: version=%q providerSource=%q providerVersion=%q", cfg.Version, cfg.ProviderSource, cfg.ProviderVersion)
	}
	return cfg, nil
}

func (c TerraformConfig) isComplete() bool {
	return c.Version != "" && c.ProviderSource != "" && c.ProviderVersion != ""
}

func loadTerraformVarsFromMakefile() (map[string]string, error) {
	repoRoot, err := repositoryRoot()
	if err != nil {
		return nil, err
	}

	//nolint:gosec // The path is derived from the repository root of this checked-in package.
	file, err := os.Open(filepath.Join(repoRoot, "Makefile"))
	if err != nil {
		return nil, fmt.Errorf("cannot open Makefile: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	values := map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "export ") {
			continue
		}

		assignment := strings.TrimSpace(strings.TrimPrefix(line, "export "))
		name, value, ok := splitAssignment(assignment)
		if !ok {
			continue
		}
		switch name {
		case envTerraformVersion, envTerraformProviderSource, envTerraformProviderVersion:
			values[name] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cannot read Makefile: %w", err)
	}
	return values, nil
}

func splitAssignment(line string) (string, string, bool) {
	for _, separator := range []string{"?=", ":="} {
		name, value, ok := strings.Cut(line, separator)
		if ok {
			return strings.TrimSpace(name), strings.TrimSpace(value), true
		}
	}
	return "", "", false
}

func repositoryRoot() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to locate testconfig source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..")), nil
}

package lockfile

import (
	"fmt"
	"os"
	"regexp"

	"github.com/coreos/go-semver/semver"
)

type Lock struct {
	LockVersion    string   `json:"lockfile_version"`
	PackageName    string   `json:"package_name"`
	Meta           Metadata `json:"meta"`
	PackageVersion string   `json:"version"`
	SemverVersion  *semver.Version
	Sources        map[string]string           `json:"sources"`
	ContractTypes  map[string]ContractType     `json:"contract_type"`
	Deployments    map[string]ContractInstance `json:"deployments"`
	BuildDeps      map[string]string           `json:"build_dependencies"`
}

func (l Lock) validate() (err error) {
	if l.LockVersion != "1" {
		return fmt.Errorf("unknown lockfile version: %v", l.LockVersion)
	}
	pkgNameRegexp := regexp.MustCompile("[a-zA-Z][-a-zA-Z0-9_]*")
	if !pkgNameRegexp.Match([]byte(l.PackageName)) {
		return fmt.Errorf("invalid package name: %v\npackage names must comply to the regular expression: [a-zA-Z][-a-zA-Z0-9_]*", l.PackageName)
	}

	if err = l.Meta.validate(); err != nil {
		return err
	}

	if *l.SemverVersion != (semver.Version{}) {
		return fmt.Errorf("unexpected pre-initialized semver struct")
	}
	l.SemverVersion, err = semver.NewVersion(l.PackageVersion)
	if err != nil {
		return fmt.Errorf("unexpected error in parsing semver compliant package version: %v", err)
	}

	for m, _ := range l.Sources {
		if _, err = os.Stat(m); os.IsNotExist(err) {
			return err
		}
		// todo: Should I add a checksum on the source of the files to be hashed properly
	}

	//for m, v := range l.ContractTypes {

	//}

	//for m, v := range l.Deployments {

	//}

	//for m, v := range l.BuildDeps {

	//}
	return nil
}

type Metadata struct {
	Authors     []string `json:authors`
	License     string   `json:license`
	Description string   `json:description`
	Keywords    []string `json:keywords`
}

func (metadata Metadata) validate() (err error) {
	return nil
}

type ContractType struct {
	Name            string       `json:"contract_name"`
	Bytecode        string       `json:"bytecode"`
	RuntimeBytecode string       `json:"runtime_bytecode"`
	Abi             string       `json:"abi"`
	Natspec         string       `json:"natspec"`
	Compiler        CompilerInfo `json:"compiler"`
}

type ContractInstance struct {
	Type     string      `json:"contract_type"`
	Address  string      `json:"address"`
	Tx       string      `json:"transaction"`
	Block    string      `json:"block"`
	Runtime  string      `json:"runtime_bytecode"`
	LinkDeps []LinkValue `json:"link_dependencies"`
}

type LinkValue struct {
	Offset int    `json:"offset"`
	Value  string `json:"value"`
}

type CompilerInfo struct {
	Type     string           `json:"type"`
	Version  *semver.Version  `json:"version"`
	Settings CompilerSettings `json:"settings"`
}

func (compiler CompilerInfo) validate() (err error) {
	if compiler.Type != "solc" {
		return fmt.Errorf("invalid compiler type selected: %v", compiler.Type)
	}
	return nil
}

type CompilerSettings struct {
	Optimize     bool `json:"optimize"`
	OptimizeRuns int  `json:"optimize_runs"`
}

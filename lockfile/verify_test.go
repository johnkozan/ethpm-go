package lockfile

import (
	"testing"

	"github.com/coreos/go-semver/semver"
)

func TestValidateCompilerInfo(t *testing.T) {
	for i, test := range []struct {
		compilerInfo CompilerInfo
		err          string
	}{
		{
			compilerInfo: CompilerInfo{"serpent", semver.New("0.4.13"), CompilerSettings{true, 500}},
			err:          "invalid compiler type selected: serpent",
		},
		{
			compilerInfo: CompilerInfo{"LLL", semver.New("0.4.13"), CompilerSettings{true, 500}},
			err:          "invalid compiler type selected: LLL",
		},
		{
			compilerInfo: CompilerInfo{"solc", semver.New("0.4.13"), CompilerSettings{true, 500}},
			err:          "",
		},
		{
			compilerInfo: CompilerInfo{"solc", semver.New("nightly-0.4.14-f129372245d1b4fd4ff6425e9f7cbe701247cdc1"), CompilerSettings{true, 500}},
			err:          "",
		},
		{
			compilerInfo: CompilerInfo{"solcjs", semver.New("0.4.13"), CompilerSettings{true, 500}},
			err:          "invalid compiler type selected: solcjs",
		},
	} {
		err := test.compilerInfo.validate()
		if err != nil && len(test.err) == 0 {
			t.Errorf("%d failed. Expected no err but got: %v", i, err)
			continue
		}
		if err == nil && len(test.err) != 0 {
			t.Errorf("%d failed. Expected err: %v but got none", i, test.err)
			continue
		}
		if err != nil && len(test.err) != 0 && err.Error() != test.err {
			t.Errorf("%d failed. Expected err: '%v' got err: '%v'", i, test.err, err)
			continue
		}
	}
}

package variables

import (
	"io/ioutil"
	"strings"

	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/pkg/errors"
)

// ParseVariables parses the shell expressions of the form `export FOO=bar` or  `export FOO="bar"` in the given file if it exists
// and converts them to environment variables we can use to pass into kaniko or mink
func ParseVariables(path string) (map[string]string, error) {
	m := map[string]string{}
	exists, err := files.FileExists(path)
	if err != nil {
		return m, errors.Wrapf(err, "failed to check if path exists %s", path)
	}
	if !exists {
		return m, nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return m, errors.Wrapf(err, "failed to read file %s", path)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "export") {
			continue
		}

		rest := strings.TrimSpace(strings.TrimPrefix(line, "export"))
		idx := strings.Index(rest, "=")
		if idx <= 0 {
			continue
		}
		name := strings.TrimSpace(rest[0:idx])
		value := strings.TrimSpace(rest[idx+1:])

		for _, q := range []string{"\"", "'"} {
			if strings.HasPrefix(value, q) && strings.HasSuffix(value, q) {
				value = value[1 : len(value)-1]
			}
		}
		m[name] = value
	}
	return m, nil
}

package variables_test

import (
	"path/filepath"
	"testing"

	"github.com/jenkins-x-plugins/jx-mink/pkg/variables"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVariables(t *testing.T) {
	testCases := []struct {
		file     string
		expected map[string]string
	}{
		{
			file: "simple.sh",
			expected: map[string]string{
				"CHEESE": "edam",
				"BEER":   "stella",
				"WINE":   "merlot",
			},
		},
	}

	for _, tc := range testCases {
		name := tc.file
		path := filepath.Join("test_data", name)
		require.FileExists(t, path, "should have file for %s", name)

		m, err := variables.ParseVariables(path)
		require.NoError(t, err, "should not have failed to parse %s", name)
		assert.Equal(t, tc.expected, m, "variables for %s", name)
		t.Logf("test %s got variables: %v\n", name, m)
	}
}

package plywood

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatTester(t *testing.T) {
	f := CustomFormat{
		Regex:          `([\d-]* [\d:.]* [+-]?\d* [^ ]*)`,
		TimeFormat:     "2006-01-02 15:04:05.999999999 -0700 MST",
		ExampleLine:    "2020-10-19 19:19:17.497204 +1300 NZDT line1 ",
		ExpectedOutput: "2020-10-19 07:19:17.497 [test] line1 \n",
	}

	err := TestCustomFormat(f)
	require.NoError(t, err)
}

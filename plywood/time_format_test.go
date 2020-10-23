package plywood

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTimeExtractor(t *testing.T) {
	ex := timeExtractor{
		regex:  regexp.MustCompile(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`),
		layout: "2006-01-02 15:04:05.999999999 -0700 MST",
	}

	line := "2020-10-19 19:19:17.497204 +1300 NZDT line1 "
	now, rest, err := ex.Parse(line)
	require.NoError(t, err)
	require.Equal(t, "2020-10-19 19:19:17.497204 +1300 NZDT", now.String())
	require.Equal(t, rest, " line1 ")
}
func TestDefaultFormats(t *testing.T) {
	for _, format := range defaultFormats {
		err := TestCustomFormat(format)
		require.NoError(t, err)
	}
}

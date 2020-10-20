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

func TestFmt1(t *testing.T) {
	ex := timeExtractor{
		regex:  regexp.MustCompile(`(..... ..:..:..\.\d*) `),
		layout: "02Jan 15:04:05.999999999",
	}

	line := `18Oct 01:21:25.325 - line `
	now, rest, err := ex.Parse(line)
	require.NoError(t, err)
	require.Equal(t, "2020-10-18 01:21:25.325 +0000 UTC", now.String())
	require.Equal(t, "- line ", rest)

	line = `18Oct 01:21:25.325 ERROR - line `
	now, rest, err = ex.Parse(line)
	require.NoError(t, err)
	require.Equal(t, "2020-10-18 01:21:25.325 +0000 UTC", now.String())
	require.Equal(t, "ERROR - line ", rest)

}

func TestFmt2(t *testing.T) {
	ex := timeExtractor{
		regex:  regexp.MustCompile(`(..../../.. ..:..:..) `),
		layout: "2006/01/02 15:04:05",
	}

	line := `2020/10/18 01:21:22 INFO - line `
	now, rest, err := ex.Parse(line)
	require.NoError(t, err)
	require.Equal(t, "2020-10-18 01:21:22 +0000 UTC", now.String())
	require.Equal(t, "INFO - line ", rest)
}

func TestFmt3(t *testing.T) {
	ex := timeExtractor{
		regex:  regexp.MustCompile(`(..../../.. ..:..:..\.\d*) `),
		layout: "2006/01/02 15:04:05.999999",
	}

	line := `2020/10/18 01:21:22.123456 INFO - line `
	now, rest, err := ex.Parse(line)
	require.NoError(t, err)
	require.Equal(t, "2020-10-18 01:21:22.123456 +0000 UTC", now.String())
	require.Equal(t, "INFO - line ", rest)
}

func TestLogCat(t *testing.T) {
	ex := timeExtractor{
		regex:  regexp.MustCompile(`{([-\d]* ..:..:..\.\d*) `),
		layout: "01-02 15:04:05.999999",
	}

	line := `{10-20 12:01:15.531  2342  2767 line `
	now, rest, err := ex.Parse(line)
	require.NoError(t, err)
	require.Equal(t, "2020-10-20 12:01:15.531 +0000 UTC", now.String())
	require.Equal(t, " 2342  2767 line ", rest)
}

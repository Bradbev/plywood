package plywood

import (
	"io/ioutil"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var log1 = `2020-10-19 19:19:17.497204 +1300 NZDT line 1
2020-10-19 19:19:18.497204 +1300 NZDT line 2
2020-10-19 19:19:19.497204 +1300 NZDT line 3
2020-10-19 19:19:20.497204 +1300 NZDT line 4`

var log2 = `10-19-2020 19:19:17.497504 +1300 NZDT line 1a
10-19-2020 19:19:18.497504 +1300 NZDT line 2a
10-19-2020 19:19:19.497504 +1300 NZDT line 3a
10-19-2020 19:19:20.497504 +1300 NZDT line 4a`

func TestBasic(t *testing.T) {
	r1 := strings.NewReader(log1)
	p := &Plywood{}
	p.AddReader("1", r1)
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* .*?) `, "2006-01-02 15:04:05.999999999 -0700 MST")

	ply, err := ioutil.ReadAll(p)
	require.NoError(t, err)

	require.Equal(t, `2020-10-19 19:19:17.497204 +1300 NZDT [1] line 1
2020-10-19 19:19:18.497204 +1300 NZDT [1] line 2
2020-10-19 19:19:19.497204 +1300 NZDT [1] line 3
2020-10-19 19:19:20.497204 +1300 NZDT [1] line 4
`, string(ply))
}

func TestTwoWithBadFormat(t *testing.T) {
	r1 := strings.NewReader(log1)
	r2 := strings.NewReader(log2)
	p := &Plywood{}
	p.AddReader("1", r1)
	p.AddReader("2", r2)
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* .*?) `, "2006-01-02 15:04:05.999999999 -0700 MST")

	ply, err := ioutil.ReadAll(p)
	require.NoError(t, err)

	require.Equal(t, `z [2]10-19-2020 19:19:17.497504 +1300 NZDT line 1a
z [2]10-19-2020 19:19:18.497504 +1300 NZDT line 2a
z [2]10-19-2020 19:19:19.497504 +1300 NZDT line 3a
z [2]10-19-2020 19:19:20.497504 +1300 NZDT line 4a
2020-10-19 19:19:17.497204 +1300 NZDT [1] line 1
2020-10-19 19:19:18.497204 +1300 NZDT [1] line 2
2020-10-19 19:19:19.497204 +1300 NZDT [1] line 3
2020-10-19 19:19:20.497204 +1300 NZDT [1] line 4
`, string(ply))
}

func TestS(t *testing.T) {
	d, _ := time.ParseDuration("1h2m3.456s")
	require.Equal(t, "01:02:03:456", formatDuration(d))
}

func TestTwoWithGoodFormat(t *testing.T) {
	r1 := strings.NewReader(log1)
	r2 := strings.NewReader(log2)
	p := &Plywood{IncludeZeroBasis: true}
	p.AddReader("1", r1)
	p.AddReader("2", r2)
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* .*?) `, "2006-01-02 15:04:05.999999999 -0700 MST")
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* .*?) `, "01-02-2006 15:04:05.999999999 -0700 MST")

	ply, err := ioutil.ReadAll(p)
	require.NoError(t, err)

	require.Equal(t, `2020-10-19 19:19:17.497204 +1300 NZDT [00:00:00:0][1] line 1
2020-10-19 19:19:17.497504 +1300 NZDT [00:00:00:0][2] line 1a
2020-10-19 19:19:18.497204 +1300 NZDT [00:00:01:0][1] line 2
2020-10-19 19:19:18.497504 +1300 NZDT [00:00:01:0][2] line 2a
2020-10-19 19:19:19.497204 +1300 NZDT [00:00:02:0][1] line 3
2020-10-19 19:19:19.497504 +1300 NZDT [00:00:02:0][2] line 3a
2020-10-19 19:19:20.497204 +1300 NZDT [00:00:03:0][1] line 4
2020-10-19 19:19:20.497504 +1300 NZDT [00:00:03:0][2] line 4a
`, string(ply))
}

func TestTimeExtractor(t *testing.T) {
	ex := timeExtractor{
		regex:  regexp.MustCompile(`([\d-]* [\d:.]* [+-]?\d* .*?) `),
		layout: "2006-01-02 15:04:05.999999999 -0700 MST",
	}

	line := "2020-10-19 19:19:17.497204 +1300 NZDT line1 "
	now, rest, err := ex.Parse(line)
	require.NoError(t, err)
	require.Equal(t, "2020-10-19 19:19:17.497204 +1300 NZDT", now.String())
	require.Equal(t, rest, " line1 ")
}

func TestTimedLineReader(t *testing.T) {
	r := newTimedLineReader(strings.NewReader(log1))
	p := &Plywood{}
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* .*?) `, "2006-01-02 15:04:05.999999999 -0700 MST")
	r.prepare(p)

	line := func(timeStr, lineStr string) {
		t1, text1 := r.logTime(), r.logText()
		require.Equal(t, timeStr, t1.String())
		require.Equal(t, lineStr, text1)
	}
	scan := func() {
		ok, err := r.scan()
		require.True(t, ok)
		require.NoError(t, err)
	}
	line("2020-10-19 19:19:17.497204 +1300 NZDT", " line 1")
	// verify that it's OK to double scan
	line("2020-10-19 19:19:17.497204 +1300 NZDT", " line 1")

	scan()
	line("2020-10-19 19:19:18.497204 +1300 NZDT", " line 2")

	scan()
	line("2020-10-19 19:19:19.497204 +1300 NZDT", " line 3")

	scan()
	line("2020-10-19 19:19:20.497204 +1300 NZDT", " line 4")

	// past the last line
	ok, err := r.scan()
	require.False(t, ok)
	require.Error(t, err)
}

func TestTimedLineReaderNoMatchingFormat(t *testing.T) {
	r := newTimedLineReader(strings.NewReader(log1))
	p := &Plywood{}
	r.prepare(p)

	line := func(lineStr string) {
		_, text1 := r.logTime(), r.logText()
		require.Equal(t, lineStr, text1)
	}
	scan := func() {
		ok, err := r.scan()
		require.True(t, ok)
		require.NoError(t, err)
	}
	line("2020-10-19 19:19:17.497204 +1300 NZDT line 1")
	// verify that it's OK to double scan
	line("2020-10-19 19:19:17.497204 +1300 NZDT line 1")

	scan()
	line("2020-10-19 19:19:18.497204 +1300 NZDT line 2")

	scan()
	line("2020-10-19 19:19:19.497204 +1300 NZDT line 3")

	scan()
	line("2020-10-19 19:19:20.497204 +1300 NZDT line 4")

	// past the last line
	ok, err := r.scan()
	require.False(t, ok)
	require.Error(t, err)
}

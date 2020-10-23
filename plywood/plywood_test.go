package plywood

import (
	"io/ioutil"
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

var longLines = `10-19-2020 19:19:17.497504 +1300 NZDT line 1b
line 2b
line 3b`

func TestBasic(t *testing.T) {
	r1 := strings.NewReader(log1)
	p := &Plywood{IncludeRelativeTime: false, IncludeAbsoluteTime: true}
	p.AddReader("1", r1)
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`, "2006-01-02 15:04:05.999999999 -0700 MST")

	ply, err := ioutil.ReadAll(p)
	require.NoError(t, err)

	require.Equal(t, `2020-10-19 07:19:17.497 [1] line 1
2020-10-19 07:19:18.497 [1] line 2
2020-10-19 07:19:19.497 [1] line 3
2020-10-19 07:19:20.497 [1] line 4
`, string(ply))
}

func TestTwoWithBadFormat(t *testing.T) {
	r1 := strings.NewReader(log1)
	r2 := strings.NewReader(log2)
	p := &Plywood{IncludeRelativeTime: false, IncludeAbsoluteTime: true}
	p.AddReader("1", r1)
	p.AddReader("2", r2)
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`, "2006-01-02 15:04:05.999999999 -0700 MST")

	ply, err := ioutil.ReadAll(p)
	require.NoError(t, err)

	require.Equal(t, `z [2]10-19-2020 19:19:17.497504 +1300 NZDT line 1a
z [2]10-19-2020 19:19:18.497504 +1300 NZDT line 2a
z [2]10-19-2020 19:19:19.497504 +1300 NZDT line 3a
z [2]10-19-2020 19:19:20.497504 +1300 NZDT line 4a
2020-10-19 07:19:17.497 [1] line 1
2020-10-19 07:19:18.497 [1] line 2
2020-10-19 07:19:19.497 [1] line 3
2020-10-19 07:19:20.497 [1] line 4
`, string(ply))
}

func TestS(t *testing.T) {
	d, _ := time.ParseDuration("1h2m3.456s")
	require.Equal(t, "01:02:03:456", formatDuration(d))
}

func TestTwoWithGoodFormat(t *testing.T) {
	r1 := strings.NewReader(log1)
	r2 := strings.NewReader(log2)
	p := &Plywood{IncludeRelativeTime: true, IncludeAbsoluteTime: true}
	p.AddReader("1", r1)
	p.AddReader("2", r2)
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`, "2006-01-02 15:04:05.999999999 -0700 MST")
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`, "01-02-2006 15:04:05.999999999 -0700 MST")

	ply, err := ioutil.ReadAll(p)
	require.NoError(t, err)

	require.Equal(t, `2020-10-19 07:19:17.497 [00:00:00:000][1] line 1
2020-10-19 07:19:17.497 [00:00:00:000][2] line 1a
2020-10-19 07:19:18.497 [00:00:01:000][1] line 2
2020-10-19 07:19:18.497 [00:00:01:000][2] line 2a
2020-10-19 07:19:19.497 [00:00:02:000][1] line 3
2020-10-19 07:19:19.497 [00:00:02:000][2] line 3a
2020-10-19 07:19:20.497 [00:00:03:000][1] line 4
2020-10-19 07:19:20.497 [00:00:03:000][2] line 4a
`, string(ply))
}
func TestTwoWithGoodFormatBrokenLines(t *testing.T) {
	r1 := strings.NewReader(log1)
	r2 := strings.NewReader(longLines)
	p := &Plywood{IncludeRelativeTime: true, IncludeAbsoluteTime: true}
	p.AddReader("1", r1)
	p.AddReader("2", r2)
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`, "2006-01-02 15:04:05.999999999 -0700 MST")
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`, "01-02-2006 15:04:05.999999999 -0700 MST")

	ply, err := ioutil.ReadAll(p)
	require.NoError(t, err)

	require.Equal(t, `2020-10-19 07:19:17.497 [00:00:00:000][1] line 1
2020-10-19 07:19:17.497 [00:00:00:000][2] line 1b
2020-10-19 07:19:17.497 [00:00:00:000][2] line 2b
2020-10-19 07:19:17.497 [00:00:00:000][2] line 3b
2020-10-19 07:19:18.497 [00:00:01:000][1] line 2
2020-10-19 07:19:19.497 [00:00:02:000][1] line 3
2020-10-19 07:19:20.497 [00:00:03:000][1] line 4
`, string(ply))
}

func TestTimedLineReader(t *testing.T) {
	r := newTimedLineReader(strings.NewReader(log1))
	p := &Plywood{IncludeRelativeTime: false, IncludeAbsoluteTime: true}
	p.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`, "2006-01-02 15:04:05.999999999 -0700 MST")
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
	p := &Plywood{IncludeRelativeTime: false, IncludeAbsoluteTime: true}
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

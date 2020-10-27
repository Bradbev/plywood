package plywood

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"
)

type Plywood struct {
	IncludeRelativeTime bool
	IncludeAbsoluteTime bool
	readerNames         []string
	lineReaders         []*timedLineReader
	timeExtractors      []*timeExtractor
	activeReaders       int
	firstLineTime       time.Time
	ready               sync.Once
	reader              *io.PipeReader
	writer              *io.PipeWriter
}

func (p *Plywood) AddReader(readerName string, reader io.Reader) {
	p.readerNames = append(p.readerNames, readerName)
	p.lineReaders = append(p.lineReaders, newTimedLineReader(reader))
	p.activeReaders++
}

func (p *Plywood) AddTimeFormat(regex, layout string) {
	ex := timeExtractor{
		regex:  regexp.MustCompile(regex),
		layout: layout,
	}
	p.timeExtractors = append(p.timeExtractors, &ex)
}

func (p *Plywood) Read(buf []byte) (n int, err error) {
	p.ready.Do(func() {
		for _, reader := range p.lineReaders {
			reader.prepare(p)
		}
		p.reader, p.writer = io.Pipe()
		go p.readFromInputs()
	})
	return p.reader.Read(buf)
}

func (p *Plywood) readFromInputs() {
	for p.activeReaders > 0 {
		readerIndex := p.nextReader()
		reader := p.lineReaders[readerIndex]
		moreLines := true
		for moreLines {
			when := reader.logTime()
			line := reader.logText()

			p.formatLine(readerIndex, when, line)
			moreLines, _ = reader.scan()
			if !reader.active {
				p.activeReaders--
			}
			if moreLines {
				// if the time is the same, keep reading
				moreLines = when == reader.logTime()
			}
		}
	}
	p.writer.Close()
}

func (p *Plywood) formatLine(readerIndex int, when time.Time, line string) {
	name := p.readerNames[readerIndex]
	var logLine string
	if when.IsZero() {
		logLine = fmt.Sprintf("z [%v]%v\n", name, line)
	} else {
		if p.IncludeAbsoluteTime {
			logLine = fmt.Sprintf("%v ", when.Format("2006-01-02 03:04:05.000"))
		}
		if p.IncludeRelativeTime {
			if p.firstLineTime.IsZero() {
				p.firstLineTime = when
			}
			offset := formatDuration(when.Sub(p.firstLineTime))
			logLine += fmt.Sprintf("[%v]", offset)
		}
		logLine += fmt.Sprintf("[%v]%v\n", name, line)
	}
	p.writer.Write([]byte(logLine))
}

func formatDuration(d time.Duration) string {
	inMs := d.Milliseconds()
	ms := inMs % 1000
	sec := (inMs / 1000) % 60
	min := (inMs / (60 * 1000)) % 60
	hour := (inMs / (60 * 60 * 1000)) % 60

	return fmt.Sprintf("%02d:%02d:%02d:%03d", hour, min, sec, ms)
}

// nextReader finds the earliest logtime from all the readers
func (p *Plywood) nextReader() (readerIndex int) {
	// some time in the far future..
	minTime := time.Date(3000, 0, 0, 0, 0, 0, 0, time.UTC)
	for i, lineReader := range p.lineReaders {
		if lineReader.active {
			readerTime := lineReader.logTime()
			if readerTime.IsZero() || readerTime.Before(minTime) {
				minTime = lineReader.logTime()
				readerIndex = i
			}
		}
	}
	return
}

type timeExtractor struct {
	regex  *regexp.Regexp
	layout string
}

func (t *timeExtractor) Parse(line string) (time.Time, string, error) {
	if matches := t.regex.FindStringSubmatch(line); len(matches) > 1 {
		match := matches[1]
		now, err := time.Parse(t.layout, match)
		if err != nil {
			return time.Time{}, "", err
		}
		if now.Year() == 0 {
			now = now.AddDate(time.Now().Year(), 0, 0)
		}

		return now, line[len(matches[0]):], nil
	}
	return time.Time{}, "", fmt.Errorf("No regex match for %v", line)
}

type timedLineReader struct {
	scanner     *bufio.Scanner
	currentTime *time.Time
	currentLine string
	extractor   *timeExtractor
	active      bool
}

func newTimedLineReader(reader io.Reader) *timedLineReader {
	result := &timedLineReader{
		scanner: bufio.NewScanner(reader),
	}
	result.scanner.Split(bufio.ScanLines)
	result.active = true
	return result
}

func (t *timedLineReader) prepare(p *Plywood) {
	t.scanner.Scan()
	line := t.scanner.Text()
	for _, extractor := range p.timeExtractors {
		logTime, rest, err := extractor.Parse(line)
		if err == nil {
			t.extractor = extractor
			t.currentTime = &logTime
			t.currentLine = rest
			return
		}
	}
	t.currentLine = line
}

func (t *timedLineReader) logTime() time.Time {
	if t.extractor == nil {
		return time.Time{}
	}
	return *t.currentTime
}

func (t *timedLineReader) logText() string {
	return t.currentLine
}

func (t *timedLineReader) scan() (bool, error) {
	if t.scanner.Scan() {
		if t.extractor == nil {
			t.currentLine = t.scanner.Text()
			return true, nil
		}

		logTime, rest, err := t.extractor.Parse(t.scanner.Text())
		if err != nil {
			t.currentLine = " " + t.scanner.Text()
			return true, err
		}
		t.currentTime = &logTime
		t.currentLine = rest
		return true, nil
	}
	t.active = false
	return false, fmt.Errorf("No more lines to scan")
}

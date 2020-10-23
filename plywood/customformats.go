package plywood

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type CustomFormat struct {
	Regex          string
	TimeFormat     string
	ExampleLine    string
	ExpectedOutput string
}

func TestCustomFormat(format CustomFormat) error {
	regex, err := regexp.Compile(format.Regex)
	if err != nil {
		return err
	}
	ex := timeExtractor{
		regex:  regex,
		layout: format.TimeFormat,
	}

	_, _, err = ex.Parse(format.ExampleLine)
	if err != nil {
		return err
	}

	p := &Plywood{IncludeAbsoluteTime: true, IncludeRelativeTime: false}
	p.AddTimeFormat(format.Regex, format.TimeFormat)
	reader := strings.NewReader(format.ExampleLine)
	p.AddReader("test", reader)

	outb, err := ioutil.ReadAll(p)
	if err != nil {
		return err
	}
	out := string(outb)

	if out != format.ExpectedOutput {
		return fmt.Errorf("Expected line of\n'%v'\nDid not match\n'%v'\n", format.ExpectedOutput, out)
	}
	return nil
}

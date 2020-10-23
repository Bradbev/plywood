package plywood

import (
	"encoding/json"
	"fmt"
	"os"
)

var defaultFormats = []CustomFormat{
	{
		Regex:          `^(..../../.. ..:..:..) `,
		TimeFormat:     "2006/01/02 15:04:05",
		ExampleLine:    `2020/10/18 01:21:22 INFO - line `,
		ExpectedOutput: "2020-10-18 01:21:22.000 [test]INFO - line \n",
	},
	{
		Regex:          `^([\d-]* [\d:.]* [+-]?\d* [^ ]*)`,
		TimeFormat:     "2006-01-02 15:04:05.999999999 -0700 MST",
		ExampleLine:    "2020-10-19 19:19:17.497204 +1300 NZDT line1 ",
		ExpectedOutput: "2020-10-19 07:19:17.497 [test] line1 \n",
	},
	{
		Regex:          `^(..... ..:..:......) `,
		TimeFormat:     "02Jan 15:04:05.999999999",
		ExampleLine:    `18Oct 01:21:25.325 - line `,
		ExpectedOutput: "2020-10-18 01:21:25.325 [test]- line \n",
	},
	{
		Regex:          `^(...... ..:..:..) `,
		TimeFormat:     "Jan 02 15:04:05",
		ExampleLine:    `Oct 18 01:21:25 - line `,
		ExpectedOutput: "2020-10-18 01:21:25.000 [test]- line \n",
	},
	{
		Regex:          `^(..../../.. ..:..:..\.\d*) `,
		TimeFormat:     "2006/01/02 15:04:05.999999",
		ExampleLine:    `2020/10/18 01:21:22.123456 INFO - line `,
		ExpectedOutput: "2020-10-18 01:21:22.123 [test]INFO - line \n",
	},
	{
		Regex:          `^{([-\d]* ..:..:..\.\d*) `,
		TimeFormat:     "01-02 15:04:05.999999",
		ExampleLine:    `{10-20 12:01:15.531  2342  2767 line `,
		ExpectedOutput: "2020-10-20 12:01:15.531 [test] 2342  2767 line \n",
	},
}

func DefaultPlywood() *Plywood {
	return PlywoodFromFormatters(defaultFormats)
}

func PlywoodFromFormatters(formatters []CustomFormat) *Plywood {
	result := &Plywood{IncludeAbsoluteTime: false, IncludeRelativeTime: true}
	for _, format := range formatters {
		err := TestCustomFormat(format)
		if err != nil {
			fmt.Printf("Error in timestamp formats\n%v\n", err)
			os.Exit(1)
		}
		result.AddTimeFormat(format.Regex, format.TimeFormat)
	}
	return result
}

func PrintDefaults() {
	d, _ := json.MarshalIndent(defaultFormats, "", "  ")
	fmt.Println(string(d))
}

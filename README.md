# Plywood

Plywood is a log combiner. It accepts as input multiple files and produces a
single output stream on stdout. Each line of input should start with a
timestamp so that logs can be sorted. The first line of a log file determines
the timestamp format for the entire file. If a timestamp is not understood,
then that entire file will be at the beginning of output, and each line will
start with 'z'.  
If the environment variable PLYWOOD is set, then that file name will be loaded
to enable custom runtime formatters.

The -showDefaults flag will print the default formatters in the correct format
to be loaded.
See https://golang.org/pkg/time/#Parse for details on the time layout format.
When defining custom formatters you must provide a sample line and the expected
output.
Plywood will always format the output time in UTC.

### Flags

- -a Show absolute timestamp. Defaults to false
- -r Hide relative timestamp. Defaults to false
- -showDefaults Print the default timestamp formats

### Usage

`plywood one.log two.log > plywood.log`

## Example custom time format

```json
[
  {
    "Name": "Use the optional name key to give commentary on this format",
    "Regex": "([\\d-]* [\\d:.]* [+-]?\\d* [^ ]*)",
    "TimeFormat": "2006-01-02 15:04:05.999999999 -0700 MST",
    "ExampleLine": "2020-10-19 19:19:17.497204 +1300 NZDT line1 ",
    "ExpectedOutput": "2020-10-19 07:19:17.497 [test] line1 \n"
  }
]
```

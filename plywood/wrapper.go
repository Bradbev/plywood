package plywood

func DefaultPlywood() *Plywood {
	result := &Plywood{ExcludeOriginalTime: true}
	result.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* .*?) `, "2006-01-02 15:04:05.999999999 -0700 MST")
	return result
}

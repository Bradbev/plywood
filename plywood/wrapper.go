package plywood

func DefaultPlywood() *Plywood {
	result := &Plywood{IncludeAbsoluteTime: false, IncludeRelativeTime: true}
	result.AddTimeFormat(`(..../../.. ..:..:..) `, "2006/01/02 15:04:05")
	result.AddTimeFormat(`([\d-]* [\d:.]* [+-]?\d* [^ ]*)`, "2006-01-02 15:04:05.999999999 -0700 MST")
	result.AddTimeFormat(`(..... ..:..:......) `, "02Jan 15:04:05.999999999")
	result.AddTimeFormat(`(..../../.. ..:..:..\.\d*) `, "2006/01/02 15:04:05.999999")
	result.AddTimeFormat(`{([-\d]* ..:..:..\.\d*) `, "01-02 15:04:05.999999")
	return result
}

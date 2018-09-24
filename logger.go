package itbit

type logger interface {
	Debugf(string, ...interface{})
	Debug(a ...interface{})
	Printf(string, ...interface{})
}

var lg logger

// SetLogger sets a logger on the package that will print
// messages. Must have Printf and Debugf.
func SetLogger(l logger) {
	lg = l
}

func debugf(f string, a ...interface{}) {
	if lg == nil {
		return
	}
	lg.Debugf(f, a...)
}

func debug(a ...interface{}) {
	if lg == nil {
		return
	}
	lg.Debug(a...)
}

func logf(f string, a ...interface{}) {
	if lg == nil {
		return
	}
	lg.Printf(f, a...)
}

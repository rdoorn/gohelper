package logging

type Wrapper struct {
	Log    SimpleLogger
	Prefix []interface{}
	Level  Level
}

func (w *Wrapper) Println(i ...interface{}) {
	if w.Level > InfoLevel {
		return
	}
	i = w.addPrefix(i)
	w.Log.Infof(i...)
}

func (w *Wrapper) Debugf(i ...interface{}) {
	if w.Level > DebugLevel {
		return
	}
	i = w.addPrefix(i)
	w.Log.Debugf(i...)
}

func (w *Wrapper) Infof(i ...interface{}) {
	if w.Level > InfoLevel {
		return
	}
	i = w.addPrefix(i)
	w.Log.Infof(i...)
}

func (w *Wrapper) Warnf(i ...interface{}) {
	if w.Level > WarnLevel {
		return
	}
	i = w.addPrefix(i)
	w.Log.Warnf(i...)
}

func (w *Wrapper) Errorf(i ...interface{}) {
	if w.Level > ErrorLevel {
		return
	}
	i = w.addPrefix(i)
	w.Log.Errorf(i...)
}

func (w *Wrapper) Fatalf(i ...interface{}) {
	if w.Level > FatalLevel {
		return
	}
	i = w.addPrefix(i)
	w.Log.Fatalf(i...)
}

func (w *Wrapper) Panicf(i ...interface{}) {
	if w.Level > PanicLevel {
		return
	}
	i = w.addPrefix(i)
	w.Log.Panicf(i...)
}

func (w *Wrapper) addPrefix(i []interface{}) []interface{} {
	return append(i, w.Prefix...)
}

package telemetry

// create new arguments avoid concurrent map writes
func LargeArguments() Arguments {
	return Arguments{
		"large": "...",
	}
}

type Arguments map[string]interface{}

func (arg Arguments) ToMap() map[string]interface{} {
	return arg
}

func (arg Arguments) Truncate() {
	truncateArgs(arg)
}

func truncateArgs(arg interface{}) interface{} {
	switch x := arg.(type) {
	case Arguments:
		for k, v := range x {
			x[k] = truncateArgs(v)
		}
		return x
	case map[string]interface{}:
		for k, v := range x {
			x[k] = truncateArgs(v)
		}
		return x
	case []string:
		if len(x) > 3 {
			x = x[:3]
		}
		for i, v := range x {
			x[i] = truncateArgs(v).(string)
		}
		return x
	case []interface{}:
		if len(x) > 3 {
			x = x[:3]
		}
		for i, v := range x {
			x[i] = truncateArgs(v)
		}
		return x
	case string:
		if len(x) < 32 {
			return x
		}
		return x[:32] + "..."
	}
	return arg
}

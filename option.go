package eflag

type option uint64

const (
	WithEnv option = 1 << iota
	WithColor
)

func newOption(providedOptions ...option) option {
	var resultingOption option
	for _, opt := range providedOptions {
		resultingOption |= opt
	}
	return resultingOption
}

func (t option) isSet(providedOption option) bool {
	if t == option(0) {
		return false
	}
	return providedOption&t == t
}

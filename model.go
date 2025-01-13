package eflag

import (
	"errors"
	"unsafe"
)

const (
	efName  = "efName"
	efUsage = "efUsage"

	OptionEmpty option = 1 << iota
	OptionContinueOnUnknownKind
)

var (
	ErrKindIsNotSupported = errors.New("kind is not supported")
	ErrInvalidInput       = errors.New("invalid input")
)

type option uint64

func (t option) isSet(providedOptions option) bool {
	return providedOptions&t == t
}

type carrier struct {
	ptr     any
	value   any
	uptr    unsafe.Pointer
	efName  string
	efUsage string
}

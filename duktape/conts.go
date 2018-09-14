package duktape

import "fmt"

type Error struct {
	Type       string
	Message    string
	FileName   string
	LineNumber int
	Stack      string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

type Type int

const (
	CompileEval uint = 1 << iota
	CompileFunction
	CompileStrict
	CompileSafe
	CompileNoResult
	CompileNoSource
	CompileStrlen
)

const (
	TypeNone Type = iota
	TypeUndefined
	TypeNull
	TypeBoolean
	TypeNumber
	TypeString
	TypeObject
	TypeBuffer
	TypePointer
	TypeLightFunc
)

const (
	TypeMaskNone uint = 1 << iota
	TypeMaskUndefined
	TypeMaskNull
	TypeMaskBoolean
	TypeMaskNumber
	TypeMaskString
	TypeMaskObject
	TypeMaskBuffer
	TypeMaskPointer
	TypeMaskLightFunc
)

const (
	DUK_ENUM_INCLUDE_NONENUMERABLE uint = 1 << iota
	DUK_ENUM_INCLUDE_HIDDEN
	DUK_ENUM_INCLUDE_SYMBOLS
	DUK_ENUM_EXCLUDE_STRINGS
	DUK_ENUM_OWN_PROPERTIES_ONLY
	DUK_ENUM_ARRAY_INDICES_ONLY
	DUK_ENUM_SORT_ARRAY_INDICES
	DUK_ENUM_NO_PROXY_BEHAVIOR
)

const (
	DUK_DEFPROP_WRITABLE uint = 1 << iota
	DUK_DEFPROP_ENUMERABLE
	DUK_DEFPROP_CONFIGURABLE
	DUK_DEFPROP_HAVE_WRITABLE
	DUK_DEFPROP_HAVE_ENUMERABLE
	DUK_DEFPROP_HAVE_CONFIGURABLE
	DUK_DEFPROP_HAVE_VALUE
	DUK_DEFPROP_HAVE_GETTER
	DUK_DEFPROP_HAVE_SETTER
	DUK_DEFPROP_FORCE
)

const (
	ErrNone int = 0

	// Internal to Duktape
	ErrUnimplemented int = 50 + iota
	ErrUnsupported
	ErrInternal
	ErrAlloc
	ErrAssertion
	ErrAPI
	ErrUncaughtError
)

const (
	// Common prototypes
	ErrError int = 1 + iota
	ErrEval
	ErrRange
	ErrReference
	ErrSyntax
	ErrType
	ErrURI
)

const (
	// Returned error values
	ErrRetUnimplemented int = -(ErrUnimplemented + iota)
	ErrRetUnsupported
	ErrRetInternal
	ErrRetAlloc
	ErrRetAssertion
	ErrRetAPI
	ErrRetUncaughtError
)

const (
	ErrRetError int = -(ErrError + iota)
	ErrRetEval
	ErrRetRange
	ErrRetReference
	ErrRetSyntax
	ErrRetType
	ErrRetURI
)

const (
	ExecSuccess = iota
	ExecError
)

const (
	LogTrace int = iota
	LogDebug
	LogInfo
	LogWarn
	LogError
	LogFatal
)

const (
	BufobjDuktapeAuffer     = 0
	BufobjNodejsAuffer      = 1
	BufobjArraybuffer       = 2
	BufobjDataview          = 3
	BufobjInt8array         = 4
	BufobjUint8array        = 5
	BufobjUint8clampedarray = 6
	BufobjInt16array        = 7
	BufobjUint16array       = 8
	BufobjInt32array        = 9
	BufobjUint32array       = 10
	BufobjFloat32array      = 11
	BufobjFloat64array      = 12
)

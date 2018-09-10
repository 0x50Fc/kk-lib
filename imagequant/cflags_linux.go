//+build !windows

package imagequant

/*
#cgo CFLAGS: -O3  -Wunknown-pragmas -fomit-frame-pointer -Wall -Wno-attributes -std=c99 -DNDEBUG -DUSE_SSE=1 -msse
#cgo LDFLAGS: -lm
*/
import "C"

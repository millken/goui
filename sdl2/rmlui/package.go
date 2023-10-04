package rmlui

/*

#cgo CXXFLAGS: -mmacosx-version-min=10.4 -std=c++14 -I${SRCDIR} -I${SRCDIR}/Include/
// #cgo darwin LDFLAGS: -L${SRCDIR} -lc++ -framework Cocoa -framework Metal -lSDL2
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR} -lRmlCore_darwin_arm64 -lshell_darwin_arm64 -lRmlDebugger_darwin_arm64
#include <stdlib.h>
#include <string.h>
#include "rmlui.h"
*/
import "C"

//swig -go -cgo -c++ -intgosize 64 rmlui.swigcxx

func RmlGetVersion() string {
	s := C.RmlGetVersion()
	//defer C.free(unsafe.Pointer(s))
	return C.GoString(s)
}

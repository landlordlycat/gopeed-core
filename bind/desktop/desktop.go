// go build -ldflags="-w -s" -buildmode=c-shared -o bin/libgopeed.dll ./bind/desktop
package main

import "C"
import (
	"github.com/monkeyWie/gopeed-core/bind"
)

func main() {}

//export Start
func Start(ip *C.char, port int) *C.char {
	return C.CString(bind.Start(C.GoString(ip), port))
}

//export Stop
func Stop() *C.char {
	return C.CString(bind.Stop())
}

// go build -ldflags="-w -s" -buildmode=c-shared -o bin/libgopeed.dll ./bind/desktop
package main

import "C"
import (
	"github.com/monkeyWie/gopeed-core/bind"
)

func main() {

}

//export Listen
func Listen() *C.char {
	return C.CString(bind.Listen())
}

//export Create
func Create(url string, opts string) *C.char {
	if err := bind.Create(url, opts); err != nil {
		return C.CString(err.Error())
	}
	return nil
}

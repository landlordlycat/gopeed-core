// gomobile bind -ldflags="-w -s" -o bin/libgopeed.aar -target=android -javapkg=com.gopeed.core ./bind/mobile
package mobile

import "C"
import (
	"github.com/monkeyWie/gopeed-core/bind"
)

func Listen() string {
	return bind.Listen()
}

func Create(url string, opts string) error {
	return bind.Create(url, opts)
}

// gomobile bind -ldflags="-w -s" -o bin/libgopeed.aar -target=android -javapkg=com.gopeed.core ./bind/mobile
package mobile

import "C"
import (
	"github.com/monkeyWie/gopeed-core/bind"
)

func Start(ip string, port int) string {
	return bind.Start(ip, port)
}

func Stop() string {
	return bind.Stop()
}

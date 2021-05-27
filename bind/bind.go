package bind

import (
	"encoding/json"
	"github.com/monkeyWie/gopeed-core/pkg/rest"
)

func Start(ip string, port int) string {
	var r StartResult
	r.Port, r.Err = rest.Start(ip, port)
	return toJSON(r)
}

func Stop() string {
	if err := rest.Stop(); err != nil {
		return err.Error()
	}
	return ""
}

type StartResult struct {
	Port int   `json:"port"`
	Err  error `json:"err"`
}

func toJSON(v interface{}) string {
	buf, _ := json.Marshal(v)
	return string(buf)
}

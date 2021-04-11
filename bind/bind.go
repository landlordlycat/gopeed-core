package bind

import (
	"encoding/json"
	"github.com/monkeyWie/gopeed-core/pkg/base"
	"github.com/monkeyWie/gopeed-core/pkg/download"
	"sync"
)

var (
	l sync.Mutex
	q = make([]*download.Event, 0)
)

func init() {
	download.Boot().Listener(func(event *download.Event) {
		l.Lock()
		defer l.Unlock()
		q = append(q, event)
	})
}

func Listen() string {
	l.Lock()
	defer l.Unlock()
	if len(q) == 0 {
		return ""
	}
	event := q[0]
	buf, _ := json.Marshal(event)
	q = q[1:]
	return string(buf)
}

func Create(url string, opts string) error {
	t, err := toOptions(opts)
	if err != nil {
		return err
	}
	return download.Boot().
		URL(url).
		Create(t)
}

func toOptions(str string) (*base.Options, error) {
	var opts base.Options
	if err := json.Unmarshal([]byte(str), &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

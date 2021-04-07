package bind

import (
	"encoding/json"
	"github.com/monkeyWie/gopeed-core/pkg/base"
	"github.com/monkeyWie/gopeed-core/pkg/download"
)

var eventCh = make(chan *download.Event)

func init() {
	download.Boot().Listener(func(event *download.Event) {
		eventCh <- event
	})
}

func Listen() string {
	buf, _ := json.Marshal(<-eventCh)
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

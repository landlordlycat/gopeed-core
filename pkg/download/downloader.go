package download

import (
	"errors"
	"github.com/google/uuid"
	"github.com/monkeyWie/gopeed-core/internal/controller"
	"github.com/monkeyWie/gopeed-core/internal/fetcher"
	"github.com/monkeyWie/gopeed-core/internal/protocol/http"
	"github.com/monkeyWie/gopeed-core/pkg/base"
	"github.com/monkeyWie/gopeed-core/pkg/util"
	"net/url"
	"strings"
	"sync"
	"time"
)

type FetcherBuilder func() (protocols []string, builder func() fetcher.Fetcher)
type Listener func(event *Event)

type Task struct {
	ID       string         `json:"id"`
	Res      *base.Resource `json:"res"`
	Opts     *base.Options  `json:"opts"`
	Status   base.Status    `json:"status"`
	Progress *Progress      `json:"progress"`

	fetcher fetcher.Fetcher
	timer   *util.Timer
	locker  *sync.Mutex
}

type Progress struct {
	// 下载耗时(纳秒)
	Used int64 `json:"used"`
	// 每秒下载字节数
	Speed int64 `json:"speed"`
	// 已下载的字节数
	Downloaded int64 `json:"downloaded"`
}

type Downloader struct {
	*controller.DefaultController
	fetchBuilders map[string]func() fetcher.Fetcher
	tasks         map[string]*Task
	listener      Listener
	eventCh       chan *Event
}

func NewDownloader(fbs ...FetcherBuilder) *Downloader {
	d := &Downloader{DefaultController: controller.NewController()}
	d.fetchBuilders = make(map[string]func() fetcher.Fetcher)
	for _, f := range fbs {
		protocols, builder := f()
		for _, p := range protocols {
			d.fetchBuilders[strings.ToUpper(p)] = builder
		}
	}
	d.tasks = make(map[string]*Task)

	// 每秒统计一次下载速度
	go func() {
		for {
			if len(d.tasks) > 0 {
				for _, task := range d.tasks {
					if task.Status == base.DownloadStatusDone ||
						task.Status == base.DownloadStatusError ||
						task.Status == base.DownloadStatusPause {
						continue
					}
					current := task.fetcher.Progress().TotalDownloaded()
					task.Progress.Used = task.timer.Used()
					task.Progress.Speed = current - task.Progress.Downloaded
					task.Progress.Downloaded = current
					d.emit(EventKeyProgress, task)
				}
			}
			time.Sleep(time.Second)
		}
	}()
	return d
}

func (d *Downloader) buildFetcher(URL string) (fetcher.Fetcher, error) {
	url, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	if fetchBuilder, ok := d.fetchBuilders[strings.ToUpper(url.Scheme)]; ok {
		fetcher := fetchBuilder()
		fetcher.Setup(d.DefaultController)
		return fetcher, nil
	}
	return nil, errors.New("unsupported protocol")
}

func (d *Downloader) Resolve(req *base.Request) (*base.Resource, error) {
	fetcher, err := d.buildFetcher(req.URL)
	if err != nil {
		return nil, err
	}
	return fetcher.Resolve(req)
}

func (d *Downloader) Create(res *base.Resource, opts *base.Options) (err error) {
	fetcher, err := d.buildFetcher(res.Req.URL)
	if err != nil {
		return
	}
	if opts == nil {
		opts = &base.Options{}
	}
	if !res.Range || opts.Connections < 1 {
		opts.Connections = 1
	}
	err = fetcher.Create(res, opts)
	if err != nil {
		return
	}
	id := uuid.New().String()
	task := &Task{
		ID:       id,
		Res:      res,
		Opts:     opts,
		Status:   base.DownloadStatusReady,
		Progress: &Progress{},
		fetcher:  fetcher,
		timer:    &util.Timer{},
		locker:   new(sync.Mutex),
	}
	d.tasks[id] = task
	err = fetcher.Start()
	if err != nil {
		return
	}
	go func() {
		task.timer.Start()
		d.emit(EventKeyStart, task)
		err = fetcher.Wait()
		if err != nil {
			task.Status = base.DownloadStatusError
			d.emit(EventKeyError, task, err)
		} else {
			task.Progress.Used = task.timer.Used()
			if task.Res.TotalSize == 0 {
				task.Res.TotalSize = task.fetcher.Progress().TotalDownloaded()
			}
			used := task.Progress.Used / int64(time.Second)
			if used == 0 {
				used = 1
			}
			task.Progress.Speed = task.Res.TotalSize / used
			task.Progress.Downloaded = task.Res.TotalSize
			task.Status = base.DownloadStatusDone
			d.emit(EventKeyDone, task)
		}
		d.emit(EventKeyFinally, task, err)
	}()
	return
}

func (d *Downloader) Pause(id string) {
	task := d.tasks[id]
	task.locker.Lock()
	defer task.locker.Unlock()
	task.timer.Pause()
	task.fetcher.Pause()
	d.emit(EventKeyPause, task)
}

func (d *Downloader) Continue(id string) {
	task := d.tasks[id]
	task.locker.Lock()
	defer task.locker.Unlock()
	task.timer.Continue()
	task.fetcher.Continue()
	d.emit(EventKeyContinue, task)
}

func (d *Downloader) SetListener(fn Listener) {
	d.listener = fn
}

func (d *Downloader) emit(eventKey EventKey, task *Task, errs ...error) {
	if d.listener != nil {
		var err error
		if len(errs) > 0 {
			err = errs[0]
		}
		d.listener(&Event{
			Key:  eventKey,
			Task: task,
			Err:  err,
		})
	}
}

var defaultDownloader = NewDownloader(http.FetcherBuilder)

func Resolve(request *base.Request) (*base.Resource, error) {
	return defaultDownloader.Resolve(request)
}

func SetListener(listener Listener) {
	defaultDownloader.SetListener(listener)
}

func Create(request *base.Request, res *base.Resource, opts *base.Options) (err error) {
	if res == nil {
		res, err = defaultDownloader.Resolve(request)
		if err != nil {
			return err
		}
	}
	return defaultDownloader.Create(res, opts)
}

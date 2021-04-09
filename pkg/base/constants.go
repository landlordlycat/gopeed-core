package base

import "errors"

type Status string

const (
	DownloadStatusReady = "ready"
	DownloadStatusStart = "start"
	DownloadStatusPause = "pause"
	DownloadStatusError = "error"
	DownloadStatusDone  = "done"
)

const (
	HttpCodeOK             = 200
	HttpCodePartialContent = 206

	HttpHeaderRange              = "Range"
	HttpHeaderContentLength      = "Content-Length"
	HttpHeaderContentRange       = "Content-Range"
	HttpHeaderContentDisposition = "Content-Disposition"

	HttpHeaderRangeFormat = "bytes=%d-%d"
)

var (
	DeleteErr = errors.New("delete")
)

package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/monkeyWie/gopeed-core/pkg/base"
	"github.com/monkeyWie/gopeed-core/pkg/download"
	"github.com/monkeyWie/gopeed-core/pkg/rest/dto"
	"io/ioutil"
	"net"
	"net/http"
)

var listener net.Listener

func Start(ip string, port int) (int, error) {
	var err error
	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return 0, err
	}
	r := mux.NewRouter()
	addRouters(r)
	server := &http.Server{}
	server.Handler = r
	go func() {
		download.SetListener(func(event *download.Event) {

		})
		server.Serve(listener)
	}()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func Stop() error {
	return listener.Close()
}

func addRouters(r *mux.Router) {
	r.Methods("POST").
		Path("/tasks").
		HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			var d dto.CreateTask
			handleJSON(writer, request, &d, func() error {
				return download.Create(&base.Request{
					URL: d.URL,
				}, nil, d.Options)
			})
		})
}

func handleJSON(writer http.ResponseWriter, request *http.Request, v interface{}, handle func() error) {
	defer request.Body.Close()
	buf, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handleError(writer, err)
		return
	}
	if err = json.Unmarshal(buf, v); err != nil {
		handleError(writer, err)
		return
	}
	if err = handle(); err != nil {
		handleError(writer, err)
		return
	}
	return
}

func handleError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(500)
	writer.Write([]byte(fmt.Sprintf("request fail:%s", err.Error())))
}

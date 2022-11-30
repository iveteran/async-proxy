package receiver

import (
	"net/http"
	"net/http/httputil"

	"matrix.works/fmx-async-proxy/conf"
	"matrix.works/fmx-async-proxy/mq"
)

type RequestHandler func(http.ResponseWriter, *http.Request)

var mqProducer *mq.MqProducer

func proxyHandler() RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-REQUEST-ID")
		if requestId != "" {
			conf.Logger.Printf("%s [%s] X-REQUEST-ID: %s\n", r.URL.Path, r.Method, requestId)
		}

		backend := r.Header.Get("PROXY-BACKEND")
		if backend == "" {
			conf.Logger.Println("Can not get backend from request header")
			w.WriteHeader(421) // 421 Misdirected Request
			return
		}

		conf.Logger.Printf("%s [%s] redirect to %s\n", r.URL.Path, r.Method, backend)

		requestBytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			conf.Logger.Println("Dump request failed: %s", err)
			w.WriteHeader(421) // 421 Misdirected Request
			return
		}
		//fmt.Printf("dumped request: %+v\n", requestBytes)

		topic := r.Header.Get("proxy-topic")
		mqProducer.Enqueue(topic, backend, requestBytes)
	}
}

func SetupProxyHandlers() {
	mqProducer = mq.NewMqProducer()
	http.Handle("/", http.HandlerFunc(proxyHandler()))
}

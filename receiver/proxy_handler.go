package receiver

import (
	"net/http"
	"net/http/httputil"

	"matrix.works/async-proxy/logger"
	"matrix.works/async-proxy/mq"
)

type RequestHandler func(http.ResponseWriter, *http.Request)

var mqProducer *mq.MqProducer

func init() {
	mqProducer = mq.NewMqProducer()
}

func proxyHandler(rtBackend string) RequestHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		requestId := r.Header.Get("X-REQUEST-ID")
		if requestId != "" {
			logger.Logger.Printf("%s [%s] X-REQUEST-ID: %s\n", r.URL.Path, r.Method, requestId)
		}

		backend := r.Header.Get("PROXY-BACKEND")
		if backend == "" && rtBackend != "" {
			backend = rtBackend
		}
		if backend == "" {
			logger.Logger.Println("Can not get backend for this request")
			w.WriteHeader(421) // 421 Misdirected Request
			return
		}

		logger.Logger.Printf("%s [%s] redirect to %s\n", r.URL.Path, r.Method, backend)

		requestBytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			logger.Logger.Println("Dump request failed: %s", err)
			w.WriteHeader(421) // 421 Misdirected Request
			return
		}
		//fmt.Printf("dumped request: %+v\n", requestBytes)

		topic := r.Header.Get("proxy-topic")
		mqProducer.Enqueue(topic, backend, requestBytes)
	}
}

func SetupProxyHandlers(routeMap map[string]string) {
	backend := ""
	http.Handle("/", http.HandlerFunc(proxyHandler(backend)))

	for pathPrefix, backend := range routeMap {
		http.Handle(pathPrefix, http.HandlerFunc(proxyHandler(backend)))
		pathLen := len(pathPrefix)
		if pathPrefix[pathLen-1] == '/' {
			http.Handle(pathPrefix[:pathLen-1], http.HandlerFunc(proxyHandler(backend)))
		} else {
			http.Handle(pathPrefix+"/", http.HandlerFunc(proxyHandler(backend)))
		}
	}
}

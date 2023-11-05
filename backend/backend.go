package backend

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"matrix.works/async-proxy/logger"
	"matrix.works/async-proxy/messages"
)

func SendMessage(backend string, requestBytes []byte) error {

	code, rsp, err := sendRequest(backend, requestBytes)
	if err != nil {
		logger.Logger.Printf("[SendMessage] sendRequest error: %s\n", err.Error())
		return err
	}
	if code != 200 {
		err := errors.New("Remote server is not available")
		logger.Logger.Printf("[SendMessage] sendRequest error: %s\n", err.Error())
		return err
	}
	if rsp.Code != 0 {
		logger.Logger.Printf("[SendMessage] sendRequest code: %d, error: %s\n", rsp.Code, rsp.Msg)
		return errors.New(rsp.Msg)
	}

	return nil
}

func sendRequest(backend string, reqeustBytes []byte,
) (int, *messages.Response, error) {

	reader := bufio.NewReader(bytes.NewReader(reqeustBytes))
	req, err := http.ReadRequest(reader)
	beUrl, err := url.Parse(backend)
	url := req.URL
	url.Scheme = beUrl.Scheme
	url.Host = beUrl.Host
	req.RequestURI = ""
	req.Host = beUrl.Host
	req.URL = url

	fmt.Printf("request: %+v\n", req)
	//fmt.Printf("request.body: %+v\n", req.Body)

	client := &http.Client{Timeout: 0} // wait forever

	rsp, err := client.Do(req)
	if err != nil {
		logger.Logger.Printf("[sendRequest] Error: %s\n", err.Error())
		return -1, nil, err
	}
	defer rsp.Body.Close()

	fmt.Println("response status: ", rsp.Status)
	//fmt.Println("response Headers:", rsp.Header)
	//fmt.Println("response body: ", rsp.Body)
	if rsp.StatusCode != 200 {
		return rsp.StatusCode, nil, nil
	}

	result := &messages.Response{}
	err1 := json.NewDecoder(rsp.Body).Decode(result)
	if err1 != nil {
		logger.Logger.Println(err1)
	}
	//fmt.Printf(">> rsp body: %+v\n", result)

	return rsp.StatusCode, result, nil
}

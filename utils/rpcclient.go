package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type HttpClient struct {
	*http.Client
}

// NewHttpClient to get http client instance
func NewHttpClient() *HttpClient {
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}
	return &HttpClient{
		httpClient,
	}
}

func buildHttpServerAddress(protocol string, host string, port int) string {
	return fmt.Sprintf("%s://%s:%d", protocol, host, port)
}

func (client *HttpClient) RPCCall(
	method string,
	params interface{},
	rpcResponse interface{},
) (err error) {
	rpcProtocol := GetENV("RPC_PROTOCOL", "http")
	rpcHost := GetENV("RPC_HOST", "127.0.0.1")
	rpcPort, _ := strconv.Atoi(GetENV("RPC_PORT", "9334"))
	rpcEndpoint := buildHttpServerAddress(rpcProtocol, rpcHost, rpcPort)

	payload := map[string]interface{}{
		"method": method,
		"params": params,
		"id":     0,
	}
	payloadInBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := client.Post(rpcEndpoint, "application/json", bytes.NewBuffer(payloadInBytes))

	if err != nil {
		return err
	}
	respBody := resp.Body
	defer respBody.Close()

	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, rpcResponse)
	if err != nil {
		return err
	}
	return nil
}

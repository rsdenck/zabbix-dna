package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ZabbixClient struct {
	URL   string
	Token string
	HTTP  *http.Client
}

type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	ID      int         `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
	ID      int             `json:"id"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func NewClient(url, token string, timeout int) *ZabbixClient {
	return &ZabbixClient{
		URL:   url,
		Token: token,
		HTTP: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *ZabbixClient) Login(user, password string) error {
	params := map[string]string{
		"username": user,
		"password": password,
	}

	result, err := c.Call("user.login", params)
	if err != nil {
		return err
	}

	var token string
	if err := json.Unmarshal(result, &token); err != nil {
		return err
	}

	c.Token = token
	return nil
}

func (c *ZabbixClient) Call(method string, params interface{}) (json.RawMessage, error) {
	reqBody := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	// Only add auth if we have a token and the method is not public
	if c.Token != "" && method != "apiinfo.version" && method != "user.login" {
		reqBody.Auth = c.Token
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Post(c.URL, "application/json-rpc", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if rpcResp.Error != nil {
		return nil, fmt.Errorf("zabbix error %d: %s (%s)", rpcResp.Error.Code, rpcResp.Error.Message, rpcResp.Error.Data)
	}

	return rpcResp.Result, nil
}

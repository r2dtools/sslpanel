package agent

import (
	"backend/internal/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agentintegration"
)

const (
	refreshCommand             = "refresh"
	getVhodstsCommand          = "getVhostsCommand"
	getVhostCertificateCommand = "getVhostCertificate"
)

type Agent struct {
	token  string
	client *client
	logger logger.Logger
}

type Response struct {
	Status,
	Error string
	Data interface{}
}

type requestData struct {
	Token,
	Command string
	Data interface{}
}

func (a *Agent) Refresh() (*agentintegration.ServerData, error) {
	data, err := a.Request(refreshCommand, nil)

	if err != nil {
		return nil, err
	}

	var serverData agentintegration.ServerData
	err = mapstructure.Decode(data, &serverData)

	if err != nil {
		return nil, errors.New("invalid server agent data")
	}

	return &serverData, nil
}

func (a *Agent) GetVhosts() ([]agentintegration.VirtualHost, error) {
	data, err := a.Request(getVhodstsCommand, nil)

	if err != nil {
		return nil, err
	}

	var vhosts []agentintegration.VirtualHost
	err = mapstructure.Decode(data, &vhosts)

	if err != nil {
		return nil, errors.New("invalid vhosts data")
	}

	return vhosts, nil
}

func (a *Agent) GetVhostCertificate(vhsotName string) (*agentintegration.Certificate, error) {
	data, err := a.Request(getVhostCertificateCommand, map[string]string{
		"vhostName": vhsotName,
	})

	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, nil
	}

	var certificate agentintegration.Certificate
	err = mapstructure.Decode(data, &certificate)

	if err != nil {
		return nil, fmt.Errorf("invalid certificate data: %v", err)
	}

	return &certificate, nil
}

func (a *Agent) Request(command string, data interface{}) (interface{}, error) {
	reqData := requestData{
		Token:   a.token,
		Command: command,
		Data:    data,
	}
	reqByteData, err := json.Marshal(reqData)

	if err != nil {
		return nil, fmt.Errorf("could not encode data: %v", err)
	}

	respData, err := a.client.Request(reqByteData)

	if err != nil {
		return nil, err
	}

	a.logger.Debug("response from agent: " + string(respData))

	var resp Response

	if err = json.Unmarshal(respData, &resp); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}

	if resp.Status != "ok" {
		message := resp.Error

		if message == "" {
			message = "unknown error"
		}

		return nil, errors.New(message)
	}

	return resp.Data, nil
}

func NewAgent(ipv4, ipv6, token string, port int) (*Agent, error) {
	if ipv4 == "" && ipv6 == "" {
		return nil, errors.New("ipv4 or ipv6 address must be specified")
	}

	if port <= 0 {
		return nil, errors.New("invalid port")
	}

	if token == "" {
		return nil, errors.New("invalid token")
	}

	ip := ipv4

	if ip == "" {
		ip = ipv6
	}

	tcpClient := client{
		ip:      ip,
		port:    port,
		timeout: defaultTimeout,
	}

	return &Agent{
		token:  token,
		client: &tcpClient,
	}, nil
}

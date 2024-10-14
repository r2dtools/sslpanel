package remote

import (
	"backend/logger"
	"backend/models"
	"backend/tcp"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agentintegration"
)

// Agent manages remote server agents
type Agent struct {
	Server *models.Server
}

// Response from the server agent
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

// Refresh sends request to refresh the remote server agent data
func (a *Agent) Refresh() (*agentintegration.ServerData, error) {
	data, err := a.Request("refresh", nil)
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

// GetVhosts returns informations about server virtual hosts
func (a *Agent) GetVhosts() ([]agentintegration.VirtualHost, error) {
	data, err := a.Request("getVhosts", nil)

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

// GetVhostCertificate returns vhost certificate
func (a *Agent) GetVhostCertificate(vhsotName string) (*agentintegration.Certificate, error) {
	data, err := a.Request("getVhostCertificate", map[string]string{
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

func (a *Agent) getClient() (*tcp.Client, error) {
	ip := a.Server.Ipv4Address
	if ip == "" {
		ip = a.Server.Ipv6Address
	}
	if ip == "" {
		return nil, errors.New("server agent ip address is not specified")
	}

	port := a.Server.AgentPort
	if port == 0 {
		return nil, errors.New("server agent port is not specified")
	}

	tcpClient := tcp.Client{
		IP:   ip,
		Port: port,
	}

	return &tcpClient, nil
}

func (a *Agent) Request(command string, data interface{}) (interface{}, error) {
	tcpClient, err := a.getClient()
	if err != nil {
		return nil, err
	}

	token := a.Server.Token
	if token == "" {
		return nil, errors.New("server agent token is not specified")
	}

	reqData := requestData{
		Token:   token,
		Command: command,
		Data:    data,
	}
	reqByteData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("could not encode data: %v", err)
	}

	respData, err := tcpClient.Request(reqByteData)
	if err != nil {
		return nil, err
	}

	logger.Debug("response from agent: " + string(respData))

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

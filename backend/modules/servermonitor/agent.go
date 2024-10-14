package servermonitor

import (
	"backend/remote"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/r2dtools/agentintegration"
)

type agent struct {
	serverAgent *remote.Agent
}

func (a *agent) loadStatisticsData(requestData *agentintegration.ServerMonitorStatisticsRequestData) (*agentintegration.ServerMonitorStatisticsResponseData, error) {
	responseData, err := a.serverAgent.Request("servermonitor.loadStatisticsData", requestData)
	if err != nil {
		return nil, err
	}
	if responseData == nil {
		return nil, nil
	}

	var data agentintegration.ServerMonitorStatisticsResponseData
	err = mapstructure.Decode(responseData, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid response data: %v", err)
	}

	return &data, nil
}

func (a *agent) loadDiskStatisticsData(requestData *agentintegration.ServerMonitorStatisticsRequestData) (*agentintegration.ServerMonitorDiskResponseData, error) {
	responseData, err := a.serverAgent.Request("servermonitor.loadStatisticsData", requestData)
	if err != nil {
		return nil, err
	}
	if responseData == nil {
		return nil, nil
	}

	var data agentintegration.ServerMonitorDiskResponseData
	err = mapstructure.Decode(responseData, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid response data for disk statistics: %v", err)
	}

	return &data, nil
}

func (a *agent) loadNetworkStatisticsData(requestData *agentintegration.ServerMonitorStatisticsRequestData) (*agentintegration.ServerMonitorNetworkResponseData, error) {
	responseData, err := a.serverAgent.Request("servermonitor.loadStatisticsData", requestData)
	if err != nil {
		return nil, err
	}
	if responseData == nil {
		return nil, nil
	}

	var data agentintegration.ServerMonitorNetworkResponseData
	err = mapstructure.Decode(responseData, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid response data for network statistics: %v", err)
	}

	return &data, nil
}

func (a *agent) loadProcessStatisticsData(requestData *agentintegration.ServerMonitorStatisticsRequestData) (*agentintegration.ServerMonitorProcessResponseData, error) {
	responseData, err := a.serverAgent.Request("servermonitor.loadStatisticsData", requestData)
	if err != nil {
		return nil, err
	}
	if responseData == nil {
		return nil, nil
	}

	var data agentintegration.ServerMonitorProcessResponseData
	err = mapstructure.Decode(responseData, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid response data for processes statistics: %v", err)
	}

	return &data, nil
}

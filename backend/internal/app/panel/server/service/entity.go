package service

import (
	domainService "backend/internal/app/panel/domain/service"
	"time"
)

type Server struct {
	ID           int       `json:"id"`
	Guid         string    `json:"guid"`
	Name         string    `json:"name"`
	OsCode       string    `json:"os_code"`
	OsVersion    string    `json:"os_version"`
	Ipv4Address  string    `json:"ipv4_address"`
	Ipv6Address  string    `json:"ipv6_address"`
	AgentVersion string    `json:"agent_version"`
	AgentPort    int       `json:"agent_port"`
	IsActive     int       `json:"is_active"`
	IsRegistered int       `json:"is_registered"`
	AccountID    int       `json:"account_id"`
	CreatedAt    time.Time `json:"created_at"`
	Token        string    `json:"token"`
}

type ServerDetails struct {
	Server

	HostName       string                 `json:"hostname"`
	Os             string                 `json:"os"`
	PlatformFamily string                 `json:"platform_family"`
	KernelVersion  string                 `json:"kernal_version"`
	KernelArch     string                 `json:"kernal_arch"`
	Virtualization string                 `json:"virtualization"`
	Uptime         uint64                 `json:"uptime"`
	BootTime       uint64                 `json:"boottime"`
	Domains        []domainService.Domain `json:"domains"`
}

type NewServerRequest struct {
	Name        string `json:"name" validate:"nonzero"`
	Ipv4Address string `json:"ipv4_address"`
	Ipv6Address string `json:"ipv6_address"`
	AgentPort   int    `json:"agent_port" validate:"nonzero"`
	Token       string `json:"token" validate:"nonzero"`
	AccountID   int
}

type UpdateServerRequest struct {
	ID          int    `json:"id" validate:"nonzero"`
	Name        string `json:"name" validate:"nonzero"`
	Ipv4Address string `json:"ipv4_address"`
	Ipv6Address string `json:"ipv6_address"`
	AgentPort   int    `json:"agent_port" validate:"nonzero"`
	Token       string `json:"token" validate:"nonzero"`
}

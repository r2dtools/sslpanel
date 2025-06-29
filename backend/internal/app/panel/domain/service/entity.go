package service

type DomainSetting struct {
	ID           int    `json:"id"`
	SettingName  string `json:"settingname"`
	SettingValue string `json:"settingvalue"`
}

type DomainConfigRequest struct {
	ServerGuid string
	DomainName string
	WebServer  string `form:"webserver"`
}

type DomainRequest struct {
	ServerGuid string
	DomainName string
	WebServer  string `form:"webserver"`
}

type DomainCertificateRequest struct {
	ServerGuid string
	DomainName string
	WebServer  string `form:"webserver"`
}

type DomainSettingsRequest struct {
	ServerGuid string
	DomainName string
}

type ChangeDomainSettingRequest struct {
	ServerGuid   string
	DomainName   string
	SettingName  string `json:"settingname"`
	SettingValue string `json:"settingvalue"`
}

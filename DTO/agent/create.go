package agent

import "github.com/gabemanfroi/midgard/domain/models"

type CreateAgentDTO struct {
	ElasticsearchId   string            `json:"elasticsearchId"`
	ElasticsearchName string            `json:"elasticsearchName"`
	Name              string            `json:"name"`
	Ip                string            `json:"ip,omitempty"`
	UserId            uint8             `json:"UserId"`
	CompanyId         uint8             `json:"companyId"`
	GroupId           uint8             `json:"groupId,omitempty"`
	DeviceType        models.DeviceType `json:"deviceType"`
}

package agent

import "github.com/gabemanfroi/midgard/domain/models"

type UpdateAgentDTO struct {
	ElasticsearchId   *string            `json:"elasticsearchId,omitempty"`
	ElasticsearchName *string            `json:"elasticsearchName,omitempty"`
	Name              *string            `json:"name,omitempty"`
	Ip                *string            `json:"ip,omitempty"`
	UserId            *uint8             `json:"UserId,omitempty"`
	CompanyId         *uint8             `json:"companyId,omitempty"`
	GroupId           *uint8             `json:"groupId,omitempty"`
	DeviceType        *models.DeviceType `json:"deviceType,omitempty"`
}

package agent

import "github.com/gabemanfroi/midgard/domain/models"

type AgentGeneralData struct {
	Id         uint              `json:"id"`
	Alias      string            `json:"alias"`
	Name       string            `json:"name"`
	Ip         *string           `json:"ip,omitempty"`
	DeviceType models.DeviceType `json:"deviceType"`
}

type ReadAgentDTO struct {
	AgentGeneralData *AgentGeneralData `json:"generalData"`
}

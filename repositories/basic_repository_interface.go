package repositories

import (
	"github.com/gabemanfroi/showcaseme/domain/DTO/basic"
)

type IBasicRepository interface {
	Create(dto *basic.CreateBasicDTO) *basic.ReadBasicDTO
	GetAll() ([]*basic.ReadBasicDTO, error)
	GetById(id uint) (*basic.ReadBasicDTO, error)
	Delete(id uint) error
	Update(id uint, dto *basic.UpdateBasicDTO) (*basic.ReadBasicDTO, error)
}

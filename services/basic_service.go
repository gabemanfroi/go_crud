package services

import (
	"github.com/gabemanfroi/midgard/domain/DTO/basic"
	"github.com/gabemanfroi/midgard/domain/interfaces/repositories"
	"github.com/gabemanfroi/midgard/internal/utils"
	"github.com/golobby/container/v3"
)

type BasicService struct {
	repository repositories.IBasicRepository
}

func CreateBasicService() *BasicService { return &BasicService{repository: getBasicService()} }

func (service BasicService) Create(dto *basic.CreateBasicDTO) *basic.ReadBasicDTO {
	return service.repository.Create(dto)
}

func (service BasicService) GetAll() ([]*basic.ReadBasicDTO, error) {
	return service.repository.GetAll()
}

func (service BasicService) GetById(id uint) (*basic.ReadBasicDTO, error) {
	return service.repository.GetById(id)
}

func (service BasicService) Delete(id uint) error {
	return service.repository.Delete(id)
}

func (service BasicService) Update(id uint, dto *basic.UpdateBasicDTO) (*basic.ReadBasicDTO, error) {
	return service.repository.Update(id, dto)
}

func getBasicService() repositories.IBasicRepository {
	var injector repositories.IBasicRepository
	utils.Check(container.Resolve(&injector), "Error while retrieving BasicRepository instance")
	return injector
}

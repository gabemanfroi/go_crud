package controllers

import (
	"encoding/json"
	"github.com/gabemanfroi/showcaseme/domain/DTO/basic"
	"github.com/gabemanfroi/showcaseme/domain/interfaces/services"
	"github.com/gabemanfroi/showcaseme/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golobby/container/v3"
	"strconv"
)

type BasicController struct {
	service services.IBasicService
}

func (controller BasicController) Create(c *fiber.Ctx) error {
	var dto basic.CreateBasicDTO
	if err := c.BodyParser(&dto); err != nil {
		return err
	}

	return c.Status(200).JSON(controller.service.Create(&dto))
}

func (controller BasicController) GetAll(c *fiber.Ctx) error {
	basics, err := controller.service.GetAll()
	if err != nil {
		return err
	}
	utils.Check(json.NewEncoder(c).Encode(&basics), "failed to encode basics")
	return c.Status(200).JSON(basics)
}

func (controller BasicController) GetById(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	utils.Check(err, "failed to get basic id")
	b, err := controller.service.GetById(uint(id))
	if err != nil {
		return err
	}
	utils.Check(json.NewEncoder(c).Encode(&b), "failed to encode basic")
	return c.Status(200).JSON(b)
}

func (controller BasicController) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	utils.Check(err, "failed to get basic id")
	err = controller.service.Delete(uint(id))
	if err != nil {
		return err
	}

	return c.Status(204).JSON("user deleted")
}

func (controller BasicController) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	utils.Check(err, "failed to get basic id")
	var dto basic.UpdateBasicDTO
	if err := c.BodyParser(&dto); err != nil {
		return err
	}
	updatedBasic, err := controller.service.Update(uint(id), &dto)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(updatedBasic)
}

func CreateBasicController() *BasicController { return &BasicController{service: getBasicService()} }

func getBasicService() services.IBasicService {
	var injector services.IBasicService
	utils.Check(container.Resolve(&injector), "Error while retrieving BasicService instance ")
	return injector
}

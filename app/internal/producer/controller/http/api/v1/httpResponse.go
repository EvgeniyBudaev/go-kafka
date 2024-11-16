package v1

import (
	"github.com/gofiber/fiber/v2"
)

func ResponseError(ctx *fiber.Ctx, err error, httpStatusCode int) error {
	return ctx.Status(httpStatusCode).JSON(err.Error())
}

func ResponseOk(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(data)
}

func ResponseCreated(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(fiber.StatusCreated).JSON(data)
}

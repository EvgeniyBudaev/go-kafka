package controller

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/EvgeniyBudaev/kafka-go/app/internal/producer/controller/http/api/v1"
	"github.com/EvgeniyBudaev/kafka-go/app/internal/producer/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"net/http"
	"strconv"
)

type ProfileController struct {
	kafkaWriter *kafka.Writer
}

func NewProfileController(kw *kafka.Writer) *ProfileController {
	return &ProfileController{
		kafkaWriter: kw,
	}
}

func (pc *ProfileController) AddLike() fiber.Handler {
	return func(ctf *fiber.Ctx) error {
		fmt.Println("POST /gateway/api/v1/profiles/likes")
		hc := &entity.HubContent{
			Message: "Hello, World",
			UserId:  1,
		}
		hubContentJson, err := json.Marshal(hc)
		if err != nil {
			return v1.ResponseError(ctf, err, http.StatusInternalServerError)
		}
		err = pc.kafkaWriter.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(strconv.FormatUint(1, 10)),
				Value: hubContentJson,
			},
		)
		if err != nil {
			fmt.Println("WriteMessages err: ", err)
			return v1.ResponseError(ctf, err, http.StatusInternalServerError)
		}
		return v1.ResponseCreated(ctf, "OK")
	}
}

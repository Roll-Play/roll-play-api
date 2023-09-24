package handler

import (
	"net/http"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/labstack/echo/v4"
)

func GenerateNPCHandler(context echo.Context) error {
	npc := entities.NewNPC()

	return context.JSON(http.StatusOK, npc)
}

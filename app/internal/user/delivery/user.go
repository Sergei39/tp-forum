package delivery

import (
	"net/http"

	userModel "github.com/forums/app/internal/user"
	"github.com/forums/app/models"
	"github.com/forums/utils/errors"
	"github.com/forums/utils/logger"
	"github.com/labstack/echo/v4"
)

type handler struct {
	userUsecase userModel.UserUsecase
}

func NewUserHandler(usecase userModel.UserUsecase) userModel.UserHandler {
	return &handler{
		userUsecase: usecase,
	}
}

func (h *handler) CreateUser(c echo.Context) error {
	ctx := models.GetContext(c)

	nickname := c.Param("nickname")
	newUser := new(models.User)
	if err := c.Bind(newUser); err != nil {
		sendErr := errors.New(http.StatusBadRequest, err.Error())
		logger.Delivery().Error(ctx, sendErr)
		return c.NoContent(sendErr.Code())
	}
	newUser.Nickname = nickname

	response, err := h.userUsecase.CreateUser(ctx, newUser)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *handler) GetUser(c echo.Context) error {
	ctx := models.GetContext(c)

	nickname := c.Param("nickname")

	response, err := h.userUsecase.GetUserByName(ctx, nickname)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(response.Code(), response.Body())
}

func (h *handler) UpdateUser(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

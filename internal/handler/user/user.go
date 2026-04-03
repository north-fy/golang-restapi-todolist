package user

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/pkg/write"
	"github.com/sirupsen/logrus"
)

const op = "handler/user/user"

type ServiceUser interface {
	CreateUser(ctx context.Context, firstName, lastName, numberPhone string) (int, error)
	GetUser(ctx context.Context, id int) (models.User, error)
	GetUsersWithPagination(ctx context.Context, pt models.Pagination) ([]models.User, error)
	EditInfoUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type HandlerUser struct {
	log     *logrus.Logger
	service ServiceUser
}

func NewHandlerUser(log *logrus.Logger, user ServiceUser) *HandlerUser {
	return &HandlerUser{
		log:     log,
		service: user,
	}
}

func (h *HandlerUser) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestUser := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		h.log.Errorf("%s: %s", op, models.ErrBadRequest.Error())
		http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.CreateUser(ctx, requestUser.FirstName, requestUser.LastName, requestUser.NumberPhone)
	if err != nil {
		if errors.As(err, &models.ErrTargetExist) || errors.As(err, &models.ErrBadRequest) {
			h.log.Errorf("%s: %s", op, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, models.ErrInternal.Error(), http.StatusInternalServerError)
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *HandlerUser) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/users/"))
	if err != nil {
		http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	userModel, err := h.service.GetUser(ctx, userID)
	if err != nil {
		if errors.As(err, &models.ErrInvalidID) {
			http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
			h.log.Errorf("%s: %s", op, err.Error())
			return
		}

		http.Error(w, models.ErrInternal.Error(), http.StatusInternalServerError)
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, userModel)
}

func (h *HandlerUser) HandleGetUsersWithPagination(w http.ResponseWriter, r *http.Request) {
	// TODO: переделать функцию, так не должно работать!
	ctx := r.Context()

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	if limitStr == "" || offsetStr == "" {
		h.log.Errorf("%s: %s", op, models.ErrInvalidLimitOffset)
		http.Error(w, models.ErrInvalidLimitOffset.Error(), http.StatusBadRequest)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	pt := models.Pagination{
		Limit:  limit,
		Offset: offset,
	}

	users, err := h.service.GetUsersWithPagination(ctx, pt)
	if err != nil {
		h.log.Errorf("%s: %s", op, err.Error())
		http.Error(w, models.ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	write.WriteJSON(w, http.StatusOK, []any{pt, users})
}

func (h *HandlerUser) HandleEditUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestUser := models.User{}
	if err := json.NewDecoder(r.Body).Decode(&requestUser); err != nil {
		http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	requestUser.ID = userID

	if err = h.service.EditInfoUser(ctx, requestUser); err != nil {
		if models.IsErrValidate(err) || errors.As(err, &models.ErrInvalidID) {
			http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
			h.log.Errorf("%s: %s", op, err.Error())
			return
		}

		http.Error(w, models.ErrInternal.Error(), http.StatusInternalServerError)
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, nil)
}

func (h *HandlerUser) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	if err = h.service.DeleteUser(ctx, userID); err != nil {
		if errors.As(err, &models.ErrInvalidID) {
			http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
			h.log.Errorf("%s: %s", op, err.Error())
			return
		}

		http.Error(w, models.ErrInternal.Error(), http.StatusInternalServerError)
		h.log.Errorf("%s: %s", op, err.Error())
		return
	}

	write.WriteJSON(w, http.StatusOK, nil)
}

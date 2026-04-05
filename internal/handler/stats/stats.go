package stats

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/pkg/write"
	"github.com/sirupsen/logrus"
)

const op = "handler/stats/"

type ServiceStats interface {
	GetStatisticsAll(ctx context.Context) (models.ResponseStats, error)
	GetStatisticsByUser(ctx context.Context, id int) (models.ResponseStats, error)
	GetStatisticsByTime(ctx context.Context, timeFlt models.FilterTime) (models.ResponseStats, error)
}

type HandlerStats struct {
	log     *logrus.Logger
	service ServiceStats
}

func NewHandlerStats(log *logrus.Logger, serv ServiceStats) *HandlerStats {
	return &HandlerStats{
		log:     log,
		service: serv,
	}
}

func (h *HandlerStats) HandleGetStatistics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	model := r.URL.Query().Get("filter")
	switch model {
	case "all":
		resp, err := h.service.GetStatisticsAll(ctx)
		if err != nil {
			h.log.Error(err)
			http.Error(w, models.ErrInternal.Error(), http.StatusInternalServerError)
			return
		}

		write.WriteJSON(w, http.StatusOK, resp)

	case "user":
		id, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			h.log.Errorf("%s: %s", op, err.Error())
			http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		resp, err := h.service.GetStatisticsByUser(ctx, id)
		if err != nil {
			if errors.As(err, &models.ErrInvalidID) {
				h.log.Error(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			h.log.Error(err)
			http.Error(w, models.ErrInternal.Error(), http.StatusInternalServerError)
			return
		}

		write.WriteJSON(w, http.StatusOK, resp)

	case "time":
		var err error
		filter := models.FilterTime{}

		if startStr := r.URL.Query().Get("start_time"); startStr != "" {
			filter.Start, err = time.Parse(time.RFC3339, startStr)
			if err != nil {
				h.log.Errorf("%s: %s", op, err.Error())
				http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
				return
			}
		} else {
			h.log.Errorf("%s: empty field 'start_time' with time statistics", op)
			http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
			return
		}

		if endStr := r.URL.Query().Get("end_time"); endStr != "" {
			filter.End, err = time.Parse(time.RFC3339, endStr)
			if err != nil {
				h.log.Errorf("%s: %s", op, err.Error())
				http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
				return
			}
		} else {
			filter.End, err = time.Parse(time.RFC3339, time.Now().String())
			if err != nil {
				h.log.Errorf("%s: %s", op, err.Error())
				http.Error(w, models.ErrBadRequest.Error(), http.StatusBadRequest)
				return
			}
		}

		resp, err := h.service.GetStatisticsByTime(ctx, filter)
		if err != nil {
			h.log.Error(err)
			http.Error(w, models.ErrInternal.Error(), http.StatusInternalServerError)
			return
		}

		write.WriteJSON(w, http.StatusOK, resp)
	}
}

package stats

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/north-fy/golang-restapi-todolist/internal/domain/models"
	"github.com/north-fy/golang-restapi-todolist/internal/storage/redis"
	"github.com/sirupsen/logrus"
)

/*
type ResponseStats struct {
	Count           int
	CountCompleted  int
	PerCompleted    float64
	AvgTimeComplete time.Time
}

type ServiceStats interface {
	GetStatisticsAll(ctx context.Context) (ResponseStats, error)
	GetStatisticsByUser(ctx context.Context, id int) (ResponseStats, error)
	GetStatisticsByTime(ctx context.Context, time models.FilterTime) (ResponseStats, error)
}
*/

const op = "service/stats/"

type StorageStats interface {
	SelectAllTasks(ctx context.Context) ([]models.Task, error)
	SelectTasksByUser(ctx context.Context, id int) ([]models.Task, error)
	SelectTasksByTime(ctx context.Context, time models.FilterTime) ([]models.Task, error)
}

type ServiceStats struct {
	log     *logrus.Logger
	storage StorageStats
	cache   *redis.Storage
}

func NewServiceStats(log *logrus.Logger, storage StorageStats, cache *redis.Storage) *ServiceStats {
	return &ServiceStats{
		log:     log,
		storage: storage,
	}
}

func (s *ServiceStats) GetStatisticsAll(ctx context.Context) (models.ResponseStats, error) {
	resp := models.ResponseStats{}
	if err := s.cache.Get(ctx, models.KeyFormatALl, &resp); err != nil {
		s.log.Error(err)

		tasks, err := s.storage.SelectAllTasks(ctx)
		if err != nil {
			s.log.Error(err)
			return models.ResponseStats{}, err
		}

		// Main logic generate resp model
		for _, task := range tasks {
			resp.Count++

			// Count completed tasks and calculate average time complete
			if task.Completed {
				resp.CountCompleted++

				resp.AvgTimeComplete += task.CreatedAt.Sub(task.CompletedAt)
			}

		}

		resp.PerCompleted = math.Round(float64(resp.CountCompleted) / float64(resp.Count) * 100)

		// TODO: check error if time.Duration is not variable
		resp.AvgTimeComplete = resp.AvgTimeComplete / time.Duration(resp.Count)

		if err = s.cache.Set(ctx, models.KeyFormatALl, resp); err != nil {
			s.log.Error(err)
			return models.ResponseStats{}, err
		}
	}

	return resp, nil
}

func (s *ServiceStats) GetStatisticsByUser(ctx context.Context, id int) (models.ResponseStats, error) {
	resp := models.ResponseStats{}
	if err := s.cache.Get(ctx, fmt.Sprintf("%s%d", models.KeyFormatUser, id), &resp); err != nil {
		s.log.Error(err)

		tasks, err := s.storage.SelectTasksByUser(ctx, id)
		if err != nil {
			s.log.Error(err)
			return models.ResponseStats{}, err
		}

		// Main logic generate resp model
		for _, task := range tasks {
			resp.Count++

			// Count completed tasks and calculate average time complete
			if task.Completed {
				resp.CountCompleted++

				resp.AvgTimeComplete += task.CreatedAt.Sub(task.CompletedAt)
			}

		}

		resp.PerCompleted = math.Round(float64(resp.CountCompleted) / float64(resp.Count) * 100)

		// TODO: check error if time.Duration is not variable
		resp.AvgTimeComplete = resp.AvgTimeComplete / time.Duration(resp.Count)

		if err = s.cache.Set(ctx, fmt.Sprintf("%s%d", models.KeyFormatUser, id), resp); err != nil {
			s.log.Error(err)
			return models.ResponseStats{}, err
		}
	}

	return resp, nil
}

func (s *ServiceStats) GetStatisticsByTime(ctx context.Context, timeFlt models.FilterTime) (models.ResponseStats, error) {
	resp := models.ResponseStats{}
	if err := s.cache.Get(ctx, fmt.Sprintf("%s%s-%s", models.KeyFormatTime, timeFlt.Start, timeFlt.End), &resp); err != nil {
		s.log.Error(err)

		tasks, err := s.storage.SelectTasksByTime(ctx, timeFlt)
		if err != nil {
			s.log.Error(err)
			return models.ResponseStats{}, err
		}

		// Main logic generate resp model
		for _, task := range tasks {
			resp.Count++

			// Count completed tasks and calculate average time complete
			if task.Completed {
				resp.CountCompleted++

				resp.AvgTimeComplete += task.CreatedAt.Sub(task.CompletedAt)
			}

		}

		resp.PerCompleted = math.Round(float64(resp.CountCompleted) / float64(resp.Count) * 100)

		// TODO: check error if time.Duration is not variable
		resp.AvgTimeComplete = resp.AvgTimeComplete / time.Duration(resp.Count)

		if err = s.cache.Set(ctx, fmt.Sprintf("%s%s-%s", models.KeyFormatTime, timeFlt.Start, timeFlt.End), resp); err != nil {
			s.log.Error(err)
			return models.ResponseStats{}, err
		}
	}

	return resp, nil
}

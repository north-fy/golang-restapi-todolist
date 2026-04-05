package restapi

import (
	"net/http"

	"github.com/north-fy/golang-restapi-todolist/internal/handler/stats"
	taskhandler "github.com/north-fy/golang-restapi-todolist/internal/handler/task"
	userhandler "github.com/north-fy/golang-restapi-todolist/internal/handler/user"
)

type Route struct {
	mu *http.ServeMux
}

func newRouter() *Route {
	return &Route{
		mu: http.NewServeMux(),
	}
}

func (rt *Route) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.mu.ServeHTTP(w, r)
}

func (r *Route) ConfigureRouter(user *userhandler.HandlerUser, task *taskhandler.HandlerTask, stats *stats.HandlerStats) {
	// endpoints users
	r.mu.Handle("POST /api/users", MethodHandler("POST", user.HandleCreateUser))
	r.mu.Handle("GET /api/users/{id}", MethodHandler("GET", user.HandleGetUser))
	r.mu.Handle("GET /api/users", MethodHandler("GET", user.HandleGetUsersWithPagination))
	r.mu.Handle("PATCH /api/users/{id}", MethodHandler("PATCH", user.HandleEditUser))
	r.mu.Handle("DELETE /api/users/{id}", MethodHandler("DELETE", user.HandleDeleteUser))

	r.mu.Handle("GET /api/users/{id}/tasks", MethodHandler("GET", task.HandleGetTasksByUserID))

	// endpoints task
	r.mu.Handle("POST /api/tasks", MethodHandler("POST", task.HandleCreateTask))
	r.mu.Handle("GET /api/tasks/{id}", MethodHandler("GET", task.HandleGetTaskByID))
	r.mu.Handle("GET /api/tasks", MethodHandler("GET", task.HandleGetPaginationTasks))
	r.mu.Handle("PATCH /api/tasks/{id}", MethodHandler("PATCH", task.HandleEditTask))
	r.mu.Handle("DELETE /api/tasks/{id}", MethodHandler("DELETE", task.HandleDeleteTask))

	// endpoint statistics
	r.mu.Handle("GET /api/statistics", MethodHandler("GET", stats.HandleGetStatistics))
}

func MethodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}

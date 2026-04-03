package restapi

import (
	"net/http"

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

func (r *Route) ConfigureRouter(user *userhandler.HandlerUser, task *taskhandler.HandlerTask) {
	// endpoints users
	r.mu.Handle("POST /users", MethodHandler("POST", user.HandleCreateUser))
	r.mu.Handle("GET /users/{id}", MethodHandler("GET", user.HandleGetUser))
	r.mu.Handle("GET /users", MethodHandler("GET", user.HandleGetUsersWithPagination))
	r.mu.Handle("PATCH /users/{id}", MethodHandler("PATCH", user.HandleEditUser))
	r.mu.Handle("DELETE /users/{id}", MethodHandler("DELETE", user.HandleDeleteUser))

	r.mu.Handle("GET /users/{id}/tasks", MethodHandler("GET", task.HandleGetTasksByUserID))

	// endpoints task
	r.mu.Handle("POST /tasks", MethodHandler("POST", task.HandleCreateTask))
	r.mu.Handle("GET /tasks/{id}", MethodHandler("GET", task.HandleGetTaskByID))
	r.mu.Handle("GET /tasks", MethodHandler("GET", task.HandleGetPaginationTasks))
	r.mu.Handle("PATCH /tasks/{id}", MethodHandler("PATCH", task.HandleEditTask))
	r.mu.Handle("DELETE /tasks/{id}", MethodHandler("DELETE", task.HandleDeleteTask))
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

package restapi

import (
	"net/http"

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

func (r *Route) ConfigureRouter(user *userhandler.HandlerUser) {
	r.mu.Handle("POST /users", MethodHandler("POST", user.HandleCreateUser))
	r.mu.Handle("GET /users/{id}", MethodHandler("GET", user.HandleGetUser))
	r.mu.Handle("GET /users/tasks/{id}", MethodHandler("GET", user.HandleGetTasks))
	r.mu.Handle("PATCH /users/{id}", MethodHandler("PATCH", user.HandleEditUser))
	r.mu.Handle("DELETE /users/{id}", MethodHandler("DELETE", user.HandleDeleteUser))
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

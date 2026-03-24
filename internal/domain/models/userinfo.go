package models

type UserInfo struct {
	user  User
	tasks []Task
}

func (u UserInfo) GetUser() User {
	return u.user
}

func (u UserInfo) GetTasks() []Task {
	return u.tasks
}

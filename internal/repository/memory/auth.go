package memory

import (
	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/alipourhabibi/urlshortener/internal/core/messages"
)

type MemUsers struct {
	users []entity.User
}

func NewMemUser() *MemUsers {
	return &MemUsers{
		users: []entity.User{},
	}
}

func (m *MemUsers) Add(u entity.User) error {
	m.users = append(m.users, u)
	return nil
}
func (m *MemUsers) Exists(username string) (bool, error) {
	for _, v := range m.users {
		if v.Username == username {
			return true, nil
		}
	}
	return false, nil
}
func (m *MemUsers) Get(username string) (entity.User, error) {
	for _, v := range m.users {
		if v.Username == username {
			return v, nil
		}
	}
	return entity.User{}, messages.ErrNotFound
}

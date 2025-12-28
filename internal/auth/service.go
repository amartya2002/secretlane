package auth

import "errors"

type User struct {
	ID       int
	Username string
	Password string
}

type AuthService struct {
	repo *Repository
}

func NewAuthService() *AuthService {
	return &AuthService{repo: NewDefaultRepository()}
}

func (s *AuthService) Authenticate(username, password string) (*User, error) {
	u, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if u.Password != password {
		return nil, errors.New("invalid username or password")
	}

	return u, nil
}

// Signup creates a new user with the given username and password.
func (s *AuthService) Signup(username, password string) (*User, error) {
	exists, err := s.repo.UserExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	return s.repo.CreateUser(username, password)
}

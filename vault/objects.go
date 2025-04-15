package vault

import "fmt"

// Write wrapper for the vault client

// SimplePassword is a struct that holds a username and password
type SimplePassword struct {
	Username string
	Password string
}

func (s *SimplePassword) Serialize() (map[string]any, error) {
	return map[string]any{
		"username": s.Username,
		"password": s.Password,
	}, nil
}

func (s *SimplePassword) Deserialize(data map[string]any) error {
	username, ok := data["username"].(string)
	if !ok {
		return fmt.Errorf("username is not a string")
	}

	password, ok := data["password"].(string)
	if !ok {
		return fmt.Errorf("password is not a string")
	}

	s.Username = username
	s.Password = password

	return nil
}

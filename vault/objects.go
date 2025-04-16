package vault

import (
	"encoding/json"
	"fmt"
)

// Write wrapper for the vault client

// SimplePassword is a struct that holds a username and password
type SimplePassword struct {
	Username string
	Password string
	Number   int64
}

func (s *SimplePassword) Serialize() (map[string]any, error) {
	return map[string]any{
		"username": s.Username,
		"password": s.Password,
		"aged":     s.Number,
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

	aged, ok := data["aged"].(json.Number)
	if !ok {
		return fmt.Errorf("aged is not a int64")
	}

	s.Username = username
	s.Password = password
	s.Number, _ = aged.Int64()

	return nil
}

func (s *SimplePassword) Init() *SimplePassword {
	if s == nil {
		return &SimplePassword{}
	}

	return s
}

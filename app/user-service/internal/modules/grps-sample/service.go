package echo

import "context"

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Echo(ctx context.Context, message string) (string, error) {
	// some logic
	return "user-service returned:" + message, nil
}

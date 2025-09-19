// Package echo содержит бизнес-логику echo-модуля user-service
package echo

import "context"

type Service struct{}

func NewService() *Service { return &Service{} }

func (s *Service) Echo(ctx context.Context, message string) (string, error) {
    return "user-service:" + message, nil
}



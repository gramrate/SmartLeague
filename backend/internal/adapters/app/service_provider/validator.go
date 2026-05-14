package service_provider

import "SmartLeague/internal/adapters/controller/api/validator"

func (s *ServiceProvider) Validator() *validator.Validator {
	if s.validator == nil {
		s.validator = validator.New()
	}

	return s.validator
}

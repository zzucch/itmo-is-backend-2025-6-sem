package sell

import "github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) FindAllPhones() ([]general.Phone, error) {
	return s.repository.FindAllPhones()
}

func (s *Service) CreatePhone(phone *general.Phone) error {
	if err := general.ValidateCondition(phone.Condition); err != nil {
		return err
	}

	return s.repository.CreatePhone(phone)
}

func (s *Service) DeletePhone(id uint) error {
	return s.repository.DeletePhone(id)
}

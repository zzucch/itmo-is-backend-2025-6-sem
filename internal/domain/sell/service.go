package sell

import (
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general"
	"github.com/is-web-y26/m3302-milovatskiy/internal/domain/general/storage/s3"
)

type Service struct {
	repository Repository
	s3Client   *s3.S3Client
}

func NewService(repository Repository, s3Client *s3.S3Client) *Service {
	return &Service{
		repository: repository,
		s3Client:   s3Client,
	}
}

func (s *Service) FindAllPhones() ([]general.Phone, error) {
	return s.repository.FindAllPhones()
}

func (s *Service) GetPhoneByID(id uint) (*general.Phone, error) {
	return s.repository.FindPhoneByID(id)
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

func (s *Service) UpdatePhone(phone *general.Phone) error {
	if err := general.ValidateCondition(phone.Condition); err != nil {
		return err
	}

	return s.repository.UpdatePhone(phone)
}

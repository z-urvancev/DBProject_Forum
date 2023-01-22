package impl

import (
	"DBProject/internal/models"
	"DBProject/internal/repositories"
)

type ServiceUsecaseImpl struct {
	serviceRepository repositories.ServiceRepository
}

func NewServiceUseCase(service repositories.ServiceRepository) *ServiceUsecaseImpl {
	return &ServiceUsecaseImpl{serviceRepository: service}
}

func (suc *ServiceUsecaseImpl) ClearService() error {
	return suc.serviceRepository.ClearService()
}

func (suc *ServiceUsecaseImpl) GetService() (*models.Status, error) {
	return suc.serviceRepository.GetService()
}

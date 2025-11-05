package services

import (
	"errors"
	"invoicer-go/m/src/dto"
	"invoicer-go/m/src/models"

	"gorm.io/gorm"
)

type UserService struct {
	database *gorm.DB
}

func NewUserService(database *gorm.DB) *UserService {
	return &UserService{
		database: database,
	}
}

func (s *UserService) UpdateUser(id string, payload dto.UpdateUserDto) (*models.User, error) {
	user, err := s.GetUser(id)
	if err != nil {
		return nil, err
	}

	if payload.Name != nil {
		user.Name = *payload.Name
	}
	if payload.Email != nil {

		if err := s.checkEmailUniqueness(*payload.Email, id); err != nil {
			return nil, err
		}
		user.Email = *payload.Email
	}
	if payload.Phone != nil {
		user.Phone = *payload.Phone
	}
	if payload.RcNumber != nil {
		user.RcNumber = *payload.RcNumber
	}
	if payload.CompanyLogo != nil {
		user.CompanyLogo = *payload.CompanyLogo
	}
	if payload.CompanyName != nil {
		user.CompanyName = *payload.CompanyName
	}
	if payload.Website != nil {
		user.Website = *payload.Website
	}
	if payload.TaxId != nil {
		user.TaxId = *payload.TaxId
	}

	if payload.BankInformation != nil {
		user.BankInformation = &models.BankInformation{
			AccountName:   *payload.BankInformation.AccountName,
			AccountNumber: *payload.BankInformation.AccountNumber,
			BankName:      *payload.BankInformation.BankName,
			BankSwiftCode: *payload.BankInformation.BankSwiftCode,
			Iban:          *payload.BankInformation.Iban,
			RoutingNumber: *payload.BankInformation.RoutingNumber,
		}
	}

	if err := s.database.Save(user).Error; err != nil {
		return nil, err
	}

	return s.GetUser(id)
}

func (s *UserService) DeleteUser(id string) error {
	user, err := s.GetUser(id)
	if err != nil {
		return err
	}

	return s.database.Transaction(func(tx *gorm.DB) error {
		if user.BankInformation != nil {
			if err := tx.Delete(user.BankInformation).Error; err != nil {
				return err
			}
		}

		return tx.Delete(user).Error
	})
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	user := &models.User{}

	if err := s.database.First(user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *UserService) checkEmailUniqueness(email, excludeID string) error {
	var count int64
	query := s.database.Model(&models.User{}).Where("email = ?", email)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	if err := query.Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("email already exists")
	}

	return nil
}
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := s.database.Preload("BankInformation").Where("email = ?", email).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}
func (s *UserService) GetAllUsers(limit, page int) ([]models.User, int64, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	if limit > 100 {
		limit = 100
	}

	var users []models.User
	var total int64

	offset := (page - 1) * limit

	if err := s.database.Preload("BankInformation").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	if err := s.database.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

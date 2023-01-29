package configurator

import (
	"fmt"

	"gorm.io/gorm"
)

type ZeroAffectedRows struct {
}

func (z ZeroAffectedRows) Error() string {
	return "Zero affected rows. Expected more."
}

/////////////////////

type StorageErrorExists struct {
	items []string
}

func (s StorageErrorExists) Error() string {
	return fmt.Sprintf("database already has those items=%v", s.items)
}

///////////////////////

type ConfiguratorError struct {
	rowAffected []int
	errors      []error
}

func (s *ConfiguratorError) AppendSQLError(res *gorm.DB) *ConfiguratorError {
	s.rowAffected = append(s.rowAffected, int(res.RowsAffected))
	s.errors = append(s.errors, res.Error)
	return s
}

func (s *ConfiguratorError) AppenError(err error) *ConfiguratorError {
	s.errors = append(s.errors, err)
	return s
}

func (s *ConfiguratorError) GetError() error {
	for _, row := range s.rowAffected {
		if row != 0 {
			return nil
		}
	}

	for _, err := range s.errors {
		if err != nil {
			return err
		}
	}

	return ZeroAffectedRows{}
}

func (s *ConfiguratorError) GetErrorWithAllowedZeroRows() error {
	for _, err := range s.errors {
		if err != nil {
			return err
		}
	}

	return nil
}

package configurator

import (
	"darkbot/configurator/models"
	"darkbot/utils"
	"fmt"

	"gorm.io/gorm"
)

type IConfiguratorTags interface {
	TagsAdd(channelID string, tags ...string) *ConfiguratorError
	TagsRemove(channelID string, tags ...string) *ConfiguratorError
	TagsList(channelID string) ([]string, *ConfiguratorError)
	TagsClear(channelID string) *ConfiguratorError
}

type ConfiguratorBase struct {
	Configurator
}

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

type ZeroAffectedRows struct {
}

func (z ZeroAffectedRows) Error() string {
	return "Zero affected rows. Expected more."
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

type StorageErrorExists struct {
	items []string
}

func (s StorageErrorExists) Error() string {
	return fmt.Sprintf("database already has those items=%v", s.items)
}

func (c ConfiguratorBase) TagsAdd(channelID string, tags ...string) *ConfiguratorError {
	objs := []models.TagBase{}
	errors := &ConfiguratorError{}
	for _, tag := range tags {
		objs = append(objs, models.TagBase{
			TagTemplate: models.TagTemplate{
				ChannelShared: models.ChannelShared{

					ChannelID: channelID,
				},
				Tag: tag,
			},
		})
	}

	presentTags, _ := c.TagsList(channelID)
	for _, tag := range presentTags {
		for _, newtag := range tags {
			if tag == newtag {
				fmt.Printf("TagsAdd. Tag %s is already present in channelID=%s\n", tag, channelID)
				errors.AppenError(StorageErrorExists{items: []string{tag}})
				return errors
			}
		}
	}

	res := c.db.Create(objs)
	utils.CheckWarn(res.Error, "unsuccesful result of c.db.Create")
	errors.AppendSQLError(res)
	return errors
}

func (c ConfiguratorBase) TagsRemove(channelID string, tags ...string) *ConfiguratorError {
	errors := &ConfiguratorError{}
	for _, tag := range tags {
		result := c.db.Where("channel_id = ? AND tag = ?", channelID, tag).Delete(&models.TagBase{})
		utils.CheckWarn(result.Error, "unsuccesful result of c.db.Delete")
		errors.AppendSQLError(result)
	}
	return errors
}

func (c ConfiguratorBase) TagsList(channelID string) ([]string, *ConfiguratorError) {
	objs := []models.TagBase{}
	result := c.db.Where("channel_id = ?", channelID).Find(&objs)
	utils.CheckWarn(result.Error, "unsuccesful result of c.db.Find")

	return utils.CompL(objs,
		func(x models.TagBase) string { return x.Tag }), (&ConfiguratorError{}).AppendSQLError(result)
}

func (c ConfiguratorBase) TagsClear(channelID string) *ConfiguratorError {
	var tags []models.TagBase
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	fmt.Println("Clear.Find.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Find.err=", result.Error)
	result = c.db.Unscoped().Delete(&tags)
	fmt.Println("Clear.Delete.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Delete.err=", result.Error)
	return (&ConfiguratorError{}).AppendSQLError(result)
}

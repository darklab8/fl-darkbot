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

func (s *ConfiguratorError) AppendAll(res *gorm.DB) *ConfiguratorError {
	s.rowAffected = append(s.rowAffected, int(res.RowsAffected))
	s.errors = append(s.errors, res.Error)
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

	res := c.db.Create(objs)
	utils.CheckWarn(res.Error, "unsuccesful result of c.db.Create")
	errors.AppendAll(res)
	return errors
}

func (c ConfiguratorBase) TagsRemove(channelID string, tags ...string) *ConfiguratorError {
	errors := &ConfiguratorError{}
	for _, tag := range tags {
		result := c.db.Where("channel_id = ? AND tag = ?", channelID, tag).Delete(&models.TagBase{})
		utils.CheckWarn(result.Error, "unsuccesful result of c.db.Delete")
		errors.AppendAll(result)
	}
	return errors
}

func (c ConfiguratorBase) TagsList(channelID string) ([]string, *ConfiguratorError) {
	objs := []models.TagBase{}
	result := c.db.Where("channel_id = ?", channelID).Find(&objs)
	utils.CheckWarn(result.Error, "unsuccesful result of c.db.Find")

	return utils.CompL(objs,
		func(x models.TagBase) string { return x.Tag }), (&ConfiguratorError{}).AppendAll(result)
}

func (c ConfiguratorBase) TagsClear(channelID string) *ConfiguratorError {
	var tags []models.TagBase
	result := c.db.Unscoped().Where("channel_id = ?", channelID).Find(&tags)
	fmt.Println("Clear.Find.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Find.err=", result.Error)
	result = c.db.Unscoped().Delete(&tags)
	fmt.Println("Clear.Delete.rowsAffected=", result.RowsAffected)
	fmt.Println("Clear.Delete.err=", result.Error)
	return (&ConfiguratorError{}).AppendAll(result)
}

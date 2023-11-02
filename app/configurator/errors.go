package configurator

import (
	"fmt"
)

type ErrorZeroAffectedRows struct {
	ExtraMsg string
}

const ErrorZeroAffectedRowsMsg = "Zero affected rows. not found records. Expected more."

func (z ErrorZeroAffectedRows) Error() string {
	return ErrorZeroAffectedRowsMsg + z.ExtraMsg
}

/////////////////////

type StorageErrorExists struct {
	items []string
}

func (s StorageErrorExists) Error() string {
	return fmt.Sprintf("database already has those items=%v", s.items)
}

///////////////////////

type errorAggregator struct {
	errors []error
}

func NewErrorAggregator() errorAggregator {
	return errorAggregator{}
}

func (e *errorAggregator) Append(err error) {
	e.errors = append(e.errors, err)
}

func (e errorAggregator) TryToGetError() error {
	for _, err := range e.errors {
		if err != nil {
			return err
		}
	}
	return nil
}

package storage

import "darkbot/scrappy/shared/api"

type IStorage interface {
	Update()
	Length() int
	API() api.APIinterface
}

/*
Reusable API code to request []byte code of smth. Reusable for player and base.
*/

package apiRawData

type APIinterface interface {
	GetData() []byte
	New() APIinterface
}

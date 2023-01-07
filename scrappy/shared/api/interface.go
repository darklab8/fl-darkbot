/*
Reusable API code to request []byte code of smth. Reusable for player and base.
*/

package api

type APIinterface interface {
	GetData() []byte
	New() APIinterface
}

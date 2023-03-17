package option

import (
	"github.com/bww/go-dbx/v1/persist"
)

type Config = persist.Config
type Option = persist.Option

func UseConfig(c Config) Option {
	return persist.UseConfig(c)
}
func Cascade(on bool) Option {
	return persist.Cascade(on)
}
func FetchRelated(on bool) Option {
	return persist.FetchRelated(on)
}
func StoreRelated(on bool) Option {
	return persist.StoreRelated(on)
}
func DeleteRelated(on bool) Option {
	return persist.DeleteRelated(on)
}
func Params(p map[string]interface{}) Option {
	return persist.Params(p)
}

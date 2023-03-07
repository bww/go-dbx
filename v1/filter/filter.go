// DEPRECATED: everything in this package is deprecated and
// will be removed in a future release. Use the counterpart
// types and functions from the 'query' package instead.
package filter

import (
	"github.com/bww/go-dbx/v1/query"
)

type Option = query.ReadOption
type Filter = query.ReadConfig
type Range = query.Range
type Order = query.Order
type Timeframe = query.Timeframe

func UseFilter(f Filter) Option {
	return query.UseReadConfig(f)
}
func WithLimit(r Range) Option {
	return query.WithLimit(r)
}
func WithOrder(o Order) Option {
	return query.WithOrder(o)
}
func WithTimeframe(t Timeframe) Option {
	return query.WithTimeframe(t)
}

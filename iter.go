package page

import (
	"context"
	"iter"
)

type Read[I, P any] func(ctx context.Context, nextPage *P) ([]I, *P, error)

func Iter[I, P any](ctx context.Context, readPage Read[I, P]) iter.Seq2[I, error] {
	var nextPage *P
	var i int
	var ts []I

	var err error
	ts, nextPage, err = readPage(ctx, nextPage)
	return func(yield func(I, error) bool) {
		for {
			if err != nil {
				var t I
				if !yield(t, err) {
					return
				}
			}
			for ; i < len(ts); i++ {
				t := ts[i]
				if !yield(t, nil) {
					i++
					return
				}
			}
			if nextPage == nil {
				return
			}
			ts, nextPage, err = readPage(ctx, nextPage)
			i = 0
		}
	}
}

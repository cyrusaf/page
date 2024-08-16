package page

import (
	"context"
	"iter"
)

type Read[I, P any] func(ctx context.Context, nextPage *P) ([]I, *P, error)

func Iter[I, P any](ctx context.Context, readPage Read[I, P]) iter.Seq2[I, error] {
	return func(yield func(I, error) bool) {
		var nextPage *P
		for {
			var err error
			var ts []I
			ts, nextPage, err = readPage(ctx, nextPage)
			if err != nil {
				var t I
				if !yield(t, err) {
					return
				}
			}
			for _, t := range ts {
				if !yield(t, nil) {
					return
				}
			}
			if nextPage == nil {
				return
			}
		}
	}
}

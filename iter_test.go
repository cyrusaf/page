package page_test

import (
	"context"
	"testing"

	"github.com/cyrusaf/page"
)

func TestIter(t *testing.T) {
	ctx := context.Background()
	f := FakePager{
		Items: 10,
		Limit: 2,
	}

	i := 0
	for j, err := range page.Iter(ctx, f.Read) {
		if err != nil {
			t.Fatalf("expected err to be nil, but got %s instead", err.Error())
		}
		if i != j {
			t.Fatalf("expected iterator value to be %d but got %d instead", i, j)
		}
		i++
	}
	if f.Invocations != 5 {
		t.Fatalf("expected 5 paginator invocations, but got %d instead", f.Invocations)
	}
}

type FakePager struct {
	Invocations int
	Cursor      int
	Items       int
	Limit       int
}

func (f *FakePager) Read(ctx context.Context, nextPage *int) ([]int, *int, error) {
	page := 0
	if nextPage != nil {
		page = *nextPage
	}

	f.Invocations += 1
	items := []int{}
	for i := 0; i < f.Limit && f.Cursor < f.Items; i++ {
		items = append(items, page+i)
		f.Cursor += 1
	}
	if f.Cursor < f.Items {
		return items, &f.Cursor, nil
	}
	return items, nil, nil
}

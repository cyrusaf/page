package page_test

import (
	"context"
	"sync"
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
	if f.invocations != 5 {
		t.Fatalf("expected 5 paginator invocations, but got %d instead", f.invocations)
	}
}

func TestIterInterrupt(t *testing.T) {
	ctx := context.Background()
	f := FakePager{
		Items: 10,
		Limit: 2,
	}

	i := 0
	iter := page.Iter(ctx, f.Read)
	for j, err := range iter {
		if err != nil {
			t.Fatalf("expected err to be nil, but got %s instead", err.Error())
		}
		if i != j {
			t.Fatalf("expected iterator value to be %d but got %d instead", i, j)
		}
		i++
		break
	}
	for j, err := range iter {
		if err != nil {
			t.Fatalf("expected err to be nil, but got %s instead", err.Error())
		}
		if i != j {
			t.Fatalf("expected iterator value to be %d but got %d instead", i, j)
		}
		i++
	}
	if f.invocations != 5 {
		t.Fatalf("expected 5 paginator invocations, but got %d instead", f.invocations)
	}
}

type FakePager struct {
	mu          sync.Mutex
	invocations int
	Cursor      int
	Items       int
	Limit       int
}

func (f *FakePager) Read(ctx context.Context, nextPage *int) ([]int, *int, error) {
	page := 0
	if nextPage != nil {
		page = *nextPage
	}

	f.mu.Lock()
	f.invocations += 1
	f.mu.Unlock()
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

func (f *FakePager) Invocations() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.invocations
}

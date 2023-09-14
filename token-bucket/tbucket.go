// Package tbucket is an implementation of a token bucket algorithm.
package tbucket

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Bucket is a token bucket.
// Tokens are added to the bucket at a fixed interval.
type Bucket struct {
	ticker *time.Ticker
	bucket chan struct{}
}

// NewBucket creates a new token bucket.
// A ticker is created to add tokens to the bucket at a fixed interval.
func NewBucket(size int, refillInterval time.Duration) (*Bucket, error) {
	var err error
	if size <= 0 {
		err = errors.New("size must be greater than 0")
	}
	if refillInterval <= 0 {
		err = errors.Join(err, errors.New("refillInterval must be greater than 0"))
	}
	if err != nil {
		return nil, err
	}
	b := &Bucket{
		ticker: time.NewTicker(refillInterval),
		bucket: make(chan struct{}, size),
	}

	// fill the bucket
	for i := 0; i < size; i++ {
		b.bucket <- struct{}{}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		for {
			_, ok := <-b.ticker.C
			if !ok {
				close(b.bucket)
				return
			}
			if len(b.bucket) < size {
				b.bucket <- struct{}{}
			}
		}
	}()
	wg.Wait()

	return b, nil
}

// Get returns a channel that will sent a token when one is available.
func (b *Bucket) Get(ctx context.Context) error {
	select {
	case <-b.bucket:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

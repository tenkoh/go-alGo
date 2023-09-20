// Package tbucket is an implementation of a token bucket algorithm.
package tbucket

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrBucketClosed is returned when the bucket is closed.
var ErrBucketClosed = errors.New("bucket is closed")

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
			<-b.ticker.C
			if len(b.bucket) < size {
				b.bucket <- struct{}{}
			}
		}
	}()
	wg.Wait()

	return b, nil
}

// Get returns block the response until a token is available.
// If the context is canceled, it returns an error.
// If the bucket is closed, it returns an error.
func (b *Bucket) Get(ctx context.Context) error {
	select {
	case _, ok := <-b.bucket:
		if !ok {
			return ErrBucketClosed
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close closes the bucket. After closing, Get returns an error.
func (b *Bucket) Close() {
	// stop refilling the bucket
	b.ticker.Stop()
	close(b.bucket)
	// remove all tokens from the bucket
	for range b.bucket {
		continue
	}
}

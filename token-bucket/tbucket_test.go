package tbucket_test

import (
	"context"
	"tbucket"
	"testing"
	"time"
)

// 1. bucketに空きがある場合はtokenを取り出せる
// 2. bucketに空きがない場合はtokenを取り出せない(ブロックされる)
// 3. bucketからの取り出しはキャンセルできる
// 4. bucketには決めた間隔でtokenが補充される
// 5. bucketには最大値を超えてtokenを補充できない
// 6. bucketの最大サイズは1以上でなければならない
// 7. bucketの補充間隔は0より大きくなければならない

func TestBucket_Get(t *testing.T) {
	bucketSize := 3
	interval := 500 * time.Millisecond
	tolerance := interval / 10

	ctx := context.Background()

	b, err := tbucket.NewBucket(bucketSize, interval)
	if err != nil {
		t.Error(err)
		return
	}

	// Confirm that the bucket is not blocked
	t0 := time.Now()
	for i := 0; i < bucketSize; i++ {
		<-b.Get(ctx)
	}
	t1 := time.Now()
	if t1.Sub(t0) > tolerance {
		t.Errorf("bucket is blocked from being taken out")
	}

	// Confirm that the bucket is blocked
	<-b.Get(ctx)
	t1 = time.Now()
	if dt := t1.Sub(t0); dt > interval+tolerance || dt < interval-tolerance {
		t.Errorf("bucket is not blocked properly")
	}
}

package tbucket_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/tenkoh/go-alGo/tbucket"
)

// 1. ok bucketに空きがある場合はtokenを取り出せる
// 2. ok bucketに空きがない場合はtokenを取り出せない(ブロックされる)
// 3. ok bucketからの取り出しはキャンセルできる
// 4. ok bucketには決めた間隔でtokenが補充される
// 5. 自明 bucketには最大値を超えてtokenを補充できない
// 6. ok bucketの最大サイズは1以上でなければならない
// 7. ok bucketの補充間隔は0より大きくなければならない
// 8. ok bucketはクローズできる。クローズしたbucketからの取り出しはエラーになる。

func TestNewBucket(t *testing.T) {
	// 6. bucketの最大サイズは1以上でなければならない
	if _, err := tbucket.NewBucket(0, 1*time.Second); err == nil {
		t.Error("want error: bucket size is not larger than 0, but got nil")
	}
	// 7. bucketの補充間隔は0より大きくなければならない
	if _, err := tbucket.NewBucket(1, 0); err == nil {
		t.Error("want error: interval is not larger than 0, but got nil")
	}

	if _, err := tbucket.NewBucket(1, 1*time.Nanosecond); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func TestBucket_Close(t *testing.T) {
	interval := 10 * time.Second // 補充されない十分大きいインターバルを設定する
	ctx := context.Background()

	// bucketをクローズした後に新たなGetができないことを確認する
	b, err := tbucket.NewBucket(1, interval)
	if err != nil {
		t.Error(err)
		return
	}

	b.Close()
	if err := b.Get(ctx); !errors.Is(err, tbucket.ErrBucketClosed) {
		t.Errorf("want error: bucket is closed, but got %v", err)
	}

	// bucketをクローズした後に既存のGetがエラーを返すことを確認する
	b, err = tbucket.NewBucket(1, interval)
	if err != nil {
		t.Error(err)
		return
	}

	var errch = make(chan error, 1)

	b.Get(ctx) // bucketを空にする

	go func() {
		errch <- b.Get(ctx)
	}()

	// wait for the goroutine to start
	// [TODO] sleepによる待ち合わせは不確実なため、より確実な方法を検討する
	<-time.After(100 * time.Millisecond)
	b.Close()

	if err := <-errch; !errors.Is(err, tbucket.ErrBucketClosed) {
		t.Errorf("want error: bucket is just closed, but got %v", err)
	}

}

func TestBucket_Get(t *testing.T) {
	bucketSize := 3
	interval := 500 * time.Millisecond
	tolerance := interval / 10 // 許容する誤差.time.Tickerのテストにも同様の考えがある

	ctx := context.Background()

	b, err := tbucket.NewBucket(bucketSize, interval)
	if err != nil {
		t.Error(err)
		return
	}

	// Confirm that the bucket is not blocked
	t0 := time.Now()
	for i := 0; i < bucketSize; i++ {
		b.Get(ctx)
	}
	t1 := time.Now()
	if t1.Sub(t0) > tolerance {
		t.Errorf("bucket is blocked from being taken out")
	}

	// Confirm that the bucket is blocked, and the token is reffiled after the interval
	b.Get(ctx)
	t1 = time.Now()
	if dt := t1.Sub(t0); dt > interval+tolerance || dt < interval-tolerance {
		t.Errorf("bucket is not blocked properly")
	}

	// Confirm that the bucket can be canceled
	ctx1, cancel1 := context.WithTimeout(ctx, interval/10)
	if err := b.Get(ctx1); err == nil {
		t.Errorf("bucket is not canceled")
	}
	defer cancel1()
}

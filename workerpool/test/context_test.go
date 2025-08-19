package test

import (
	"context"
	"fmt"
	"testing"
)

func TestRequestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	useCanceledContext(ctx)

	fmt.Println("err")

	fmt.Println("ctx")
}

func useCanceledContext(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}

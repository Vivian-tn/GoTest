package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	//客户端
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	serve(ctx)
	time.Sleep(3 * time.Second)
	fmt.Println(3)
}

//服务端
func serve(ctx context.Context) {
	fmt.Println(1, ctx.Err())
	time.Sleep(2 * time.Second)
	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		return
	}
	fmt.Println(2)
}

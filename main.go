package main

import (
	"context"
	"fmt"
	"github.com/codingWhat/ditributed-lock/mysql"
)

func main() {
	//ctx := context.Background()
	//
	//var st int //sleep time = logic cost time
	//flag.IntVar(&st, "s", 0, "-s xxx")
	//flag.Parse()
	//
	//holder := fmt.Sprintf("h:%d", time.Now().Unix()%10)
	//l, err := mysql.NewLockerV1(holder)
	//if err != nil {
	//	panic(err)
	//}
	////1678440722
	//for {
	//	err := l.Lock(ctx, "dsadas", 10)
	//	if err != nil {
	//		fmt.Printf("time:%d, %s lock failed. err: %+v \n", time.Now().Unix(), holder, err)
	//		time.Sleep(100 * time.Millisecond)
	//		continue
	//	}
	//	//t := rand.Int()%5 + 1
	//	fmt.Printf("time:%d, %s lock success. sleep: %d \n", time.Now().Unix(), holder, st)
	//	time.Sleep(time.Duration(st) * time.Second)
	//	_ = l.UnLock(ctx, "dsadas")
	//}

	l, _ := mysql.NewLockerV2()

	ctx := context.Background()

	if err := l.Acquire(ctx, "aaas"); err != nil {
		fmt.Println("acquire -->", err.Error())
		return
	}
	fmt.Println("---->")
	// do business logic
	if err := l.Release(ctx); err != nil {
		fmt.Println("release -->", err.Error())
		return
	}
}

package chanpprof

import (
	"fmt"
	"os"
	"time"
)

// 未关闭的 channel 的 pprof
// GODEBUG='gctrace=1' go run main.go > UncloseChanPPROF.log
func UncloseChanPPROF(exitChan chan os.Signal, openCorrect bool) {
	index := 0
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-exitChan:
			fmt.Println("receive exit signal")
			goto END
		case <-ticker.C:
			fmt.Println("start go routine index", index)
			index++
			go func(i int) {
				count := 0
				goRoutineTimer := time.NewTimer(time.Second)
				receiverChan := make(chan int)
				go func() {
					for v := range receiverChan {
						fmt.Println("go routine index", i, "receive chan value", v)
					}
					fmt.Println("go routine index", i, "receiver end")
				}()

				for count < 30 {
					<-goRoutineTimer.C
					receiverChan <- count
					goRoutineTimer.Reset(time.Second)
					count++
				}
				if openCorrect {
					// 如果不关闭 channel，会造成 监听 receiverChan 的协程泄露
					// 协程未泄露，在 http://127.0.0.1:9527/debug/pprof/ 中 goroutine 的数量会保持不变
					// 协程泄露，在 http://127.0.0.1:9527/debug/pprof/ 中 goroutine 的数量会持续增长
					close(receiverChan)
				}
				fmt.Println("go routine index", i, "sender end")
			}(index)
		}
	}
END:
	fmt.Println("sleep 5s to check pprof")
	time.Sleep(time.Second * 5)
	fmt.Println("sleep end, exit")
}

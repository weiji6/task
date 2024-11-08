1.思考以下程序在并发中出现panic的可能是什么？并给出解决方案
```go
package main

import (
	"fmt"
	"time"
)

type message struct {
	Topic     string
	Partition int32
	Offset    int64
}

type FeedEventDM struct {
	Type    string
	UserID  int
	Title   string
	Content string
}

type MSG struct {
	ms        message
	feedEvent FeedEventDM
}

const ConsumeNum = 5

func main() {
	var consumeMSG []MSG
	var lastConsumeTime time.Time // 记录上次消费的时间
	msgs := make(chan MSG)

	//这里源源不断的生产信息
	go func() {
		for i := 0; ; i++ {
			msgs <- MSG{
				ms: message{
					Topic:     "消费主题",
					Partition: 0,
					Offset:    0,
				},
				feedEvent: FeedEventDM{
					Type:    "grade",
					UserID:  i,
					Title:   "成绩提醒",
					Content: "您的成绩是xxx",
				},
			}
			//每次发送信息会停止0.01秒以模拟真实的场景
			time.Sleep(100 * time.Millisecond)
		}
	}()

	//不断接受消息进行消费
	for msg := range msgs {
		// 添加新的值到events中
		consumeMSG = append(consumeMSG, msg)
		// 如果数量达到额定值就批量消费
		if len(consumeMSG) >= ConsumeNum {
			//进行异步消费
			go func() {
				m := consumeMSG[:ConsumeNum]
				fn(m)
			}()
			// 更新上次消费时间
			lastConsumeTime = time.Now()
			// 清除插入的数据
			consumeMSG = consumeMSG[ConsumeNum:]
		} else if !lastConsumeTime.IsZero() && time.Since(lastConsumeTime) > 5*time.Minute {
			// 如果距离上次消费已经超过5分钟且有未处理的消息
			if len(consumeMSG) > 0 {
				//进行异步消费
				go func() {
					m := consumeMSG[:ConsumeNum]
					fn(m)
				}()
				// 更新上次消费时间
				lastConsumeTime = time.Now()
				// 清空插入的数据
				consumeMSG = consumeMSG[ConsumeNum:]
			}
		}
	}
}

func fn(m []MSG) {
	fmt.Printf("本次消费了%d条消息\n", len(m))
}
```
在主协程和消费协程共享一个consumeMSG，在添加新值时，会改变消费切片，导致报错，应该加一个互斥锁将两者的行为分开

修改后如下
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type message struct {
	Topic     string
	Partition int32
	Offset    int64
}

type FeedEventDM struct {
	Type    string
	UserID  int
	Title   string
	Content string
}

type MSG struct {
	ms        message
	feedEvent FeedEventDM
}

const ConsumeNum = 5

var lock sync.Mutex

func main() {
	var consumeMSG []MSG
	var lastConsumeTime time.Time
	msgs := make(chan MSG)

	go func() {
		for i := 0; ; i++ {
			lock.Lock()
			msgs <- MSG{
				ms: message{
					Topic:     "消费主题",
					Partition: 0,
					Offset:    0,
				},
				feedEvent: FeedEventDM{
					Type:    "grade",
					UserID:  i,
					Title:   "成绩提醒",
					Content: "您的成绩是xxx",
				},
			}
			lock.Unlock()
			time.Sleep(100 * time.Millisecond)
		}
	}()

	for msg := range msgs {
		lock.Lock()
		consumeMSG = append(consumeMSG, msg)
		if len(consumeMSG) >= ConsumeNum {
			m := consumeMSG[:ConsumeNum]
			lastConsumeTime = time.Now()
			consumeMSG = consumeMSG[ConsumeNum:]
			lock.Unlock()
			go func() {
				fn(m)
			}()
		} else if !lastConsumeTime.IsZero() && time.Since(lastConsumeTime) > 5*time.Minute {
			if len(consumeMSG) > 0 {
				m := consumeMSG[:ConsumeNum]
				lastConsumeTime = time.Now()
				consumeMSG = consumeMSG[ConsumeNum:]
				lock.Unlock()
				go func() {
					fn(m)
				}()
			}
		} else {
			lock.Unlock()
		}
	}
}

func fn(m []MSG) {
	fmt.Printf("本次消费了%d条消息\n", len(m))
}
```

2.使用for循环生成20个goroutine，并向一个channel传入随机数和goroutine编号，等待这些goroutine都生成完后，想办法给这些goroutine按照编号进行排序(输出排序前和排序后的结果,要求不使用额外的空间存储着20个数据)
```go
package main

import (
	"fmt"
	"golang.org/x/exp/rand"
	"sync"
)

type nums struct {
	id  int
	num int
}

var wg sync.WaitGroup

func main() {
	ch := make(chan nums, 20)
	for i := 1; i < 21; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			v := rand.Intn(100)
			ch <- nums{i, v}
		}(i)
		fmt.Printf("%v ", <-ch)
	}
	fmt.Println()
	wg.Wait()
}
```
3.经典老题：交叉打印下面两个字符串（要求一个打印完，另一个会继续打印）
"ABCDEFGHIJKLMNOPQRSTUVWXYZ" "0123..."
得到："AB01CD23EF34..."
```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	str1 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str2 := "0123456789"

	ch1 := make(chan byte)
	ch2 := make(chan byte)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := range str1 {
			ch1 <- str1[i]
		}
		close(ch1)
	}()

	go func() {
		defer wg.Done()
		for i := range str2 {
			ch2 <- str2[i]
		}
		close(ch2)
	}()
	go func() {
		wg.Wait()
	}()

	var a, b byte
	var o1, o2 bool
	for {
		select {
		case a, o1 = <-ch1:
			if !o1 {
				ch1 = nil
			} else {
				fmt.Printf(string(a))
			}
		case b, o2 = <-ch2:
			if !o2 {
				ch2 = nil
			} else {
				fmt.Printf(string(b))
			}
		}

		if ch1 == nil && ch2 == nil {
			break
		}
	}
}
```
package homework

import (
	"fmt"
	"log"
	"net/rpc"
	"sync"
	"time"
)

type Coordinator struct {
	lock     sync.RWMutex
	stage    string      //运行阶段
	toDo     chan Update //需要更新的数据
	Deadline time.Time   //Deadline
	Workers  [2]int      //workers列表
	Data     [10000]int  //存储的数据，默认值为0
}

func MakeCoordinator() *Coordinator {
	c := Coordinator{
		toDo:    make(chan Update),
		Workers: [2]int{1, 2},
	}
	fmt.Println("coordinator创建完成")
	for _, id := range c.Workers {
		fmt.Println("循环启动worker服务")
		go func() {
			MakeWork(id)
		}()
	}
	for {
		time.Sleep(500 * time.Millisecond)
		order := OrderInfo{
			Stage: PREPARE,
		}
		reply := Reply{}
		workername := "worker"
		callPrepare(workername+"GetAndReturn", &order, &reply)
		if reply.Stage == Checked {
			commitOrder := OrderInfo{
				Stage: COMMIT,
				Args:  c.Getdata(reply.IndexI),
			}
			callCommit(workername+".doCommit", &commitOrder, &reply)
		} else {
			rollbackOrder := OrderInfo{
				Stage: ROLLBACK,
			}
			callRollBack(workername+".doROLLBACK", &rollbackOrder, &reply)
		}

	}
	return &c
}

//读取
func (c *Coordinator) Getdata(i int) [3]int {
	c.lock.RLock()
	args := [3]int{c.Data[i], c.Data[(i+1)%10000], c.Data[(i+2)%10000]}
	c.lock.RUnlock()
	return args
}

//写入
func (c *Coordinator) Updata() bool {
	c.lock.Lock()
	for {
		c.lock.Lock()
		task := <-c.toDo
		c.Data[task.Index] = task.Number
		c.lock.Unlock()
	}
	c.lock.Unlock()
	return true
}

//告诉workers准备好，得到传来的index
func callPrepare(rpcname string, args interface{}, reply interface{}) bool {
	sockname := coordinatorSock()
	// 拨号服务
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

//把第一阶段得到的数据给workers  ，然后workers计算结束后返回
func callCommit(rpcname string, args interface{}, reply interface{}) bool {
	sockname := coordinatorSock()
	// 拨号服务
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

//第一阶段不成功，那么告诉workers要ROLLBACK
func callRollBack(rpcname string, args interface{}, reply interface{}) bool {
	sockname := coordinatorSock()
	// 拨号服务
	c, err := rpc.DialHTTP("unix", sockname)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

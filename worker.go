package homework

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type KeyValue struct {
	Key   int
	Value int
}

type Worker struct {
	I         int
	J         int
	Data      [3]int
	WorkTimes int
	WorkerId  int
}

//开始work
func MakeWork(Id int) {
	w := Worker{
		WorkTimes: 0,
		I:         rand.Int() % 10000,
		J:         rand.Int() % 10000,
		WorkerId:  os.Getegid(),
	}
	fmt.Printf("%d工作！", w.WorkerId)
	id := os.Getpid()
	log.Println("工号：%d 开始工作", id)
	w.server() //启动服务开始监听来自coordinator的信息
	//每次工作完自动更新 I 、 J 的值 ACK

	log.Println("worker Id :%d 完成10000次", id)

}

func (w *Worker) GetAndReturn(args *OrderInfo, reply *Reply) error {
	reply.IndexI = w.I
	reply.Stage = Checked
	w.Data = args.Args
	w.Data[0] = w.Data[0] + w.Data[1] + w.Data[2]
	return nil
}

func (w *Worker) doCommit(args *OrderInfo, reply *Reply) error {
	if w.WorkTimes < 10000 {
		reply.Data.Number = w.Data[0]
		reply.Data.Index = w.J
		reply.Stage = ACK
		w.WorkTimes++
		log.Println("本%d 次传递完成", w.WorkTimes)
		w.I = rand.Int() % 10000
		w.J = rand.Int() % 10000
	} else {
		reply.Stage = DONE
	}
	return nil
}

func (w *Worker) doRollBack(args *OrderInfo, reply *Reply) error {
	w.WorkTimes--
	reply.Stage = ACK
	w.I = rand.Int() % 10000
	w.J = rand.Int() % 10000
	log.Println("本%d 次无效", w.WorkTimes)
	return nil
}

//启动服务 监听coordinator的消息
func (w *Worker) server() {
	rpc.Register(w)
	rpc.HandleHTTP()
	sockname := "homework"
	os.Remove(sockname) //删除文件
	//监听来自homework的信息
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

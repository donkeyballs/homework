package homework

import (
	"os"
	"strconv"
	"time"
)

type Update struct {
	Index    int
	Number   int
	WorkerId int
	Deadline time.Time
}

const (
	PREPARE  = "PREPARE"
	COMMIT   = "COMMIT"
	ROLLBACK = "ROLLBACK"
	ACK      = "ACK"
	Checked  = "YES"
	DONE     = "DONE"
)

//命令
type OrderInfo struct {
	WorkerId int
	Stage    string //PREPARE , COMMIT, ROLLBACK
	Args     [3]int //传递的参数
}

//回复
type Reply struct {
	Stage  string // ACK  PREPARED
	IndexI int
	Data   Update
}

func coordinatorSock() string {
	s := "homework"
	s += strconv.Itoa(os.Getuid())
	return s
}

// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0
/**
For each incoming connection form a TCP client, the TCP server will start
a new goroutine to handle that request.
*/
import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
)

var CloseInterrupt = errors.New("Closing server")
type multiEchoServer struct {
	// TODO: implement this!
	mq []chan string
	close chan error
	cnt int
	clients chan map[int]*client
	ln      *net.TCPListener
	broadChan chan []byte

}


type client struct {
	id int
	conn *net.TCPConn

	recvChan chan []byte
	sendChan chan []byte
	server *multiEchoServer
	toCloseRecv chan int
	toCloseSend chan int
}

func handleConnection(c net.Conn, mes *multiEchoServer) error {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	idx := len(mes.mq)
	mes.mq = append(mes.mq, make(chan string, 100))
	for {
		netData, err:= bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			mes.cnt -= 1
			return err
		}
		for _, ele := range mes.mq {
			ele<-netData
		}
		go print(c, mes.mq[idx])

	}

	c.Close()
	return nil
}


func print(c net.Conn, msg chan string) bool {
	for {
		line, err:= <-msg
		if err {
			return err
		}
		c.Write([]byte(line))
	}
	return true

}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!

	return &multiEchoServer{

		mq:    make([]chan string, 1),
		close: make(chan error),
		cnt:   0,


	}
}

func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!
	PORT := ":" + strconv.Itoa(port)
	addr, err := net.ResolveTCPAddr("tcp", PORT)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()
	//enable multiple servers
	go func() {
		for {
			fmt.Println("I am waiting")
			c, err := l.Accept()
			if err != nil {
				fmt.Println("I have error")
				fmt.Println(err)
			}
			fmt.Println("I am free")
			closeErr := <-mes.close
			if closeErr != nil {
				return closeErr
			}
			mes.cnt += 1
			go handleConnection(c, mes)
		}
	}()



}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!

}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	return mes.cnt
}

// TODO: add additional methods/functions below!

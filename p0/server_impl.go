// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0
/**
For each incoming connection form a TCP client, the TCP server will start
a new goroutine to handle that request.
*/
import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

type multiEchoServer struct {
	// TODO: implement this!
	//mq []chan s
	//close chan error
	clients map[int]*client
	ln *net.TCPListener
	broadChan chan []byte
	cnt int

}

type client struct {
	id int
	conn *net.TCPConn
	//used to store the msg received
	recvChan chan[] byte
	//used to store the msg sent
	sendChan chan[] byte
	server *multiEchoServer
	toCloseRecv chan int
	toCloseSend chan int
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	// TODO: implement this!
	s := &multiEchoServer{
		cnt: 0,
		clients:     make(map[int]*client),
		broadChan:   make(chan []byte, 1)}
	//s.clients <- make(map[int]*client, 1)
	return s
}



func (mes *multiEchoServer) Start(port int) error {
	// TODO: implement this!
	PORT := ":"+strconv.Itoa(port)
	addr, err := net.ResolveTCPAddr("tcp", PORT)
	if err != nil {
		return err
	}
	mes.ln, err = net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	go func() {
		for {
			conn, err := mes.ln.AcceptTCP()
			if err != nil {
				fmt.Println(err)
				return
			}
			c := &client{
				id: mes.cnt,
				conn: conn,
				recvChan: make(chan []byte, 1),
				sendChan: make(chan []byte, 10000),
				toCloseSend: make(chan int, 1),
				toCloseRecv: make(chan int, 1),
				server: mes,
			}
			mes.cnt++
			//update the map for the

			mes.clients[c.id] = c
			//cli := <-mes.clients
			//cli[c.id] = c
			//mes.clients<-cli



			go c.recvLoop()
			go c.sendLoop()
		}

	}()
	go mes.broadCastLoop()
	return nil

}

func (c *client) recvLoop() {
	br := bufio.NewReader(c.conn)
	for {
		msg, err := br.ReadBytes('\n')
		if err != nil {
			c.toCloseRecv <- 1
			return
		}
		select {
		case c.server.broadChan <- []byte(msg):
			break
		//if the server shuts down, then the server stops reading the message from
		//the client
		case <- c.toCloseSend:
			c.toCloseRecv <- 1
			return
		}

	}
}


func (c *client) sendLoop() {
	for {
		select {
		case data, ok := <-c.sendChan:
			if !ok {
				return
			}
			_, err := c.conn.Write(data)
			if err != nil {
				return
			}
			//stop sending msg to this client, delete the client
			//from the map, so the server stops broadcasting msg
			//to client
		case <- c.toCloseRecv:

			delete(c.server.clients, c.id)
			//cli :=<- c.server.clients
			//delete(cli, c.id)
			//c.server.clients <- cli
			return
		}

	}
}

func (mes *multiEchoServer) broadCastLoop() {
	for {
		select {
		case data, ok :=<- mes.broadChan:
			if !ok {
				return
			}
			cli := mes.clients
			//cli :=<- mes.clients
			for _, c := range cli {
				select {
				case c.sendChan <- data:
					break
				default:
					break
				}
			}
			mes.clients = cli
		}
	}
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!
	mes.ln.Close()
	cli := mes.clients
	for _, c := range cli {
		c.conn.Close()
		c.toCloseSend <- 1
	}
	mes.clients = cli
	close(mes.broadChan)

}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!

	return len(mes.clients)
}

// TODO: add additional methods/functions below!

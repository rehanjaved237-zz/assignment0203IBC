package main

import (
	"bufio"
//	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
  "strings"
//	"encoding/gob"
	a1 "github.com/assignment01IBC_i160140"
)

var addr = flag.String("addr", "", "The address to listen to; default is \"\" (all interfaces).")
var port = flag.Int("port", 8000, "The port to listen on; default is 8000.")

var node []string

type Account struct {
	Name string
	Balance float32
	MiningStatus bool
}

var acct Account

var chainHead *a1.Block

func main() {
	flag.Parse()

	fmt.Print("[+] Enter your Name: ")
	fmt.Scanln(&acct.Name)

	src := *addr + ":" + strconv.Itoa(*port)
	listener, _ := net.Listen("tcp", src)


	defer listener.Close()

	for i:=0;i<5;i++{
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go handleConnection(conn)
	}

  time.Sleep(500 * time.Millisecond)
	acct.MiningStatus = true

  // establish connections now

		for i:=0;i<len(node);i++ {

			conn2, err2 := net.Dial("tcp", node[i])

  	if err2 != nil {
  		if _, t := err2.(*net.OpError); t {
  			fmt.Println("Some problem connecting.")
  		} else {
  			fmt.Println("Unknown error: " + err2.Error())
  		}
  		os.Exit(1)
  	}

		writeAddr(conn2, i)

		conn2.Close()
  }

	time.Sleep(1 * time.Second)

	for {
		input := ""

		fmt.Printf("\n<---- Welcome to BitCoin Network, %s ---->\n", acct.Name)
		fmt.Printf("[1]- Press 1 to view Account balance\n")
		fmt.Printf("[2]- Press 2 to make a transaction\n")
		fmt.Printf("[3]- Press 3 to view connected peers\n")
		fmt.Printf("[4]- Press 4 to view blockchain\n")

		fmt.Scanln(&input)

		if input == "1" {
			fmt.Printf("%s's Account Balance: %f\n",acct.Name, acct.Balance)
		} else if input == "2" {
			fmt.Println("Format(SenderName-send-ReceiverName-amount-coins)")
			trans1 := a1.Transaction{}
			trans1.Input()

			writeTransaction(trans1)
		} else if input == "3" {

			fmt.Printf("\nConnected Nodes:\n")
			for i:=0;i<len(node);i++{
				fmt.Println(node[i])
			}

		} else if input == "4" {
			a1.ListBlocks(chainHead)
		}
	}

}

func writeTransaction(msg a1.Transaction) {
	for i:=0;i<len(node);i++ {

		conn, err := net.Dial("tcp", node[i])

		if err != nil {
			if _, t := err.(*net.OpError); t {
				fmt.Println("Some problem connecting.")
			} else {
				fmt.Println("Unknown error: " + err.Error())
			}
			continue
		}

		wrStr := "/trans," + msg.Trans[0] + "\n"

		n, err := conn.Write([]byte(wrStr))
		if err != nil {
			fmt.Println("Error in here.")
		}
		fmt.Println("Bytes written: ", n)

		conn.Close()
	}
}

func writeAddr(conn net.Conn, i int) {
	writeString := "/addr2,localhost:8000"
	for i < len(node) {
		writeString += "+"
		writeString += node[i]
		i++
	}
	writeString += "\n"

	_, err := conn.Write([]byte(writeString))
	if err != nil {
		fmt.Println("Error in write connection")
	}
}

var noofclients int
// Read
func handleConnection(conn net.Conn) {
	noofclients += 1
	fmt.Printf("%d Client connected.\n", noofclients)

  for {
    message, err := bufio.NewReader(conn).ReadString('\n')
    if err != nil {
      break
    }
		fmt.Println(message)

    handleMessage(message, conn)
  }

  fmt.Println("Disconnected")
}

func handleMessage(message string, conn net.Conn) {
	message = message[:len(message)-1]
  s := strings.Split(message, ",")
	fmt.Println("success")

	if len(s[0]) > 0 && s[0][0] == '/' {
		switch {
    case s[0] == "/addr":
      node = append(node, s[1])

		case s[0] == "/time":
			resp := "It is " + time.Now().String() + "\n"
			fmt.Print("< " + resp)
			conn.Write([]byte(resp))

		case s[0] == "/trans":
			fmt.Println(s[1])

			if acct.MiningStatus == true {
				var t a1.Transaction
				t.Add(s[1])
				a1.InsertBlock(t, chainHead)
			}

		case s[0] == "/quit":
			fmt.Println("Quitting.")
			conn.Write([]byte("I'm shutting down now.\n"))
			fmt.Println("< " + "%quit%")
			conn.Write([]byte("%quit%\n"))
			os.Exit(0)

		default:
			conn.Write([]byte("Unrecognized command.\n"))
		}
	}
}

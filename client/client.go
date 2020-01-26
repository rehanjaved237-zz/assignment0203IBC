
package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
//	"encoding/gob"
	"time"
//	"bytes"
	a1 "../assignment01IBC_i160140"
)

var host = flag.String("host", "localhost", "The hostname or IP to connect to; defaults to \"localhost\".")
var port = flag.Int("port", 8002, "The port to connect to; defaults to 8000.")

var node []string

var chainHead *a1.Block

type Account struct {
	Name string
	Balance float32
	MiningStatus bool
}

var acct Account

func main() {
	flag.Parse()

	ownAddr := *host + ":" + strconv.Itoa(*port)

	fmt.Print("[+] Enter your Name: ")
	fmt.Scanln(&acct.Name)
	acct.MiningStatus = false

	conn, err := net.Dial("tcp", "localhost:8000")

	if err != nil {
		if _, t := err.(*net.OpError); t {
			fmt.Println("Some problem connecting.")
		} else {
			fmt.Println("Unknown error: " + err.Error())
		}
		os.Exit(1)
	}

  _, err = conn.Write([]byte("/addr," + ownAddr + "\n"))
  if err != nil {
    fmt.Println("Error writing to stream.")
  }

  listener, err := net.Listen("tcp", ownAddr)
  if err != nil {
    fmt.Println("Error : ", err)
  }

	defer listener.Close()

	go acceptClient(listener)

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

		}
	}
}

func acceptClient(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Some connection error: %s\n", err)
		}

		go readConnection(conn)
	}
}

func readConnection(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			conn.Close()
			break
		}

		handleMessage(message, conn)
	}
}

func writeTransaction(msg a1.Transaction) {
	for i:=0;i<len(node);i++ {

		fmt.Println("In here")
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

func handleMessage(message string, conn net.Conn) {
	message = message[:len(message)-1]
  s := strings.Split(message, ",")

	if len(s[0]) > 0 && s[0][0] == '/' {
		switch {
    case s[0] == "/addr":
      node = append(node, s[1])

		case s[0] == "/addr2":
			//fmt.Println(s[1])
			str1 := strings.Split(s[1], "+")
			for _, msg := range str1 {
				node = append(node, msg)
			}

		case s[0] == "/trans":
			fmt.Println(s[1])

			if acct.MiningStatus == true {
				var t a1.Transaction
				t.Add(s[1])
				a1.InsertBlock(t, chainHead)
			}


/*			fmt.Printf("\nReceived Nodes:\n")
			for i:=0;i<len(node);i++{
				fmt.Println(node[i])
			}*/


		case s[0] == "/time":
			resp := "It is " + time.Now().String() + "\n"
			fmt.Print("< " + resp)
			conn.Write([]byte(resp))

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

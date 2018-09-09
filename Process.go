package main

import (
	"encoding/json"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

//Variáveis globais interessantes para o processo
var err string
var myPort string          // porta do meu servidor
var nServers int           // qtde de outros processo
var CliConn []*net.UDPConn // vetor com conexões para os servidores dos outros processos
var ServConn *net.UDPConn  // conexão do meu servidor (onde recebo mensagens dos outros processos)
var CriticalSessionConn *net.UDPConn // vetor com conexões para os servidores dos outros processos

var ID int
var myClock int
var state string
var numOfReplies int
var requestQueue []request

type jsonMSg struct {
	MyId int
	MyClock int
	MyMsg string
}

type request struct {
	Id int
	Clock int
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

func MaxFunc(a int, b int) int {
	if a > b { return a }
	return b
}

func doServerJob() {

	// Ler (uma vez somente) da conexão UDP a mensagem
	buf := make([]byte, 1024) // Buffer de tamanho 1024

	n, addr, err := ServConn.ReadFromUDP(buf) // Escuta a mensagem
	CheckError(err)

	var message jsonMSg
	err = json.Unmarshal(buf[:n], &message)
	CheckError(err)

	fmt.Println("Received ", string(buf[0:n]), " from ", addr, " - Process ", message.MyId) // Imprime a mensagem lida

	// Atualiza clock
	myClock = 1 + MaxFunc(myClock, message.MyClock)

	if message.MyMsg == "Request" {
		processRequest(message.MyId, message.MyClock);	
	} else if message.MyMsg == "Reply"{
		processReply()
	}
}

func doClientJob(otherProcess int, content string) {	
	// Atualiza clock
	myClock++

	// Envia uma mensagem para o servidor do processo otherServer
	msg := jsonMSg { 
					ID, 
					myClock, 
					content,
				}

	jsonSerialized, err := json.Marshal(msg) // Serializar o JSON
	CheckError(err)

	_, err = CliConn[otherProcess -1].Write(jsonSerialized)
	CheckError(err)

	time.Sleep(time.Second * 1) // Espera 1 segundo
}

func initConnections() {
	ID, _ = strconv.Atoi(os.Args[1])
	myPort = os.Args[ID+1]
	nServers = len(os.Args) - 2 // Tira o nome (no caso Process) e tira a primeira porta(que é a minha). As demais portas são dos outros processos

	//	Outros códigos para deixar ok a conexão do meu servidor
	ServerAddr, err := net.ResolveUDPAddr("udp", myPort)
	CheckError(err)

	ServConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(	err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	//	Outros códigos para deixar ok as conexões com os servidores dos outros processos
	for i := 2; i < len(os.Args); i++ {
		ServerAddr, err := net.ResolveUDPAddr("udp", os.Args[i])
		CheckError(err)

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		CliConn = append(CliConn, Conn)
		CheckError(err)
	}

	// Conectando com o recurso compartilhado
	ServerAddr, err = net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)

	CriticalSessionConn, err = net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	// Inicia clocks
	myClock = 0

	// Inicia estado
	state = "RELEASED"

	// Inicia fila de requisições
	requestQueue = make([]request, 0)
}

func readInput(ch chan string) {
	// Non-blocking async routine to listen for terminal input
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		ch <- string(text)
	}
}

func requestCriticalSession() {
	state = "WANTED"
	for i := 0; i < len(CliConn); i++ {
		if i+1 != ID { doClientJob(i + 1, "Request") }
	}
}

func processReply() {
	numOfReplies++;
	if numOfReplies == nServers-1 {
		numOfReplies = 0
		accessCriticalSession()
	}
}

func processRequest(id int, clock int) {
	if state == "HELD" || (state == "WANTED" && myClock < clock) {
		requestQueue = append(requestQueue, request {id, clock})
	} else {
		doClientJob(id, "Reply")
	}
}

func accessCriticalSession() {
	state = "HELD"
	fmt.Println("Entrando na CS")

	msg := "Processo: " + strconv.Itoa(ID) + " - " + "Clock: " + strconv.Itoa(myClock) + " na CS"
	buf := []byte(msg)
	_, err := CriticalSessionConn.Write(buf)
	CheckError(err)

	time.Sleep(time.Second * 20)
	fmt.Println("Saindo da CS")

	state = "RELEASED"
	clearQueue()
}

func clearQueue() {
	for len(requestQueue) > 0 {
		doClientJob(requestQueue[0].Id, "Reply")
		requestQueue = requestQueue[1:] // Deque
	}
}

func main() {
	initConnections()
	ch := make(chan string)

	//O fechamento de conexões devem ficar aqui, assim só fecha
	//conexão quando a main morrer
	defer ServConn.Close()

	for i := 0; i < nServers; i++ {
		defer CliConn[i].Close()
	}
	//Todo Process fará a mesma coisa: ouvir msg e mandar infinitos i’s para os outros processos
	for {
		//Server
		go doServerJob()

		go readInput(ch)

		// When there is a request (from stdin). Do it!
		select {
		case x, valid := <-ch:
			if valid {
				fmt.Printf("Recebi do teclado: %s \n", x)

				if x == strconv.Itoa(ID) {
					myClock++
					fmt.Println("Meu Clock: ", myClock)
				} else {
					if state == "RELEASED" {
						fmt.Println("Solicitando recurso compartilhado!")
						requestCriticalSession()
					} else {
						fmt.Println(x, "ignorado. Processo está no estado: ", state)
					}
				}
			} else {
				fmt.Println("Channel closed!")
			}
		default:
			// Do nothing in the non-blocking approach.
			time.Sleep(time.Second * 1)
		}

		// Wait a while
		time.Sleep(time.Second * 1)
	}
}

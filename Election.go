package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"strings"
)

//Variáveis globais interessantes para o processo
var err string
var myPort string          // porta do meu servidor
var nServers int           // qtde de outros processo
var CliConn []*net.UDPConn // vetor com conexões para os servidores dos outros processos
var ServConn *net.UDPConn  // conexão do meu servidor (onde recebo mensagens dos outros processos)

var myId string
var ID int

func maxInt( a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

func election(msgType string, receivedId int) {
	nextID := maxInt( (ID + 1) % (nServers + 1), 1)

	if msgType == "S" {
		if receivedId == ID {
			fmt.Println("I'm the leader!")
			doClientJob(nextID, "F" + myId)
		} else {
			msg := "S" + strconv.Itoa(maxInt(receivedId, ID))
			doClientJob(nextID, msg)
		}
	} 
	if msgType == "F" {
		if receivedId == ID {
			fmt.Println("Election finished!")
		} else {
			msg := "F" + strconv.Itoa(receivedId)
			doClientJob(nextID, msg)
		}
		
	}
}

func doServerJob() {
	// Ler (uma vez somente) da conexão UDP a mensagem
	// Escreve na tela a msg recebida
	buf := make([]byte, 1024) // Buffer de tamanho 1024

	n, addr, err := ServConn.ReadFromUDP(buf) // Escuta a mensagem
	CheckError(err)

	// Mensagem
	msgType := string(buf[0])
	id,_ := strconv.Atoi(string(buf[1:n]))

	fmt.Println("Recebeu ", msgType, id, "de ", addr) // Imprime a mensagem lida
	
	election(msgType, id)
}

func doClientJob(otherProcess int, msg string) {

	time.Sleep(time.Second * 5) // Espera 5 segundo

	// Envia uma mensagem para o servidor do processo otherServer
	buf := []byte(msg)
	_, err := CliConn[otherProcess-1].Write(buf) // Escreve no canal

	if err != nil {
		CheckError(err)
	} else {
		fmt.Println("Enviou", msg)
	}
}

func initConnections() {
	myId = os.Args[1]
	ID, _ = strconv.Atoi(myId)
	myPort = os.Args[ID+1]
	nServers = len(os.Args) - 2 // Tira o nome (no caso Process) e tira a primeira porta(que é a minha). As demais portas são dos outros processos

	//	Outros códigos para deixar ok a conexão do meu servidor
	ServerAddr, err := net.ResolveUDPAddr("udp", myPort)
	CheckError(err)

	ServConn, err = net.ListenUDP("udp", ServerAddr)
	CheckError(err)

	//	Outros códigos para deixar ok as conexões com os servidores dos outros processos
	for i := 2; i < len(os.Args); i++ {
		ServerAddr, err := net.ResolveUDPAddr("udp", os.Args[i])
		CheckError(err)

		LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
		CheckError(err)

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		CliConn = append(CliConn, Conn)
		CheckError(err)
	}
}

func readInput(ch chan string) {
	// Non-blocking async routine to listen for terminal input
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _, _ := reader.ReadLine()
		ch <- string(text)
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
				
				s := strings.Split(x, " ")

				otherProcess,_ := strconv.Atoi(s[0])
				msg := s[1]

				doClientJob(otherProcess, msg)
				
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


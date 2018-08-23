package main

import (
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

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

func PrintError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
	}
}

func doServerJob() {
	// Ler (uma vez somente) da conexão UDP a mensagem
	// Escreve na tela a msg recebida

	buf := make([]byte, 1024) // Buffer de tamanho 1024

	n, addr, err := ServConn.ReadFromUDP(buf)                  // Escuta a mensagem
	fmt.Println("Received ", string(buf[0:n]), " from ", addr) // Imprime a mensagem lida

	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func doClientJob(otherProcess int, i int) {
	// Envia uma mensagem (com valor i) para o servidor do processo
	// otherServer

	msg := strconv.Itoa(i)
	buf := []byte(msg)
	_, err := CliConn[otherProcess].Write(buf) // Escreve no canal

	if err != nil {
		fmt.Println(msg, err)
	}

	time.Sleep(time.Second * 1) // Espera 1 segundo
}

func initConnections() {
	myPort = os.Args[1]
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

package main

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func main() {
	/*
	   Endereço da conexão - Servidor
	       - UDP
	       - porta 10001

	   Utiliza o pacote net
	*/
	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:10001")
	CheckError(err)

	/*
	   Endereço da conexão - Local
	       - UDP
	       - porta 10001

	   Utiliza o pacote net
	*/
	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	/*
	   Conexão dial-up
	       - Endereço do Cliente:  127.0.0.1:0
	       - Endereço do Servidor: 127.0.0.1:10001
	*/
	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	/*
	   Fecha a conexão após o servidor ler as mensagens
	*/
	defer Conn.Close()
	i := 0
	for {
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_, err := Conn.Write(buf) // Escreve no canal
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1) // Espera 1 segundo
	}
}

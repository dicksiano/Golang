package main

import (
	"fmt"
	"net"
	"os"
)

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	/*
	   Endereço da conexão
	       - UDP
	       - porta 10001

	   Utiliza o pacote net
	*/
	ServerAddr, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)

	/*
	   Servidor começa a escutar pacotes da porta 10001
	*/
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)

	/*
	   Fecha a conexão após o servidor ler as mensagens
	*/
	defer ServerConn.Close()

	buf := make([]byte, 1024) // Buffer de tamanho 1024

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)                // Escuta a mensagem
		fmt.Println("Received ", string(buf[0:n]), " from ", addr) // Imprime a mensagem lida

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

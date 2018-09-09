package main

import (
	"fmt"
	"net"
	"os"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Erro: ", err)
		os.Exit(0)
	}
}

func main() {
	Address, err := net.ResolveUDPAddr("udp", ":10001")
	CheckError(err)
	Connection, err := net.ListenUDP("udp", Address)
	CheckError(err)
	defer Connection.Close()
	for {
		//Loop infinito para receber mensagem e escrever todo
		//conteúdo (processo que enviou, seu relógio e texto)
		//na tela	
		buf := make([]byte, 1024) // Buffer de tamanho 1024	
		n, addr, err := Connection.ReadFromUDP(buf)                // Escuta a mensagem
		fmt.Println("Received ", string(buf[0:n]), " from ", addr) // Imprime a mensagem lida

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}
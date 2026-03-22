package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func StartServer(porta string) error {
	ln, err := net.Listen("tcp", ":"+porta)
	if err != nil {
		return err
	}
	defer ln.Close()

	log.Printf("Servidor à escuta na porta %s\n", porta)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			r := bufio.NewReader(c)

			msg, err := r.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					log.Println("read error:", err)
				}
				return
			}
			fmt.Printf("Mensagem de %s: %s", c.RemoteAddr(), msg)
		}(conn)
	}
}

func connectToPeer(ip_destino, porta string) error {
	addr := net.JoinHostPort(ip_destino, porta)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = fmt.Fprintln(conn, "Olá da máquina A!")
	return err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go [server|client]")
		return
	}

	modo := os.Args[1]

	if modo == "server" {
		// Inicia o servidor na porta 5555
		err := StartServer("5555")
		if err != nil {
			log.Fatal("Erro ao iniciar servidor:", err)
		}
	} else if modo == "client" {
		// Liga-se ao nosso próprio computador (IP local: 127.0.0.1) na porta 5555
		err := connectToPeer("127.0.0.1", "5555")
		if err != nil {
			log.Fatal("Erro ao ligar ao servidor:", err)
		}
		fmt.Println("Mensagem enviada com sucesso!")
	} else {
		fmt.Println("Modo desconhecido. Escolhe 'server' ou 'client'.")
	}
}

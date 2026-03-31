package main

import (
	"fmt"
	"net/http"
	"os"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// var clientes = make(map[*websocket.Conn]bool)

type Mensagem struct {
	Tipo string `json:"tipo"`
	NomeSala string `json:"nomeSala"`
	Password string `json:"password"`
	Dados string `json:"dados"`
}

type Sala struct {
	Password string
	Clientes map[*websocket.Conn]bool
}

var salas = make(map[string]*Sala)

func lidarComLigacao(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro", err)
		return
	}

	defer ws.Close()

	clientes[ws] = true

	for {
		tipo, mensagem, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Erro", err)
			delete(clientes, ws)
			break
		}

		fmt.Println("Recebido:", string(mensagem))

		for cliente := range clientes {
			if cliente != ws {
				cliente.WriteMessage(tipo, mensagem)
			}
		}
	}
}

func main() {

	porta := os.Getenv("PORT")

	if porta == "" {
		porta = "8080"
	}

	http.HandleFunc("/sinal", lidarComLigacao)

	fmt.Println("o servidor está a trabalhar na porta"+porta)

	err := http.ListenAndServe(":"+porta, nil)
	if err != nil {
		fmt.Println("Erro", err)
	}
}

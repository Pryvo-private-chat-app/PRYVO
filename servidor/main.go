package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clientes = make(map[*websocket.Conn]bool)

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
	http.HandleFunc("/sinal", lidarComLigacao)

	fmt.Println("o servidor está a trabalhar na porta 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erro", err)
	}
}

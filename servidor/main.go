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
	Oferta string
	Clientes map[*websocket.Conn]bool
}

var salas = make(map[string]*Sala)

// func lidarComLigacao(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println("Erro", err)
// 		return
// 	}

// 	defer ws.Close()

// 	clientes[ws] = true

// 	for {
// 		tipo, mensagem, err := ws.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Erro", err)
// 			delete(clientes, ws)
// 			break
// 		}

// 		fmt.Println("Recebido:", string(mensagem))

// 		for cliente := range clientes {
// 			if cliente != ws {
// 				cliente.WriteMessage(tipo, mensagem)
// 			}
// 		}
// 	}
// }

func lidarComLigacao(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro", err)
		return
	}
	defer ws.Close()

	for {
		_, mensagemRecebida, err := ws.ReadMessage()
		if err != nil {
			break // Se alguém fechar a app, sai do ciclo
		}

		// 1. A Rececionista abre o envelope JSON
		var envelope Mensagem
		err = json.Unmarshal(mensagemRecebida, &envelope)
		if err != nil {
			continue // Se não for um envelope JSON válido, ignora
		}

		// 2. O que é que a pessoa quer fazer?
		if envelope.Tipo == "criar" {
			// Cria a gaveta da sala com a password e guarda logo a Oferta
			salas[envelope.NomeSala] = &Sala{
				Password: envelope.Password,
				Oferta:   envelope.Dados,
				Clientes: make(map[*websocket.Conn]bool),
			}
			// Coloca o anfitrião dentro da sala
			salas[envelope.NomeSala].Clientes[ws] = true
			fmt.Println("Sala criada:", envelope.NomeSala)

		} else if envelope.Tipo == "entrar" {
			// Procura a sala no registo
			salaExiste := salas[envelope.NomeSala]

			if salaExiste != nil {
				// Verifica a password!
				if salaExiste.Password == envelope.Password {
					// Password certa! Entra na sala.
					salaExiste.Clientes[ws] = true
					fmt.Println("Alguém entrou na sala:", envelope.NomeSala)

					// A rececionista entrega a Oferta do anfitrião diretamente ao convidado
					ws.WriteMessage(websocket.TextMessage, []byte(salaExiste.Oferta))
				} else {
					// Password errada - Barra a entrada
					ws.WriteMessage(websocket.TextMessage, []byte("ERRO: Password incorreta!"))
				}
			} else {
				// Sala não existe
				ws.WriteMessage(websocket.TextMessage, []byte("ERRO: Essa sala não existe!"))
			}

		} else if envelope.Tipo == "resposta" {
			// O convidado gerou a Resposta e quer enviá-la para o Anfitrião
			salaExiste := salas[envelope.NomeSala]
			if salaExiste != nil {
				// Envia a resposta para o Anfitrião (que é a outra pessoa na sala)
				for cliente := range salaExiste.Clientes {
					if cliente != ws {
						cliente.WriteMessage(websocket.TextMessage, []byte(envelope.Dados))
					}
				}
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

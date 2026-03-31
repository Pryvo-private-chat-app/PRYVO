package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var conexaoServidor *websocket.Conn

func SetupWebRTC(db *sql.DB) (*webrtc.PeerConnection, *webrtc.DataChannel, error) {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, nil, err
	}

	dataChannel, err := peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		return nil, nil, err
	}

	dataChannel.OnOpen(func() {
		fmt.Println("\n[+] Canal de dados aberto! A encriptação DTLS está ativa. Podes começar a falar.")
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("\n[Amigo]: %s\n", string(msg.Data))
		GravarMensagem(db, "Amigo", string(msg.Data))
	})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("\n[Amigo]: %s\n", string(msg.Data))
			GravarMensagem(db, "Amigo", string(msg.Data))
		})
	})

	return peerConnection, dataChannel, nil
}

func Encode(obj interface{}) string {
	BytesdoJson, err := json.Marshal(obj)

	if err != nil {
		panic(err)
	}

	textoBase64 := base64.StdEncoding.EncodeToString(BytesdoJson)
	return textoBase64
}

func Decode(texto string, obj interface{}) {
	bytesDoJson, err := base64.StdEncoding.DecodeString(texto)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bytesDoJson, obj)
	if err != nil {
		panic(err)
	}

}

// func escutarServidor() {
// 	for {
// 		_, mensagem, err := conexaoServidor.ReadMessage()
// 		if err != nil {
// 			fmt.Println("Erro, ", err)
// 			break
// 		}
// 		fmt.Println("Central diz: ", string(mensagem))
// 	}
// }

func main() {

	var err error

	conexaoServidor, _, err = websocket.DefaultDialer.Dial("wss://pryvo-central.onrender.com/sinal", nil)
	if err != nil {
		fmt.Println("Erro, ", err)
		return
	}

	db := InitDB()

	defer db.Close()

	var escolha string
	fmt.Println("Escolhe: [1] Criar Sala, [2] Entrar numa Sala ou [3] Limpar Histórico de Conversas")
	fmt.Scanln(&escolha)
	if escolha == "1" {
		peerConnection, dataChannel, err := SetupWebRTC(db)
		if err != nil {
			log.Fatal("Erro fatal ao iniciar WebRTC:", err)
		}

		gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

		offer, err := peerConnection.CreateOffer(nil)
		if err != nil {
			log.Fatal(err)
		}

		err = peerConnection.SetLocalDescription(offer)
		if err != nil {
			log.Fatal(err)
		}

		<-gatherComplete

		offerBase64 := Encode(peerConnection.LocalDescription())

		fmt.Println("\n=== A ENVIAR CONVITE PARA A CENTRAL ===")
		conexaoServidor.WriteMessage(websocket.TextMessage, []byte(offerBase64))

		fmt.Println("À espera que um amigo entre na sala...")

		_, mensagem, err := conexaoServidor.ReadMessage()
		if err != nil {
			log.Fatal("Erro a ler da central:", err)
		}

		codigoAmigo := string(mensagem)
		fmt.Println("Resposta do amigo recebida! A ligar P2P...")

		var answer webrtc.SessionDescription
		Decode(codigoAmigo, &answer)

		err = peerConnection.SetRemoteDescription(answer)
		if err != nil {
			log.Fatal(err)
		}

		leitorChat := bufio.NewReader(os.Stdin)

		for {
			mensagem, _ := leitorChat.ReadString('\n')
			mensagem = strings.TrimSpace(mensagem)

			if mensagem == "" {
				continue
			}

			dataChannel.SendText(mensagem)
			GravarMensagem(db, "Eu", mensagem)

		}

	} else if escolha == "2" {

		peerConnection, dataChannel, err := SetupWebRTC(db)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("À espera do convite do anfitrião na Central...")
		_, mensagem, err := conexaoServidor.ReadMessage()
		if err != nil {
			fmt.Println("Erro, ", err)
		}
		codigoAmigo := string(mensagem)

		fmt.Println("Código recebido! A gerar resposta...")

		var offer webrtc.SessionDescription
		Decode(codigoAmigo, &offer)

		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			log.Fatal(err)
		}

		gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			log.Fatal(err)
		}

		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			log.Fatal(err)
		}

		<-gatherComplete

		answerBase64 := Encode(peerConnection.LocalDescription())

		fmt.Println("A enviar a nossa resposta para o Anfitrião através da Central...")

		conexaoServidor.WriteMessage(websocket.TextMessage, []byte(answerBase64))

		fmt.Println("Resposta enviada! Ligação P2P direta estabelecida. Podes começar a teclar!")

		leitorChat := bufio.NewReader(os.Stdin)

		for {
			mensagem, _ := leitorChat.ReadString('\n')
			mensagem = strings.TrimSpace(mensagem)

			if mensagem == "" {
				continue
			}

			dataChannel.SendText(mensagem)
			GravarMensagem(db, "Eu", mensagem)
		}
	} else if escolha == "3" {
		LimparHistorico(db)
	}
}

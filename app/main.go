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

	"github.com/pion/webrtc/v3"
)

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

func main() {

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

		fmt.Println("\n=== SALA CRIADA COM SUCESSO ===")
		fmt.Println("Copia o código abaixo e envia ao teu amigo:")
		fmt.Println("--------------------------------------------------")
		fmt.Println(offerBase64)
		fmt.Println("--------------------------------------------------")
		fmt.Println("Fico à espera do do teu Amigo. Cola aqui o código e prime Enter")

		leitor := bufio.NewReader(os.Stdin)
		codigoAmigo, _ := leitor.ReadString('\n')
		codigoAmigo = strings.TrimSpace(codigoAmigo)

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

		fmt.Println("Cola o código que o teu amigo enviou e prime Enter:")
		leitor := bufio.NewReader(os.Stdin)
		codigoAmigo, _ := leitor.ReadString('\n')
		codigoAmigo = strings.TrimSpace(codigoAmigo)

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

		fmt.Println("\n=== SUCESSO: RESPOSTA GERADA ===")
		fmt.Println("Copia o código abaixo e devolve ao teu amigo (Host):")
		fmt.Println("--------------------------------------------------")
		fmt.Println(answerBase64)
		fmt.Println("--------------------------------------------------")

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

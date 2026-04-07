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

type Mensagem struct {
	Tipo     string `json:"tipo"`
	NomeSala string `json:"nomeSala"`
	Password string `json:"password"`
	Dados    string `json:"dados"`
}

type MensagemChat struct {
	Nome  string `json:"nome"`
	Foto  string `json:"foto"`
	Texto string `json:"texto"`
}

var conexaoServidor *websocket.Conn

var meuNome, minhaFoto string

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
		var msgRecebida MensagemChat
		err := json.Unmarshal(msg.Data, &msgRecebida)

		if err != nil {
			fmt.Printf("\n[Desconhecido]: %s\n", string(msg.Data))
			GravarMensagem(db, "Desconhecido", string(msg.Data))
			return
		}
		fmt.Printf("\n[%s]: %s\n", msgRecebida.Nome, msgRecebida.Texto)
		GravarMensagem(db, msgRecebida.Nome, msgRecebida.Texto)
	})

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			var msgRecebida MensagemChat
			err := json.Unmarshal(msg.Data, &msgRecebida)

			if err != nil {
				fmt.Printf("\n[Desconhecido]: %s\n", string(msg.Data))
				GravarMensagem(db, "Desconhecido", string(msg.Data))
				return
			}

			fmt.Printf("\n[%s]: %s\n", msgRecebida.Nome, msgRecebida.Texto)
			GravarMensagem(db, msgRecebida.Nome, msgRecebida.Texto)
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

	var err error

	conexaoServidor, _, err = websocket.DefaultDialer.Dial("wss://pryvo-central.onrender.com/sinal", nil)
	if err != nil {
		fmt.Println("Erro, ", err)
		return
	}

	db := InitDB()

	nomeLido, fotoLida, existe := LerPerfil(db)

	switch existe {
	case false:
		fmt.Println("Bem vindo ao PRYVO.\nQual é o teu nome?")
		fmt.Scanln(&meuNome)
		fmt.Println("Escolhe uma foto de perfil:")
		fmt.Scanln(&minhaFoto)
		GravarPerfil(db, meuNome, minhaFoto)
	case true:
		meuNome = nomeLido
		minhaFoto = fotoLida
		fmt.Println("Bem-vindo/a, " + meuNome + "!")
	}

	defer db.Close()

	var escolha string
	fmt.Println("Escolhe: [1] Criar Sala, [2] Entrar numa Sala, [3] Limpar Histórico de Conversas, [4] Ler Histórico de Conversa")
	fmt.Scanln(&escolha)
	if escolha == "1" {
		var password string
		var sala string
		fmt.Println("Cria um nome para a tua sala:")
		fmt.Scanln(&sala)
		fmt.Println("Cria uma password:")
		fmt.Scanln(&password)
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
		envelope := Mensagem{
			Tipo:     "criar",
			NomeSala: sala,
			Password: password,
			Dados:    offerBase64,
		}
		conexaoServidor.WriteJSON(envelope)

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

			envelopeP2P := MensagemChat{
				Nome:  meuNome,
				Foto:  minhaFoto,
				Texto: mensagem,
			}

			bytesDoEnvelope, err := json.Marshal(envelopeP2P)
			if err != nil {
				fmt.Println("Erro", err)
				continue
			}

			dataChannel.SendText(string(bytesDoEnvelope))
			GravarMensagem(db, "Eu", mensagem)

		}

	} else if escolha == "2" {
		var password string
		var sala string
		fmt.Println("Qual é o nome da sala?")
		fmt.Scanln(&sala)
		fmt.Println("Qual é a password?")
		fmt.Scanln(&password)
		peerConnection, dataChannel, err := SetupWebRTC(db)
		if err != nil {
			log.Fatal(err)
		}
		envelopeEntrada := Mensagem{
			Tipo:     "entrar",
			NomeSala: sala,
			Password: password,
			Dados:    "",
		}
		conexaoServidor.WriteJSON(envelopeEntrada)

		fmt.Println("À espera do convite do anfitrião na Central...")
		_, mensagem, err := conexaoServidor.ReadMessage()
		if err != nil {
			fmt.Println("Erro, ", err)
		}

		codigoAmigo := string(mensagem)

		if strings.Contains(codigoAmigo, "ERRO") {
			fmt.Println("o Servidor recusou a entrada:", codigoAmigo)
			return

		}

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

		envelopeResposta := Mensagem{
			Tipo:     "resposta",
			NomeSala: sala,
			Password: password,
			Dados:    answerBase64,
		}
		conexaoServidor.WriteJSON(envelopeResposta)

		fmt.Println("Resposta enviada! Ligação P2P direta estabelecida. Podes começar a teclar!")

		leitorChat := bufio.NewReader(os.Stdin)

		for {
			mensagem, _ := leitorChat.ReadString('\n')
			mensagem = strings.TrimSpace(mensagem)

			if mensagem == "" {
				continue
			}

			envelopeP2P := MensagemChat{
				Nome:  meuNome,
				Foto:  minhaFoto,
				Texto: mensagem,
			}

			bytesDoEnvelope, err := json.Marshal(envelopeP2P)
			if err != nil {
				fmt.Println("Erro", err)
				continue
			}

			dataChannel.SendText(string(bytesDoEnvelope))
			GravarMensagem(db, "Eu", mensagem)
		}
	} else if escolha == "3" {
		LimparHistorico(db)
	} else if escolha == "4" {
		LerHistorico(db)
	}
}

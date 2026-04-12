// package main

// import (
// 	"context"
// 	"database/sql"
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strings"

// 	"github.com/gorilla/websocket"
// 	"github.com/pion/webrtc/v3"
// 	"github.com/wailsapp/wails/v2/pkg/runtime"
// )

// type App struct {
// 	ctx       context.Context
// 	db        *sql.DB
// 	ws        *websocket.Conn
// 	pc        *webrtc.PeerConnection
// 	dc        *webrtc.DataChannel
// 	meuNome   string
// 	minhaFoto string
// }

// type SinalMensagem struct {
// 	Tipo     string `json:"tipo"`
// 	NomeSala string `json:"nomeSala"`
// 	Password string `json:"password"`
// 	Dados    string `json:"dados"`
// }

// type ChatMensagem struct {
// 	Nome  string `json:"nome"`
// 	Foto  string `json:"foto"`
// 	Texto string `json:"texto"`
// }

// func Encode(obj interface{}) string {
// 	BytesdoJson, _ := json.Marshal(obj)
// 	return base64.StdEncoding.EncodeToString(BytesdoJson)
// }

// func Decode(texto string, obj interface{}) {
// 	bytesDoJson, _ := base64.StdEncoding.DecodeString(texto)
// 	json.Unmarshal(bytesDoJson, obj)
// }

// func (a *App) SetupWebRTC() error {
// 	config := webrtc.Configuration{
// 		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
// 	}

// 	peerConnection, err := webrtc.NewPeerConnection(config)
// 	if err != nil {
// 		return err
// 	}
// 	a.pc = peerConnection

// 	dataChannel, err := peerConnection.CreateDataChannel("chat", nil)
// 	if err != nil {
// 		return err
// 	}
// 	a.dc = dataChannel

// 	configurarCanal := func(d *webrtc.DataChannel) {
// 		d.OnMessage(func(msg webrtc.DataChannelMessage) {
// 			var msgRecebida ChatMensagem
// 			json.Unmarshal(msg.Data, &msgRecebida)

// 			GravarMensagem(a.db, msgRecebida.Nome, msgRecebida.Texto)

// 			runtime.EventsEmit(a.ctx, "mensagem_recebida", msgRecebida)
// 		})
// 	}

// 	dataChannel.OnOpen(func() {
// 		runtime.EventsEmit(a.ctx, "room_ready", "P2P conection established! You can start chatting.")
// 	})

// 	configurarCanal(dataChannel)
// 	// peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
// 	// 	configurarCanal(d)
// 	// })
// 	// return nil
// 	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
// 		a.dc = d

// 		d.OnOpen(func() {
// 			runtime.EventsEmit(a.ctx, "room_ready", "P2P conection established! You can start chatting.")
// 		})

// 		configurarCanal(d)
// 	})

// 	return nil
// }

// func NewApp() *App {
// 	return &App{}
// }

// func (a *App) startup(ctx context.Context) {
// 	a.ctx = ctx
// 	a.db = InitDB()
// }

// // Greet returns a greeting for the given name
// func (a *App) Greet(name string) string {
// 	return fmt.Sprintf("Hello %s, It's show time!", name)
// }

// func (a *App) shutdown(ctx context.Context) {
// 	if a.db != nil {
// 		a.db.Close()
// 	}
// }

// func (a *App) VerificarPerfil() map[string]interface{} {
// 	nome, foto, existe := LerPerfil(a.db)

// 	return map[string]interface{}{
// 		"nome":   nome,
// 		"foto":   foto,
// 		"existe": existe,
// 	}
// }

// func (a *App) GravarPerfil(nome, foto string) bool {
// 	GravarPerfil(a.db, nome, foto)
// 	return true
// }

// func (a *App) EscolherFoto() string {
// 	caminhoFicheiro, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
// 		Title: "Escolhe a tua foto",
// 		Filters: []runtime.FileFilter{
// 			{DisplayName: "Imagens", Pattern: "*.jpg;*.jpeg;*.png"},
// 		},
// 	})

// 	if err != nil || caminhoFicheiro == "" {
// 		return ""
// 	}

// 	bytes, err := os.ReadFile(caminhoFicheiro)
// 	if err != nil {
// 		fmt.Println("Erro a ler ficheiro:", err)
// 		return ""
// 	}

// 	mimeType := http.DetectContentType(bytes)

// 	base64Str := base64.StdEncoding.EncodeToString(bytes)

// 	return "data:" + mimeType + ";base64," + base64Str
// }
// func (a *App) EnviarMensagemNet(texto string) {
// 	if a.dc != nil {
// 		mensagem := ChatMensagem{
// 			Nome:  a.meuNome,
// 			Foto:  a.minhaFoto,
// 			Texto: texto,
// 		}

// 		bytesDoJson, err := json.Marshal(mensagem)
// 		if err == nil {
// 			a.dc.SendText(string(bytesDoJson))
// 		}
// 	}
// }

// func (a *App) CriarSalaNet(codigoSala string, nome string, foto string) string {
// 	a.meuNome = nome
// 	a.minhaFoto = foto
// 	conn, _, err := websocket.DefaultDialer.Dial("wss://pryvo-central.onrender.com/sinal", nil)
// 	if err != nil {
// 		return "Erro al ligar à central"
// 	}

// 	a.ws = conn

// 	a.SetupWebRTC()

// 	offer, err := a.pc.CreateOffer(nil)
// 	if err != nil {
// 		return "Erro ao criar convite"
// 	}

// 	gatherComplete := webrtc.GatheringCompletePromise(a.pc)
// 	a.pc.SetLocalDescription(offer)
// 	<-gatherComplete

// 	textoBase64 := Encode(a.pc.LocalDescription())

// 	sinal := SinalMensagem{
// 		Tipo:     "criar",
// 		NomeSala: codigoSala,
// 		Password: "",
// 		Dados:    textoBase64,
// 	}

// 	a.ws.WriteJSON(sinal)

// 	go func() {
// 		_, msg, err := a.ws.ReadMessage()
// 		if err != nil {
// 			return
// 		}

// 		var answer webrtc.SessionDescription
// 		Decode(string(msg), &answer)
// 		a.pc.SetRemoteDescription(answer)
// 	}()

// 	return "ok"
// }

// func (a *App) EntrarSalaNet(codigoSala, nome, foto string) string {
// 	a.meuNome = nome
// 	a.minhaFoto = foto
// 	conn, _, err := websocket.DefaultDialer.Dial("wss://pryvo-central.onrender.com/sinal", nil)
// 	if err != nil {
// 		return "Erro ao legar à central"
// 	}

// 	a.ws = conn

// 	a.SetupWebRTC()

// 	sinal := SinalMensagem{
// 		Tipo:     "entrar",
// 		NomeSala: codigoSala,
// 		Password: "",
// 		Dados:    "",
// 	}

// 	a.ws.WriteJSON(sinal)
// 	_, msg, err := a.ws.ReadMessage()
// 	if err != nil || strings.Contains(string(msg), "ERRO") {
// 		return "Sala não encontrada"
// 	}

// 	var offer webrtc.SessionDescription
// 	Decode(string(msg), &offer)

// 	err = a.pc.SetRemoteDescription(offer)
// 	if err != nil {
// 		return "Erro"
// 	}

// 	gatherComplete := webrtc.GatheringCompletePromise(a.pc)

// 	answer, err := a.pc.CreateAnswer(nil)
// 	if err != nil {
// 		return "Erro"
// 	}

// 	err = a.pc.SetLocalDescription(answer)
// 	if err != nil {
// 		return "Erro"
// 	}
// 	<-gatherComplete

// 	answerBase64 := Encode(a.pc.LocalDescription())
// 	envelopeResposta := SinalMensagem{
// 		Tipo:     "resposta",
// 		NomeSala: codigoSala,
// 		Password: "",
// 		Dados:    answerBase64,
// 	}
// 	a.ws.WriteJSON(envelopeResposta)

// 	return "ok"

// }

// func (a *App) SairDaSalaNet() {
// 	if a.pc != nil {
// 		a.pc.Close()
// 	}
// 	if a.ws != nil {
// 		a.ws.Close()
// 	}
// }

package main

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx       context.Context
	db        *sql.DB
	ws        *websocket.Conn
	pc        *webrtc.PeerConnection
	dc        *webrtc.DataChannel
	meuNome   string
	minhaFoto string
}

type SinalMensagem struct {
	Tipo     string `json:"tipo"`
	NomeSala string `json:"nomeSala"`
	Password string `json:"password"`
	Dados    string `json:"dados"`
}

type ChatMensagem struct {
	Nome  string `json:"nome"`
	Foto  string `json:"foto"`
	Texto string `json:"texto"`
}

func Encode(obj interface{}) string {
	BytesdoJson, _ := json.Marshal(obj)
	return base64.StdEncoding.EncodeToString(BytesdoJson)
}

func Decode(texto string, obj interface{}) {
	bytesDoJson, _ := base64.StdEncoding.DecodeString(texto)
	json.Unmarshal(bytesDoJson, obj)
}

func (a *App) SetupWebRTC() error {
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{{URLs: []string{"stun:stun.l.google.com:19302"}}},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}
	a.pc = peerConnection

	dataChannel, err := peerConnection.CreateDataChannel("chat", nil)
	if err != nil {
		return err
	}
	a.dc = dataChannel

	configurarCanal := func(d *webrtc.DataChannel) {
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Println("📥 MENSAGEM RECEBIDA DA NET:", string(msg.Data))

			var msgRecebida ChatMensagem
			err := json.Unmarshal(msg.Data, &msgRecebida)
			if err != nil {
				fmt.Println("❌ ERRO a traduzir JSON:", err)
			}

			// ⚠️ Gravação desligada temporariamente para testar!
			// GravarMensagem(a.db, msgRecebida.Nome, msgRecebida.Texto)

			fmt.Println("📣 A enviar a mensagem para o ecrã do Wails...")
			runtime.EventsEmit(a.ctx, "mensagem_recebida", msgRecebida)
		})
	}

	dataChannel.OnOpen(func() {
		runtime.EventsEmit(a.ctx, "room_ready", "P2P conection established! You can start chatting.")
	})

	configurarCanal(dataChannel)

	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		a.dc = d

		d.OnOpen(func() {
			runtime.EventsEmit(a.ctx, "room_ready", "P2P conection established! You can start chatting.")
		})

		configurarCanal(d)
	})

	return nil
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.db = InitDB()
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

func (a *App) VerificarPerfil() map[string]interface{} {
	nome, foto, existe := LerPerfil(a.db)

	return map[string]interface{}{
		"nome":   nome,
		"foto":   foto,
		"existe": existe,
	}
}

func (a *App) GravarPerfil(nome, foto string) bool {
	GravarPerfil(a.db, nome, foto)
	return true
}

func (a *App) EscolherFoto() string {
	caminhoFicheiro, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Escolhe a tua foto",
		Filters: []runtime.FileFilter{
			{DisplayName: "Imagens", Pattern: "*.jpg;*.jpeg;*.png"},
		},
	})

	if err != nil || caminhoFicheiro == "" {
		return ""
	}

	bytes, err := os.ReadFile(caminhoFicheiro)
	if err != nil {
		fmt.Println("Erro a ler ficheiro:", err)
		return ""
	}

	mimeType := http.DetectContentType(bytes)

	base64Str := base64.StdEncoding.EncodeToString(bytes)

	return "data:" + mimeType + ";base64," + base64Str
}

// func (a *App) EnviarMensagemNet(texto string) {
// 	fmt.Println("👉 A preparar para enviar:", texto)

// 	if a.dc != nil {
// 		mensagem := ChatMensagem{
// 			Nome:  a.meuNome,
// 			Foto:  a.minhaFoto,
// 			Texto: texto,
// 		}

// 		bytesDoJson, err := json.Marshal(mensagem)
// 		if err == nil {
// 			fmt.Println("👉 A enviar para o túnel WebRTC...")
// 			a.dc.SendText(string(bytesDoJson))
// 		} else {
// 			fmt.Println("❌ ERRO a empacotar JSON:", err)
// 		}
// 	} else {
// 		fmt.Println("❌ ERRO GRAVE: O túnel (a.dc) está desligado/nil!")
// 	}
// }

func (a *App) EnviarMensagemNet(texto string) {
	fmt.Println("👉 A preparar para enviar:", texto)

	if a.dc != nil {
		mensagem := ChatMensagem{
			Nome:  a.meuNome,
			Foto:  "", // ⚠️ O TRUQUE ESTÁ AQUI: Não enviamos a foto para não entupir!
			Texto: texto,
		}

		bytesDoJson, err := json.Marshal(mensagem)
		if err == nil {
			fmt.Println("👉 A enviar para o túnel WebRTC... Tamanho:", len(bytesDoJson), "bytes")

			// Atiramos para o túnel, mas agora VERIFICAMOS se ele aceitou!
			errEnvio := a.dc.SendText(string(bytesDoJson))
			if errEnvio != nil {
				fmt.Println("❌ O TÚNEL REJEITOU A MENSAGEM:", errEnvio)
			} else {
				fmt.Println("👉 O túnel engoliu a mensagem com sucesso!")
			}
		} else {
			fmt.Println("❌ ERRO a empacotar JSON:", err)
		}
	} else {
		fmt.Println("❌ ERRO GRAVE: O túnel (a.dc) está desligado/nil!")
	}
}

func (a *App) CriarSalaNet(codigoSala string, nome string, foto string) string {
	a.meuNome = nome
	a.minhaFoto = foto
	conn, _, err := websocket.DefaultDialer.Dial("wss://pryvo-central.onrender.com/sinal", nil)
	if err != nil {
		return "Erro ao ligar à central"
	}

	a.ws = conn

	a.SetupWebRTC()

	offer, err := a.pc.CreateOffer(nil)
	if err != nil {
		return "Erro ao criar convite"
	}

	gatherComplete := webrtc.GatheringCompletePromise(a.pc)
	a.pc.SetLocalDescription(offer)
	<-gatherComplete

	textoBase64 := Encode(a.pc.LocalDescription())

	sinal := SinalMensagem{
		Tipo:     "criar",
		NomeSala: codigoSala,
		Password: "",
		Dados:    textoBase64,
	}

	a.ws.WriteJSON(sinal)

	go func() {
		_, msg, err := a.ws.ReadMessage()
		if err != nil {
			return
		}

		var answer webrtc.SessionDescription
		Decode(string(msg), &answer)
		a.pc.SetRemoteDescription(answer)
	}()

	return "ok"
}

func (a *App) EntrarSalaNet(codigoSala, nome, foto string) string {
	a.meuNome = nome
	a.minhaFoto = foto
	conn, _, err := websocket.DefaultDialer.Dial("wss://pryvo-central.onrender.com/sinal", nil)
	if err != nil {
		return "Erro ao ligar à central"
	}

	a.ws = conn

	a.SetupWebRTC()

	sinal := SinalMensagem{
		Tipo:     "entrar",
		NomeSala: codigoSala,
		Password: "",
		Dados:    "",
	}

	a.ws.WriteJSON(sinal)
	_, msg, err := a.ws.ReadMessage()
	if err != nil || strings.Contains(string(msg), "ERRO") {
		return "Sala não encontrada"
	}

	var offer webrtc.SessionDescription
	Decode(string(msg), &offer)

	err = a.pc.SetRemoteDescription(offer)
	if err != nil {
		return "Erro"
	}

	gatherComplete := webrtc.GatheringCompletePromise(a.pc)

	answer, err := a.pc.CreateAnswer(nil)
	if err != nil {
		return "Erro"
	}

	err = a.pc.SetLocalDescription(answer)
	if err != nil {
		return "Erro"
	}
	<-gatherComplete

	answerBase64 := Encode(a.pc.LocalDescription())
	envelopeResposta := SinalMensagem{
		Tipo:     "resposta",
		NomeSala: codigoSala,
		Password: "",
		Dados:    answerBase64,
	}
	a.ws.WriteJSON(envelopeResposta)

	return "ok"
}

func (a *App) SairDaSalaNet() {
	if a.pc != nil {
		a.pc.Close()
	}
	if a.ws != nil {
		a.ws.Close()
	}
}

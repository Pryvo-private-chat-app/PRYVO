package main

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// )

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	db  *sql.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.db = InitDB()
}

// Greet returns a greeting for the given name
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

// EscolherFoto abre o explorador nativo do Linux/Windows/Mac
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

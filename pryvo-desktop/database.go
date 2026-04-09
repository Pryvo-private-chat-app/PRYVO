package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "./chat.db")
	if err != nil {
		log.Fatal("Erro ao abrir a base de dados: ", err)
	}

	db.SetMaxOpenConns(1)

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS mensagens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		remetente TEXT,
		texto TEXT,
		data_hora DATETIME DEFAULT CURRENT_TIMESTAMP
	);
		CREATE TABLE IF NOT EXISTS perfil (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT,
		foto TEXT
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Erro ao criar a tabela: ", err)
	}

	return db
}

func LerPerfil(db *sql.DB) (string, string, bool) {
	var nome, foto string

	err := db.QueryRow("SELECT nome, foto FROM perfil ORDER BY id DESC LIMIT 1").Scan(&nome, &foto)

	if err != nil {
		return "", "", false
	}

	return nome, foto, true
}

func GravarPerfil(db *sql.DB, nome string, foto string) {

	instrucao := `Insert into perfil (nome, foto) VALUES (?, ?)`

	_, err := db.Exec(instrucao, nome, foto)

	if err != nil {
		log.Println("Erro", err)
	}
}

func GravarMensagem(db *sql.DB, remetente string, texto string) {

	instrucao := `Insert Into mensagens (remetente, texto) VALUES (?, ?)`

	_, err := db.Exec(instrucao, remetente, texto)

	if err != nil {
		log.Println("Erro ao gravar mensagem:", err)
	}
}

func LimparHistorico(db *sql.DB) {
	_, err := db.Exec("DELETE FROM mensagens")
	if err != nil {
		log.Println("Erro ao apagar histórico:", err)
	} else {
		fmt.Println("\n[!] O histórico de mensagens foi apagado com sucesso.")
	}
}

func LerHistorico(db *sql.DB) {
	linhas, err := db.Query("SELECT data_hora, remetente, texto FROM mensagens ORDER BY id ASC")
	if err != nil {
		log.Println("Erro", err)
		return
	}
	defer linhas.Close()

	fmt.Println("\n=== HISTÓRICO DE MENSAGENS ===")

	for linhas.Next() {
		var data, remetente, msg string

		err := linhas.Scan(&data, &remetente, &msg)
		if err != nil {
			fmt.Println("Erro", err)
			continue
		}

		fmt.Printf("[%s] %s: %s\n", data, remetente, msg)
	}
	fmt.Println("==============================\n")
}

package main

import (
	"fmt"
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "./chat.db")
	if err != nil {
		log.Fatal("Erro ao abrir a base de dados: ", err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS mensagens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		remetente TEXT,
		texto TEXT,
		data_hora DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("Erro ao criar a tabela: ", err)
	}

	return db
}

func GravarMensagem (db *sql.DB, remetente string, texto string) {

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
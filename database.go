package main

import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./chat.db")
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
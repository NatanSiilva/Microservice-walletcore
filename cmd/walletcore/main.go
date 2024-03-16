package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/NatanSiilva/ms-wallet/internal/database"
	"github.com/NatanSiilva/ms-wallet/internal/event"
	"github.com/NatanSiilva/ms-wallet/internal/usecase/create_account"
	"github.com/NatanSiilva/ms-wallet/internal/usecase/create_client"
	"github.com/NatanSiilva/ms-wallet/internal/usecase/create_transaction"
	"github.com/NatanSiilva/ms-wallet/internal/web"
	"github.com/NatanSiilva/ms-wallet/internal/web/webserver"
	"github.com/NatanSiilva/ms-wallet/pkg/events"
	_ "github.com/go-sql-driver/mysql"
)

func connectToDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:root@tcp(172.18.0.2:3306)/wallet?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}

func main() {
	db, err := connectToDatabase()
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreated()
	// eventDispatcher.Register("TransactionCreated", handle)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	transactionDb := database.NewTransactionDB(db)

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(transactionDb, accountDb, eventDispatcher, transactionCreatedEvent)

	webserver := webserver.NewWebServer(":3000")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Server is running")

	webserver.Start()

}
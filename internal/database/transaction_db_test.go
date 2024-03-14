package database

import (
	"database/sql"
	"testing"

	"github.com/NatanSiilva/ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	client        *entity.Client
	client2       *entity.Client
	accountFrom   *entity.Account
	accountTo     *entity.Account
	transactionDB *TransactionDB
}

func (su *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	su.Nil(err)
	su.db = db

	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date)")
	db.Exec("CREATE TABLE transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)")

	client, err := entity.NewClient("john", "j@j.com")
	su.Nil(err)
	su.client = client

	client2, err := entity.NewClient("john2", "j@jj.com")
	su.Nil(err)
	su.client2 = client2

	accountFrom := entity.NewAccount(su.client)
	accountFrom.Balance = 1000
	su.accountFrom = accountFrom

	accountTo := entity.NewAccount(su.client2)
	accountTo.Balance = 1000
	su.accountTo = accountTo

	su.transactionDB = NewTransactionDB(db)
}

func (suite *TransactionDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE clients")
	suite.db.Exec("DROP TABLE accounts")
	suite.db.Exec("DROP TABLE transactions")
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}

func (su *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(su.accountFrom, su.accountTo, 100)
	su.Nil(err)

	err = su.transactionDB.Create(transaction)
	su.Nil(err)
}

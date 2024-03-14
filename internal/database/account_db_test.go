package database

import (
	"database/sql"
	"testing"

	"github.com/NatanSiilva/ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
	client    *entity.Client
}

func (su *AccountDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	su.Nil(err)
	su.db = db

	db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date, updated_at date)")

	su.accountDB = NewAccountDB(db)
	su.client, _ = entity.NewClient("john", "j@j.com")
}

func (suite *AccountDBTestSuite) TearDownSuite() {
	defer suite.db.Close()
	suite.db.Exec("DROP TABLE clients")
	suite.db.Exec("DROP TABLE accounts")
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}

func (suite *AccountDBTestSuite) TestSave() {
	account := entity.NewAccount(suite.client)
	err := suite.accountDB.Save(account)
	suite.Nil(err)
}

func (su *AccountDBTestSuite) TestFindById() {
	su.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?,?,?,?)",
		su.client.ID, su.client.Name, su.client.Email, su.client.CreatedAt,
	)

	account := entity.NewAccount(su.client)
	err := su.accountDB.Save(account)
	su.Nil(err)

	accountDB, err := su.accountDB.FindById(account.ID)

	su.Nil(err)
	su.Equal(account.ID, accountDB.ID)
	su.Equal(account.Client.ID, accountDB.Client.ID)
	su.Equal(account.Balance, accountDB.Balance)
	su.Equal(account.Client.Name, accountDB.Client.Name)
	su.Equal(account.Client.Email, accountDB.Client.Email)

}

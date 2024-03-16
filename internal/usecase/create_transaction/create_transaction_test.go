package createtransaction

import (
	"testing"

	"github.com/NatanSiilva/ms-wallet/internal/entity"
	"github.com/NatanSiilva/ms-wallet/internal/event"
	"github.com/NatanSiilva/ms-wallet/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (mo *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := mo.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (mock *AccountGatewayMock) Save(account *entity.Account) error {
	args := mock.Called(account)
	return args.Error(0)
}

func (mock *AccountGatewayMock) FindById(id string) (*entity.Account, error) {
	args := mock.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func TestCreateTransactionUseCase_Execute(t *testing.T) {
	t.Run("should create a transaction", func(t *testing.T) {
		client1, _ := entity.NewClient("client1", "j@j1.com")
		account1 := entity.NewAccount(client1)
		account1.Credit(1000)

		client2, _ := entity.NewClient("client2", "j@j2.com")
		account2 := entity.NewAccount(client2)
		account2.Credit(1000)

		accountGatewayMock := &AccountGatewayMock{}
		accountGatewayMock.On("FindById", account1.ID).Return(account1, nil)
		accountGatewayMock.On("FindById", account2.ID).Return(account2, nil)

		transactionGatewayMock := &TransactionGatewayMock{}
		transactionGatewayMock.On("Create", mock.Anything).Return(nil)

		inputDTO := CreateTransactionInputDTO{
			AccountIDFrom: account1.ID,
			AccountIDTo:   account2.ID,
			Amount:        100,
		}

		dispatcher := events.NewEventDispatcher()
		event := event.NewTransactionCreated()

		useCase := NewCreateTransactionUseCase(transactionGatewayMock, accountGatewayMock, dispatcher, event)
		outputDTO, err := useCase.Execute(inputDTO)

		assert.Nil(t, err)
		assert.NotNil(t, outputDTO)
		accountGatewayMock.AssertExpectations(t)
		transactionGatewayMock.AssertExpectations(t)
		accountGatewayMock.AssertNumberOfCalls(t, "FindById", 2)
		transactionGatewayMock.AssertNumberOfCalls(t, "Create", 1)
	})
}

package create_account

import (
	"testing"

	"github.com/NatanSiilva/ms-wallet/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientGatewayMock struct {
	mock.Mock
}

func (mock *ClientGatewayMock) Save(client *entity.Client) error {
	args := mock.Called(client)
	return args.Error(0)
}

func (mock *ClientGatewayMock) Get(id string) (*entity.Client, error) {
	args := mock.Called(id)
	return args.Get(0).(*entity.Client), args.Error(1)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (mock *AccountGatewayMock) Save(account *entity.Account) error {
	args := mock.Called(account)
	return args.Error(0)
}

func (mock *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := mock.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func TestCreateAccountUseCase_Execute(t *testing.T) {
	t.Run("should create account", func(t *testing.T) {
		client, _ := entity.NewClient("John Doe", "j@j.com")
		clientGatewayMock := &ClientGatewayMock{}
		clientGatewayMock.On("Get", client.ID).Return(client, nil)

		accountGatewayMock := &AccountGatewayMock{}
		accountGatewayMock.On("Save", mock.Anything).Return(nil)

		useCase := NewCreateAccountUseCase(accountGatewayMock, clientGatewayMock)
		inputDTO := CreateAccountInputDTO{
			ClientID: client.ID,
		}
		outputDTO, err := useCase.Execute(inputDTO)

		assert.Nil(t, err)
		assert.NotNil(t, outputDTO)
		clientGatewayMock.AssertExpectations(t)
		accountGatewayMock.AssertExpectations(t)
		clientGatewayMock.AssertNumberOfCalls(t, "Get", 1)
		accountGatewayMock.AssertNumberOfCalls(t, "Save", 1)
	})
}

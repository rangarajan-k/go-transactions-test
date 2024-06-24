package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go-transactions-test/config"
	"go-transactions-test/datastore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := datastore.NewMockIPostgressStore(ctrl)
	mockClient := datastore.NewMockIPgClientInterface(ctrl)
	config := &config.MainConfig{}

	//r := gin.Default()
	mockClient.EXPECT().GetDbClient().Return(&pg.DB{}).AnyTimes()
	c := NewTransactionServiceController(config, mockClient, mockStore)
	// Test case: Successful account creation
	t.Run("Successful account creation", func(t *testing.T) {
		createAccountRequest := CreateAccountRequest{
			AccountId:      1,
			InitialBalance: 100.0,
		}
		body, _ := json.Marshal(createAccountRequest)

		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		mockStore.EXPECT().CreateAccountQuery(gomock.Any()).Return(nil)

		c.CreateAccount(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test case: Missing account ID
	t.Run("Missing account ID", func(t *testing.T) {
		createAccountRequest := CreateAccountRequest{
			InitialBalance: 100.0,
		}
		body, _ := json.Marshal(createAccountRequest)

		req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		c.CreateAccount(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestQueryAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := datastore.NewMockIPostgressStore(ctrl)
	mockClient := datastore.NewMockIPgClientInterface(ctrl)
	config := &config.MainConfig{}

	mockClient.EXPECT().GetDbClient().Return(&pg.DB{}).AnyTimes()
	c := NewTransactionServiceController(config, mockClient, mockStore)

	// Test case: Successful account query
	t.Run("Successful account query", func(t *testing.T) {
		account := &datastore.Account{AccountId: 1, BalanceAmount: 100.0}

		req, _ := http.NewRequest("GET", "/accounts/1", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{gin.Param{Key: "account_id", Value: "1"}}
		ctx.Request = req

		mockStore.EXPECT().GetAccountDetailsQuery(gomock.Any()).Return(account, nil)

		c.QueryAccount(ctx)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test case: Account not found
	t.Run("Account not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/accounts/2", nil)
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Params = gin.Params{gin.Param{Key: "account_id", Value: "2"}}
		ctx.Request = req

		mockStore.EXPECT().GetAccountDetailsQuery(gomock.Any()).Return(nil, errors.New("pg: no rows in result set"))

		c.QueryAccount(ctx)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestSubmitTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := datastore.NewMockIPostgressStore(ctrl)
	mockClient := datastore.NewMockIPgClientInterface(ctrl)

	config := &config.MainConfig{}

	mockClient.EXPECT().GetDbClient().Return(&pg.DB{}).AnyTimes()
	c := NewTransactionServiceController(config, mockClient, mockStore)

	// Test case: Successful transaction submission
	t.Run("Successful transaction submission", func(t *testing.T) {
		submitTransactionRequest := datastore.Transaction{
			SourceAccountId:      1,
			DestinationAccountId: 2,
			Amount:               50.0,
		}
		body, _ := json.Marshal(submitTransactionRequest)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		mockStore.EXPECT().GetMultipleAccountDetailsQuery(gomock.Any(), gomock.Any()).Return([]*datastore.Account{
			{AccountId: 1, BalanceAmount: 100.0},
			{AccountId: 2, BalanceAmount: 50.0},
		}, nil)
		mockStore.EXPECT().UpdateAccountBalanceQuery(gomock.Any(), gomock.Any()).Return(nil)
		mockStore.EXPECT().SubmitTransactionQuery(gomock.Any(), gomock.Any()).Return(nil)

		c.SubmitTransaction(ctx)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test case: Insufficient balance
	t.Run("Insufficient balance", func(t *testing.T) {
		submitTransactionRequest := datastore.Transaction{
			SourceAccountId:      1,
			DestinationAccountId: 2,
			Amount:               200.0,
		}
		body, _ := json.Marshal(submitTransactionRequest)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		mockStore.EXPECT().GetMultipleAccountDetailsQuery(gomock.Any(), gomock.Any()).Return([]*datastore.Account{
			{AccountId: 1, BalanceAmount: 100.0},
			{AccountId: 2, BalanceAmount: 50.0},
		}, nil)

		c.SubmitTransaction(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test case: Both source and destination accounts must exist
	t.Run("Both source and destination accounts must exist", func(t *testing.T) {
		submitTransactionRequest := datastore.Transaction{
			SourceAccountId:      1,
			DestinationAccountId: 2,
			Amount:               50.0,
		}
		body, _ := json.Marshal(submitTransactionRequest)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		// Return only one account in the query result
		mockStore.EXPECT().GetMultipleAccountDetailsQuery(gomock.Any(), gomock.Any()).Return([]*datastore.Account{
			{AccountId: 1, BalanceAmount: 100.0},
		}, nil)

		c.SubmitTransaction(ctx)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Source or Destination Account Invalid")
	})
}

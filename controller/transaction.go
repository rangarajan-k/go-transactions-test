package controller

import (
	"github.com/gin-gonic/gin"
	"go-transactions-test/config"
	"go-transactions-test/datastore"
	"go-transactions-test/dicontainer"
	"net/http"
	"strconv"
)

type ITransactionServiceController interface {
	CreateAccount(*gin.Context)
	QueryAccount(*gin.Context)
	SubmitTransaction(*gin.Context)
	HandleErrorResponse(*gin.Context, int, string)
	HandleSuccessResponse(*gin.Context, int, interface{})
}

type transactionServiceController struct {
	config    *config.MainConfig
	container dicontainer.IDiContainer
}

func NewTransactionServiceController(config *config.MainConfig, container dicontainer.IDiContainer) ITransactionServiceController {
	return &transactionServiceController{config, container}
}

func (c *transactionServiceController) CreateAccount(ctx *gin.Context) {
	var createAccountRequest datastore.Account
	err := ctx.BindJSON(&createAccountRequest)
	if err != nil {
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Error unmarshalling request")
	}

	if createAccountRequest.AccountId == 0 || createAccountRequest.BalanceAmount == 0 {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Mandatory params not present")
	}

	err = datastore.CreateAccountQuery(c.container, createAccountRequest)
	if err != nil {
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Error when creating account")
	}
	c.HandleSuccessResponse(ctx, http.StatusCreated, createAccountRequest)
}

func (c *transactionServiceController) QueryAccount(ctx *gin.Context) {
	req := ctx.Request
	accountId := req.URL.Query().Get("account_id")
	var AccountModel datastore.Account
	if accountId == "" {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Mandatory params not present")
	}
	id, err := strconv.Atoi(accountId)
	if err != nil {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid account id")
	}
	result, readErr := datastore.GetAccountDetailsQuery(c.container, AccountModel, id)
	if readErr != nil {
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Error when querying account")
	}
	c.HandleSuccessResponse(ctx, http.StatusOK, result)

}

func (c *transactionServiceController) SubmitTransaction(ctx *gin.Context) {
	var submitTransactionRequest datastore.Transaction
	err := ctx.BindJSON(&submitTransactionRequest)
	if err != nil {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Error unmarshalling request")
	}

	if submitTransactionRequest.SourceAccountId == 0 || submitTransactionRequest.DestinationAccountId == 0 || submitTransactionRequest.SourceAccountId == submitTransactionRequest.DestinationAccountId {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid params in the request")
	}

	err = datastore.SubmitTransactionQuery(c.container, submitTransactionRequest)
	if err != nil {
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Error when submitting transaction")
	}
	c.HandleSuccessResponse(ctx, http.StatusCreated, submitTransactionRequest)
}

func (c *transactionServiceController) HandleErrorResponse(ctx *gin.Context, statusCode int, error_message string) {
	ctx.IndentedJSON(statusCode, gin.H{"error": error_message})
}
func (c *transactionServiceController) HandleSuccessResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.IndentedJSON(statusCode, data)
}

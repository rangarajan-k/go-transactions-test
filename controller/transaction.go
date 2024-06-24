package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"go-transactions-test/config"
	"go-transactions-test/datastore"
	"log"
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
	config   *config.MainConfig
	pgClient *pg.DB
}

type CreateAccountRequest struct {
	AccountId      int     `json:"account_id"`
	InitialBalance float32 `json:"initial_balance"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewTransactionServiceController(config *config.MainConfig, pgClient *pg.DB) ITransactionServiceController {
	return &transactionServiceController{config, pgClient}
}

// CreateAccount godoc
// @Summary Creates an account for the customer
// @Accept json
// @Product json
// @Param payload body CreateAccountRequest true "Request Payload"
// @Success 201 {object} datastore.Account
// @Failure 400	{object} ErrorResponse{}
// @Failure 500 {object} ErrorResponse{}
// @Router /accounts [post]
func (c *transactionServiceController) CreateAccount(ctx *gin.Context) {
	var createAccountRequest *CreateAccountRequest
	err := ctx.BindJSON(&createAccountRequest)
	if err != nil {
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if createAccountRequest.AccountId == 0 {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Mandatory params not present")
		return
	}

	var insertAccount = &datastore.Account{
		AccountId:     createAccountRequest.AccountId,
		BalanceAmount: createAccountRequest.InitialBalance,
	}

	err = datastore.CreateAccountQuery(c.pgClient, insertAccount)
	if err != nil {
		log.Printf("%+v", err)
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Account already exists")
		return
	}
	c.HandleSuccessResponse(ctx, http.StatusCreated, insertAccount)
}

// QueryAccount godoc
// @Summary Queries an exiting customer account based on account id
// @Accept json
// @Product json
// @Param account_id path string true "Example: 12121"
// @Success 200 {object} datastore.Account
// @Failure 400	{object} ErrorResponse{}
// @Failure 404	{object} ErrorResponse{}
// @Failure 500 {object} ErrorResponse{}
// @Router /accounts/:account_id [get]
func (c *transactionServiceController) QueryAccount(ctx *gin.Context) {
	//req := ctx.Request
	accountId := ctx.Param("account_id")
	fmt.Println(accountId)
	if accountId == "" {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Mandatory params not present")
		return
	}
	id, err := strconv.Atoi(accountId)
	if err != nil {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Account Id missing/invalid")
		return
	}
	AccountModel := &datastore.Account{
		AccountId: id,
	}
	fmt.Println("Making query")
	result, readErr := datastore.GetAccountDetailsQuery(c.pgClient, AccountModel)
	if readErr != nil {
		errMessage := fmt.Sprintf("%v", readErr)
		if errMessage == "pg: no rows in result set" {
			c.HandleErrorResponse(ctx, http.StatusNotFound, "Account not found")
			return
		}
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}
	c.HandleSuccessResponse(ctx, http.StatusOK, result)
}

// SubmitTransaction godoc
// @Summary Posts a transaction against a source account and destination account
// @Accept json
// @Product json
// @Param payload body datastore.Transaction true "Request Payload"
// @Success 200 {object} datastore.Transaction
// @Failure 400	{object} ErrorResponse{}
// @Failure 404	{object} ErrorResponse{}
// @Failure 500 {object} ErrorResponse{}
// @Router /transactions [post]
func (c *transactionServiceController) SubmitTransaction(ctx *gin.Context) {
	var submitTransactionRequest *datastore.Transaction
	err := ctx.BindJSON(&submitTransactionRequest)
	if err != nil {
		log.Printf("%+v", err)
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid params in the request")
		return
	}

	if submitTransactionRequest.SourceAccountId == 0 || submitTransactionRequest.DestinationAccountId == 0 || submitTransactionRequest.SourceAccountId == submitTransactionRequest.DestinationAccountId {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Invalid params in the request")
		return
	}

	//get the accounts first check if they are valid
	var ids = []int{submitTransactionRequest.SourceAccountId, submitTransactionRequest.DestinationAccountId}

	/* since this is a multi table update operation use transaction to rollback in case of any errors
	we can also implement some form of locking here either using dblock / distributed redis lock
	to avoid simultaneous updates on same resource */
	tx, _ := c.pgClient.Begin()
	defer tx.Close()
	results, err := datastore.GetMultipleAccountDetailsQuery(tx, ids)
	if err != nil {
		log.Printf("%+v", err)
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(results) < 2 {
		c.HandleErrorResponse(ctx, http.StatusNotFound, "Source account / Destination Account Invalid")
		return
	}
	//check if source account has sufficient balance
	sourceAccount := results[0]
	destinationAccount := results[1]
	if sourceAccount.BalanceAmount < submitTransactionRequest.Amount {
		c.HandleErrorResponse(ctx, http.StatusBadRequest, "Insufficient balance in source account")
		return
	}

	sourceAccount.BalanceAmount = sourceAccount.BalanceAmount - submitTransactionRequest.Amount
	destinationAccount.BalanceAmount = destinationAccount.BalanceAmount + submitTransactionRequest.Amount
	//update source and destination accounts
	var accountsToUpdate = []*datastore.Account{sourceAccount, destinationAccount}
	err = datastore.UpdateAccountBalanceQuery(tx, accountsToUpdate)
	if err != nil {
		log.Printf("%+v", err)
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		tx.Rollback()
		return
	}
	//write a transaction log
	err = datastore.SubmitTransactionQuery(tx, submitTransactionRequest)
	if err != nil {
		log.Printf("%+v", err)
		c.HandleErrorResponse(ctx, http.StatusInternalServerError, "Something went wrong")
		tx.Rollback()
		return
	}
	tx.Commit()
	c.HandleSuccessResponse(ctx, http.StatusCreated, submitTransactionRequest)
}

func (c *transactionServiceController) HandleErrorResponse(ctx *gin.Context, statusCode int, errorMessage string) {
	err := ErrorResponse{Error: errorMessage}
	ctx.AbortWithStatusJSON(statusCode, err)
}
func (c *transactionServiceController) HandleSuccessResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.IndentedJSON(statusCode, data)
}

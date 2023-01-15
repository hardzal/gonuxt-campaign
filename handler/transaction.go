package handler

import (
	"crowdfounding/helper"
	"crowdfounding/transaction"
	"crowdfounding/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// parameter di uri
// tangkap parameter mapping input struct
// panggil service, input struct sebagai parameter
// service, berbekal campaign id bisa digunakan untuk memanggil repo
// repo mencari data transacation suatu campaign
type transactionsHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionsHandler {
	return &transactionsHandler{service}
}

func (h *transactionsHandler) GetCampaignTransactions(c *gin.Context) {
	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign Transactions detail", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))
	c.JSON(http.StatusOK, response)
}

//GetUserTransactions
//Handler
//ambil nilai user dari jwt/middleware
//service
//repo => ambil data transactions (preload campaign)

func (h *transactionsHandler) GetUserTransactions(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transaction, err := h.service.GetTransactionByUserID(userID)

	if err != nil {
		response := helper.APIResponse("Failed to get user's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("User Transactions", http.StatusOK, "success", transaction)
	c.JSON(http.StatusOK, response)
}

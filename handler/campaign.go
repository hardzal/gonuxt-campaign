package handler

import (
	"crowdfounding/campaign"
	"crowdfounding/helper"
	"crowdfounding/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// tangkap parameter di handler
// handler ke service
// service yang menentukan repository mana yang dicall
// repository : FindAll, FindByUserID
// db

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", campaign.FormatCampaigns(campaigns))

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of all campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

// api/v1/campaign/:id
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	// handler : mapping id yang berasal dari url ke struct input => service, call formatter
	// service : inputnya struct, input => menangkap id di url manggil repo
	// repository: get campaign by id

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	// tangkap parameter dari user ke input struct
	// ambil current user dari jwt/handler
	// panggil service, parameternya input struct (buat slug juga)
	// panggil repository untuk simpan data campaign baru

	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("failed to create a campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	NewCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to create a campaign", http.StatusOK, "success", campaign.FormatCampaign(NewCampaign))
	c.JSON(http.StatusOK, response)
}

// user masukkan input
// handler
// mapping dari input ke input struct
// input dari user, termasuk input dari uri
// service
// repository update data campaign
//func
func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	var campaignID campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&campaignID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	updateCampaign, err := h.service.UpdateCampaign(campaignID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updateCampaign))
	c.JSON(http.StatusOK, response)
}

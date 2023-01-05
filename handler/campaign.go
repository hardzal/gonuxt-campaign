package handler

import (
	"crowdfounding/campaign"
	"crowdfounding/helper"
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

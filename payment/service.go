package payment

import (
	"crowdfounding/campaign"
	"crowdfounding/transaction"
	"crowdfounding/user"
	midtrans "github.com/veritrans/go-midtrans"
	"strconv"
)

type service struct {
	transactionRepository transaction.Repository
	campaignRepository    campaign.Repository
}

type Service interface {
	GetPaymentURL(transaction Transaction, user user.User) (string, error)
	ProcessPayment(input transaction.NotificationTransactionInput) error
}

func NewService(transactionRepository transaction.Repository, campaignRepository campaign.Repository) *service {
	return &service{transactionRepository, campaignRepository}
}

func (s *service) GetPaymentURL(transaction Transaction, user user.User) (string, error) {
	midclient := midtrans.NewClient()
	midclient.ServerKey = "SB-Mid-server-Uf7KU4uqToLoB2ufPMeox2ev"
	midclient.ClientKey = "SB-Mid-client-Cgqmx2GXhxAg4klA"
	midclient.APIEnvType = midtrans.Sandbox

	snapGateway := midtrans.SnapGateway{
		Client: midclient,
	}

	snapReq := &midtrans.SnapReq{
		CustomerDetail: &midtrans.CustDetail{
			Email: user.Email,
			Phone: user.Name,
		},
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(transaction.ID),
			GrossAmt: int64(transaction.Amount),
		},
	}

	snapTokenResp, err := snapGateway.GetToken(snapReq)
	if err != nil {
		return "", err
	}

	return snapTokenResp.RedirectURL, nil
}

func (s *service) ProcessPayment(input transaction.NotificationTransactionInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID) // integer to string
	transactionPayment, err := s.transactionRepository.GetByID(transaction_id)

	if err != nil {
		return err
	}

	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transactionPayment.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transactionPayment.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transactionPayment.Status = "cancelled"
	}

	updatedTransaction, err := s.transactionRepository.Update(transactionPayment)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}
	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}

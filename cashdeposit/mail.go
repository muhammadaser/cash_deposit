package cashdeposit

import (
	"net/smtp"
	"strconv"

	"github.com/muhammadaser/cash_deposit/accounts"
)

// Mail of products
type Mail interface {
	SendReceiptNotif(deposit CashDeposit) error
}

// NewMail return struct that implement store interface
func NewMail(store accounts.Store) Mail {
	return &setMail{
		store: store,
	}
}

type setMail struct {
	store accounts.Store
}

func (m *setMail) SendReceiptNotif(deposit CashDeposit) error {
	account, err := m.store.GetAccount(deposit.AccountID)
	if err != nil {
		return err
	}
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		"5ab4b484649aa3",
		"18a343ea576383",
		"smtp.mailtrap.io",
	)
	to := []string{account.Email}
	msg := []byte("To: " + account.Email + "\r\n" +
		"Subject: Cash Deposit\r\n" +
		"\r\n" +
		"Dear " + account.FirstName + " " + account.LastName + ", \r\n" +
		"\r\n" +
		"Transaksi Cash Deposit dengan ID transaksi " + deposit.DepositID + " sejumlah IDR " +
		strconv.FormatInt(deposit.DepositAmount, 10) + " telah berhasil di lakukan pada " + deposit.DepositDate.String() +
		"\r\n \r\n" +
		"Terimakasih")
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err = smtp.SendMail(
		"smtp.mailtrap.io:2525",
		auth,
		"noreply@testbank.org",
		to,
		msg,
	)
	if err != nil {
		return err
	}
	return nil
}

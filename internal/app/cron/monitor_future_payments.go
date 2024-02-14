package cron

import (
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/nighostchris/everytrack-backend/internal/database"
)

func (cj *CronJob) MonitorFuturePayments() {
	go func() {
		for {
			cj.Logger.Info("starts")

			// Get all future payments of all users in database
			futurePayments, getFuturePaymentsError := database.GetAllFuturePayments(cj.Db)
			if getFuturePaymentsError != nil {
				cj.Logger.Error(
					fmt.Sprintf("failed to get all future payment records from database. %s", getFuturePaymentsError.Error()),
				)
				return
			}

			for _, payment := range futurePayments {
				if payment.ScheduledAt.Unix() < time.Now().Unix() {
					cj.Logger.Info(fmt.Sprintf("going to process future payment %s of amount %s for account %s", payment.Id, payment.Amount, payment.AccountId))

					// Get original balance for target account
					accountBalance, getAccountBalanceError := database.GetAccountBalance(cj.Db, payment.AccountId)
					if getAccountBalanceError != nil {
						cj.Logger.Error(
							fmt.Sprintf(
								"failed to get account balance for %s from database. %s",
								payment.AccountId,
								getAccountBalanceError.Error(),
							),
						)
						return
					}

					// Converting balance to float number
					balanceInFloat, parseBalanceError := strconv.ParseFloat(accountBalance, 64)
					if parseBalanceError != nil {
						cj.Logger.Error(fmt.Sprintf("failed to parse account balance into float. %s", parseBalanceError.Error()))
						return
					}

					// Converting amount to float number
					amountInFloat, parseAmountError := strconv.ParseFloat(payment.Amount, 64)
					if parseAmountError != nil {
						cj.Logger.Error(fmt.Sprintf("failed to parse payment amount into float. %s", parseAmountError.Error()))
						return
					}

					// Calculate the final account balance after spending the expense amount
					amount := big.NewFloat(amountInFloat)
					if !payment.Income {
						amount = big.NewFloat(0).Neg(amount)
					}
					newAccountBalance := big.NewFloat(0).Add(big.NewFloat((balanceInFloat)), amount)
					newAccountBalanceInFloat, _ := newAccountBalance.Float64()
					cj.Logger.Debug(fmt.Sprintf("going to change balance for account %s from %s to %s", payment.AccountId, accountBalance, strconv.FormatFloat(newAccountBalanceInFloat, 'f', -1, 64)))

					// Update the account balance for after receiving or paying scheduled payment
					_, updateAccountBalanceError := database.UpdateAccountBalance(
						cj.Db,
						strconv.FormatFloat(newAccountBalanceInFloat, 'f', -1, 64),
						payment.AccountId,
					)
					if updateAccountBalanceError != nil {
						cj.Logger.Error(
							fmt.Sprintf("failed to update latest balance after merging with payment amount. %s", updateAccountBalanceError.Error()),
						)
						return
					}

					// Update the next schedule date according to frequency if the payment is on rolling basis
					// Otherwise delete the payment
					if payment.Rolling {
						nextScheduledDate := payment.ScheduledAt.Add(time.Duration(payment.Frequency.Int64))
						_, updateFuturePaymentScheduleError := database.UpdateFuturePaymentSchedule(cj.Db, nextScheduledDate, payment.Id)
						if updateFuturePaymentScheduleError != nil {
							// TODO Later: need to revert balance update in previous step upon failure to delete record
							cj.Logger.Error(
								fmt.Sprintf("failed to update schedule for future payment. %s", updateFuturePaymentScheduleError.Error()),
							)
							return
						}
					} else {
						_, deleteFuturePaymentError := database.DeleteFuturePayment(cj.Db, payment.Id, payment.ClientId)
						if deleteFuturePaymentError != nil {
							// TODO LATER: need to revert balance update in previous step upon failure to delete record
							cj.Logger.Error(
								fmt.Sprintf("failed to delete paid future payment. %s", deleteFuturePaymentError.Error()),
							)
							return
						}
					}
				}
			}

			time.Sleep(1 * time.Hour)
		}
	}()
}

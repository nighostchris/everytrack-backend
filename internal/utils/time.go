package utils

type PaymentFrequency struct {
	Days   int
	Months int
	Years  int
}

func CalculateActualPaymentFrequency(seconds int) PaymentFrequency {
	paymentFrequency := PaymentFrequency{
		Days:   0,
		Months: 0,
		Years:  0,
	}

	days := seconds / 86400
	if days < 29 {
		paymentFrequency.Days = days
		return paymentFrequency
	}

	months := days / 30
	if months == 0 {
		months += 1
	}
	paymentFrequency.Months = months
	return paymentFrequency
}

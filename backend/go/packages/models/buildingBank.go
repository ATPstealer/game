package models

type Bank struct {
	LoansAmount       float64 `json:"loansAmount" bson:"loansAmount"`
	LoansLimit        float64 `json:"loansLimit" bson:"loansLimit"`
	BorrowedFromState float64 `json:"borrowedFromState" bson:"borrowedFromState"`
	BorrowedLimit     float64 `json:"borrowedLimit" bson:"borrowedLimit"`
}

type CreditTerms struct {
	Limit  float64 `json:"limit" bson:"limit"`
	Rate   float64 `json:"rate" bson:"rate"`
	Rating float64 `json:"rating" bson:"rating"`
}

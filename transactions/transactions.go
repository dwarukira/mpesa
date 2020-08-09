package transactions

import (
	"log"
	"strings"
	"time"
)

type Transaction struct {
	ID              string    `json:"id"`
	Message         string    `json:"message"`
	Type            string    `json:"type"`
	Entity          string    `json:"entity"`
	Amount          string    `json:"amount"`
	Balance         string    `json:"balance"`
	TransactionCost string    `json:"transactionCost"`
	Date            string    `json:"date"`
	DateTime        time.Time `json:"dateTime"`
	Contact         string    `json:"contact"`
}

func NewFromMessage(message string) (Transaction, error) {
	tokens := strings.Split(strings.ToLower(message), " ")
	transactionID := tokens[0]
	log.Println(transactionID)
	transactionType := GetTransactionType(message)
	entity := GetEntityFromMessage(message, transactionType)
	amount := ""
	balance := ""
	transactionCost := ""

	replaceMoney := func(value string) string {
		return strings.ReplaceAll(value, "ksh", "")
	}

	transactionTime, _ := GetDateFroMessage(message, transactionType)
	contact, _ := GetContactFromEntity(entity)

	if len(contact) > 1 {
		entity = strings.ReplaceAll(entity, contact, "")
	}

	count := 0
	for _, token := range tokens {
		if strings.HasPrefix(token, "ksh") {
			if count == 0 {
				amount = replaceMoney(token)

			} else if count == 1 {
				balance = strings.TrimSuffix(replaceMoney(token), ".")

			} else if count == 2 {
				transactionCost = strings.TrimSuffix(replaceMoney(token), ".")

			}
			count++

		}

	}
	transaction := Transaction{
		Message:         message,
		ID:              transactionID,
		Type:            transactionType,
		Entity:          entity,
		Amount:          amount,
		Balance:         balance,
		TransactionCost: transactionCost,
		DateTime:        transactionTime,
		Contact:         contact,
	}

	return transaction, nil

}

package transactions

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/metakeule/fmtdate"
)

const (
	WITHDRAW = "withdraw"
	SEND     = "send"
	PAYBILL  = "paybill"
	AIRTIME  = "airtime"
	RECEIVED = "received"
	FULIZA   = "fuliza"
)

func GetTransactionType(message string) string {
	transactionType := "unknown"
	lowerCaseMessage := strings.ToLower(message)

	contains := func(text string) bool {
		return strings.Contains(lowerCaseMessage, text)

	}
	if contains("sent to") {
		transactionType = SEND

	} else if contains("withdraw") {
		transactionType = WITHDRAW

	} else if contains("paid to") {
		transactionType = PAYBILL

	} else if contains("you bought") {
		transactionType = AIRTIME

	} else if contains("you have received") {
		transactionType = RECEIVED

	} else if contains("fuliza m-pesa amount") {
		transactionType = FULIZA

	}

	return transactionType
}

func GetMatchFromRegex(reg, message string) ([]string, error) {
	lowerCaseMessage := strings.ToLower(message)
	valid, err := regexp.Compile(reg)
	if err != nil {
		return nil, err
	}
	match := valid.FindStringSubmatch(lowerCaseMessage)
	return match, nil
}

func GetEntityFromMessage(message string, transactionType string) string {
	lowerCaseMessage := strings.ToLower(message)

	extractEntity := func(reg string) string {
		valid, _ := regexp.Compile(reg)
		match := valid.FindStringSubmatch(lowerCaseMessage)
		return match[1]

	}

	if transactionType == WITHDRAW {
		return extractEntity(`from(.*) new m-pesa`)

	} else if transactionType == RECEIVED {
		return extractEntity("from(.*) on [1-9]")

	} else if transactionType == SEND {
		return extractEntity("sent to (.*) on [1-9]")

	} else if transactionType == AIRTIME {
		return "Safaricom Airtime"

	} else if transactionType == PAYBILL {
		return extractEntity("paid to(.*). on")

	}

	return ""

}

func GetDateFroMessage(message string, transactionType string) (time.Time, error) {
	// timeRegex := "([0-1]?[0-9]|2[0-3]):[0-5][0-9] pm|am"
	dateRegex := "([0-1][0-9+]/[0-9]/[0-9]+)"
	// matchTime, _ := GetMatchFromRegex(timeRegex, message)
	matchDate, _ := GetMatchFromRegex(dateRegex, message)
	log.Println(matchDate)
	datetime := " "
	transTime, err := fmtdate.Parse("D/M/YY hh:mmam", datetime)

	if err != nil {
		fmt.Println(err.Error())
		return transTime, err
	}
	return transTime, nil
}

func GetContactFromEntity(entity string) (string, error) {
	match, err := GetMatchFromRegex(`07(\d{8})`, entity)
	if len(match) <= 1 {
		return "", err
	}
	return match[0], err
}

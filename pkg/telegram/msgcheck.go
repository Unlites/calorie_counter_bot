package telegram

import (
	"log"
	"regexp"
)

func CheckText(text string) bool {
	matched, err := regexp.MatchString(`^\s*[a-z-а-яё]+(?:\s+[a-z-а-яё]+){0,}\s*$`, text)
	if err != nil {
		log.Print(err)
	}
	if !matched {
		return true
	}
	return false
}

func CheckDigits(callories string) bool {
	matched, err := regexp.MatchString(`^[0-9]*$`, callories)
	if err != nil {
		log.Print(err)
	}
	if !matched {
		return true
	}
	return false
}

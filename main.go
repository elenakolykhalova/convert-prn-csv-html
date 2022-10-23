package main

import (
	"fmt"
	"log"
	"regexp"
)

var regexBody *regexp.Regexp // Регулярка для данных

func init() {
	// Инициализация регулярного выражения
	patternBody := `(?P<First_name>\"(.*?)\")\,(?P<Address>(?:\"(.*?)\")|(.*?))\,(?P<Postcode>(.*?))\,(?P<Mobile>(.*?))\,(?P<Limit>(.*?))\,(?P<Birthday>(.*))`
	thisRegexBody, err := regexp.Compile(patternBody)
	// Обработка ошибки при использовании регулярки
	if err != nil {
		log.Fatal(err)
	}
	regexBody = thisRegexBody
}

func main() {
	fmt.Println("Start program...")
	GetUsers(csvFileType, GetHTML(csvFileType, "data.csv"))
	GetUsers(prnFileType, GetHTML(prnFileType, "data.prn"))
	fmt.Println("End program...")
}

package main

import (
	"fmt"
	"golang.org/x/text/encoding/korean"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var menuPrices map[string]int = map[string]int{
	"\xe2\x93\x9f": 1700,
	"\xe2\x93\x91": 2000,
	"\xe2\x93\x92": 2500,
	"\xe2\x93\x93": 3000,
	"\xe2\x93\x94": 3500,
	"\xe2\x93\x95": 4000,
	"\xe2\x93\x96": 4500,
	"\xe2\x93\x97": 0,
}

var menuUrlPrefix string = "http://www.snuco.com/html/restaurant/"
var menuUrls = [...]string{menuUrlPrefix + "restaurant_menu1.asp",
	menuUrlPrefix + "restaurant_menu2.asp"}

func test_textsInHtml() {
	data := "abc<def>ghi<ddd>"
	res := textsInHtml(data)
	if res[0] != "abc" {
		log.Fatal("test fail!", res)
		return
	}
	if res[1] != "ghi" {
		log.Fatal("test fail!", res)
		return
	}
	log.Println("success!")
}

func textsInHtml(s string) []string {
	start, end := 0, len(s)
	var res []string

	for i, c := range s {
		if trimmed := strings.TrimSpace(s[start:i]); trimmed != "" &&
			trimmed != "&nbsp;" && c == '<' {
			end = i
			res = append(res, strings.TrimSpace(s[start:end]))
		} else if c == '>' {
			start = i + 1
		}
	}
	if start != len(s) {
		res = append(res, s[start:])
	}
	return res
}

func toUtf8(src []byte) []byte {
	buffer := make([]byte, len(src)*2)
	transformer := korean.EUCKR.NewDecoder()
	_, nSrc, err := transformer.Transform(buffer, src, true)
	if err != nil {
		log.Fatal("error while encoding transform...", err)
	}
	if nSrc < len(src) {
		log.Printf("only %d B of web page(%d B) read")
	}

	return buffer
}

func printMenu(texts []string) {
	// breakfast, lunch, dinner 3 fields
	for i, _ := range []string{"breakfast", "lunch", "dinner"} {
		splitted := strings.Split(texts[i], "/")
		for _, menu := range splitted {
			key := string([]rune(menu)[0])
			if price, exist := menuPrices[key]; exist {
				trimmed := strings.Trim(menu, key)
				trimmed = strings.Trim(trimmed, "&nbsp;")
				fmt.Printf("%s (%d 원) / ", trimmed, price)
			}
		}
		fmt.Printf("\n")
	}
}

func getMenu(cafes []string) {
	for _, menuUrl := range menuUrls {
		resp, err := http.Get(menuUrl)
		if err != nil {
			log.Fatal("error while get ", menuUrl, err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Fatal("failed to read body", err)
		}

		bodyString := string(toUtf8(body))
		texts := textsInHtml(bodyString)

		for _, cafe := range cafes {
			for i, s := range texts {
				if s == cafe {
					fmt.Printf("[%s]\n", cafe)
					printMenu(texts[i+1:])
				}
			}
		}
	}
}

func main() {
	cafes := []string{"302동"}
	if len(os.Args) > 1 {
		cafes = os.Args[1:]
	}
	getMenu(cafes)
}

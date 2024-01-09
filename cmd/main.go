package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func getAppID(name string) string {
	resp, err := http.Get("http://store.steampowered.com/search/results/?term=" + name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	re := regexp.MustCompile(`<a href="https://store.steampowered.com/app/([0-9]+)/`)
	matches := re.FindStringSubmatch(string(body))
	if matches != nil {
		return matches[1]
	} else {
		fmt.Println("No matches found")
	}
	defer resp.Body.Close()
	return ""
}

func getPrice(id string) string {
	resp, err := http.Get("http://store.steampowered.com/api/appdetails?appids=" + id)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return ""
	}
	price := map[string]interface{}{}
	err = json.Unmarshal(body, &price)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	nestedValue := price[id].(map[string]interface{})
	finalPrice := nestedValue["data"].(map[string]interface{})["price_overview"].(map[string]interface{})["final_formatted"]
	defer resp.Body.Close()
	return finalPrice.(string)
}

func main() {
	var name string
	flag.StringVar(&name, "name", "Cyberpunk 2077", "Name of the game you want to search for")
	flag.Parse()
	fmt.Printf("Fetching price data for %v...\n", name)
	properString := strings.ReplaceAll(name, " ", "_")
	resp := getAppID(properString)
	price := getPrice(resp)
	fmt.Printf("Price of %v: %v\n", name, price)
}

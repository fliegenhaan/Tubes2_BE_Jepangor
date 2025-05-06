package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const URL = "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	result := make(map[string][][]string)

	doc.Find("table").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			cells := row.Find("td").Not("td table td")
			if cells.Length() < 2 {
				return
			}

			var element string
			var recipes [][]string
			valid := false

			cells.Each(func(k int, cell *goquery.Selection) {
				text := strings.TrimSpace(cell.Text())
				if element == "" {
					element = text
				}

				for _, line := range strings.Split(text, "\n") {
					if strings.Contains(line, "+") {
						valid = true
						parts := strings.Split(line, "+")
						var cleanedParts []string
						for _, part := range parts {
							cleanedParts = append(cleanedParts, strings.TrimSpace(part))
						}
						recipes = append(recipes, cleanedParts)
					}
				}
			})

			if valid {
				result[element] = recipes
			}
		})
	})

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		log.Fatal("Failed to save JSON file:", err)
	}
}

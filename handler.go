package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/anaskhan96/soup"
	// utils "github.com/aosousa/golang-utils"
)

const (
	baseURL = "https://howlongtobeat.com/search_results?page=1"
	version = "1.0.0"
)

// Prints the list of accepted commands
func printHelp() {
	fmt.Printf("How Long to Beat Lookup (version %s)\n", version)
	fmt.Println("Available commands:")
	fmt.Println("* -h | --help\t Prints the list of available commands")
	fmt.Println("* -v | --version Prints the version of the application")
	fmt.Println("* <title> Prints the How Long to Beat statistics for the first game found with the title received (e.g. go-hltb-lookup.exe Final Fantasy X)")
}

// Prints the current version of the application
func printVersion() {
	fmt.Printf("How Long to Beat version %s\n", version)
}

/* Handles a request to look up a game's statistics from the How Long to Beat website
 * Receives:
 * * args ([]string) - Arguments passed in the terminal by the user
 */
func handleOptions(args []string) {
	// create a string from the elements in the slice of arguments
	lenArgs := len(args)
	gameTitle := strings.Join(args[1:lenArgs], " ")

	formData := url.Values{}
	formData.Add("queryString", gameTitle)
	formData.Add("t", "games")
	formData.Add("sorthead", "popular")
	formData.Add("sortd", "Normal Order")
	formData.Add("length_type", "main")

	res, err := http.Post(baseURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	checkResponse(res.StatusCode)

	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	defer res.Body.Close()

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	printGameStatistics(gameTitle, string(content))
}

/* Print a game's statistics
 * Receives:
 * * title (string) - Title of the game
 * * document (string) - HowLongTobeat HTML document
 */
func printGameStatistics(title, document string) {
	fmt.Printf("%s Stats\n\n", title)

	doc := soup.HTMLParse(document)
	gameList := doc.Find("ul")
	gameLiTag := gameList.Children()[3]
	gameInfoDiv := gameLiTag.Children()[3]
	gameDetailsDiv := gameInfoDiv.Children()[3].Children()[1]

	for index, row := range gameDetailsDiv.Children() {
		if row.NodeValue == "div" {
			fmt.Printf(row.Text())
		}

		if index%2 != 0 {
			fmt.Println()
		}
	}
}

/* Check if there was an error in the request
 * Receives:
 * * code (int) - Response status code
 */
func checkResponse(code int) {
	if code != 200 {
		fmt.Println("Error: An error occurred while performing your request. Please try again later.")
		os.Exit(1)
	}
}

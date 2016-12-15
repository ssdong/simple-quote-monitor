package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var usage = `
  Usage: simple-quote-monitor [options]

  -se <stock exchange>     The code of stock exchange
  -ss <stock symbol>       The unique series of letters of a security
  -min <number>            The minimum price watching for
  -max <number>            The maximum price watching for
`

func main() {
	if len(os.Args) < 2 {
		man()
	}

	var symbol string
	var exchange string
	var min float64
	var max float64

	flag.StringVar(&symbol, "ss", "", "The code of stock exchange")
	flag.StringVar(&exchange, "se", "", "The series of letters of a security")
	flag.Float64Var(&min, "min", -0.1, "The minimum price watching for")
	flag.Float64Var(&max, "max", -0.1, "The maximum price watching for")

	flag.Parse()

	if symbol == "" || exchange == "" || min == -0.1 || max == -0.1 {
		log.Fatalf("Please provide all necessary arguments\n %s\n", usage)
	}

	if min > max {
		log.Fatalf("Please make sure the minimum price %f is smaller than the maximum price %f", min, max)
	}

	exchange = strings.ToUpper(exchange)
	symbol = strings.ToUpper(symbol)
	url := fmt.Sprintf("https://www.google.com/finance?q=%s:%s", exchange, symbol)

	// Ping and get the price data in every 30 - 40 seconds
	for {
		// Get the response from google finance search
		response, getError := http.Get(url)
		if getError != nil {
			log.Fatalln(getError)
		}

		defer response.Body.Close()

		// Parse the fetched html file into nodes
		node, parseError := html.Parse(response.Body)
		if parseError != nil {
			log.Fatalln(parseError)
		}

		// Starting from this line, the logic is tailored to match the id
		// on the google finance page. Right now, we are focusing on the
		// following HTML snippet structure which contains the stock price
		//
		//<div id="price-panel" class="id-price-panel goog-inline-block">
		//   <div>
		//       <span class="pr">
		//           <span id="ref_121280306767104_l">9.89</span>
		//       </span>
		//       <div class="id-price-change nwp">
		//           <span class="ch bld">
		//              <span class="chr" id="ref_121280306767104_c">-0.13</span>
		//              <span class="chr" id="ref_121280306767104_cp">(-1.30%)</span>
		//           </span>
		//       </div>
		//   </div>
		// </div>
		result := findNodes(node, "class", "pr")

		// Notice, the Node type within package net/html is really referring to a node
		// in the tree-based API for XML, the DOM. I.e. An opening tag, text or a comment
		// is also a valid node. After we find the opening tag <span class="pr">, we need
		// to further looking for the number within the its <span> child.

		// This line of code is really specific for the above DOM structure. We are getting
		// the number's text node. Notice, the trickest part is that the newline character
		// after <span class="pr"> is also another valid text node. I.e.
		// <span class="pr">
		//    (A "\n" text node here)<span id="ref_121280306767104_l">(span element node)9.89(text node)</span>

		// If the returned and parsed HTML file does not contain any node with the specified
		// (name, value) pair, we assume the provided stock symbol is actually invalid
		if len(result) == 0 {
			log.Fatalf("The stock exchange code %s or the stock symbol %s is invalid\n", exchange, symbol)
		}
		priceNode := result[0].FirstChild.NextSibling.FirstChild
		price, _ := strconv.ParseFloat(priceNode.Data, 64)

		var title string
		var message string
		if price <= min {
			title = "Price drop alert"
			message = fmt.Sprintf("The stock %s's price is now $%f\n", symbol, price)
		}
		if price >= max {
			title = "Price rise alert"
			message = fmt.Sprintf("The stock %s's price is now $%f\n", symbol, price)
		}
		fmt.Printf(message)
		alert := exec.Command("terminal-notifier", "-title", title, "-message", message)
		alert.Run()
		// Sleep randomly from 30 seconds to 40 seconds
		time.Sleep(time.Duration((30000 + rand.Intn(11)*1000)) * time.Millisecond)
	}
}

// This will print out the flag options
func man() {
	fmt.Printf("%s\n", usage)
	os.Exit(0)
}

// Parse the node(and its children) and find nodes with the passed-in (name, value) pair
func findNodes(parent *html.Node, name string, value string) []*html.Node {
	var result []*html.Node = nil

	// Check the attribute array on current node
	attributes := parent.Attr
	for _, attribute := range attributes {
		if attribute.Key == name && attribute.Val == value {
			result = append(result, parent)
		}
	}
	// Run Depth-first Search on its children and then siblings
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		result = append(result, findNodes(child, name, value)...)
	}
	return result
}

// Copyright © 2018 Sysdig
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tembleking/weavesocks-loadtest/pkg/client"
	"github.com/tembleking/weavesocks-loadtest/pkg/types"
)

var rootCmd = &cobra.Command{
	Use:   "weavesocks-loadtest",
	Short: "Creates some fake load in the weavesocks demo application",
	Long:  `Create some fake load in the weavesocks demo application available at https://microservices-demo.github.io`,
	Run:   Run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	numOfClients     *int
	delayBeforeStart *int
	hostName         *string
	numOfRequests    *int
)

func init() {
	numOfClients = rootCmd.Flags().IntP("clients", "c", 2, "Number of concurrent clients")
	delayBeforeStart = rootCmd.Flags().IntP("delay", "d", 0, "Delay before start")
	hostName = rootCmd.Flags().StringP("hostname", "n", "", "Target host url (eg. http://localhost:8080)")
	numOfRequests = rootCmd.Flags().IntP("requests", "r", 10, "Number of requests per client")
}

func Run(cmd *cobra.Command, args []string) {
	if *hostName == "" {
		fmt.Fprintln(os.Stderr, "Hostname can't be empty")
		return
	}

	if *delayBeforeStart > 0 {
		fmt.Println("Waiting", *delayBeforeStart, "seconds before starting...")
		time.Sleep(time.Duration(*delayBeforeStart) * time.Second)
	}

	wg := sync.WaitGroup{}

	fmt.Println("Running", *numOfClients, "clients with", *numOfRequests, "requests per clientRoutine to", *hostName, "...")
	catalog, err := getCatalog(*hostName)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < *numOfClients; i++ {
		wg.Add(1)
		go clientRoutine(&wg, *numOfRequests, *hostName, catalog)
	}

	wg.Wait()
}

func getCatalog(host string) (catalog types.Catalog, err error) {
	hostUrl, _ := url.Parse(host)
	catalogUrl, _ := hostUrl.Parse("/catalogue")
	response, err := http.Get(catalogUrl.String())
	if err != nil {
		err = errors.Wrap(err, "error retrieving the catalog")
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("response error: %s", response.Status))
		return
	}

	err = json.NewDecoder(response.Body).Decode(&catalog)

	if err != nil {
		err = errors.Wrap(err, "error decoding the catalog response")
		return
	}
	return
}

func getRandomItemInCatalog(catalog types.Catalog) (id string) {
	choosenCatalogElement := catalog[rand.Intn(len(catalog))]
	return choosenCatalogElement.ID
}

func clientRoutine(wg *sync.WaitGroup, maxRequests int, host string, catalog types.Catalog) {
	defer wg.Done()

	c, err := client.New(host)
	if err != nil {
		log.Println(err)
		return
	}

	for numRequests := 0; numRequests < maxRequests; {
		randomItemId := getRandomItemInCatalog(catalog)

		err = c.Get("")
		if err != nil {
			log.Println(err)
		}
		numRequests++

		err = c.Login("user", "password")
		if err != nil {
			log.Println(err)
		}
		numRequests++

		err = c.Get("category.html")
		if err != nil {
			log.Println(err)
		}
		numRequests++

		err = c.Get(fmt.Sprintf("detail.html?id=%d", randomItemId))
		if err != nil {
			log.Println(err)
		}
		numRequests++

		err = c.Delete("cart")
		if err != nil {
			log.Println(err)
		}
		numRequests++

		cartRequest := types.CartRequest{
			ID:       randomItemId,
			Quantity: 1,
		}
		buff := bytes.NewBuffer([]byte{})
		json.NewEncoder(buff).Encode(cartRequest)
		err = c.Post("cart", "application/json", buff)
		if err != nil {
			log.Println(err)
		}
		numRequests++

		err = c.Get("basket.html")
		if err != nil {
			log.Println(err)
		}
		numRequests++

		err = c.Post("orders", "", nil)
		if err != nil {
			log.Println(err)
		}
		numRequests++

	}
}

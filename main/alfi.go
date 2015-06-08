package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/codegangsta/cli"
	"net/http"
	"os"
	"strings"
)

type Result struct {
	Header 		Header		`json:"responseHeader"`
	Response 	Response	`json:"response"`
	SpellCheck 	SpellCheck	`json:"spellcheck"`
}

type Header struct {

}

type Response struct {
	NumFound 		int				`json:"numFound"`
	Repositories 	[]Repository	`json:"docs"`
}

type Repository struct {
	Id 				string	`json:"id"`
	LatestVersion 	string	`json:"latestVersion"`
}

type SpellCheck struct {
	Suggestions	[]string	`json:"suggestions"`
}

func main() {
	app := cli.NewApp()
	app.Name = "alfi"
	app.Usage = "alfi [search query]"
	app.Action = func (c *cli.Context) {
		if (c.Args().Present()) {
			request(c.Args()[0])
		} else {
			println("Missing query parameter")
			fmt.Printf("%s\n", app.Usage)
		}

	}
	app.Run(os.Args)
}

func request(query string) {
	response, err := http.Get("http://search.maven.org/solrsearch/select?q=" + query + "&rows=25&wt=json")
	println("Searching...")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", string(contents))

		var result Result;
		err = json.Unmarshal(contents, &result)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		for _, repo := range result.Response.Repositories {
			if (strings.Contains(repo.Id, "processor") || strings.Contains(repo.Id, "compiler")) {
				fmt.Printf("provided '%s:%s'\n", repo.Id, repo.LatestVersion)
			} else {
				fmt.Printf("compile '%s:%s'\n", repo.Id, repo.LatestVersion)
			}
		}
		fmt.Printf("Found %d result(s)\n", result.Response.NumFound)
	}
}
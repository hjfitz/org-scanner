package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"

	"github.com/google/go-github/github"
	gitleaks "github.com/zricethezav/gitleaks/src"
)

var wg sync.WaitGroup

type issues struct {
	passwordLeaks   int
	vulnerabilities []ScanIssue
}

func checkRepos(repos []*github.Repository) {
	for _, repo := range repos {
		repoURL := *repo.GitURL
		name := *repo.Name
		options := &gitleaks.Options{GithubURL: repoURL}

		// check for passcodes in git history
		leakCount, err := gitleaks.Run(options)

		// check for vulnerabilities in repo
		NodeScan(repoURL)

		if err != nil {
			log.Fatal(err)
		}

		if leakCount == 0 {
			fmt.Printf("%d leaks found for '%s'\n", leakCount, name)
		} else {
			fmt.Printf("!! %d leaks found for '%s'\n", leakCount, name)
		}
	}
	wg.Done()
}

func main() {
	fmt.Printf("There are %d threads available\n", runtime.NumCPU())
	client := github.NewClient(nil)
	// list public repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, err := client.Repositories.ListByOrg(context.Background(), "github", opt)

	if err != nil {
		log.Fatal(err)
	}

	// split repos in between threads
	threads := runtime.NumCPU()
	groupCount := len(repos) / threads
	fmt.Printf("Found %d repos\n", len(repos))
	fmt.Printf("Giving %d repos to each thread\n", groupCount)

	NodeScan(*repos[0].GitURL)

	// // calculate a deficit (rounding errors when splitting work between threads)
	// deficit := len(repos) - (groupCount * threads)
	// fmt.Printf("There is a deficit of %d groups. Adding to thread...\n", deficit)
	// wg.Add(1)
	// go checkRepos(repos[deficit:])
	// repos = repos[deficit:]

	// beginning := 0
	// end := groupCount
	// for i := 0; i < threads; i++ {
	// 	wg.Add(1)
	// 	go checkRepos(repos[beginning:end])
	// 	repos = repos[end:]
	// }

	// // wait for checkRepos to terminate
	// wg.Wait()
	// fmt.Printf("there are %d repos remaining\n", len(repos))
}

package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"

	. "github.com/logrusorgru/aurora"
	"github.com/google/go-github/github"
	gitleaks "github.com/zricethezav/gitleaks/src"
)

var wg sync.WaitGroup

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
			fmt.Printf("%s %s %d leaks found for '%s'\n", Green(" ✓"), Bold("[leaks]"), leakCount, name)
		} else {
			fmt.Printf("%s %s %d leaks found for '%s'\n", Red("!!"), Bold("[leaks]"), leakCount, name)
		}
	}
	wg.Done()
}

func main() {
	org := "github"
	fmt.Printf("%s Initialising scan on %s\n", Bold("   [main]"), org)

	fmt.Printf("%s There are %d threads available\n", Bold("   [main]"), runtime.NumCPU())
	client := github.NewClient(nil)
	// list public repositories for org "github"
	opt := &github.RepositoryListByOrgOptions{Type: "public"}
	repos, _, err := client.Repositories.ListByOrg(context.Background(), org, opt)

	if err != nil {
		log.Fatal(err)
	}


	// get all node repos
	// go doesn't have map :(
	var jsRepos []*github.Repository
	for _, repo := range repos {
		if repo.Language == nil {
			continue
		}
		if *repo.Language == "JavaScript" || *repo.Language == "TypeScript" {
			jsRepos = append(jsRepos, repo)
		}
	}

	// split repos in between threads
	threads := runtime.NumCPU()
	groupCount := len(jsRepos) / threads
	fmt.Printf("%s Found %d JavaScript/TypeScript repos\n", Bold("   [main]"), len(jsRepos))
	fmt.Printf("%s Giving %d repo(s) to each thread\n", Bold("   [main]"), groupCount)

	// // calculate a deficit (rounding errors when splitting work between threads)
	deficit := len(jsRepos) - (groupCount * threads)
	if deficit > 0 {
		fmt.Printf("%s There is a deficit of %d groups. Adding to thread...\n", Yellow("   [main (warn)]") ,deficit)
		wg.Add(1)
		go checkRepos(jsRepos[deficit:])
		jsRepos = jsRepos[deficit:]
	}

	beginning := 0
	end := groupCount
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go checkRepos(jsRepos[beginning:end])
		jsRepos = jsRepos[end:]
	}

	// wait for checkRepos to terminate
	wg.Wait()
}

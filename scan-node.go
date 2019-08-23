package main

import (
	"os"
	"fmt"
	"time"
	"path"
	"strings"
	"os/exec"
	"math/rand"
	"encoding/json"

	"gopkg.in/src-d/go-git.v4"
	. "github.com/logrusorgru/aurora"
)

type SecIssue struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Tag         string `json:"tag"`
	Line        int 	`json:"line"`
	Lines       string `json:"lines"`
	Filename    string `json:"filename"`
	Path        string `json:"path"`
	Sha2        string `json:"sha2"`
}

type HeaderIssue struct {
	Description	string `json:"description"`
	Tag			string `json:"tag"`
	Title		string `json:"title"`
}

type RepoIssue struct {
	Files 			[]string		`json:"files"`
	SecIssues 		[]SecIssue		`json:"security_issues"`
	HeaderIssues	[]HeaderIssue	`json:"header_issues"`
}

// happens to be the same for both yarn and npm
type AuditMeta struct {
	Info		string	`json:info`
	Low			string	`json:low`
	Moderate	string	`json:moderate`
	High		string	`json:high`
	Critical	string	`json:critical`
}

func Exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}

func runAudit(pkg string, dir string) AuditMeta {
	var cmd *exec.Cmd
	var auditResult AuditMeta
	if (pkg == "yarn") {
		cmd = exec.Command("yarn", "audit", "--json", "|", "jq", "'select(.data.vulnerabilities != null) | .data.vulnerabilities'", "-r", "-M")
	} else { // if pkg == npm
		cmd = exec.Command("npm", "audit", "--json", "|", "jq", "'select(.metadata.vulnerabilities != null) | .metadata.vulnerabilities'", "-r", "-M")
	}
	cmd.Dir = dir
	pkgIssues, err := cmd.Output()
	handleErr(err)
	fmt.Println(string(pkgIssues))
	json.Unmarshal(pkgIssues, &auditResult)
	return auditResult
}

func auditPkgs(directory string) {
	var auditResult AuditMeta
	// none? npm install --package-lock-only
	// yarn: yarn audit --json | jq 'select(.data.vulnerabilities != null)|.data.vulnerabilities' -r -M
	// npm: npm audit --json | jq 'select(.metadata.vulnerabilities != null)|.metadata.vulnerabilities' -r -M
	yarnDir := path.Join(directory, "yarn.lock")
	npmDir := path.Join(directory, "package-lock.json")

	// var pkgAudit AuditMeta

	// check yarn.lock
	if Exists(yarnDir) {
		auditResult = runAudit("yarn", directory)
	// check package-lock
	} else if Exists(npmDir) {
		auditResult = runAudit("npm", directory)
	} else {
		// create a lockfile
		fmt.Printf("%s creating lockfile\n", Bold("   [package audit]"))
		cmd := exec.Command("yarn")
		cmd.Dir = directory
		_, err := cmd.Output()
		handleErr(err)
		auditResult = runAudit("yarn", directory)
		// create a lockfile
	}

	fmt.Println(auditResult)
}

func NodeScan(repo string) {
	rand.Seed(time.Now().UnixNano())

	directory := path.Join("/tmp/", randStringRunes(25))
	outfile := path.Join(directory, (randStringRunes(15) + ".json"))

	fmt.Printf("%s cloning %s for SSA into %s\n", Bold("   [source analysis]"), repo, directory)

	_, err := git.PlainClone(directory, false, &git.CloneOptions{URL: repo})

	fmt.Printf("%s cloning %s complete. Scanning...\n", Bold("   [source analysis]"), repo)

	handleErr(err)

	cmd := exec.Command("nodejsscan", "-d", directory, "-o", outfile)

	_, err = cmd.Output()

	handleErr(err)

	parser := exec.Command("node", "./nodejsscan-parser.js", outfile)
	parserOut, err := parser.Output()


	var result RepoIssue
	json.Unmarshal(parserOut, &result)

	issues := len(result.SecIssues)

	if (issues > 0) {	
		fmt.Printf("%s %s Found %d issues in \"%s\"\n", Red("!!"), Bold("[source analysis]"), len(result.SecIssues), repo)
		for _, issue := range result.SecIssues {
			filename := strings.Replace(issue.Path, directory, "", 1)
			fmt.Printf("%s %s Found %s in file %s\n", Red("!!"), Bold("[source analysis]"), issue.Title, filename)
		}
	} else {
		fmt.Printf("%s %s Found %d issues in \"%s\"\n", Green(" ✓"), Bold("[source analysis]"), len(result.SecIssues), repo)
	}

	// handle audits
	// check for yarn.lock or package-lock.json

	auditPkgs(directory)

	// attempt to run npm audit (reason 2 for docker container)

	// attempt to remove the dir
	fmt.Printf("%s Removing %s\n", Bold("   [source analysis]"), directory)
	os.RemoveAll(directory)

}

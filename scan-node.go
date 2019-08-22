package main

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

type ScanIssue struct {
	title       string
	description string
	tag         string
	line        string
	lines       string
	filename    string
	path        string
	sha2        string
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func NodeScan(repo string) {
	rand.Seed(time.Now().UnixNano())

	directory := "/tmp/" + randStringRunes(10)
	fmt.Printf("Cloning %s for SSA into %s", repo, directory)

	_, err := git.PlainClone(directory, false, &git.CloneOptions{URL: repo})

	fmt.Printf("Clone complete. Scanning...")

	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("nodejsscan", "-d", directory)

	out, err := cmd.Output()

	fmt.Println(string(out))

}

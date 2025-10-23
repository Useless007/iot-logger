package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"time"
)

func run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// ตรวจว่ามีการเปลี่ยนแปลงใน repo ไหม
func changed() bool {
	out, _ := run("status", "--porcelain")
	return len(out) > 0
}

func main() {
	commitInterval := 6 * time.Hour     // commit ทุก 6 ชั่วโมง
	checkEvery := 10 * time.Minute      // ตรวจทุก 10 นาที
	lastCommit := time.Now().Add(-commitInterval)

	log.Printf("Auto committer started: check every %v, commit interval %v", checkEvery, commitInterval)

	for {
		time.Sleep(checkEvery)

		if changed() && time.Since(lastCommit) >= commitInterval {
			now := time.Now()
			msg := fmt.Sprintf("data: update logs at %s", now.Format("2006-01-02 15:04:05"))

			if _, err := run("add", "-A"); err != nil {
				log.Println("git add failed:", err)
				continue
			}
			if _, err := run("commit", "-m", msg); err != nil {
				log.Println("git commit failed:", err)
				continue
			}
			if out, err := run("push"); err != nil {
				log.Println("git push failed:", err, out)
				continue
			}

			lastCommit = now
			log.Printf("Committed and pushed: %s", msg)
		}
	}
}

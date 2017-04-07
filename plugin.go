package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	Repo struct {
		Owner   string
		Name    string
		Link    string
		Avatar  string
		Branch  string
		Private bool
		Trusted bool
	}

	Build struct {
		Number   int
		Event    string
		Status   string
		Deploy   string
		Created  int64
		Started  int64
		Finished int64
		Link     string
	}

	Commit struct {
		Remote  string
		Sha     string
		Ref     string
		Link    string
		Branch  string
		Message string
		Author  Author
	}

	Author struct {
		Name   string
		Email  string
		Avatar string
	}

	Config struct {
		// plugin-specific parameters and secrets
		Registry string
		Token    string
		Repo     string
		Tag      StrSlice
		File     string
		Context  string
		Storage  string
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Commit Commit
		Config Config
	}
)

func (p Plugin) Exec() error {
	go func() {
		args := []string{"daemon"}

		if len(p.Config.Storage) != 0 {
			args = append(args, "-s", p.Config.Storage)
		}

		cmd := exec.Command("/usr/bin/docker", args...)
		if os.Getenv("DOCKER_LAUNCH_DEBUG") == "true" {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		} else {
			cmd.Stdout = ioutil.Discard
			cmd.Stderr = ioutil.Discard
		}
		trace(cmd)
		cmd.Run()
	}()

	// ping Docker until available
	for i := 0; i < 3; i++ {
		cmd := exec.Command("/usr/bin/docker", "info")
		cmd.Stdout = ioutil.Discard
		cmd.Stderr = ioutil.Discard
		err := cmd.Run()
		if err == nil {
			break
		}
		time.Sleep(time.Second * 5)
	}

	// Login to Docker
	cmd := exec.Command("/usr/bin/docker", "login", "-u", "_json_key", "-p", p.Config.Token, "-e", "chunkylover53@aol.com", p.Config.Registry)
	// cmd.Dir = workspace.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("Login failed.")
		os.Exit(1)
	}

	// Build the container
	cmd = exec.Command("/usr/bin/docker", "build", "--pull=true", "--rm=true", "-f", p.Config.File, "-t", p.Commit.Sha, p.Config.Context)
	// cmd.Dir = workspace.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Unable to build docker image: %v", err)
	}

	// Creates image tags
	for _, tag := range p.Config.Tag.Slice() {
		// create the full tag name
		tag_ := fmt.Sprintf("%s:%s", p.Config.Repo, tag)
		if tag == "latest" {
			tag_ = p.Config.Repo
		}

		// tag the build image sha
		cmd = exec.Command("/usr/bin/docker", "tag", p.Commit.Sha, tag_)
		// cmd.Dir = workspace.Path
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		trace(cmd)
		err = cmd.Run()
		if err != nil {
			return fmt.Errorf("Unable to tag docker image: %v", err)
		}
	}

	// Push the image and tags to the registry
	cmd = exec.Command("/usr/bin/docker", "push", p.Config.Repo)
	// cmd.Dir = workspace.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Unable to push docker image: %v", err)
	}

	// plugin logic goes here
	return nil
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}

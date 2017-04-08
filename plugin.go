package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type (
	Commit struct {
		Sha string
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
		Commit Commit
		Config Config
	}
)

//Keep docker paths as constant to highlight shift from dcoker -d to dockerd for docker daemon invokation
const (
	//Docker daemon and client path are the same for this version of docker.
	dockerClientFullPath = "/usr/bin/docker"
	dockerDaemonFullPath = "/usr/bin/docker"
)

func (p Plugin) Exec() error {
	go func() {
		//NOTE: No arguments needed if using dockerd. Could not get docker login to work with newer docker release (non-interactive login error)
		//args := []string{}
		args := []string{"daemon"}

		if len(p.Config.Storage) != 0 {
			args = append(args, "-s", p.Config.Storage)
		}

		cmd := exec.Command(dockerDaemonFullPath, args...)
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

	var dockerInfoSucceeded = false

	// ping Docker until available
	for i := 0; i < 3; i++ {
		cmd := exec.Command(dockerClientFullPath, "info")
		cmd.Stdout = ioutil.Discard
		cmd.Stderr = ioutil.Discard
		err := cmd.Run()
		if err == nil {
			dockerInfoSucceeded = true
			break
		}
		time.Sleep(time.Second * 5)
	}

	//Added test to make sure that docker info succeeded before proceeding.
	if dockerInfoSucceeded {
		log.Println("Docker Daemon is running.")
	} else {
		return errors.New("Unable to detect that docker daemon is running. Docker info loop timed out.")
	}

	// Login to Docker
	cmd := exec.Command(dockerClientFullPath, "login", "-u", "_json_key", "-p", p.Config.Token, "-e", "chunkylover53@aol.com", p.Config.Registry)
	// cmd.Dir = workspace.Path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Login failed: %v\n", err)
	}

	// Build the container
	cmd = exec.Command(dockerClientFullPath, "build", "--pull=true", "--rm=true", "-f", p.Config.File, "-t", p.Commit.Sha, p.Config.Context)
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
		cmd = exec.Command(dockerClientFullPath, "tag", p.Commit.Sha, tag_)
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
	cmd = exec.Command(dockerClientFullPath, "push", p.Config.Repo)
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

package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"os/exec"
// 	"strings"
// 	"time"
// )

// type Docker struct {
// 	Registry string   `json:"registry"`
// 	Storage  string   `json:"storage_driver"`
// 	Token    string   `json:"token"`
// 	Repo     string   `json:"repo"`
// 	Tag      StrSlice `json:"tag"`
// 	File     string   `json:"file"`
// 	Context  string   `json:"context"`
// }

// var (
// 	buildCommit string
// )

// func main() {
// 	fmt.Printf("Drone GCR Plugin built from %s\n", buildCommit)

// 	// workspace := plugin.Workspace{}
// 	// build := plugin.Build{}
// 	// vargs := Docker{}

// 	// plugin.Param("workspace", &workspace)
// 	// plugin.Param("build", &build)
// 	// plugin.Param("vargs", &vargs)
// 	// plugin.MustParse()

// 	vargs := Docker{}

// 	// Repository name should have gcr prefix
// 	vargs.Registry = os.Getenv("PLUGIN_REGISTRY")
// 	if len(vargs.Registry) == 0 {
// 		vargs.Registry = "gcr.io"
// 	}
// 	// Set the Dockerfile name
// 	vargs.File = os.Getenv("PLUGIN_FILE")
// 	if len(vargs.File) == 0 {
// 		vargs.File = "Dockerfile"
// 	}
// 	// Set the Context value
// 	vargs.Context = os.Getenv("PLUGIN_CONTEXT")
// 	if len(vargs.Context) == 0 {
// 		vargs.Context = "."
// 	}
// 	// Set the Tag value
// 	vargs.Tag = os.Getenv("PLUGIN_TAG")
// 	if vargs.Tag.Len() == 0 {
// 		vargs.Tag = StrSlice{[]string{"latest"}}
// 	}
// 	// Concat the Registry URL and the Repository name if necessary
// 	vargs.Repo = os.Getenv("PLUGIN_REPO")
// 	if strings.Count(vargs.Repo, "/") == 1 {
// 		vargs.Repo = fmt.Sprintf("%s/%s", vargs.Registry, vargs.Repo)
// 	}
// 	// Trim any spaces or newlines from the token
// 	vargs.Token = os.Getenv("PLUGIN_TOKEN")
// 	vargs.Token = strings.TrimSpace(vargs.Token)

// 	vargs.Storage = os.Getenv("PLUGIN_STORAGE")

// 	go func() {
// 		args := []string{"daemon"}

// 		if len(vargs.Storage) != 0 {
// 			args = append(args, "-s", vargs.Storage)
// 		}

// 		cmd := exec.Command("/usr/bin/docker", args...)
// 		if os.Getenv("DOCKER_LAUNCH_DEBUG") == "true" {
// 			cmd.Stdout = os.Stdout
// 			cmd.Stderr = os.Stderr
// 		} else {
// 			cmd.Stdout = ioutil.Discard
// 			cmd.Stderr = ioutil.Discard
// 		}
// 		trace(cmd)
// 		cmd.Run()
// 	}()

// 	// ping Docker until available
// 	for i := 0; i < 3; i++ {
// 		cmd := exec.Command("/usr/bin/docker", "info")
// 		cmd.Stdout = ioutil.Discard
// 		cmd.Stderr = ioutil.Discard
// 		err := cmd.Run()
// 		if err == nil {
// 			break
// 		}
// 		time.Sleep(time.Second * 5)
// 	}

// 	// Login to Docker
// 	cmd := exec.Command("/usr/bin/docker", "login", "-u", "_json_key", "-p", vargs.Token, "-e", "chunkylover53@aol.com", vargs.Registry)
// 	cmd.Dir = workspace.Path
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println("Login failed.")
// 		os.Exit(1)
// 	}

// 	// Build the container
// 	cmd = exec.Command("/usr/bin/docker", "build", "--pull=true", "--rm=true", "-f", vargs.File, "-t", build.Commit, vargs.Context)
// 	cmd.Dir = workspace.Path
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	trace(cmd)
// 	err = cmd.Run()
// 	if err != nil {
// 		os.Exit(1)
// 	}

// 	// Creates image tags
// 	for _, tag := range vargs.Tag.Slice() {
// 		// create the full tag name
// 		tag_ := fmt.Sprintf("%s:%s", vargs.Repo, tag)
// 		if tag == "latest" {
// 			tag_ = vargs.Repo
// 		}

// 		// tag the build image sha
// 		cmd = exec.Command("/usr/bin/docker", "tag", build.Commit, tag_)
// 		cmd.Dir = workspace.Path
// 		cmd.Stdout = os.Stdout
// 		cmd.Stderr = os.Stderr
// 		trace(cmd)
// 		err = cmd.Run()
// 		if err != nil {
// 			os.Exit(1)
// 		}
// 	}

// 	// Push the image and tags to the registry
// 	cmd = exec.Command("/usr/bin/docker", "push", vargs.Repo)
// 	cmd.Dir = workspace.Path
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	trace(cmd)
// 	err = cmd.Run()
// 	if err != nil {
// 		os.Exit(1)
// 	}
// }

// // Trace writes each command to standard error (preceded by a ‘$ ’) before it
// // is executed. Used for debugging your build.
// func trace(cmd *exec.Cmd) {
// 	fmt.Println("$", strings.Join(cmd.Args, " "))
// }

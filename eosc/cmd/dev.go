// Copyright © 2018 EOS Canada <info@eoscanada.com>

package cmd

import (
	"context"
	"fmt"
	"os"

	dexec "github.com/ahmetalpbalkan/go-dexec"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "EOS Smart Contract Development tools",
}

func init() {
	RootCmd.AddCommand(devCmd)
}

var devDockerRepo = "gcr.io/eoscanada-public/"
var devDockerImageName = "eosio-cdt-slim"
var devDockerImageTag = "v1.2.1"
var devDockerImageID = "cc0eb6ea938115cf587842f1b41c0f2f84d1baf9ae815da3e68fce9755d6fa5f"
var devDockerImageFullName = devDockerRepo + devDockerImageName
var devDockerImageTaggedName = devDockerImageFullName + ":" + devDockerImageTag

func sanityCheckDevEnvironment() {
	fmt.Println("-- Performing sanity check of your Docker environment")
	cl, err := docker.NewClientFromEnv()

	errorCheck("Unable to communicate with Docker.", err)

	_, err = cl.InspectImage(devDockerImageID)
	if err != nil {
		ctx, cancel := context.WithCancel(context.Background())
		go setupSignalHandler(func(err error) {
			fmt.Println("Stopping early during while pulling Docker image for EOS Smart Contract Development")
			cancel()
		})

		pullImagesOptions := docker.PullImageOptions{
			Repository:        devDockerImageFullName,
			Tag:               devDockerImageTag,
			InactivityTimeout: 0,
			Context:           ctx,
			OutputStream:      os.Stdout,
		}
		noAuth := docker.AuthConfiguration{}
		fmt.Println("-- We need to pull the Docker image from a remote registry - this may take a while...")
		err = cl.PullImage(pullImagesOptions, noAuth)
		errorCheck("Unable to pull the Docker image for EOS Smart Contract Development", err)
	}

	var pwd string
	pwd, err = os.Getwd()
	errorCheck("Unable to get current working directory", err)

	if _, err := os.Stat(pwd + "/build.sh"); !os.IsNotExist(err) {
		d := dexec.Docker{cl}
		cco := docker.CreateContainerOptions{
			Config: &docker.Config{
				Image: devDockerImageTaggedName,
				Volumes: map[string]struct{}{
					"/workspace": {},
				},
			},
			HostConfig: &docker.HostConfig{
				Binds: []string{pwd + ":/workspace"},
			},
		}
		m, _ := dexec.ByCreatingContainer(cco)
		cmd := d.Command(m, "./build.sh")
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		errorCheck("Unable to exec the EOS Smart Contract Development Docker image", err)
	} else {
		errorCheck("Make sure the current working directory contains an EOS Smart Contract project", err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"time"

	vmfork "github.com/lox/govmomi-vmfork"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

const (
	envURL      = "VSPHERE_HOST"
	envUserName = "VSPHERE_USERNAME"
	envPassword = "VSPHERE_PASSWORD"
)

var (
	urlDescription      = fmt.Sprintf("ESX or vCenter URL [%s]", envURL)
	userNameDescription = fmt.Sprintf("ESX or vCenter Username [%s]", envUserName)

	urlFlag      = flag.String("url", getEnvString(envURL, "https://username:password@host/sdk"), urlDescription)
	userNameFlag = flag.String("username", getEnvString(envUserName, ""), userNameDescription)
	passwordFlag = flag.String("password", "", fmt.Sprintf("ESX or vCenter Username [%s]", envPassword))
	parentName   = flag.String("parent-vm", "", "The vm to fork")
	childName    = flag.String("child-vm", "", "The name of the child vm that is created")
	script       = flag.String("script", "", "The script to execute")
)

// getEnvString returns string from environment variable.
func getEnvString(v string, def string) string {
	r := os.Getenv(v)
	if r == "" {
		return def
	}

	return r
}

// NewClient creates a govmomi.Client
func newClient(ctx context.Context) (*govmomi.Client, error) {
	flag.Parse()

	// Parse URL from string
	u, err := soap.ParseURL(*urlFlag)
	if err != nil {
		return nil, err
	}

	password := getEnvString(envPassword, "")
	if *passwordFlag != "" {
		password = *passwordFlag
	}

	u.User = url.UserPassword(*userNameFlag, password)
	log.Printf("Connecting to %s", u.String())

	// Connect and log in to ESX or vCenter
	return govmomi.NewClient(ctx, u, true)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[random.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := newClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	vm, err := vmfork.FindVirtualMachine(ctx, *parentName, client.Client)
	if err != nil {
		log.Fatal(err)
	}

	// generate a name if one doesn't exist
	if *childName == "" {
		generatedName := fmt.Sprintf("%s-child-%s", *parentName, randSeq(10))
		childName = &generatedName
	}

	err = vm.Fork(ctx, vmfork.CreateChildSpec{
		Name:   *childName,
		Script: *script,
	})
	if err != nil {
		log.Fatal(err)
	}
}

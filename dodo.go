package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"fmt"
	"github.com/colebrumley/dodo/actions"
	"github.com/digitalocean/godo"
	"os"
)

var (
	cmdToken      string
	cmdArgs       []string
	cmdAction     string
	cmdActionArgs []string
	doClient      *godo.Client
)

func init() {
	flag.StringVar(&cmdToken, "token", "", "DO access token (or set $DO_TOKEN environment variable)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "DOdo - v.1 The Unofficial DigitalOcean CLI tool\nBasic usage: dodo [global args] action [action args]\n")
		fmt.Fprintf(os.Stderr, "  Examples:\n\tdodo list drops\n\tdodo list ips\n\tdodo create droplet name=awesomedroplet.com size=512mb image=fedora-21-x64\nGlobal Args:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Actions:\n  list\n  create\n  delete\n")
		fmt.Fprintf(os.Stderr, "For help on an action, add 'help' after it (i.e. dodo list help)\n")
	}
	flag.Parse()
	if cmdToken == "" {
		cmdToken = os.Getenv("DO_TOKEN")
	}
	cmdArgs = flag.Args()
	cmdAction, cmdActionArgs = shift(cmdArgs)
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: cmdToken},
	}
	doClient = godo.NewClient(t.Client())
}

func main() {
	switch cmdAction {
	case "list":
		if len(cmdActionArgs) <= 0 {
			actions.List("help", doClient)
			os.Exit(1)
		}
		actions.List(cmdActionArgs[0], doClient)
	case "create":
		if len(cmdActionArgs) <= 0 {
			actions.Create("help", cmdActionArgs, doClient)
			os.Exit(1)
		}
		actions.Create(cmdActionArgs[0], cmdActionArgs[1:], doClient)
	case "delete":
		if len(cmdActionArgs) <= 0 {
			actions.Delete("help", cmdActionArgs, doClient)
			os.Exit(1)
		}
		actions.Delete(cmdActionArgs[0], cmdActionArgs[1:], doClient)
	}
}

func shift(slice []string) (string, []string) {
	switch len(slice) {
	case 0:
		return "", nil
	case 1:
		return slice[0], nil
	}
	return slice[0], slice[1:]
}


package main

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"fmt"
	"github.com/digitalocean/godo"
	"os"
	"github.com/colebrumley/dodo/actions"
)

var (
	cmdToken      string
	cmdArgs       []string
	cmdAction     string
	cmdActionArgs []string
	doClient      *godo.Client
)

func init() {
	flag.StringVar(&cmdToken, "token", "db422c4066e6cccf038e2f139f3ed4101268d1ab8a120bf607e011951ea1edcc", "DO access token")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "DOdo - v.1 The Unofficial DigitalOcean CLI tool\nBasic usage: dodo [global args] action [action args]\n")
		fmt.Fprintf(os.Stderr, "  Examples:\n\tdodo list drops\n\tdodo list ips\n\tdodo create droplet name=awesomedroplet.com memory=512mb\nGlobal Args:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Actions:\n  list\n  create\n")
		fmt.Fprintf(os.Stderr, "For help on an action, add 'help' after it (i.e. dodo list help)\n")
	}
	flag.Parse()
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
			fmt.Printf("lolwut: list what?\n")
			os.Exit(1)
		}
		actions.List(cmdActionArgs[0], doClient)
	case "create":
		if len(cmdActionArgs) <= 0 {
			fmt.Printf("lolwut: create what?\n")
			os.Exit(1)
		}
		actions.Create(cmdActionArgs[0], cmdActionArgs[1:], doClient)
	case "delete":
		if len(cmdActionArgs) <= 0 {
			fmt.Printf("lolwut: delete what?\n")
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


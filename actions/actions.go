package actions

import (
	"fmt"
	"os"
	"strings"
	"log"
	"text/tabwriter"
	"github.com/digitalocean/godo"
)

func Delete(deleteItem string, args []string, client *godo.Client) {
	switch deleteItem {
	case "help":
		fmt.Printf("USAGE: 'dodo delete droplet [name]'\n")
		os.Exit(2)
	case "droplet":
			var myId int
			myName := args[0]
			drops, _ := DropletList(client)
			for _, d := range drops {
				if d.Name == myName {
					myId = d.ID
				}
			}
			if myId > 0 {
				_, err := client.Droplets.Delete(myId)
				if err != nil {
					fmt.Printf("Could not delete droplet: %s\n", err.Error())
				}
				fmt.Printf("Successfully deleted droplet %s\n", myName)
			}
	}
}

func Create(createItem string, args []string, client *godo.Client) {
	switch createItem {
	case "help":
		fmt.Printf("USAGE: 'dodo create droplet [var=value]'\nDroplet vars:\n  name\n  size\n  image\n  region\n  userdata\n  keys\n  backups\n  ipv6\n  privatenetworking\n")
		os.Exit(2)
	case "droplet":
		createRequest := &godo.DropletCreateRequest{
			Region: "nyc3",
		}
		if len(args) > 1 {
			for _, arg := range args[0:] {
				split := strings.Split(arg, "=")
				//fmt.Printf("%s=%s\n", split[0], split[1])
				switch split[0] {
				case "name", "Name":
					createRequest.Name = split[1]
				case "size", "Size":
					createRequest.Size = split[1]
				case "image", "Image":
					createRequest.Image = split[1]
				case "region", "Region":
					createRequest.Region = split[1]
				case "userdata", "UserData":
					createRequest.UserData = split[1]
				case "keys", "sshkeys":
					allKeys, _ := SshKeyList(client)
					for _, key := range strings.Split(split[1], ",") {
						for _, k := range allKeys {
							if key == k.Name || key == fmt.Sprintf("%v", k.ID) {
								createRequest.SSHKeys = append(createRequest.SSHKeys, k.ID)								
							}
						}
					}
				case "backups", "Backups":
					if split[1] == "true" {
						createRequest.Backups = true
					}
				case "ipv6", "IPv6":
					if split[1] == "true" {
						createRequest.IPv6 = true
					}
				case "privatenetworking", "PrivateNetworking":
					if split[1] == "true" {
						createRequest.PrivateNetworking = true
					}
				}
			}
			_, _, err := client.Droplets.Create(createRequest)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}
			fmt.Printf("Successfully created droplet %s\n", createRequest.Name)
			return
		}
	case "image":
		fmt.Printf("Not quite there yet\n")
	}
	fmt.Printf("Not enough info!\n")
}

func List(listItem string, client *godo.Client) {
	switch listItem {
	case "help", "-h", "--help":
		fmt.Printf("USAGE: 'dodo list ITEM' where ITEM is one of:\n  droplets or drops\n  ips\n  keys or sshkeys\n  images or distros\n")
		os.Exit(2)
	case "droplets", "drops", "dr":
		drops, err := DropletList(client)
		if err != nil {
			log.Fatal("Couldn't list droplets")
		}
		PrettyPrintDroplets(drops)
	case "ips", "ip":
		drops, err := DropletList(client)
		if err != nil {
			log.Fatal("Couldn't list droplets")
		}
		PrettyPrintIPs(drops)
	case "keys", "sshkeys":
		keys, err := SshKeyList(client)
		if err != nil {
			log.Fatal("Couldn't list SSH keys")
		}
		PrettyPrintKeys(keys)
	case "images", "distros", "img", "dist":
		distros, err := DistroList(client)
		if err != nil {
			log.Fatal("Couldn't list distros")
		}
		PrettyPrintDistros(distros)
	}
}

func DropletList(client *godo.Client) ([]godo.Droplet, error) {
	list := []godo.Droplet{}
	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := client.Droplets.List(opt)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, d := range droplets {
			list = append(list, d)
		}
		if resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			log.Fatal(err.Error())
		}
		opt.Page = page + 1
	}
	return list, nil
}

func DistroList(client *godo.Client) ([]godo.Image, error) {
	list := []godo.Image{}
	opt := &godo.ListOptions{}
	for {
		images, resp, err := client.Images.List(opt)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, i := range images {
			list = append(list, i)
		}
		if resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			log.Fatal(err.Error())
		}
		opt.Page = page + 1
	}
	return list, nil
}

func SshKeyList(client *godo.Client) ([]godo.Key, error) {
	list := []godo.Key{}
	opt := &godo.ListOptions{}
	for {
		keys, resp, err := client.Keys.List(opt)
		if err != nil {
			log.Fatal(err.Error())
		}
		for _, k := range keys {
			list = append(list, k)
		}
		if resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			log.Fatal(err.Error())
		}
		opt.Page = page + 1
	}
	return list, nil
}

func PrettyPrintDroplets(drops []godo.Droplet) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 5, 1, ' ', 0)
	//w.Init(output, minwidth, tabwidth, padding, padchar, flags)
	fmt.Fprintln(w, "NAME\t", "CORES\t", "RAM\t", "DISK\t", "STATUS\t")
	for _, drop := range drops {
		fmt.Fprintln(w, drop.Name, "\t", drop.Vcpus, "\t", drop.Memory, "MB\t", drop.Disk, "GB\t", drop.Status)
	}
	w.Flush()
}

func PrettyPrintDistros(distros []godo.Image) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 5, 1, ' ', 0)
	//w.Init(output, minwidth, tabwidth, padding, padchar, flags)
	fmt.Fprintln(w, "NAME\t", "SLUG\t", "ID")
	for _, distro := range distros {
		fmt.Fprintln(w, distro.Name, "\t", distro.Slug, "\t", distro.ID)
	}
	w.Flush()
}

func PrettyPrintIPs(drops []godo.Droplet) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 5, 1, ' ', 0)
	//w.Init(output, minwidth, tabwidth, padding, padchar, flags)
	fmt.Fprintln(w, "DROPLET\t", "PROTO\t", "TYPE\t", "ADDRESS")
	for _, drop := range drops {
		for _, net := range drop.Networks.V4 {
			fmt.Fprintln(w, drop.Name, "\t", "IPv4\t", net.Type, "\t", net.IPAddress)
		}
		for _, net := range drop.Networks.V6 {
			fmt.Fprintln(w, drop.Name, "\t", "IPv6\t", net.Type, "\t", net.IPAddress)
		}
	}
	w.Flush()
}

func PrettyPrintKeys(keys []godo.Key) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 20, 5, 1, ' ', 0)
	//w.Init(output, minwidth, tabwidth, padding, padchar, flags)
	fmt.Fprintln(w, "NAME\t", "ID\t", "FINGERPRINT")
	for _, key := range keys {
		fmt.Fprintln(w, key.Name, "\t", key.ID, "\t", key.Fingerprint)
	}
	w.Flush()
}


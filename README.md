# DOdo
### Intro

DOdo is a lightweight CLI interface for the DigitalOcean API.  It adheres to many of the same usage standards as modern Golang apps like Etcd and Docker.  This 
early version is still being fleshed out, but right now it can:
 - list
   - droplets
   - ssh keys
   - IPs
   - distros
 - create droplets
 - delete droplets

Generally, if the syntax is wrong or the feature isn't implemented the program will exit with no output.  If an action was attempted, you should see something.

### Usage
```sh
DOdo - v.1 The Unofficial DigitalOcean CLI tool Basic usage: dodo [global args] action [action args]
  Examples:
        dodo list drops
        dodo list ips
        dodo create droplet name=awesomedroplet.com memory=512mb
 Global Args:
  -token="": DO access token (or set $DO_TOKEN environment variable)
 Actions:
  list
  create
  delete For help on an action, add 'help' after it (i.e. dodo list help)
  ```
  
### Examples

List droplets: `dodo list droplets` 

Create a basic CoreOS droplet: `dodo create droplet name=test image=coreos-alpha size=512mb keys=main`

Delete that droplet: `dodo delete droplet test`

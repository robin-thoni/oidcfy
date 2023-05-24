package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/robin-thoni/oidcfy/internal/config"
	"github.com/robin-thoni/oidcfy/internal/profiles"
	"github.com/robin-thoni/oidcfy/internal/server"
)

// type a struct {
// 	a string
// }

// type b struct {
// 	a string
// }

// type c struct {
// 	a
// 	b
// }

// func (d a) test() {

// }

// func (d b) test() {

// }

// func test(d a) {

// }

// func test1() {
// 	myC := c{}
// 	myC.test()
// }

func main() {

	rootConfig := config.RootConfig{}
	rootConfigStr, err := os.ReadFile("./configs/locals/example.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}
	err = yaml.Unmarshal([]byte(rootConfigStr), &rootConfig)
	if err != nil {
		log.Fatal(err)
		return
	}

	profiles := profiles.Profiles{}
	errs := profiles.FromConfig(&rootConfig)
	if len(errs) > 0 {
		log.Println(errs)
	}

	for name, profile := range profiles.MatchProfiles {
		if !profile.IsValid() {
			log.Printf("Match profile %s is invalid", name)
		}
	}
	for name, profile := range profiles.AuthenticationProfiles {
		if !profile.IsValid() {
			log.Printf("Authentication profile %s is invalid", name)
		}
	}
	for name, profile := range profiles.AuthorizationProfiles {
		if !profile.IsValid() {
			log.Printf("Authorization profile %s is invalid", name)
		}
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", rootConfig.Http.Address, rootConfig.Http.Port))
	if err != nil {
		log.Fatal(err)
		ln.Close()
	}

	server := server.NewServer(&rootConfig, &profiles)

	err = http.Serve(ln, server)
	if err != nil {
		log.Fatal(err)
		ln.Close()
	}
}

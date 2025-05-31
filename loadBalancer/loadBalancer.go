package loadbalancer

import (
	"encoding/json"
	"making-loadbalancer/server"
	"os"
)

type LoadBalancer struct {
	PORT    uint16
	Servers []server.Server
}

type configFile struct {
	Port    uint16   `json:"PORT"`
	Servers []string `json:"Servers"`
}

func Initialize(configFilePath string) (*LoadBalancer, error) {

	//read config.json file and store marshalled data in "data"
	data, fileReadError := os.ReadFile((configFilePath))
	if fileReadError != nil {
		return nil, fileReadError
	}

	//unmarshal json data from "data" into cf struct
	cf := new(configFile)
	parsingError := json.Unmarshal(data, cf)
	if parsingError != nil {
		return nil, parsingError
	}
	// fmt.Println("config file details")
	// fmt.Printf("PORT %d\n", cf.Port)
	// for index, server := range cf.Servers {
	// 	fmt.Printf("Server %d is %s\n", index, server)
	// }

	//intialize loadBalancer
	lb := new(LoadBalancer)
	lb.PORT = cf.Port
	for _, url := range cf.Servers {
		s := server.NewServer(url)
		lb.Servers = append(lb.Servers, *s)
	}

	return lb, nil
}

func (lb *LoadBalancer) ReqHandler() {

}

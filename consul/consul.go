package consul

import (
	"log"

	"github.com/hashicorp/consul/api"
) 

var (
	client *api.Client

	// The name of the service in Consul
	ServiceName string = "unkown-service"
	
	// The name of the proxy-path
	Proxy string = "unkown"
)

// Init Consul client
func init() {
	var err error
	client, err = api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Panicf("cannot init consul client: %v\n", err)
	}
}

func unregisterProxy() {
	log.Println("deregister service...")
	if err := client.Agent().ServiceDeregister(ServiceName); err != nil {
		log.Printf("cannot deregister proxy: %v\n", err)
	}
}

func registerProxy() error {
	log.Printf("Registering service '%s'\n", ServiceName)
	check := &api.AgentServiceCheck{
		HTTP:     "http://localhost:8080/health",
		Interval: "1s",
		Timeout:  "100ms",
	}
	return client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name:  ServiceName,
		Port:  8080,
		Check: check,
		Tags:  []string{"proxy://" + Proxy},
	})
}

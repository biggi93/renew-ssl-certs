package hetzner

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"github.com/biggi93/simple-file-server/config"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type hetznerData struct {
	HetznerClient         *hcloud.Client
	LoadBalancerId        int
	LoadBalancer          *hcloud.LoadBalancer
	TargetServiceIdForSSL int
	ListenPort            int
	TargetPort            int
	Domain                string
}

func GetHetznerObject() *hetznerData {
	config := config.Get()
	client := hcloud.NewClient(hcloud.WithToken(config.Hetzner.Token))

	h := &hetznerData{
		LoadBalancerId: config.Hetzner.LoadbalanceID,
		ListenPort:     config.Hetzner.LoadBalancerListenerPort,
		TargetPort:     config.Hetzner.LoadBalancerTargetPort,
		HetznerClient:  client,
		Domain:         config.Domain,
	}

	h.setLoadBalancer()
	return h
}

func (h *hetznerData) setLoadBalancer() {
	lb, _, err := h.HetznerClient.LoadBalancer.GetByID(context.Background(), h.LoadBalancerId)

	if err != nil {
		panic(err)
	}

	h.LoadBalancer = lb
}

func (h *hetznerData) AddTcpPortServiceForSSL() error {

	opts := hcloud.LoadBalancerAddServiceOpts{
		Protocol:        hcloud.LoadBalancerServiceProtocolTCP,
		ListenPort:      Ptr(h.ListenPort),
		DestinationPort: Ptr(h.TargetPort),
		Proxyprotocol:   Ptr(false),
		HTTP:            nil,
		HealthCheck: &hcloud.LoadBalancerAddServiceOptsHealthCheck{
			Protocol: hcloud.LoadBalancerServiceProtocolTCP,
			Port:     Ptr(h.ListenPort),
			Interval: Ptr(15 * time.Second),
			Timeout:  Ptr(10 * time.Second),
			Retries:  Ptr(3),
			HTTP:     nil,
		},
	}

	action, _, err := h.HetznerClient.LoadBalancer.AddService(context.Background(), h.LoadBalancer, opts)

	if err != nil {
		return err
	}

	if action.Status != "success" {
		return fmt.Errorf("New Service for http could not be created")
	}

	h.TargetServiceIdForSSL = action.ID

	log.Printf("New Loadbalance Service for SSL Certification were added with hetzner id %d", action.ID)

	return nil

}

func (h *hetznerData) WaitForTcp80LBService() error {
	for {
		log.Println("Starting Waiting Setting Up Loadbalancer Service")
		lb, _, err := h.HetznerClient.LoadBalancer.GetByID(context.Background(), h.LoadBalancer.ID)
		if err != nil {
			return err
		}

		for _, service := range lb.Targets[0].HealthStatus {
			if service.ListenPort != h.ListenPort {
				continue
			}
			if service.Status == hcloud.LoadBalancerTargetHealthStatusStatusHealthy {
				log.Printf("Load Balancer successfully listens on port %d\n", h.ListenPort)
				return nil
			} else {
				log.Printf("Waiting for Load Balancer. Status: %s\n", service.Status)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (h *hetznerData) TestLbService() {
	for {

		response, err := http.Get(fmt.Sprintf("http://%s:%d/test", h.Domain, h.ListenPort))
		if err != nil {
			continue
		}
		if response.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func (h *hetznerData) DeleteTcpPort80ServiceForSSL() error {
	if h.TargetServiceIdForSSL == 0 {
		return fmt.Errorf("No Loadbalancer Target Service available")
	}

	action, _, err := h.HetznerClient.LoadBalancer.DeleteService(context.Background(), h.LoadBalancer, h.ListenPort)

	if err != nil {
		return err
	}

	if action.Status != "success" {
		return fmt.Errorf("New Service for http could not be deleted")
	}

	log.Printf("Loadbalance Service for SSL Certification with hetzner id %d were deleted", h.TargetServiceIdForSSL)
	h.TargetServiceIdForSSL = 0

	return nil
}

func Ptr[T any](value T) *T {
	return &value
}

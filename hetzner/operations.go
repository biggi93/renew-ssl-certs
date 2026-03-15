package hetzner

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/biggi93/simple-file-server/config"

	"github.com/hetznercloud/hcloud-go/hcloud"
)

type hetznerData struct {
	HetznerClient         *hcloud.Client
	LoadBalancerId        int
	FireWallName			 string
	LoadBalancer          *hcloud.LoadBalancer
	FireWall					 *hcloud.Firewall
	TargetServiceIdForSSL int
	ListenPort            int
	TargetPort            int
	PortSSHIn				 string
	PortBWIn					 string
	PortBWOut				 string
	PortCertBotIn			 string
	PortCertBotOut			 string
	PortBWMailOut			 string
	Port443					 string
	Domain                string
}

func GetHetznerObject() *hetznerData {
	config := config.Get()
	client := hcloud.NewClient(hcloud.WithToken(config.Hetzner.Token))

	h := &hetznerData{
		LoadBalancerId: config.Hetzner.LoadbalanceID,
		ListenPort:     config.Hetzner.LoadBalancerListenerPort,
		TargetPort:     config.Hetzner.LoadBalancerTargetPort,
		PortSSHIn: 		 config.Hetzner.PortSshIn,
		PortBWIn: 		 config.Hetzner.PortBWIn,
		PortBWOut: 		 config.Hetzner.PortBWOut,
		PortBWMailOut:	 config.Hetzner.PortBwMailOut,
		PortCertBotIn:  config.Hetzner.PortCbIn,
		PortCertBotOut:  config.Hetzner.PortCbOut,
		Port443: config.Hetzner.Port443,
		FireWallName: 	 config.Hetzner.FireWallName,
		HetznerClient:  client,
		Domain:         config.Domain,
	}

	h.setLoadBalancer()
	h.setFireWall()
	return h
}

func (h *hetznerData) setFireWall() {

	fw,_,  err := h.HetznerClient.Firewall.GetByName(context.Background(), h.FireWallName)

	if err != nil {
		panic(err)
	}

	h.FireWall = fw
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
	time.Sleep(15 * time.Second)
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
		println("testing lb")
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



func(h *hetznerData) OpenAddFWRule() error {
	if h.FireWall == nil {
		return fmt.Errorf("no firewall object available")
	}

	rules := hcloud.FirewallSetRulesOpts{
		Rules: []hcloud.FirewallRule{
				{
					Direction: hcloud.FirewallRuleDirectionIn,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortSSHIn),
					SourceIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("ssh"),
				},
				{
					Direction: hcloud.FirewallRuleDirectionIn,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortCertBotIn),
					SourceIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("certbot"),
				},
				{
					Direction: hcloud.FirewallRuleDirectionIn,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortBWIn),
					SourceIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("bitwarden"),

				},
				{
					Direction: hcloud.FirewallRuleDirectionOut,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortCertBotOut),
					DestinationIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("certbot"),
				},
				{
					Direction: hcloud.FirewallRuleDirectionOut,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortBWOut),
					DestinationIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("bitwarden"),

				},
				{
					Direction: hcloud.FirewallRuleDirectionOut,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortBWMailOut),
					DestinationIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("bw_mail"),
				},
				{
					Direction: hcloud.FirewallRuleDirectionOut,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.Port443),
					DestinationIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("bw_update"),
				},
		},
	}

	_, _, err := h.HetznerClient.Firewall.SetRules(context.Background(), h.FireWall, rules)

	if err != nil {
		return  err
	}

	return  nil
}


func(h *hetznerData) CloseAddFWRule() error {
	if h.FireWall == nil {
		return fmt.Errorf("no firewall object available")
	}

	rules := hcloud.FirewallSetRulesOpts{
		Rules: []hcloud.FirewallRule{
				{
					Direction: hcloud.FirewallRuleDirectionIn,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortSSHIn),
					SourceIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("ssh"),
				},
				{
					Direction: hcloud.FirewallRuleDirectionIn,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortBWIn),
					SourceIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("bitwarden"),

				},
				{
					Direction: hcloud.FirewallRuleDirectionOut,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortBWOut),
					DestinationIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("bitwarden"),

				},
				{
					Direction: hcloud.FirewallRuleDirectionOut,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.PortBWMailOut),
					DestinationIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("bw_mail"),
				},
				{
					Direction: hcloud.FirewallRuleDirectionOut,
					Protocol: hcloud.FirewallRuleProtocolTCP,
					Port: Ptr(h.Port443),
					DestinationIPs: []net.IPNet{
						{
							IP:   net.ParseIP("0.0.0.0"),
							Mask: net.CIDRMask(0, 32),
						},
            	},
					Description: Ptr("bw_update"),
				},
		},
	}

	_, _, err := h.HetznerClient.Firewall.SetRules(context.Background(), h.FireWall, rules)

	if err != nil {
		return  err
	}

	return  nil
}
package keystore

import (
	"encoding/json"
	"log"

	"github.com/coreos/go-etcd/etcd"
)

// Structs used for decoding json received by watch command.

type Port struct {
	ContainerPort int    `json:"containerPort,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

type Condition struct {
	Status string `json:"status,omitempty"`
	Type   string `json:"type,omitempty"`
}

type Container struct {
	Image string `json:"image,omitempty"`
	Name  string `json:"name,omitempty"`
	Ports []Port `json:"ports,omitempty"`
}

type Pod struct {
	Spec struct {
		Containers []Container `json:"containers,omitempty"`
		Host       string      `json:"host,omitempty"`
	} `json:"spec,omitempty"`

	//
	Status_ struct {
		Conditions []Condition `json:"Condition,omitempty"`
		Phase      string      `json:"phase"`
		PodIP      string      `json:"podIP,omitempty"`
	} `json:"status"`
}

func (pod *Pod) Status() string {
	return pod.Status_.Phase
}

func (pod *Pod) PodIP() string {
	return pod.Status_.PodIP
}

func (pod *Pod) Host() string {
	return pod.Spec.Host
}

func (pod *Pod) String() string {
	if b, err := json.MarshalIndent(pod, "", "\t"); err != nil {
		return "{ERROR}"
	} else {
		return string(b)
	}
}

// etcd.Watch is blocking so use this helper function to monitor.
func watch(client *etcd.Client, receiver chan *etcd.Response, stop chan bool) {
	if _, err := client.Watch("/registry/pods", 0, true, receiver, stop); err != nil {
		log.Fatal(err)
	}
}

//  Connect to master etcd instance and wait for pods to come online/offline.
func Watch(master string) error {

	machines := []string{master}
	client := etcd.NewClient(machines)

	receiver := make(chan *etcd.Response, 1)
	stop := make(chan bool, 1)

	go watch(client, receiver, stop)

	log.Println("Watching and waiting for pods to come online...")
	for {
		select {
		case resp := <-receiver:
			if resp == nil {
				log.Printf("Got nil resp in watch channel.")
			} else {

				var pod Pod
				json.Unmarshal([]byte(resp.Node.Value), &pod)

				switch resp.Action {
				case "create":
					log.Printf("> Pod %s created.", resp.Node.Key)
				case "compareAndSwap":
					if pod.Status() == "Running" {
						log.Printf("> Pod %s status changed to %s with ip %s", resp.Node.Key, pod.Status(), pod.PodIP())
					} else {
						log.Printf("> Pod %s status changed to %s", resp.Node.Key, pod.Status())
					}
					// log.Printf("\n> %s\n\t%s\n\t%s", resp.Action, resp.Pod.Value, resp.Pod.Key)
				case "delete":
					log.Printf("> Pod %s offline.", resp.Node.Key)
				}

				log.Printf(pod.String())
			}
		case <-stop:
			log.Printf("Exiting!")
			return nil
		}
	}
}

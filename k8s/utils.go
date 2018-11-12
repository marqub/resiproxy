package k8s

import (
	"strconv"

	"net"

	"github.com/marqub/resiproxy/log"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// CreateK8sMapping opens a new port at the ToxiProxy k8s service level
func CreateK8sMapping(listen string) error {
	_, listeningPort, err := net.SplitHostPort(listen)
	if err != nil {
		return err
	}

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	// get the toxyproxy service
	service, err := clientset.CoreV1().Services(Config.Namespace).Get(Config.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// add a new port mapping to the service
	log.Logger().Debug("%+v", service.Spec.Ports)
	newPort, _ := strconv.Atoi(listeningPort)
	if !doesPortExist(service.Spec.Ports, listeningPort) {
		service.Spec.Ports = append(service.Spec.Ports, v1.ServicePort{Name: strconv.Itoa(newPort), Port: int32(newPort), TargetPort: intstr.FromInt(newPort)})
		log.Logger().Debug("%+v", service.Spec.Ports)
		service, err = clientset.CoreV1().Services(Config.Namespace).Update(service)
		if err != nil {
			return err
		}
	}
	return nil
}

func doesPortExist(existingPorts []v1.ServicePort, value string) bool {
	for _, port := range existingPorts {
		if port.Name == value {
			return true
		}
	}
	return false
}

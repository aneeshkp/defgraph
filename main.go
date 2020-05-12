package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	report "github.com/aneeshkp/depgraph/internal/report"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//apps "k8s.io/api/apps/v1"
	//v1 "k8s.io/api/apps/v1"
)

var (
	clientset *kubernetes.Clientset
	n         = report.NewNodeTable()
)

func main() {
	// creates the in-cluster config

	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig *string
		if home, _ := os.UserHomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
		fmt.Println(kubeconfig)

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	// creates the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	for _, pod := range pods.Items {
		ors := pod.GetOwnerReferences()
		if len(ors) > 0 {
			refPod := report.RefResult{Name: pod.Name, Namespace: pod.Namespace, Kind: "Pod", OwnerReference: ors[0].Name, Ownerkind: ors[0].Kind}
			genTree(&refPod)
			n.AddNode(pod.Name, pod.Namespace, "Pod", ors[0].Name)
		} else {
			n.AddNode(pod.Name, pod.Namespace, "Pod", "")
		}

	}
	n.ShowAll()

}

func genTree(node *report.RefResult) {
	fmt.Println("node.Ownerkind", node.Ownerkind)
	if node.OwnerReference == "" {
		n.AddNode(node.Name, node.Namespace, node.Kind, node.OwnerReference)
	}
	err, item := GetResourse(node.Ownerkind, node.OwnerReference, node.Namespace)
	if err != nil || item == nil {
		fmt.Printf("%#v\n", err)
	} else {
		if item.Kind == "Node" {
			n.AddNode(item.Name, item.Namespace, item.Kind, "")
		} else if item.OwnerReference != "" && item.Name == node.OwnerReference {
			//previous item references to me and I do  have parent ,  find my parent
			genTree(item)
			n.AddNode(item.Name, item.Namespace, item.Kind, item.OwnerReference)

		} else if item.OwnerReference == "" && item.Name == node.OwnerReference {
			//previous item references to me and I don;t have parent , print
			n.AddNode(item.Name, item.Namespace, item.Kind, item.OwnerReference)

		}
	}

}

func GetResourse(kind string, name string, namespace string) (err error, refSource *report.RefResult) {
	switch kind {
	case "DaemonSet":
		{
			ds, errs := clientset.AppsV1().DaemonSets(namespace).Get(name, metav1.GetOptions{})
			if errs != nil {
				err = fmt.Errorf("Error getting resource %#v", errs)
				return
			}
			ors := ds.GetOwnerReferences()
			if len(ors) > 0 {
				refSource = &report.RefResult{Name: ds.Name, Namespace: namespace, Kind: kind, OwnerReference: ors[0].Name, Ownerkind: ors[0].Kind}

			} else {
				//fmt.Println(pod.Name, pod.Namespace, pod.Kind, "")
				refSource = &report.RefResult{Name: ds.Name, Namespace:namespace,Kind:kind}
			}
		}
	case "Deployment":
		{
			ds, errs := clientset.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
			if errs != nil {
				err = fmt.Errorf("Error getting resource %#v", errs)
				return
			}
			ors := ds.GetOwnerReferences()
			if len(ors) > 0 {
				fmt.Println(ds.Name, ds.Namespace, ds.Kind, ors[0].Name)
				refSource = &report.RefResult{Name: ds.Name, Namespace: namespace, Kind: kind, OwnerReference: ors[0].Name, Ownerkind: ors[0].Kind}
			} else {
				//fmt.Println(pod.Name, pod.Namespace, pod.Kind, "")
				refSource = &report.RefResult{Name: ds.Name, Namespace: namespace, Kind: kind}
			}
		}
	case "ReplicaSet":
		{
			ds, errs := clientset.AppsV1().ReplicaSets(namespace).Get(name, metav1.GetOptions{})
			if errs != nil {
				err = fmt.Errorf("Error getting resource %#v", errs)
				return
			}
			ors := ds.GetOwnerReferences()
			if len(ors) > 0 {
				fmt.Println(ds.Name, ds.Namespace, ds.Kind, ors[0].Name)
				refSource = &report.RefResult{Name: ds.Name, Namespace: namespace, Kind: kind, OwnerReference: ors[0].Name, Ownerkind: ors[0].Kind}

			} else {
				//fmt.Println(pod.Name, pod.Namespace, pod.Kind, "")
				refSource = &report.RefResult{Name: ds.Name, Namespace: namespace, Kind: kind}
			}
		}
	default:
		{
			refSource = &report.RefResult{Name: name, Namespace: namespace, Kind: "Node"}
			fmt.Printf("Unknown Type Name  %s  Kind %s", name, kind)
		}
	}
	return
}

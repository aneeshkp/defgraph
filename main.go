package main

import (

	"flag"
	"fmt"


	"os"
	"path/filepath"
	"strings"

	report "github.com/aneeshkp/depgraph/internal/report"
 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/dynamic"
	//apps "k8s.io/api/apps/v1"
	//v1 "k8s.io/api/apps/v1"
	//"k8s.io/apimachinery/pkg/api/meta"
	//"k8s.io/apimachinery/pkg/runtime/schema"
	//"k8s.io/client-go/discovery"
	//"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	//"k8s.io/client-go/restmapper"
)

var (
	clientset *kubernetes.Clientset
	//dynamiClient *dynamic.Client
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

	//singlePod()


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
			n.AddNode(pod.Name, pod.Namespace, "Pod", ors[0].Name,nil)
		} else {
			n.AddNode(pod.Name, pod.Namespace, "Pod", "",nil)
		}

	}
	n.ShowAll()

}
func singlePod(){
	var podname ="network-operator-cf4d548b8-2d6tz"
	var namespace= "openshift-network-operator"
	pod, err := clientset.CoreV1().Pods(namespace).Get(podname,metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	ors := pod.GetOwnerReferences()
	if len(ors) > 0 {
		refPod := report.RefResult{Name: pod.Name, Namespace: pod.Namespace, Kind: "Pod", OwnerReference: ors[0].Name, Ownerkind: ors[0].Kind}
		genTree(&refPod)
		n.AddNode(pod.Name, pod.Namespace, "Pod", ors[0].Name,nil)
	} else {
		n.AddNode(pod.Name, pod.Namespace, "Pod", "",nil)
	}
	n.ShowAll()
	os.Exit(0)
}

func genTree(node *report.RefResult) {
	fmt.Println("node.Ownerkind", node.Ownerkind)
	if node.OwnerReference == "" {
		n.AddNode(node.Name, node.Namespace, node.Kind, node.OwnerReference,nil)
	}
	err, item := GetResourse(node.Ownerkind, node.OwnerReference, node.Namespace)
	if err != nil || item == nil {
		fmt.Printf("%#v\n", err)
	} else {
		if item.Kind == "Node" {
			n.AddNode(item.Name, item.Namespace, item.Kind, "",&item.Images)
		} else if item.OwnerReference != "" && item.Name == node.OwnerReference {
			//previous item references to me and I do  have parent ,  find my parent
			genTree(item)
			n.AddNode(item.Name, item.Namespace, item.Kind, item.OwnerReference,&item.Images)
		} else if item.OwnerReference == "" && item.Name == node.OwnerReference {
			//previous item references to me and I don;t have parent , print
			n.AddNode(item.Name, item.Namespace, item.Kind, item.OwnerReference,&item.Images)
		}else if item.Kind != "Pod"{//don't know what it is
			n.AddNode(item.Name, item.Namespace, item.Kind, item.OwnerReference,&item.Images)
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
			var envList []v1.EnvVar
			var depImages []string
			if len(ds.Spec.Template.Spec.Containers) > 0 {
				//deploymentImage := ds.Spec.Template.Spec.Containers[0].Image
				//fmt.Println(deploymentImage)

				envList = ds.Spec.Template.Spec.Containers[0].Env

				for _,envlist :=range envList {
					if strings.Contains(envlist.Name,"IMAGE"){
						depImages=append(depImages,envlist.Name)
					}
				}
			}

			if len(ors) > 0 {
				fmt.Println(ds.Name, ds.Namespace, ds.Kind, ors[0].Name)
				refSource = &report.RefResult{Name: ds.Name, Namespace: namespace, Kind: kind, OwnerReference: ors[0].Name, Ownerkind: ors[0].Kind,Images:depImages}
			} else {
				//fmt.Println(pod.Name, pod.Namespace, pod.Kind, "")
				refSource = &report.RefResult{Name: ds.Name, Namespace: namespace, Kind: kind,Images:depImages}
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
	case "StatefulSet":
		{
			ds, errs := clientset.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})
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
	case "ConfigMap":
		refSource = &report.RefResult{Name: name, Namespace: namespace, Kind: "ConfigMap"}
	default:
		{
			refSource = &report.RefResult{Name: name, Namespace: namespace, Kind: kind}
			fmt.Printf("Unknown Type Name  %s  Kind %s", name, kind)
		}
	}
	return
}


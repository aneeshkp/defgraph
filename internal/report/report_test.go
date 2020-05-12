package report_test

import (
	"encoding/json"
	"fmt"
	report "github.com/aneeshkp/depgraph/internal/report"
	"testing"
)



var pod1 =`[{
"name": "service-ca-pod_1",
"kind": "Pod",
"namespace": "openshift-service-ca-operator",
"ownerkind": "ReplicaSet",
"ownerReference": "service-ca-replicaset_1"
},
{
"name": "service-ca-pod_2",
"kind": "Pod",
"namespace": "openshift-service-ca",
"ownerkind": "ReplicaSet",
"ownerReference": "service-ca-replicaset_2"
},
{
"name": "service-ca-pod_3",
"kind": "Pod",
"namespace": "openshift-service-catalog-apiserver-operator",
"ownerkind": "ReplicaSet",
"ownerReference": "service-ca-replicaset_3"
}]`
 var rs1 =`[{
"name": "service-ca-replicaset_1",
"kind": "ReplicaSet",
"namespace": "openshift-service-ca",
"ownerkind": "Deployment",
"ownerReference": "service-ca-deployment-1"
},
{
"name": "service-ca-replicaset_2",
"kind": "ReplicaSet",
"namespace": "openshift-service-catalog-apiserver-operator",
"ownerkind": "Deployment",
"ownerReference": "service-ca-deployment-1"
},
{
"name": "service-ca-replicaset_3",
"kind": "ReplicaSet",
"namespace": "openshift-service-catalog-controller-manager-operator",
"ownerkind": "Deployment",
"ownerReference": "service-ca-deployment-2"
}]`
var dep1=`[
{
"name": "service-ca-deployment-1",
"namespace": "openshift-service-ca-operator",
"kind": "Deployment"
},
{
"name": "service-ca-deployment-2",
"namespace": "openshift-service-ca",
"kind": "Deployment"
}]`


var n=report.NewNodeTable()
func TestNodes(t *testing.T) {
	pod:=make([]report.RefResult,0)
	json.Unmarshal([]byte(pod1), &pod)
	for _,item :=range pod {
		fmt.Println("")
		fmt.Println(item.Name)
		fmt.Println(item.Kind)
		gen(&item)
		n.AddNode(item.Name,item.Namespace,item.Kind,item.OwnerReference)
	}
	n.ShowAll()

}

func gen(node *report.RefResult){
	if node.OwnerReference=="" {
		n.AddNode(node.Name,node.Namespace,node.Kind,node.OwnerReference)
	}else{
		  if node.Ownerkind=="Deployment"{
			  dep:=make([]report.RefResult,0)
			  if ok:= json.Unmarshal([]byte(dep1), &dep); ok!=nil  {
				  fmt.Printf("error %#v",ok)
			  }else {

				  for _, item := range dep {
					  if item.OwnerReference != "" && item.Name == node.OwnerReference {
						  //previous item references to me and I do  have parent ,  find my parent
						  fmt.Println("Calling gen for deployment")
						  gen(&item)
						  n.AddNode(item.Name, item.Namespace, item.Kind, item.OwnerReference)
						  break
					  } else if item.OwnerReference == "" && item.Name == node.OwnerReference {
					  	//previous item references to me and I don;t have parent , print
						  n.AddNode(item.Name, item.Namespace, item.Kind, item.OwnerReference)
						  break
					  }
				  }
			  }
		  }else if node.Ownerkind=="ReplicaSet"{
			rs:=make([]report.RefResult,0)
			if ok:= json.Unmarshal([]byte(rs1), &rs); ok!=nil  {
				fmt.Printf("error %#v",ok)
			}else{
				for _,item :=range rs {
					if item.OwnerReference!="" && item.Name==node.OwnerReference {
						//previous item references to me and I do  have parent ,  find my parent
						gen(&item)
						n.AddNode(item.Name,item.Namespace,item.Kind,item.OwnerReference)
						break
					}else if item.OwnerReference=="" && item.Name==node.OwnerReference  {
						//previous item references to me and I don;t have parent , print
						n.AddNode(item.Name,item.Namespace,item.Kind,item.OwnerReference)
						break
					}
				}
			}
			}
		}

	}




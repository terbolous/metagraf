/*
Copyright 2019 The metaGraf Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package modules

import (
	"context"
	"fmt"
	"github.com/openshift/api/image/docker10"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	log "k8s.io/klog"
	"metagraf/mg/k8sclient"
	"metagraf/mg/params"
	"metagraf/pkg/helpers"
	"metagraf/pkg/imageurl"
	"metagraf/pkg/metagraf"
	"os"
	"strconv"
	"strings"
)

func GenService(mg *metagraf.MetaGraf) {
	objname := Name(mg)

	var serviceports []corev1.ServicePort

	var DockerImage string
	if len(mg.Spec.BaseRunImage) > 0 {
		DockerImage = mg.Spec.BaseRunImage
	} else if len(mg.Spec.BuildImage) > 0 {
		DockerImage = mg.Spec.BuildImage
	} else {
		DockerImage = mg.Spec.Image
	}

	var imgurl imageurl.ImageURL
	_ = imgurl.Parse(DockerImage)

	ImageInfo := &docker10.DockerImage{}

	client := k8sclient.GetImageClient()
	ist := helpers.GetImageStreamTags(
		client,
		imgurl.Namespace,
		imgurl.Image+":"+imgurl.Tag)
	ImageInfo = helpers.GetDockerImageFromIST(ist)

	for k := range ImageInfo.Config.ExposedPorts {
		ss := strings.Split(k, "/")
		port, _ := strconv.Atoi(ss[0])
		ContainerPort := corev1.ServicePort{
			Name:     strings.ToLower(ss[0]) + "-" + ss[1],
			Port:     int32(port),
			Protocol: corev1.Protocol(strings.ToUpper(ss[1])),
			TargetPort: intstr.IntOrString{
				Type:   0,
				IntVal: int32(port),
				StrVal: ss[1],
			},
		}
		serviceports = append(serviceports, ContainerPort)
	}

	// Handle no serviceports
	if len(serviceports) == 0 {
		serviceports = append(
			serviceports,
			corev1.ServicePort{
				Name:     "http",
				Port:     int32(80),
				Protocol: corev1.Protocol("TCP"),
				TargetPort: intstr.IntOrString{
					Type:   0,
					IntVal: int32(8080),
					StrVal: "8080",
				},
			})
	}

	selectors := make(map[string]string)
	selectors["app"] = objname

	labels := make(map[string]string)
	labels["app"] = objname

	obj := corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   objname,
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports:           serviceports,
			Selector:        selectors,
			Type:            "ClusterIP",
			SessionAffinity: "None",
		},
	}


	if !Dryrun {
		StoreService(obj)
	}
	if Output {
		MarshalObject(obj.DeepCopyObject())
	}

	// Optinonally also create a ServiceMonitor resource.
	if params.ServiceMonitor {
		if Output && Format == "yaml" {
			fmt.Println("---")
		}
		GenServiceMonitor(mg, &obj)
	}
}

func StoreService(obj corev1.Service) {
	client := k8sclient.GetCoreClient().Services(NameSpace)
	svc, _ := client.Get(context.TODO(), obj.Name, metav1.GetOptions{})

	if len(svc.ResourceVersion) > 0 {
		obj.ResourceVersion = svc.ResourceVersion
		obj.Spec.ClusterIP = svc.Spec.ClusterIP
		_, err := client.Update(context.TODO(), &obj, metav1.UpdateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Updated Service: ", obj.Name, " in Namespace: ", NameSpace)
	} else {
		_, err := client.Create(context.TODO(), &obj, metav1.CreateOptions{})
		if err != nil {
			log.Error(err)
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("Created Service: ", obj.Name, " in Namespace: ", NameSpace)
	}
}

func DeleteService(name string) {
	client := k8sclient.GetCoreClient().Services(NameSpace)

	_, err := client.Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println("Service: ", name, "does not exist in namespace: ", NameSpace,", skipping...")
		return
	}

	err = client.Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println( "Unable to delete Service: ", name, " in namespace: ", NameSpace)
		log.Error(err)
		return
	}
	fmt.Println("Deleted Service: ", name, ", in namespace: ", NameSpace)
}
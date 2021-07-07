package main

import (
	"context"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Does a full catalogs of the images and namespaces on the current cluster
// Also builds a image->namespaces map to find easily in which namespace an image is used (and how many times)
func ClusterImageCatalog(clientset *kubernetes.Clientset) (ns map[string]VulnSummary, imgs map[string]VulnSummary, imgs_ns map[string]map[string]uint32, err error) {
	err = nil
	ns = make(map[string]VulnSummary)
	imgs = make(map[string]VulnSummary)
	imgs_ns = make(map[string]map[string]uint32)

	// Retrieve pods from kubeapi
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}

	// Populate the sets: namespaces, images, and image->namespace map
	for _, pod := range pods.Items {
		ns[pod.ObjectMeta.Namespace] = VulnSummary{}
		for _, container := range pod.Spec.Containers {
			imgs[container.Image] = VulnSummary{}
			_, ok := imgs_ns[container.Image]
			if !ok {
				imgs_ns[container.Image] = make(map[string]uint32)
			}
			imgs_ns[container.Image][pod.ObjectMeta.Namespace] += 1
			log.Logger.Debug().Uint32("counter", imgs_ns[container.Image][pod.ObjectMeta.Namespace]).Str("image", container.Image).Str("namespace", pod.ObjectMeta.Namespace).Msg("Namespace/Image collected")
		}
	}
	return
}

package apiservice

import (
	"github.com/go-fuego/fuego"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *AppAgentRepository) DiscoveryGroup(c fuego.ContextNoBody) (metav1.APIGroup, error) {
	return metav1.APIGroup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "APIGroup",
			APIVersion: "v1",
		},
		Name: "agent.app.multi.ch",
		Versions: []metav1.GroupVersionForDiscovery{
			{
				GroupVersion: "agent.app.multi.ch/v1",
				Version:      "v1",
			},
		},
		PreferredVersion: metav1.GroupVersionForDiscovery{
			GroupVersion: "agent.app.multi.ch/v1",
			Version:      "v1",
		},
	}, nil
}

func (r *AppAgentRepository) DiscoveryV1ResourceList(c fuego.ContextNoBody) (metav1.APIResourceList, error) {
	return metav1.APIResourceList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "APIResourceList",
			APIVersion: "v1",
		},
		GroupVersion: "agent.app.multi.ch/v1",
		APIResources: []metav1.APIResource{
			{
				Name:       "apps/start",
				Namespaced: true,
				Kind:       "App",
				Verbs:      []string{"get"},
			},
			{
				Name:       "apps/stop",
				Namespaced: true,
				Kind:       "App",
				Verbs:      []string{"get"},
			},
			{
				Name:       "apps/restart",
				Namespaced: true,
				Kind:       "App",
				Verbs:      []string{"get"},
			},
		},
	}, nil
}

package apiservice

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-fuego/fuego"
	"k8s.io/apimachinery/pkg/types"
	appv1 "multi.ch/app/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AppAgentRepository struct {
	cluster client.Client
}

func NewAppAgentRepository(cluster client.Client) *AppAgentRepository {
	return &AppAgentRepository{
		cluster: cluster,
	}
}

func (r *AppAgentRepository) Register(s *fuego.Server) {
	agentGroup := fuego.Group(s, "/apis/agent.app.multi.ch")
	versionedAgentGroup := fuego.Group(agentGroup, "/v1")
	namespacedAgentGroup := fuego.Group(versionedAgentGroup, "/namespaces/{namespace}/apps/{resource}/{action}")

	fuego.Get(agentGroup, "", r.DiscoveryGroup)
	fuego.Get(versionedAgentGroup, "", r.DiscoveryV1ResourceList)
	fuego.Get(namespacedAgentGroup, "", r.DoAction)
}

func (r *AppAgentRepository) DoAction(c fuego.ContextNoBody) (*APIResponse[string], error) {
	namespace := c.PathParam("namespace")
	resource := c.PathParam("resource")
	action := c.PathParam("action")

	var key = types.NamespacedName{
		Namespace: namespace,
		Name:      resource,
	}

	var app appv1.App
	if err := r.cluster.Get(c, key, &app); err != nil {
		return nil, fuego.NotFoundError{
			Detail: err.Error(),
		}
	}

	svcRef, found := app.Status.ChildResources.Get("", "Service", app.Name)
	if !found {
		return nil, fuego.NotFoundError{
			Detail: "service is not ready",
		}
	}

	queryParams := make(url.Values)
	queryParams.Set("processname", "workload")
	queryParams.Set("action", action)

	url := &url.URL{
		Scheme:   "http",
		Host:     fmt.Sprintf("%s.%s.svc:1080", svcRef.Name, namespace),
		Path:     "/index.html",
		RawQuery: queryParams.Encode(),
	}

	httpReq, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fuego.InternalServerError{
			Detail: err.Error(),
		}
	}

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fuego.InternalServerError{
			Detail: err.Error(),
		}
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fuego.InternalServerError{
			Detail: fmt.Sprintf("failed to call %s, status code: %d", url.String(), resp.StatusCode),
		}
	}

	return NewAPIResponse(action), nil
}

package api

import (
  "fmt"
	"github.com/emicklei/go-restful"
	"github.com/go-redis/redis"
	"net/http"

  "github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/pkg/tracing"
  "github.com/opentracing/opentracing-go"
)

type Questions struct {
	Content []Question
}

type Question struct {
	Body string
}

type QuestionsResource struct {
	backingStore *redis.Client
  tracer opentracing.Tracer
}

func NewQuestionsResource(client *redis.Client, tracingClient string) *QuestionsResource {
  return &QuestionsResource{backingStore: client, tracer : tracing.Init("questionsApi", tracingClient)}
}

func (resource QuestionsResource) Register(container *restful.Container) {
	ws := new(restful.WebService)

	ws.Path("/questions")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(resource.getAll))
	ws.Route(ws.PUT("").To(resource.add))

	container.Add(ws)
}

func (resource QuestionsResource) getAll(req *restful.Request, resp *restful.Response) {
  parent := resource.tracer.StartSpan("getAllRequest")

  child := resource.tracer.StartSpan("fetchDataFromRedis", opentracing.ChildOf(parent.Context()))
	result, err := resource.backingStore.LRange("questions", 0, -1).Result()
  child.LogEvent(fmt.Sprintf("Fetching data from redis"))

	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
	}

	data := Questions{}

  child.LogEvent(fmt.Sprintf("%v", data))

	for _, res := range result {
		data.Content = append(data.Content, Question{Body: res})
	}

	resp.WriteEntity(data)
  child.Finish()
  parent.Finish()
}

func (resource QuestionsResource) add(req *restful.Request, resp *restful.Response) {
  parent := resource.tracer.StartSpan("add")
  child := resource.tracer.StartSpan("pushDataToRedis", opentracing.ChildOf(parent.Context()))

	question := Question{}
	err := req.ReadEntity(&question)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
	}
	resource.backingStore.RPush("questions", question.Body)
	resp.WriteHeaderAndEntity(http.StatusCreated, question)

  child.LogEvent(fmt.Sprintf("%v", questionBody))
  child.Finish()
  parent.Finish()
}

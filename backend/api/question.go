package api

import (
	"github.com/emicklei/go-restful"
	"github.com/go-redis/redis"
	"net/http"
)

type Questions struct {
	Content []Question
}

type Question struct {
	Body string
}

type QuestionsResource struct {
	backingStore *redis.Client
}

func NewQuestionsResource(client *redis.Client) *QuestionsResource {
	return &QuestionsResource{backingStore: client}
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
	result, err := resource.backingStore.LRange("questions", 0, -1).Result()
	if err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
	}

	data := Questions{}
	for _, res := range result {
		data.Content = append(data.Content, Question{Body: res})
	}

	resp.WriteEntity(data)
}

func (resource QuestionsResource) add(req *restful.Request, resp *restful.Response) {
	question := Question{}
	err := req.ReadEntity(&question)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
	}
	resource.backingStore.RPush("questions", question.Body)
	resp.WriteHeaderAndEntity(http.StatusCreated, question)
}

package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/danilomarques1/graphmusic/api/resolver"
	"github.com/danilomarques1/graphmusic/model"
	chi "github.com/go-chi/chi/v5"
	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server struct {
	Router *chi.Mux
	port   string
	Client *mongo.Client
}

type Error struct {
	Message string `json:"message"`
}

func NewServer(port string) (*Server, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		return nil, err
	}
	router := chi.NewRouter()
	return &Server{Router: router, port: port, Client: client}, nil
}

func (s *Server) Init() {
	musicRepository := model.NewMusicRepository(s.Client, os.Getenv("DB"))
	musicResolver := resolver.NewMusicResolver(musicRepository)
	musicType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Music",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)
	musicInput := graphql.NewInputObject(
		graphql.InputObjectConfig{
			Name: "InputMusic",
			Fields: graphql.InputObjectConfigFieldMap{
				"author": &graphql.InputObjectFieldConfig{
					Type: graphql.String,
				},
				"name": &graphql.InputObjectFieldConfig{
					Type: graphql.String,
				},
			},
		},
	)

	query := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"findAll": &graphql.Field{
					Type:    graphql.NewList(musicType),
					Resolve: musicResolver.FindAll,
				},
				"findByName": &graphql.Field{
					Type: musicType,
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: musicResolver.FindByName,
				},
			},
		})

	mutation := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "RootMutation",
			Fields: graphql.Fields{
				"addMusic": &graphql.Field{
					Type: musicType,
					Args: graphql.FieldConfigArgument{
						"name":   &graphql.ArgumentConfig{Type: graphql.String},
						"author": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: musicResolver.Save,
				},
				"removeByName": &graphql.Field{
					Type: musicType,
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: musicResolver.RemoveByName,
				},
				"removeById": &graphql.Field{
					Type: musicType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{Type: graphql.String},
					},
					Resolve: musicResolver.RemoveById,
				},
				"updateByName": &graphql.Field{
					Type: musicType,
					Args: graphql.FieldConfigArgument{
						"name":  &graphql.ArgumentConfig{Type: graphql.String},
						"music": &graphql.ArgumentConfig{Type: musicInput},
					},
					Resolve: musicResolver.UpdateByName,
				},
			},
		})

	schema, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    query,
			Mutation: mutation,
		})

	if err != nil {
		log.Fatal(err)
	}

	s.Router.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			RespondERR(w, http.StatusBadRequest, "Invalid body")
			return
		}
		result, err := s.executeQuery(string(b), schema)
		if err != nil {
			RespondERR(w, http.StatusBadRequest, err.Error())
			return
		}
		json.NewEncoder(w).Encode(result)
	})
}

func (s *Server) Start() {
	fmt.Printf("Server running on port %v\n", s.port)
	http.ListenAndServe(fmt.Sprintf(":%v", s.port), s.Router)
}

func (s *Server) executeQuery(query string, schema graphql.Schema) (interface{}, error) {
	log.Printf("Query = %v\n", query)
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})

	if len(result.Errors) > 0 {
		log.Printf("%v\n", result.Errors)
		log.Printf("Erro ao executar a query\n")
		return nil, result.Errors[0]
	}
	return result.Data, nil
}

func RespondERR(w http.ResponseWriter, statusCode int, msg string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&Error{Message: msg})
}

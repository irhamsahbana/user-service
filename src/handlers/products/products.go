package products

import (
	"encoding/json"
	"log"
	"net/http"
	"user-service/src/util/client"
	"user-service/src/util/middleware"

	"github.com/thedevsaddam/renderer"
)

type Handler struct {
	render *renderer.Render
}

func NewProductHandler(r *renderer.Render) *Handler {
	return &Handler{
		render: r,
	}
}
func (h *Handler) CreateShop(w http.ResponseWriter, r *http.Request) {
	type (
		Response struct {
			Data struct {
				Id        string `json:"id"`
				UserId    string `json:"user_id"`
				Name      string `json:"name"`
				CreatedAt string `json:"created_at"`
				UpdatedAt string `json:"updated_at"`
			} `json:"data"`
			Message string `json:"message"`
			Success bool   `json:"success"`
		}

		ResponseError struct {
			Errors  map[string]any `json:"errors"`
			Message string         `json:"message"`
			Success bool           `json:"success"`
		}

		request struct {
			Name string `json:"name"`
		}
	)
	var (
		responseSuccess Response
		ctx             = r.Context()
	)

	channel := make(chan client.Response)

	req := request{
		Name: "Shop 1",
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.render.JSON(w, http.StatusInternalServerError, err)
		return
	}
	netClient := client.NetClientRequest{
		NetClient:  &http.Client{},
		RequestUrl: "http://localhost:3000/api/shops",
		QueryParam: []client.QueryParams{
			{Param: "user_id", Value: middleware.GetUserID(ctx)},
		},
	}

	netClient.Post(req, channel)

	response := <-channel
	if response.Status != http.StatusCreated {
		log.Println("Error: ", response.Status, response.Err)
		var result ResponseError
		if err := json.Unmarshal(response.Res, &result); err != nil {
			log.Println(err)
		}

		h.render.JSON(w, response.Status, result)
		return
	}

	if err := json.Unmarshal(response.Res, &responseSuccess); err != nil {
		log.Println(err)
	}

	h.render.JSON(w, response.Status, responseSuccess)

}

package set

import (
	"encoding/json"
	"errors"
	"fmt"
	"gitlab.com/kendellfab/mtgstudio/internal/api"
	"net/http"
)

const (
	setRoute = "/sets"
)

type SetRepo interface {
	// TODO: Add pagination
	ListSets() ([]Set, error)
}

type RestSetRepo struct {
	endpoint *api.Endpoint
}

func NewRestSetRepo(endpoint *api.Endpoint) SetRepo {
	return &RestSetRepo{endpoint: endpoint}
}

func (rsr *RestSetRepo) ListSets() ([]Set, error) {

	req, err := rsr.endpoint.NewGetRequest(setRoute)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("error loading sets: %d", resp.StatusCode))
	}

	var sets api.Response[[]Set]
	err = json.NewDecoder(resp.Body).Decode(&sets)
	if err != nil {
		return nil, fmt.Errorf("error decoding sets: %w", err)
	}

	return sets.Data, nil
}

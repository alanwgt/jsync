// Copyright © 2022 Alan Weingartner <hi@alanwgt.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
//  1. Redistributions of source code must retain the above copyright notice,
//     this list of conditions and the following disclaimer.
//
//  2. Redistributions in binary form must reproduce the above copyright notice,
//     this list of conditions and the following disclaimer in the documentation
//     and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alanwgt/jsync/internal/model"
	"github.com/alanwgt/jsync/log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

type PropertiesResponse model.PaginatedResponse[[]model.Property]
type BrokersResponse model.PaginatedResponse[[]model.Broker]
type CondominiumsResponse model.PaginatedResponse[[]model.Condominium]
type BannersResponse model.PaginatedResponse[[]model.Banner]
type ActivePropertiesData struct {
	Total  int   `json:"total"`
	Result []int `json:"result"`
}
type ActivePropertiesResponse model.PaginatedResponse[ActivePropertiesData]

type Requester struct {
	webserviceKey      string
	client             *http.Client
	maxPages           int
	concurrentRequests int
}

type requestData struct {
	requester *Requester
	path      RoutePath
	startDate *time.Time
	page      int
}

func NewRequester(maxPages int, concurrentRequest int) *Requester {
	return &Requester{
		client:             &http.Client{Timeout: 10 * time.Second},
		maxPages:           maxPages,
		concurrentRequests: concurrentRequest,
	}
}

func (r *Requester) SetWebserviceKey(k string) {
	r.webserviceKey = k
}

func (r Requester) newUrl(path RoutePath, page int, startDate *time.Time) (*url.URL, error) {
	u, err := url.Parse(WebserviceEndpoint)
	if err != nil {
		return nil, err
	}

	u.Path += fmt.Sprintf("/%s/%s", r.webserviceKey, (string)(path))

	q := &url.Values{}
	q.Add("v", WebserviceVersion)
	q.Add("pageSize", "100")
	q.Add("page", strconv.Itoa(page))

	if startDate != nil {
		q.Add("start", strconv.FormatInt(startDate.Unix(), 10))
	}

	u.RawQuery = q.Encode()
	return u, nil
}

func (r Requester) get(u *url.URL) (*http.Response, error) {
	return r.client.Get(u.String())
}

func mappedGet[T any](r *Requester, path RoutePath, page int, startDate *time.Time) (T, error) {
	st := time.Now()
	u, err := r.newUrl(path, page, startDate)
	var emptyResponse T
	if err != nil {
		return emptyResponse, err
	}

	res, err := r.get(u)
	if err != nil {
		return emptyResponse, err
	}

	log.Debug().
		Str("url", u.String()).
		Int("page", page).
		Int("status_code", res.StatusCode).
		Str("duração", time.Now().Sub(st).Round(time.Millisecond).String()).
		Msg("requisição concluída")

	if res.StatusCode != 200 {
		return emptyResponse, errors.New(fmt.Sprintf("response da requisição para %s retornou status code %d", u.String(), res.StatusCode))
	}

	var mappedResponse T

	if err := json.NewDecoder(res.Body).Decode(&mappedResponse); err != nil {
		return emptyResponse, err
	}

	return mappedResponse, nil
}

func requestHandler[T any](res chan<- []T, rds <-chan requestData) {
	for rd := range rds {
		r, err := mappedGet[model.PaginatedResponse[[]T]](rd.requester, rd.path, rd.page, rd.startDate)

		if err != nil {
			res <- nil
		} else {
			res <- r.Data
		}
	}
}

func getAllPaginated[T any](r *Requester, path RoutePath, startDate *time.Time) ([]T, error) {
	// fazer a primeira requisição pra saber se precisamos paralelizar o resto
	response, err := mappedGet[model.PaginatedResponse[[]T]](r, path, 1, startDate)
	if err != nil {
		return nil, err
	}

	var maxPages int
	if response.MaxPages() < r.maxPages {
		maxPages = response.MaxPages()
	} else {
		maxPages = r.maxPages
	}

	// a jetimob manda a quantia de itens mesmo quando a resposta for vazia
	if maxPages == 1 || len(response.Data) < response.PageSize {
		return response.Data, nil
	}

	log.Debug().Int("concurrent_request", r.concurrentRequests).Str("path", string(path)).Msg("iniciando requisições em paralelo")

	items := response.Data
	wg := &sync.WaitGroup{}
	req := make(chan requestData)
	res := make(chan []T, r.concurrentRequests)
	maxConcurrentRequests := r.concurrentRequests

	if maxPages < maxConcurrentRequests {
		maxConcurrentRequests = maxPages
	}

	for w := 0; w < maxConcurrentRequests; w++ {
		go requestHandler(res, req)
	}

	go func(res chan []T, wg *sync.WaitGroup) {
		for r := range res {
			items = append(items, r...)
			wg.Done()
		}
	}(res, wg)

	for p := 2; p <= maxPages; p++ {
		wg.Add(1)
		req <- requestData{
			requester: r,
			path:      path,
			startDate: startDate,
			page:      p,
		}
	}

	close(req)
	wg.Wait()
	close(res)

	return items, nil

}

func (r Requester) GetProperties(startDate *time.Time) ([]model.Property, error) {
	return getAllPaginated[model.Property](&r, PropertiesPath, startDate)
}

func (r Requester) GetCondominiums() ([]model.Condominium, error) {
	return getAllPaginated[model.Condominium](&r, CondominiumPath, nil)
}

func (r Requester) GetActiveProperties() ([]int, error) {
	res, err := mappedGet[ActivePropertiesResponse](&r, ActivePropertiesPath, 1, nil)
	if err != nil {
		return nil, err
	}

	return res.Data.Result, nil
}

func (r Requester) GetBrokers() ([]model.Broker, error) {
	return getAllPaginated[model.Broker](&r, BrokersPath, nil)
}

func (r Requester) GetBanners() ([]model.Banner, error) {
	return getAllPaginated[model.Banner](&r, BannersPath, nil)
}

func GetFromFile(p string, into any) error {
	bs, err := os.ReadFile(p)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bs, into); err != nil {
		return err
	}

	return nil
}

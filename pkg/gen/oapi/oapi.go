// Package oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package oapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Defines values for CompanyType.
const (
	CompanyTypeCooperative        CompanyType = "cooperative"
	CompanyTypeCorporations       CompanyType = "corporations"
	CompanyTypeNonProfit          CompanyType = "nonProfit"
	CompanyTypeSoleProprietorship CompanyType = "soleProprietorship"
)

// Defines values for CompanyUpdateType.
const (
	CompanyUpdateTypeCooperative        CompanyUpdateType = "cooperative"
	CompanyUpdateTypeCorporations       CompanyUpdateType = "corporations"
	CompanyUpdateTypeNonProfit          CompanyUpdateType = "nonProfit"
	CompanyUpdateTypeSoleProprietorship CompanyUpdateType = "soleProprietorship"
)

// Defines values for ErrorCode.
const (
	CompanyNotFound     ErrorCode = "company_not_found"
	ForbiddenAccess     ErrorCode = "forbidden_access"
	Unauthorized        ErrorCode = "unauthorized"
	UnknownAuthError    ErrorCode = "unknown_auth_error"
	UnknownStorageError ErrorCode = "unknown_storage_error"
	ValidationError     ErrorCode = "validation_error"
)

// Company defines model for Company.
type Company struct {
	Description *string `json:"description,omitempty"`

	// Employees Amount of employees
	Employees  int64       `json:"employees"`
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Registered bool        `json:"registered"`
	Type       CompanyType `json:"type"`
}

// CompanyType defines model for Company.Type.
type CompanyType string

// CompanyInput defines model for CompanyInput.
type CompanyInput = Company

// CompanyUpdate defines model for CompanyUpdate.
type CompanyUpdate struct {
	Description *string            `json:"description,omitempty"`
	Employees   *int64             `json:"employees,omitempty"`
	Name        *string            `json:"name,omitempty"`
	Registered  *bool              `json:"registered,omitempty"`
	Type        *CompanyUpdateType `json:"type,omitempty"`
}

// CompanyUpdateType defines model for CompanyUpdate.Type.
type CompanyUpdateType string

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code ErrorCode `json:"code"`

	// Error Error message
	Error string `json:"error"`
}

// ErrorCode Error code
type ErrorCode string

// ParamCompanyUId defines model for paramCompanyUId.
type ParamCompanyUId = string

// PostCompaniesJSONRequestBody defines body for PostCompanies for application/json ContentType.
type PostCompaniesJSONRequestBody = CompanyInput

// PatchCompaniesUidJSONRequestBody defines body for PatchCompaniesUid for application/json ContentType.
type PatchCompaniesUidJSONRequestBody = CompanyUpdate

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// PostCompanies request with any body
	PostCompaniesWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostCompanies(ctx context.Context, body PostCompaniesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// DeleteCompaniesUid request
	DeleteCompaniesUid(ctx context.Context, uid ParamCompanyUId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetCompaniesUid request
	GetCompaniesUid(ctx context.Context, uid ParamCompanyUId, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PatchCompaniesUid request with any body
	PatchCompaniesUidWithBody(ctx context.Context, uid ParamCompanyUId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PatchCompaniesUid(ctx context.Context, uid ParamCompanyUId, body PatchCompaniesUidJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostCompaniesWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostCompaniesRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostCompanies(ctx context.Context, body PostCompaniesJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostCompaniesRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) DeleteCompaniesUid(ctx context.Context, uid ParamCompanyUId, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewDeleteCompaniesUidRequest(c.Server, uid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetCompaniesUid(ctx context.Context, uid ParamCompanyUId, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetCompaniesUidRequest(c.Server, uid)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PatchCompaniesUidWithBody(ctx context.Context, uid ParamCompanyUId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchCompaniesUidRequestWithBody(c.Server, uid, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PatchCompaniesUid(ctx context.Context, uid ParamCompanyUId, body PatchCompaniesUidJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPatchCompaniesUidRequest(c.Server, uid, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostCompaniesRequest calls the generic PostCompanies builder with application/json body
func NewPostCompaniesRequest(server string, body PostCompaniesJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostCompaniesRequestWithBody(server, "application/json", bodyReader)
}

// NewPostCompaniesRequestWithBody generates requests for PostCompanies with any type of body
func NewPostCompaniesRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/companies")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewDeleteCompaniesUidRequest generates requests for DeleteCompaniesUid
func NewDeleteCompaniesUidRequest(server string, uid ParamCompanyUId) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uid", runtime.ParamLocationPath, uid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/companies/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("DELETE", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetCompaniesUidRequest generates requests for GetCompaniesUid
func NewGetCompaniesUidRequest(server string, uid ParamCompanyUId) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uid", runtime.ParamLocationPath, uid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/companies/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPatchCompaniesUidRequest calls the generic PatchCompaniesUid builder with application/json body
func NewPatchCompaniesUidRequest(server string, uid ParamCompanyUId, body PatchCompaniesUidJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPatchCompaniesUidRequestWithBody(server, uid, "application/json", bodyReader)
}

// NewPatchCompaniesUidRequestWithBody generates requests for PatchCompaniesUid with any type of body
func NewPatchCompaniesUidRequestWithBody(server string, uid ParamCompanyUId, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "uid", runtime.ParamLocationPath, uid)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/companies/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// PostCompanies request with any body
	PostCompaniesWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCompaniesResponse, error)

	PostCompaniesWithResponse(ctx context.Context, body PostCompaniesJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCompaniesResponse, error)

	// DeleteCompaniesUid request
	DeleteCompaniesUidWithResponse(ctx context.Context, uid ParamCompanyUId, reqEditors ...RequestEditorFn) (*DeleteCompaniesUidResponse, error)

	// GetCompaniesUid request
	GetCompaniesUidWithResponse(ctx context.Context, uid ParamCompanyUId, reqEditors ...RequestEditorFn) (*GetCompaniesUidResponse, error)

	// PatchCompaniesUid request with any body
	PatchCompaniesUidWithBodyWithResponse(ctx context.Context, uid ParamCompanyUId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchCompaniesUidResponse, error)

	PatchCompaniesUidWithResponse(ctx context.Context, uid ParamCompanyUId, body PatchCompaniesUidJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchCompaniesUidResponse, error)
}

type PostCompaniesResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Company
	JSON400      *Error
	JSON401      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r PostCompaniesResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostCompaniesResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type DeleteCompaniesUidResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON401      *Error
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r DeleteCompaniesUidResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r DeleteCompaniesUidResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetCompaniesUidResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Company
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r GetCompaniesUidResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetCompaniesUidResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PatchCompaniesUidResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Company
	JSON400      *Error
	JSON401      *Error
	JSON404      *Error
	JSON500      *Error
}

// Status returns HTTPResponse.Status
func (r PatchCompaniesUidResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PatchCompaniesUidResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostCompaniesWithBodyWithResponse request with arbitrary body returning *PostCompaniesResponse
func (c *ClientWithResponses) PostCompaniesWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostCompaniesResponse, error) {
	rsp, err := c.PostCompaniesWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostCompaniesResponse(rsp)
}

func (c *ClientWithResponses) PostCompaniesWithResponse(ctx context.Context, body PostCompaniesJSONRequestBody, reqEditors ...RequestEditorFn) (*PostCompaniesResponse, error) {
	rsp, err := c.PostCompanies(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostCompaniesResponse(rsp)
}

// DeleteCompaniesUidWithResponse request returning *DeleteCompaniesUidResponse
func (c *ClientWithResponses) DeleteCompaniesUidWithResponse(ctx context.Context, uid ParamCompanyUId, reqEditors ...RequestEditorFn) (*DeleteCompaniesUidResponse, error) {
	rsp, err := c.DeleteCompaniesUid(ctx, uid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseDeleteCompaniesUidResponse(rsp)
}

// GetCompaniesUidWithResponse request returning *GetCompaniesUidResponse
func (c *ClientWithResponses) GetCompaniesUidWithResponse(ctx context.Context, uid ParamCompanyUId, reqEditors ...RequestEditorFn) (*GetCompaniesUidResponse, error) {
	rsp, err := c.GetCompaniesUid(ctx, uid, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetCompaniesUidResponse(rsp)
}

// PatchCompaniesUidWithBodyWithResponse request with arbitrary body returning *PatchCompaniesUidResponse
func (c *ClientWithResponses) PatchCompaniesUidWithBodyWithResponse(ctx context.Context, uid ParamCompanyUId, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PatchCompaniesUidResponse, error) {
	rsp, err := c.PatchCompaniesUidWithBody(ctx, uid, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePatchCompaniesUidResponse(rsp)
}

func (c *ClientWithResponses) PatchCompaniesUidWithResponse(ctx context.Context, uid ParamCompanyUId, body PatchCompaniesUidJSONRequestBody, reqEditors ...RequestEditorFn) (*PatchCompaniesUidResponse, error) {
	rsp, err := c.PatchCompaniesUid(ctx, uid, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePatchCompaniesUidResponse(rsp)
}

// ParsePostCompaniesResponse parses an HTTP response from a PostCompaniesWithResponse call
func ParsePostCompaniesResponse(rsp *http.Response) (*PostCompaniesResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostCompaniesResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Company
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseDeleteCompaniesUidResponse parses an HTTP response from a DeleteCompaniesUidWithResponse call
func ParseDeleteCompaniesUidResponse(rsp *http.Response) (*DeleteCompaniesUidResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &DeleteCompaniesUidResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetCompaniesUidResponse parses an HTTP response from a GetCompaniesUidWithResponse call
func ParseGetCompaniesUidResponse(rsp *http.Response) (*GetCompaniesUidResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetCompaniesUidResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Company
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParsePatchCompaniesUidResponse parses an HTTP response from a PatchCompaniesUidWithResponse call
func ParsePatchCompaniesUidResponse(rsp *http.Response) (*PatchCompaniesUidResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PatchCompaniesUidResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Company
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create a new company
	// (POST /companies)
	PostCompanies(ctx echo.Context) error
	// Delete a company
	// (DELETE /companies/{uid})
	DeleteCompaniesUid(ctx echo.Context, uid ParamCompanyUId) error
	// Get details of a company
	// (GET /companies/{uid})
	GetCompaniesUid(ctx echo.Context, uid ParamCompanyUId) error
	// Update a company partially
	// (PATCH /companies/{uid})
	PatchCompaniesUid(ctx echo.Context, uid ParamCompanyUId) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostCompanies converts echo context to params.
func (w *ServerInterfaceWrapper) PostCompanies(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostCompanies(ctx)
	return err
}

// DeleteCompaniesUid converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteCompaniesUid(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uid" -------------
	var uid ParamCompanyUId

	err = runtime.BindStyledParameterWithLocation("simple", false, "uid", runtime.ParamLocationPath, ctx.Param("uid"), &uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uid: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteCompaniesUid(ctx, uid)
	return err
}

// GetCompaniesUid converts echo context to params.
func (w *ServerInterfaceWrapper) GetCompaniesUid(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uid" -------------
	var uid ParamCompanyUId

	err = runtime.BindStyledParameterWithLocation("simple", false, "uid", runtime.ParamLocationPath, ctx.Param("uid"), &uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uid: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetCompaniesUid(ctx, uid)
	return err
}

// PatchCompaniesUid converts echo context to params.
func (w *ServerInterfaceWrapper) PatchCompaniesUid(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uid" -------------
	var uid ParamCompanyUId

	err = runtime.BindStyledParameterWithLocation("simple", false, "uid", runtime.ParamLocationPath, ctx.Param("uid"), &uid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uid: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PatchCompaniesUid(ctx, uid)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/companies", wrapper.PostCompanies)
	router.DELETE(baseURL+"/companies/:uid", wrapper.DeleteCompaniesUid)
	router.GET(baseURL+"/companies/:uid", wrapper.GetCompaniesUid)
	router.PATCH(baseURL+"/companies/:uid", wrapper.PatchCompaniesUid)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xWQW/jNhP9K8J831GJ5Y2cbXXrpotFgALNJacgMBhxZHErcVhy5F3X0H8vSMqWHWvR",
	"HOK0xfYoajjzZua9GW6hpNaQRs0Oii0YYUWLjHb8uqHWCL25v5X+SKIrrTKsSEMBw7+kUxJSwK+iNQ1C",
	"AYtFhj/kWXaB7358usjnMr8Q7+fXF3l+fb1Y5HmWZRmkoLwPI7iGFLRo/c3oyeLvnbIooWDbYQqurLEV",
	"Pj5vjDdzbJVeQd/3u58B8IAnYLdk0LJCdwJ7C634+gvqFddQXGVZlj53mwK2pqENTtyGn1rqNCdUJaNR",
	"ChXZVjAUoDRf57B3qTTjCq33qeREBrvMj0DNFxOQLK6UYwxl2bt5ImpQaP8/nmwBdddC8QAlWUNWeNQe",
	"oCZ9Z6lSDCmU5IsjWK0RUnDU4J0lYxUyWVcrA48n8fvDrjxAaFNAflirI5CDi9EVPX3Gkj3UoU+32nTs",
	"IYum+bWC4mEL/7dYQQH/m428nA0Nnu262z8ecG172Kq576Uv80sZOLDuY/SW3JA1x0kMBBwyOCppPyZy",
	"b6RgfG3avYBS/xrunFDgo7VkTytWksRTyQXjJPxL9xg7/ZumL3opOq6XGNyl0Gn/SVb9ERhYkX1SUqJe",
	"irJE53NZi0bJkNn+Uhm7uNTEy4o6LYOj6N0xWbHCwfZxqme7TKYwt+icWCH8lZ5GKHJKMt5c6YpOw3CN",
	"iUO7ViUmxtJaSXQJamlIaXZJFerm01NBn6y4iUwez9ZoXXQ2v8wuM58TGdTCKCjg6jK7vII0jOnQodl4",
	"1XePXFDwwAnSfk3AHTm+OYjgM0XHH0huYpM1o47KN6ZRZbg5++yiTMZp/4JpEGdIf1xPL9pw4AxpF5G+",
	"y+avHTuGnV6JpUXBKBPXBeJVXdNsfGXzLHs1GFFDEyA+CJkMNY8x5+ePeX+ovD6FxVskeqsZrRZNkADa",
	"BAfDFFzXtsJufEdCJxKRaPwyaGETTEYmz7adkn0UV4Nxkh8z+udwvuf0fdh/h6+lb+yu0WT2/DXl19gz",
	"hubffmJFYFN8+ht6m0ek5w26y1wTJ3Eq/5NYFQmRiJFRKaxwYhZ+Qj43bbK3HGwSWajGJRbZKlxPUvK7",
	"Z8cn5H2hqDpmiRFc1hM70x+fgSln27zDo/dFq/dNGdoFYN/x6v1PgBDJOQovMcKyEoEJwTJcjbLqbAMF",
	"1MymmM0aKkVTk+Pi/XyRQf/Y/xkAAP//s1HdACQRAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

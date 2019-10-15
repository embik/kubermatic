// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewAddonParams creates a new AddonParams object
// with the default values initialized.
func NewAddonParams() *AddonParams {

	return &AddonParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewAddonParamsWithTimeout creates a new AddonParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewAddonParamsWithTimeout(timeout time.Duration) *AddonParams {

	return &AddonParams{

		timeout: timeout,
	}
}

// NewAddonParamsWithContext creates a new AddonParams object
// with the default values initialized, and the ability to set a context for a request
func NewAddonParamsWithContext(ctx context.Context) *AddonParams {

	return &AddonParams{

		Context: ctx,
	}
}

// NewAddonParamsWithHTTPClient creates a new AddonParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewAddonParamsWithHTTPClient(client *http.Client) *AddonParams {

	return &AddonParams{
		HTTPClient: client,
	}
}

/*AddonParams contains all the parameters to send to the API endpoint
for the addon operation typically these are written to a http.Request
*/
type AddonParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the addon params
func (o *AddonParams) WithTimeout(timeout time.Duration) *AddonParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the addon params
func (o *AddonParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the addon params
func (o *AddonParams) WithContext(ctx context.Context) *AddonParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the addon params
func (o *AddonParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the addon params
func (o *AddonParams) WithHTTPClient(client *http.Client) *AddonParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the addon params
func (o *AddonParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *AddonParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

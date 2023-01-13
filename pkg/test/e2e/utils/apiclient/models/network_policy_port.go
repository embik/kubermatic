// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NetworkPolicyPort NetworkPolicyPort describes a port to allow traffic on
//
// swagger:model NetworkPolicyPort
type NetworkPolicyPort struct {

	// If set, indicates that the range of ports from port to endPort, inclusive,
	// should be allowed by the policy. This field cannot be defined if the port field
	// is not defined or if the port field is defined as a named (string) port.
	// The endPort must be equal or greater than port.
	// This feature is in Beta state and is enabled by default.
	// It can be disabled using the Feature Gate "NetworkPolicyEndPort".
	// +optional
	EndPort int32 `json:"endPort,omitempty"`

	// port
	Port *IntOrString `json:"port,omitempty"`

	// protocol
	Protocol Protocol `json:"protocol,omitempty"`
}

// Validate validates this network policy port
func (m *NetworkPolicyPort) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePort(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProtocol(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NetworkPolicyPort) validatePort(formats strfmt.Registry) error {
	if swag.IsZero(m.Port) { // not required
		return nil
	}

	if m.Port != nil {
		if err := m.Port.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("port")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("port")
			}
			return err
		}
	}

	return nil
}

func (m *NetworkPolicyPort) validateProtocol(formats strfmt.Registry) error {
	if swag.IsZero(m.Protocol) { // not required
		return nil
	}

	if err := m.Protocol.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("protocol")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("protocol")
		}
		return err
	}

	return nil
}

// ContextValidate validate this network policy port based on the context it is used
func (m *NetworkPolicyPort) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidatePort(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProtocol(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *NetworkPolicyPort) contextValidatePort(ctx context.Context, formats strfmt.Registry) error {

	if m.Port != nil {
		if err := m.Port.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("port")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("port")
			}
			return err
		}
	}

	return nil
}

func (m *NetworkPolicyPort) contextValidateProtocol(ctx context.Context, formats strfmt.Registry) error {

	if err := m.Protocol.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("protocol")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("protocol")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *NetworkPolicyPort) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NetworkPolicyPort) UnmarshalBinary(b []byte) error {
	var res NetworkPolicyPort
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}

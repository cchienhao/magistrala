// Copyright (c) Abstract Machines
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/absmach/magistrala/pkg/errors"
)

const (
	invitationsEndpoint = "invitations"
	acceptEndpoint      = "accept"
	rejectEndpoint      = "reject"
)

type Invitation struct {
	InvitedBy   string    `json:"invited_by"`
	UserID      string    `json:"user_id"`
	DomainID    string    `json:"domain_id"`
	Token       string    `json:"token,omitempty"`
	Relation    string    `json:"relation,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	ConfirmedAt time.Time `json:"confirmed_at,omitempty"`
	RejectedAt  time.Time `json:"rejected_at,omitempty"`
	Resend      bool      `json:"resend,omitempty"`
}

type InvitationPage struct {
	Total       uint64       `json:"total"`
	Offset      uint64       `json:"offset"`
	Limit       uint64       `json:"limit"`
	Invitations []Invitation `json:"invitations"`
}

func (sdk mgSDK) SendInvitation(invitation Invitation, token string) (err error) {
	data, err := json.Marshal(invitation)
	if err != nil {
		return errors.NewSDKError(err)
	}

	url := sdk.invitationsURL + "/" + invitation.DomainID + "/" + invitationsEndpoint

	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, data, nil, http.StatusCreated)

	return sdkerr
}

func (sdk mgSDK) Invitation(userID, domainID, token string) (invitation Invitation, err error) {
	url := sdk.invitationsURL + "/" + domainID + "/" + invitationsEndpoint + "/" + usersEndpoint + "/" + userID

	_, body, sdkerr := sdk.processRequest(http.MethodGet, url, token, nil, nil, http.StatusOK)
	if sdkerr != nil {
		return Invitation{}, sdkerr
	}

	if err := json.Unmarshal(body, &invitation); err != nil {
		return Invitation{}, errors.NewSDKError(err)
	}

	return invitation, nil
}

func (sdk mgSDK) Invitations(pm PageMetadata, domainID, token string) (invitations InvitationPage, err error) {
	endpoint := fmt.Sprintf("%s/%s", domainID, invitationsEndpoint)

	url, err := sdk.withQueryParams(sdk.invitationsURL, endpoint, pm)
	if err != nil {
		return InvitationPage{}, errors.NewSDKError(err)
	}

	_, body, sdkerr := sdk.processRequest(http.MethodGet, url, token, nil, nil, http.StatusOK)
	if sdkerr != nil {
		return InvitationPage{}, sdkerr
	}

	var invPage InvitationPage
	if err := json.Unmarshal(body, &invPage); err != nil {
		return InvitationPage{}, errors.NewSDKError(err)
	}

	return invPage, nil
}

func (sdk mgSDK) AcceptInvitation(domainID, token string) (err error) {
	url := sdk.invitationsURL + "/" + domainID + "/" + invitationsEndpoint + "/" + acceptEndpoint

	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, nil, nil, http.StatusNoContent)

	return sdkerr
}

func (sdk mgSDK) RejectInvitation(domainID, token string) (err error) {
	url := sdk.invitationsURL + "/" + domainID + "/" + invitationsEndpoint + "/" + rejectEndpoint

	_, _, sdkerr := sdk.processRequest(http.MethodPost, url, token, nil, nil, http.StatusNoContent)

	return sdkerr
}

func (sdk mgSDK) DeleteInvitation(userID, domainID, token string) (err error) {
	url := sdk.invitationsURL + "/" + domainID + "/" + invitationsEndpoint + "/" + usersEndpoint + "/" + userID

	_, _, sdkerr := sdk.processRequest(http.MethodDelete, url, token, nil, nil, http.StatusNoContent)

	return sdkerr
}

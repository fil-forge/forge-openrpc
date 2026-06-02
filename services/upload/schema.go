package upload

import (
	"fmt"
	"slices"
	"strings"

	"github.com/alanshaw/go-openrpc"
	"github.com/fil-forge/forge-openrpc/schemas"
	accesscmds "github.com/fil-forge/libforge/commands/access"
	blobcmds "github.com/fil-forge/libforge/commands/blob"
	indexcmds "github.com/fil-forge/libforge/commands/index"
	providercmds "github.com/fil-forge/libforge/commands/provider"
	spacecmds "github.com/fil-forge/libforge/commands/space"
	ucancmds "github.com/fil-forge/libforge/commands/ucan"
	uploadcmds "github.com/fil-forge/libforge/commands/upload"
	shardcmds "github.com/fil-forge/libforge/commands/upload/shard"
)

type MethodDefinition struct {
	Name        string
	Description string
	Args        any
	Result      any
}

var defs = []MethodDefinition{
	{
		Name: accesscmds.Claim.Command.String(),
		Description: `
Issuer: agent
Subject: agent

Retrieves any delegations stored via ` + "`" + `/access/delegate` + "`" + ` or ` + "`" + `/access/confirm` + "`" + ` whose _audience_ matches the invocation subject.
		`,
		Args:   accesscmds.ClaimArguments{},
		Result: accesscmds.ClaimOK{},
	},
	{
		Name: accesscmds.Confirm.Command.String(),
		Description: `
Issuer: upload service
Subject: upload service

An invocation created in response to an access request (` + "`" + `/access/request` + "`" + `).

It is created and signed by the upload service, but _not_ executed. It is sent via email to an account holder. Clicking on the link in the email sends the invocation back to the upload service for execution.

Delegations requested from the ` + "`" + `/access/request` + "`" + ` invocation are copied into the ` + "`" + `/access/confirm` + "`" + ` invocation so that when it is executed, it can create the necessary delegations from the account (` + "`" + `did:mailto:` + "`" + `) to the agent. The account delegation requires an attestation, since ` + "`" + `did:mailto:` + "`" + ` is not a crypto key and cannot sign a delegation.

Delegations are sent in the invocation receipt, but also stored by the upload service for retrieval via ` + "`" + `/access/claim` + "`" + `.
		`,
		Args:   accesscmds.ConfirmArguments{},
		Result: accesscmds.ConfirmOK{},
	},
	{
		Name: accesscmds.Delegate.Command.String(),
		Description: `
Issuer: agent
Subject: space

Send delegations to the upload service for retrieval later via ` + "`" + `/access/claim` + "`" + `. Using the space as the subject allows the service to verify a paid account is setup, so that it is not storing delegations for arbitrary audiences.
		`,
		Args:   accesscmds.DelegateArguments{},
		Result: accesscmds.DelegateOK{},
	},
	{
		Name: accesscmds.Request.Command.String(),
		Description: `
Issuer: agent
Subject: agent

Request delegations for specified commands from account (email) in the invocation arguments. Delegations are addressed to the invocation subject.

An ` + "`" + `/access/confirm` + "`" + ` invocation is created and signed by the upload service, but _not_ executed. It is sent via email to an account holder.

Receipt of the ` + "`" + `/access/confirm` + "`" + ` invocation by the upload service confirms the account has consented to the request.
		`,
		Args:   accesscmds.RequestArguments{},
		Result: accesscmds.RequestOK{},
	},
	{
		Name: blobcmds.Add.Command.String(),
		Description: `
Issuer: agent
Subject: space

Requests to add a blob (identified by multihash and size) to a space. The response includes a promise for ` + "`" + `/http/put` + "`" + ` - an task for the agent to perform and send the receipt back to the uplaod service via ` + "`" + `/ucan/conclude` + "`" + `. It also includes promise for a ` + "`" + `/blob/accept` + "`" + ` invocation which is made by the upload service to the Piri node when the ` + "`" + `/http/put` + "`" + ` receipt is received.

The ` + "`" + `/blob/accept` + "`" + ` invocation returns a location commitment for the blob that describes where the blob may be retrieved from (via a ` + "`" + `/content/retrieve` + "`" + ` or ` + "`" + `/blob/retrieve` + "`" + ` invocation).
		`,
		Args:   blobcmds.AddArguments{},
		Result: blobcmds.AddOK{},
	},
	{
		Name: blobcmds.List.Command.String(),
		Description: `
Issuer: agent
Subject: space

Lists blobs stored in a space.
		`,
		Args:   blobcmds.ListArguments{},
		Result: blobcmds.ListOK{},
	},
	{
		Name: indexcmds.Add.Command.String(),
		Description: `
Issuer: agent
Subject: space

Informs the service that one of the blobs uploaded to the space is actually an index, 
		`,
		Args:   indexcmds.AddArguments{},
		Result: indexcmds.AddOK{},
	},
	{
		Name:        providercmds.Add.Command.String(),
		Description: ``,
		Args:        providercmds.AddArguments{},
		Result:      providercmds.AddOK{},
	},
	{
		Name:        spacecmds.Info.Command.String(),
		Description: ``,
		Args:        spacecmds.InfoArguments{},
		Result:      spacecmds.InfoOK{},
	},
	{
		Name:        ucancmds.Conclude.Command.String(),
		Description: ``,
		Args:        ucancmds.ConcludeArguments{},
		Result:      ucancmds.ConcludeOK{},
	},
	{
		Name:        uploadcmds.Add.Command.String(),
		Description: ``,
		Args:        uploadcmds.AddArguments{},
		Result:      uploadcmds.AddOK{},
	},
	{
		Name:        uploadcmds.List.Command.String(),
		Description: ``,
		Args:        uploadcmds.ListArguments{},
		Result:      uploadcmds.ListOK{},
	},
	{
		Name:        shardcmds.List.Command.String(),
		Description: ``,
		Args:        shardcmds.ListArguments{},
		Result:      shardcmds.ListOK{},
	},
}

func BuildSchema() (*openrpc.Schema, error) {
	methods := []*openrpc.Method{}
	references := slices.Clone(schemas.Common)

	for _, def := range defs {
		method, refs, err := schemas.BuildMethod(
			def.Name,
			strings.TrimSpace(def.Description),
			def.Args,
			def.Result,
			references,
		)
		if err != nil {
			return nil, fmt.Errorf("building method %s: %w", def.Name, err)
		}
		methods = append(methods, method)
		references = append(references, refs...)
	}

	componentSchemas := make(map[string]*openrpc.JSONSchema, len(references))
	for _, ref := range references {
		componentSchemas[ref.ID] = &openrpc.JSONSchema{
			Schema: ref.JSONSchema,
		}
	}

	return &openrpc.Schema{
		OpenRPC: "1.4.1",
		Info: &openrpc.Info{
			Title:       "sprue",
			Description: "The Forge Upload Service.",
			Version:     "0.0.0", // TODO: fetch current version?
			License: &openrpc.License{
				Name: "Apache-2.0 OR MIT",
				URL:  "https://github.com/storacha/sprue/blob/main/LICENSE",
			},
		},
		Servers: []*openrpc.Server{
			{
				Name: "did:web:upload.forge.fil.one",
				URL:  "https://upload.forge.fil.one",
			},
		},
		Methods:    methods,
		Components: &openrpc.Components{Schemas: componentSchemas},
	}, nil
}

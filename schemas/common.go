package schemas

import (
	"reflect"

	"github.com/fil-forge/ucantone/did"
	"github.com/fil-forge/ucantone/ucan/command"
	"github.com/fil-forge/ucantone/ucan/promise"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

var AwaitOK = Schema{
	ID:   "awaitOK",
	Type: reflect.TypeOf(promise.AwaitOK{}),
	JSONSchema: &jsonschema.Schema{
		Title:       "AwaitOK",
		Description: "A promise of task completion.",
		Type:        "object",
		Properties: map[string]*jsonschema.Schema{
			"await/ok": {
				Ref: "#/components/schemas/cid",
			},
		},
	},
}

var CID = Schema{
	ID:   "cid",
	Type: reflect.TypeOf(cid.Undef),
	JSONSchema: &jsonschema.Schema{
		Title:       "CID",
		Description: "An IPLD content identifier (CID).",
		Type:        "object",
	},
}

var Command = Schema{
	ID:   "command",
	Type: reflect.TypeOf(command.Undef),
	JSONSchema: &jsonschema.Schema{
		Title:       "Command",
		Description: "A UCAN command string, e.g. `/blob/add`.",
		Type:        "string",
	},
}

var DID = Schema{
	ID:   "did",
	Type: reflect.TypeOf(did.Undef),
	JSONSchema: &jsonschema.Schema{
		Title:       "DID",
		Description: "A decentralized identifier (DID).",
		Type:        "string",
	},
}

var Multihash = Schema{
	ID:   "multihash",
	Type: reflect.TypeOf(multihash.Multihash{}),
	JSONSchema: &jsonschema.Schema{
		Title:       "Multihash",
		Description: "An IPLD multihash.",
		Type:        "object",
	},
}

var Common = []Schema{AwaitOK, CID, Command, DID, Multihash}

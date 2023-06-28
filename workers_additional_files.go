package cloudflare

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Possible script types
const (
	Wasm             = "application/wasm"
	Javascript       = "application/javascript"
	JavascriptModule = "application/javascript+module"
)

type AdditionalFile struct {
	FileName    string
	ScriptType  string
	FileContent string
}

// AdditionalFileListItem a struct representing an individual binding in a list of bindings.
// type AdditionalFileListItem struct {
// 	Name           string         `json:"name"`
// 	AdditionalFile AdditionalFile `json:"additional_file"`
// }

var jsonRes struct {
	Response
	AdditionalFiles []AdditionalFile `json:"result"` //TODO needs work
}

// AdditionalFileListResponse wrapper struct for API response to additional files list API call.
type WorkerAdditionalFileListResponse struct {
	Response
	AdditionalFiles []AdditionalFile
}

type ListAdditionalFilesParams struct {
	ScriptName string
}

// ListAdditionalFiles returns all additional files for a particular worker.
func (api *API) ListAdditionalFiles(ctx context.Context, rc *ResourceContainer, params ListAdditionalFilesParams) (WorkerAdditionalFileListResponse, error) {
	if params.ScriptName == "" {
		return WorkerAdditionalFileListResponse{}, errors.New("ScriptName is required")
	}

	if rc.Level != AccountRouteLevel {
		return WorkerAdditionalFileListResponse{}, ErrRequiredAccountLevelResourceContainer
	}

	if rc.Identifier == "" {
		return WorkerAdditionalFileListResponse{}, ErrMissingAccountID
	}

	uri := fmt.Sprintf("/accounts/%s/workers/scripts/%s/bindings", rc.Identifier, params.ScriptName)

	var jsonRes struct {
		Response
		//Bindings []workerBindingMeta `json:"result"`
		AdditionalFiles []AdditionalFile `json:"result"` //TODO needs work
	}
	var r WorkerAdditionalFileListResponse
	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(res, &jsonRes)
	if err != nil {
		return r, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	r = WorkerAdditionalFileListResponse{
		Response:        jsonRes.Response,
		AdditionalFiles: make([]AdditionalFile, 0, len(jsonRes.AdditionalFiles)),
	}
	for _, jsonAdditionalFile := range jsonRes.AdditionalFiles {
		filename := jsonAdditionalFile.FileName
		// name := jsonAdditionalFile.FileName
		// if !ok {
		// 	return r, fmt.Errorf("Binding missing name %v", jsonAdditSionalFile)
		// }
		// bType, ok := jsonAdditionalFile["type"].(string)
		// if !ok {
		// 	return r, fmt.Errorf("Binding missing type %v", jsonAdditionalFile)
		// }
		// bindingListItem := WorkerAdditionalFileItem{
		// 	Name: name,
		// }
		additionalFile := AdditionalFile{
			FileName: filename,
			//	AdditionalFile: AdditionalFile{},
		}
		r.AdditionalFiles = append(r.AdditionalFiles, additionalFile)
	}
	return r, nil
}

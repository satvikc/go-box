package box

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Folder struct {
	Id                string       `json:"id,omitempty"`                  // The folder’s ID.
	SequenceId        string       `json:"sequence_id,omitempty"`         // A unique ID for use with the /events endpoint.
	ETag              string       `json:"etag,omitempty"`                // A unique string identifying the version of this folder.
	Name              string       `json:"name,omitempty"`                // The name of this folder.
	Description       string       `json:"description,omitempty"`         // The description of this folder.
	Size              int          `json:"size,omitempty"`                // Size of this file in bytes.
	PathCollection    Collection   `json:"path_collection,omitempty"`     // The path of folders to this item, starting at the root.
	CreatedAt         BoxTime      `json:"created_at,omitempty"`          // The time the folder was created.
	ModifiedAt        BoxTime      `json:"modified_at,omitempty"`         // The time the folder or its contents were last modified.
	ThrashedAt        BoxTime      `json:"thrashed_at,omitempty"`         // The time the folder or its contents were put in the trash.
	PurgedAt          BoxTime      `json:"purged_at,omitempty"`           // The time the folder or its contents were purged from the trash.
	ContentCreatedAt  BoxTime      `json:"content_created_at,omitempty"`  // The time the folder or its contents were originally created (according to the uploader).
	ContentModifiedAt BoxTime      `json:"content_modified_at,omitempty"` // The time the folder or its contents were last modified (according to the uploader).
	CreatedBy         Entity       `json:"created_by,omitempty"`          // The user who created this folder.
	ModifiedBy        Entity       `json:"modified_by,omitempty"`         // The user who last modified this folder.
	OwnedBy           Entity       `json:"owned_by,omitempty"`            // The user who owns this file.
	SharedLink        SharedObject `json:"shared_link,omitempty"`         // The shared link for this folder. Null if not set..
	Parent            Entity       `json:"parent,omitempty"`              // The folder that contains this one.
	ItemStatus        string       `json:"item_status,omitempty"`         // Whether this item is deleted or not.
	Permissions       Permission   `json:"permissions,omitempty"`         // The permissions that the current user has on this file.
	Tags              []string     `json:"tags,omitempty"`                // All tags applied to this file.
	HasCollaborations bool         `json:"has_collaborations,omitempty"`  // Whether this folder has any collaborators.
	SyncStatus        string       `json:"sync_status,omitempty"`         // Whether this folder will be synced by the Box sync clients or not. Can be
	ItemCollection    Collection   `json:"item_collection,omitempty"`     // A collection of mini file and folder objects contained in this folder.
	FolderUploadEmail string       `json:"folder_upload_email,omitempty"` // The upload email address for this folder. Null if not set.
}

// CreateFolder creates the folder name under the folder having
// parentid as id. It also returns the created Folder.
func (box *Box) CreateFolder(parentid, name string) (*Folder, error) {
	var rv Folder
	reqBody := fmt.Sprintf(`{"name":"%s", "parent": {"id" : "%s"}}`, name, parentid)
	body, err := box.doRequest("POST", "folders", nil, reqBody)
	if err == nil {
		err = json.Unmarshal(body, &rv)
		return &rv, err
	}
	return nil, err
}

// GetFolder returns the Folder with the given id,
func (box *Box) GetFolder(id string) (*Folder, error) {
	var rv Folder
	rawurl := fmt.Sprintf("folders/%s", id)
	body, err := box.doRequest("GET", rawurl, nil, "")
	if err == nil {
		err = json.Unmarshal(body, &rv)
		return &rv, err
	}
	return nil, err
}

func (box *Box) DeleteFolder(id string) error {
	rawurl := fmt.Sprintf("folders/%s", id)
	_, err := box.doRequest("DELETE", rawurl, &url.Values{"recursive": {"true"}}, "")
	return err
}

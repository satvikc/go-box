package box

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type Folder struct {
	Id                string        `json:"id,omitempty"`                  // The folder’s ID.
	SequenceId        string        `json:"sequence_id,omitempty"`         // A unique ID for use with the /events endpoint.
	ETag              string        `json:"etag,omitempty"`                // A unique string identifying the version of this folder.
	Name              string        `json:"name,omitempty"`                // The name of this folder.
	Description       string        `json:"description,omitempty"`         // The description of this folder.
	Size              int           `json:"size,omitempty"`                // Size of this file in bytes.
	PathCollection    *Collection   `json:"path_collection,omitempty"`     // The path of folders to this item, starting at the root.
	CreatedAt         *BoxTime      `json:"created_at,omitempty"`          // The time the folder was created.
	ModifiedAt        *BoxTime      `json:"modified_at,omitempty"`         // The time the folder or its contents were last modified.
	ThrashedAt        *BoxTime      `json:"thrashed_at,omitempty"`         // The time the folder or its contents were put in the trash.
	PurgedAt          *BoxTime      `json:"purged_at,omitempty"`           // The time the folder or its contents were purged from the trash.
	ContentCreatedAt  *BoxTime      `json:"content_created_at,omitempty"`  // The time the folder or its contents were originally created (according to the uploader).
	ContentModifiedAt *BoxTime      `json:"content_modified_at,omitempty"` // The time the folder or its contents were last modified (according to the uploader).
	CreatedBy         *Entity       `json:"created_by,omitempty"`          // The user who created this folder.
	ModifiedBy        *Entity       `json:"modified_by,omitempty"`         // The user who last modified this folder.
	OwnedBy           *Entity       `json:"owned_by,omitempty"`            // The user who owns this file.
	SharedLink        *SharedObject `json:"shared_link,omitempty"`         // The shared link for this folder. Null if not set..
	Parent            *Entity       `json:"parent,omitempty"`              // The folder that contains this one.
	ItemStatus        string        `json:"item_status,omitempty"`         // Whether this item is deleted or not.
	Permissions       *Permission   `json:"permissions,omitempty"`         // The permissions that the current user has on this file.
	Tags              []string      `json:"tags,omitempty"`                // All tags applied to this file.
	HasCollaborations bool          `json:"has_collaborations,omitempty"`  // Whether this folder has any collaborators.
	SyncStatus        string        `json:"sync_status,omitempty"`         // Whether this folder will be synced by the Box sync clients or not. Can be
	ItemCollection    *Collection   `json:"item_collection,omitempty"`     // A collection of mini file and folder objects contained in this folder.
	FolderUploadEmail *UploadEmail  `json:"folder_upload_email,omitempty"` // The upload email address for this folder. Null if not set.
}

// Items returns all items (folder or files) under the given
// folder. It calls Get if the folder is not already populated.
func (f *Folder) Items(box *Box) ([]Entity, error) {
	if f.ItemCollection == nil {
		if err := f.Get(box); err != nil {
			return nil, err
		}
	}

	return f.ItemCollection.Entry, nil
}

// Create creates a sub folder under the given folder. It returns the
// created folder. Note that only Id is required apriori.
func (f *Folder) Create(box *Box, name string) (*Folder, error) {
	if f.Id == "" {
		return nil, errors.New("Empty id while using Create")
	}

	fold := Folder{Name: name, Parent: &Entity{Id: f.Id}}
	reqBody, _ := json.Marshal(fold)

	body, err := box.doRequest("POST", "folders", nil, reqBody)

	if err == nil {
		err = json.Unmarshal(body, &fold)
		return &fold, err
	}
	return nil, err
}

// Get populates the fields of the struct. Node that only Id is
// required apriori.
func (f *Folder) Get(box *Box) error {
	if f.Id == "" {
		return errors.New("Empty id while using Get")
	}
	rawurl := fmt.Sprintf("folders/%s", f.Id)
	body, err := box.doRequest("GET", rawurl, nil, nil)

	if err == nil {
		err = json.Unmarshal(body, f)
		return err
	}
	return err
}

// Delete deletes the folder. Note that only Id is required apriori.
func (f *Folder) Delete(box *Box) error {
	if f.Id == "" {
		return errors.New("Empty id while using Delete")
	}

	rawurl := fmt.Sprintf("folders/%s", f.Id)
	_, err := box.doRequest("DELETE", rawurl, &url.Values{"recursive": {"true"}}, nil)

	return err
}

// Rename renames the folder with the new name. Note that only Id is
// required apriori. The folder is populated with all the information
// after the call.
func (f *Folder) Rename(box *Box, name string) error {
	if f.Id == "" {
		return errors.New("Empty id while using Rename")
	}

	fold := Folder{Name: name}
	reqBody, _ := json.Marshal(fold)

	rawurl := fmt.Sprintf("folders/%s", f.Id)
	body, err := box.doRequest("PUT", rawurl, nil, reqBody)

	if err == nil {
		err = json.Unmarshal(body, f)
		return err
	}
	return err

}

// Move moves the current folder under the given parent. Note that
// only Id is required apriori for both parent and current folder. The
// folder is populated with all the information after the call.
func (f *Folder) Move(box *Box, parent *Folder) error {
	if f.Id == "" || parent.Id == "" {
		return errors.New("Empty id while using Move")
	}

	fold := Folder{Parent: &Entity{Id: parent.Id}}
	reqBody, _ := json.Marshal(fold)

	rawurl := fmt.Sprintf("folders/%s", f.Id)
	body, err := box.doRequest("PUT", rawurl, nil, reqBody)

	if err == nil {
		err = json.Unmarshal(body, f)
		return err
	}
	return err

}

// Copy copies the current folder under the given parent. Note that
// only Id is required apriori for both parent and current folder. The
// copied folder is returned after copy is successful.
func (f *Folder) Copy(box *Box, parent *Folder) (*Folder, error) {
	if f.Id == "" || parent.Id == "" {
		return nil, errors.New("Empty id while using Copy")
	}

	fold := Folder{Parent: &Entity{Id: parent.Id}}
	reqBody, _ := json.Marshal(fold)

	rawurl := fmt.Sprintf("folders/%s/copy", f.Id)
	body, err := box.doRequest("POST", rawurl, nil, reqBody)

	if err == nil {
		err = json.Unmarshal(body, &fold)
		return &fold, err
	}
	return nil, err

}

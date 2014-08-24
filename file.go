package box

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type File struct {
	Id                string        `json:"id,omitempty"`                  // Box’s unique string identifying this file.
	SequenceId        string        `json:"sequence_id,omitempty"`         // A unique ID for use with the /events endpoint.
	ETag              string        `json:"etag,omitempty"`                // A unique string identifying the version of this file.
	Sha1              string        `json:"sha1,omitempty"`                // The sha1 hash of this file.
	Name              string        `json:"name,omitempty"`                // The name of this file.
	Description       string        `json:"description,omitempty"`         // The description of this file.
	Size              int           `json:"size,omitempty"`                // Size of this file in bytes.
	PathCollection    *Collection   `json:"path_collection,omitempty"`     // The path of folders to this item, starting at the root.
	CreatedAt         *BoxTime      `json:"created_at,omitempty"`          // When this file was created on Box’s servers.
	ModifiedAt        *BoxTime      `json:"modified_at,omitempty"`         // When this file was last updated on the Box servers.
	ThrashedAt        *BoxTime      `json:"thrashed_at,omitempty"`         // When this file was last moved to the trash.
	PurgedAt          *BoxTime      `json:"purged_at,omitempty"`           // When this file will be permanently deleted.
	ContentCreatedAt  *BoxTime      `json:"content_created_at,omitempty"`  // When the content of this file was created.
	ContentModifiedAt *BoxTime      `json:"content_modified_at,omitempty"` // When the content of this file was last modified.
	CreatedBy         *Entity       `json:"created_by,omitempty"`          // The user who first created file.
	ModifiedBy        *Entity       `json:"modified_by,omitempty"`         // The user who last updated this file.
	OwnedBy           *Entity       `json:"owned_by,omitempty"`            // The user who owns this file.
	SharedLink        *SharedObject `json:"shared_link,omitempty"`         // The shared link object for this file.
	Parent            *Entity       `json:"parent,omitempty"`              // The folder containing this file.
	ItemStatus        string        `json:"item_status,omitempty"`         // Whether this item is deleted or not.
	VersionNumber     string        `json:"version_number,omitempty"`      // The version of the file.
	CommentCount      int           `json:"comment_count,omitempty"`       // The number of comments on a file.
	Permissions       *Permission   `json:"permissions,omitempty"`         // The permissions that the current user has on this file.
	Tags              []string      `json:"tags,omitempty"`                // All tags applied to this file.
	Lock              *BoxLock      `json:"lock,omitempty"`                // The lock held on the file.
	Extension         string        `json:"extension,omitempty"`           // Indicates the suffix, when available, on the file.
}

// Get populates the fields of the file struct. Node that only Id is
// required apriori.
func (f *File) Get(box *Box) error {
	if f.Id == "" {
		return errors.New("Empty id while using Get")
	}
	rawurl := fmt.Sprintf("files/%s", f.Id)
	body, err := box.doRequest("GET", rawurl, nil, nil)

	if err == nil {
		err = json.Unmarshal(body, f)
		return err
	}
	return err
}

// Delete deletes the file. Note that only Id is required apriori.
func (f *File) Delete(box *Box) error {
	if f.Id == "" {
		return errors.New("Empty id while using Delete")
	}

	rawurl := fmt.Sprintf("files/%s", f.Id)
	_, err := box.doRequest("DELETE", rawurl, nil, nil)

	return err
}

// Rename renames the file with the new name. Note that only Id is
// required apriori. The file object is populated with all the
// information after the call.
func (f *File) Rename(box *Box, name string) error {
	if f.Id == "" {
		return errors.New("Empty id while using Rename")
	}

	file := File{Name: name}
	reqBody, _ := json.Marshal(file)

	rawurl := fmt.Sprintf("files/%s", f.Id)
	body, err := box.doRequest("PUT", rawurl, nil, reqBody)

	if err == nil {
		err = json.Unmarshal(body, f)
		return err
	}
	return err

}

// Move moves the current file under the given parent. Note that only
// Id is required apriori for both file and parent folder. The file
// is populated with all the information after the call.
func (f *File) Move(box *Box, parent *Folder) error {
	if f.Id == "" || parent.Id == "" {
		return errors.New("Empty id while using Move")
	}

	file := File{Parent: &Entity{Id: parent.Id}}
	reqBody, _ := json.Marshal(file)

	rawurl := fmt.Sprintf("files/%s", f.Id)
	body, err := box.doRequest("PUT", rawurl, nil, reqBody)

	if err == nil {
		err = json.Unmarshal(body, f)
		return err
	}
	return err

}

// Copy copies the current file under the given parent. Note that only
// Id is required apriori for both file and parent folder. The copied
// file is returned after copy is successful.
func (f *File) Copy(box *Box, parent *Folder) (*File, error) {
	if f.Id == "" || parent.Id == "" {
		return nil, errors.New("Empty id while using Copy")
	}

	file := File{Parent: &Entity{Id: parent.Id}}
	reqBody, _ := json.Marshal(file)

	rawurl := fmt.Sprintf("files/%s/copy", f.Id)
	body, err := box.doRequest("POST", rawurl, nil, reqBody)

	if err == nil {
		err = json.Unmarshal(body, &file)
		return &file, err
	}
	return nil, err

}

// UploadFile uploads the file at the given file path. It then fills
// the information of the recently uploaded file in the file object.
// Only Id is required for the parent folder.
func (f *File) UploadFile(box *Box, path string, parent *Folder) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("filename", filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	writer.WriteField("parent_id", parent.Id)
	// only throws path error which is already checked while opening the
	// file. Only problem is if the file is deleted in between.
	stat, _ := os.Stat(path)
	modTime, _ := BoxTime(stat.ModTime()).MarshalJSON()
	writer.WriteField("content_modified_at", string(modTime))

	err = writer.Close()
	if err != nil {
		return err
	}

	rawurl := fmt.Sprintf("%s/files/content", box.APIUPLOADURL)

	// Create mutlipart request
	request, err := http.NewRequest("POST", rawurl, body)
	if err != nil {
		return err
	}

	// Get response
	var response *http.Response
	if response, err = box.client().Do(request); err != nil {
		return err
	}
	defer response.Body.Close()

	// Get response body
	var respBody []byte
	if respBody, err = getResponse(response); err != nil {
		return err
	}

	// Unmarshal returned file
	err = json.Unmarshal(respBody, f)
	return err

}

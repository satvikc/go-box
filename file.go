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

	if err == NO_CONTENT {
		return nil
	}

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

// Download downloads the file. Note that only file id is required
// apriori.
func (f *File) Download(box *Box, writer io.Writer) error {
	var request *http.Request
	var response *http.Response
	var err error

	if f.Id == "" {
		return errors.New("Empty id while using Download")
	}

	rawurl := fmt.Sprintf("%s/files/%s/content", box.APIURL, f.Id)

	if request, err = http.NewRequest("GET", rawurl, nil); err != nil {
		return err
	}

	if response, err = box.client().Do(request); err != nil {
		return err
	}

	defer response.Body.Close()

	_, err = io.Copy(writer, response.Body)

	return err

}

// Download downloads the file at the given file path. File will be
// overwritten if it already exists. Note that only file id is
// required apriori.
func (f *File) DownloadFile(box *Box, path string) error {
	out, err := os.Create("foo.txt")
	defer out.Close()
	if err != nil {
		return err
	}
	return f.Download(box, out)
}

// Upload uploads the file (given by the reader) at the given file
// path. The file name on the box server is taken from the Name
// attribute of file object. After upload, it then fills the
// information of the recently uploaded file in the file object. Note
// that Id attribute is required for the parent folder.
func (f *File) Upload(box *Box, reader io.Reader, parent *Folder) error {

	// Check is f has name attribute and parent has id attribute
	if f.Name == "" {
		return errors.New("Empty name while using Upload")
	}

	if parent.Id == "" {
		return errors.New("Empty parent id while using Upload")
	}

	// Set up multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("filename", f.Name)
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, reader); err != nil {
		return err
	}

	// Write parent id
	writer.WriteField("parent_id", parent.Id)

	// API url
	rawurl := fmt.Sprintf("%s/files/content", box.APIUPLOADURL)

	// Create mutlipart request
	request, err := http.NewRequest("POST", rawurl, body)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	if err = writer.Close(); err != nil {
		return err
	}
	// Was giving error without this as it was setting wrong content-length
	request.ContentLength = -1

	// Get response
	var response *http.Response
	if response, err = box.client().Do(request); err != nil {
		return err
	}
	defer response.Body.Close()

	// Get response body
	var respBody []byte
	if respBody, err = getResponse(response); err != nil && err != CREATED {
		return err
	}

	// All because of weird box's return format of response body
	var m map[string]json.RawMessage
	err = json.Unmarshal(respBody, &m)
	if err != nil {
		return err
	}
	var fs []json.RawMessage
	err = json.Unmarshal(m["entries"], &fs)
	if err != nil {
		return err
	}
	if len(fs) != 1 {
		return errors.New("Not enough returned argument")
	}
	err = json.Unmarshal(fs[0], f)
	if err != nil {
		return err
	}
	return nil
}

// UploadFile directly uploads the file on the box server. The name is
// taken from the Name attribute of the file object (if it is empty,
// file name is chosen). Note than only parent id is required apriori
// for the parent folder.
func (f *File) UploadFile(box *Box, path string, parent *Folder) error {
	if f.Name == "" {
		f.Name = filepath.Base(path)
	}
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	return f.Upload(box, file, parent)
}

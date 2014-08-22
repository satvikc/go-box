package box

type File struct {
	Id                string       `json:"id,omitempty"`                  // Box’s unique string identifying this file.
	SequenceId        string       `json:"sequence_id,omitempty"`         // A unique ID for use with the /events endpoint.
	ETag              string       `json:"etag,omitempty"`                // A unique string identifying the version of this file.
	Sha1              string       `json:"sha1,omitempty"`                // The sha1 hash of this file.
	Name              string       `json:"name,omitempty"`                // The name of this file.
	Description       string       `json:"description,omitempty"`         // The description of this file.
	Size              int          `json:"size,omitempty"`                // Size of this file in bytes.
	PathCollection    Collection   `json:"path_collection,omitempty"`     // The path of folders to this item, starting at the root.
	CreatedAt         BoxTime      `json:"created_at,omitempty"`          // When this file was created on Box’s servers.
	ModifiedAt        BoxTime      `json:"modified_at,omitempty"`         // When this file was last updated on the Box servers.
	ThrashedAt        BoxTime      `json:"thrashed_at,omitempty"`         // When this file was last moved to the trash.
	PurgedAt          BoxTime      `json:"purged_at,omitempty"`           // When this file will be permanently deleted.
	ContentCreatedAt  BoxTime      `json:"content_created_at,omitempty"`  // When the content of this file was created.
	ContentModifiedAt BoxTime      `json:"content_modified_at,omitempty"` // When the content of this file was last modified.
	CreatedBy         Entity       `json:"created_by,omitempty"`          // The user who first created file.
	ModifiedBy        Entity       `json:"modified_by,omitempty"`         // The user who last updated this file.
	OwnedBy           Entity       `json:"owned_by,omitempty"`            // The user who owns this file.
	SharedLink        SharedObject `json:"shared_link,omitempty"`         // The shared link object for this file.
	Parent            Entity       `json:"parent,omitempty"`              // The folder containing this file.
	ItemStatus        string       `json:"item_status,omitempty"`         // Whether this item is deleted or not.
	VersionNumber     string       `json:"version_number,omitempty"`      // The version of the file.
	CommentCount      int          `json:"comment_count,omitempty"`       // The number of comments on a file.
	Permissions       Permission   `json:"permissions,omitempty"`         // The permissions that the current user has on this file.
	Tags              []string     `json:"tags,omitempty"`                // All tags applied to this file.
	Lock              BoxLock      `json:"lock,omitempty"`                // The lock held on the file.
	Extension         string       `json:"extension,omitempty"`           // Indicates the suffix, when available, on the file.
}

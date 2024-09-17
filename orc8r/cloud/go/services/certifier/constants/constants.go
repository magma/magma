package constants

const (
	// CertInfoType is the type of CertInfo used in blobstore type fields.
	CertInfoType = "certificate_info"

	// UserType is the type of CertInfo used in blobstore type fields.
	UserType = "user"

	// PolicyType is the type of policy used in blobstore type fileds
	PolicyType = "policy"
)

type ResourceType string

const (
	Path      ResourceType = "path"
	NetworkID ResourceType = "network_id"
	TenantID  ResourceType = "tenant_id"
)

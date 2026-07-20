package aws

type ClientRegistry struct {
	// Registry for multiple AWS clients (EC2, S3, RDS, etc.)
}

func NewClientRegistry(session interface{}) *ClientRegistry {
	return &ClientRegistry{}
}

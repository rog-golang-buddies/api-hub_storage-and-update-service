package config

// GRPCConfig represents gRPC server settings
type GRPCConfig struct {
	Host string `default:""`
	Port string `default:"50051"`
}

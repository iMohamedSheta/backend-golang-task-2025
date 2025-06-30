package redis

// // Represent a redis full connections config
// type Config struct {
// 	Default     string                   `json:"default"`
// 	Connections map[string]NodeConfig    `json:"connections"` // Single nodes
// 	Clusters    map[string]ClusterConfig `json:"clusters"`    // Clusters
// }

// type NodeConfig struct {
// 	URL      string `json:"url"`
// 	Host     string `json:"host"`
// 	Port     int    `json:"port"`
// 	Database int    `json:"database"`
// 	Username string `json:"username"`
// 	Password string `json:"password"`

// 	IsActive bool   `json:"is_active"`
// 	PoolSize int    `json:"pool_size"`
// 	Timeout  string `json:"timeout"`

// 	ClientName   string `json:"client_name"`
// 	MaxRetries   int    `json:"max_retries"`
// 	MinIdleConns int    `json:"min_idle_conns"`
// 	MaxIdleConns int    `json:"max_idle_conns"`
// 	ConnIdleTime string `json:"conn_idle_time"`
// 	ConnLifetime string `json:"conn_lifetime"`
// 	// ReadOnly              bool   `json:"read_only"`
// 	IdentitySuffix        string `json:"identity_suffix"`
// 	Protocol              int    `json:"protocol"`
// 	ContextTimeoutEnabled bool   `json:"context_timeout_enabled"`

// 	TLS *TLSConfig `json:"tls"` // nil if not using TLS
// }

// type TLSConfig struct {
// 	Enabled    bool   `json:"enabled"`     // Enable TLS
// 	SkipVerify bool   `json:"skip_verify"` // Insecure
// 	CAFile     string `json:"ca_file"`     // Optional CA cert
// 	CertFile   string `json:"cert_file"`   // Optional client cert
// 	KeyFile    string `json:"key_file"`    // Optional client key
// 	ServerName string `json:"server_name"` // Optional SNI override
// }

// func (cfg NodeConfig) ToRedisOptions() (*redis.Options, error) {
// 	opt := &redis.Options{
// 		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
// 		Username: cfg.Username,
// 		Password: cfg.Password,
// 		DB:       cfg.Database,
// 	}

// 	// TLS
// 	if cfg.TLS != nil && cfg.TLS.Enabled {
// 		tlsCfg, err := buildTLSConfig(cfg.TLS)
// 		if err != nil {
// 			return nil, fmt.Errorf("TLS config error: %w", err)
// 		}
// 		opt.TLSConfig = tlsCfg
// 	}

// 	// Client name with optional suffix
// 	if cfg.ClientName != "" {
// 		name := cfg.ClientName
// 		if cfg.IdentitySuffix != "" {
// 			name += "-" + cfg.IdentitySuffix
// 		}
// 		opt.ClientName = name
// 	}

// 	if cfg.PoolSize > 0 {
// 		opt.PoolSize = cfg.PoolSize
// 	}
// 	if cfg.MinIdleConns > 0 {
// 		opt.MinIdleConns = cfg.MinIdleConns
// 	}
// 	if cfg.MaxIdleConns > 0 {
// 		opt.MaxIdleConns = cfg.MaxIdleConns
// 	}

// 	if cfg.MaxRetries >= 0 {
// 		opt.MaxRetries = cfg.MaxRetries
// 	}

// 	if cfg.Timeout != "" {
// 		if d, err := time.ParseDuration(cfg.Timeout); err == nil {
// 			opt.DialTimeout = d
// 			opt.ReadTimeout = d
// 			opt.WriteTimeout = d
// 		}
// 	}
// 	if cfg.ConnIdleTime != "" {
// 		if d, err := time.ParseDuration(cfg.ConnIdleTime); err == nil {
// 			opt.ConnMaxIdleTime = d
// 		}
// 	}
// 	if cfg.ConnLifetime != "" {
// 		if d, err := time.ParseDuration(cfg.ConnLifetime); err == nil {
// 			opt.ConnMaxLifetime = d
// 		}
// 	}

// 	if cfg.Protocol == 2 || cfg.Protocol == 3 {
// 		opt.Protocol = cfg.Protocol
// 	}

// 	opt.ContextTimeoutEnabled = cfg.ContextTimeoutEnabled

// 	return opt, nil
// }

// func (c NodeConfig) isNotActive() bool {
// 	return (c.URL == "" && c.Host == "") || !c.IsActive
// }

// // Represent a redis cluster node config
// type ClusterNode struct {
// 	URL      string `json:"url"`
// 	Host     string `json:"host"`
// 	Port     int    `json:"port"`
// 	IsActive bool   `json:"is_active"`
// }

// func (c ClusterNode) isNotActive() bool {
// 	return (c.URL == "" && c.Host == "") || !c.IsActive
// }

// // Represent a redis cluster config
// type ClusterConfig struct {
// 	Nodes    []ClusterNode `json:"nodes"`
// 	IsActive bool          `json:"is_active"`
// 	PoolSize int           `json:"pool_size"`
// 	Username string        `json:"username"`
// 	Password string        `json:"password"`
// 	Timeout  string        `json:"timeout"`
// }

// // Build tls config for go-redis client from managed config
// func buildTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
// 	if cfg == nil || !cfg.Enabled {
// 		return nil, nil
// 	}

// 	tlsCfg := &tls.Config{
// 		InsecureSkipVerify: cfg.SkipVerify,
// 		ServerName:         cfg.ServerName,
// 	}

// 	// Load CA if provided
// 	if cfg.CAFile != "" {
// 		caCert, err := os.ReadFile(cfg.CAFile)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to read CA file: %w", err)
// 		}

// 		// Create a pool with the CA certs
// 		caPool := x509.NewCertPool()
// 		if !caPool.AppendCertsFromPEM(caCert) {
// 			return nil, fmt.Errorf("failed to append CA certs")
// 		}
// 		tlsCfg.RootCAs = caPool
// 	}

// 	// Load client cert/key if provided
// 	if cfg.CertFile != "" && cfg.KeyFile != "" {
// 		cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to load client cert/key: %w", err)
// 		}
// 		tlsCfg.Certificates = []tls.Certificate{cert}
// 	}

// 	return tlsCfg, nil
// }

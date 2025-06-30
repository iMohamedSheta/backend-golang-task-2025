package redis

// func (rm *RedisManager) Initialize() error {
// 	var defaultLoaded bool

// 	// === Single node connections ===
// 	for name, cfg := range rm.config.Connections {
// 		if cfg.isNotActive() {
// 			log.Printf("Skipping Redis %s connection - inactive", name)
// 			continue
// 		}

// 		client, err := rm.createSingleClient(cfg)
// 		if err != nil {
// 			log.Printf("Failed to create Redis client %s: %v", name, err)
// 			continue
// 		}

// 		rm.mu.Lock()
// 		rm.nodeConnections[name] = client
// 		rm.mu.Unlock()

// 		log.Printf("Redis node '%s' connection established", name)
// 		if name == rm.config.Default {
// 			defaultLoaded = true
// 		}
// 	}

// 	// === Cluster connections ===
// 	for name, clusterCfg := range rm.config.Clusters {
// 		if !clusterCfg.IsActive {
// 			log.Printf("Skipping Redis cluster %s - inactive", name)
// 			continue
// 		}

// 		cluster, err := rm.createClusterClient(clusterCfg)
// 		if err != nil {
// 			log.Printf("Failed to create Redis cluster client %s: %v", name, err)
// 			continue
// 		}

// 		rm.mu.Lock()
// 		rm.clusterConnections[name] = cluster
// 		rm.mu.Unlock()

// 		log.Printf("Redis cluster '%s' connection established", name)
// 	}

// 	if rm.config.Default != "" && !defaultLoaded {
// 		log.Printf("Warning: Default Redis connection '%s' not found or not active", rm.config.Default)
// 	}

// 	return nil
// }

// func (rm *RedisManager) createSingleClient(cfg NodeConfig) (*redis.Client, error) {
// 	var opt *redis.Options
// 	if cfg.URL != "" {
// 		o, err := redis.ParseURL(cfg.URL)
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid URL: %w", err)
// 		}
// 		opt = o
// 	} else {
// 		opt = &redis.Options{
// 			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
// 			Username: cfg.Username,
// 			Password: cfg.Password,
// 			DB:       cfg.Database,
// 		}
// 	}

// 	if cfg.PoolSize > 0 {
// 		opt.PoolSize = cfg.PoolSize
// 	}

// 	if cfg.Timeout != "" {
// 		if d, err := time.ParseDuration(cfg.Timeout); err == nil {
// 			opt.DialTimeout = d
// 			opt.ReadTimeout = d
// 			opt.WriteTimeout = d
// 		}
// 	}

// 	client := redis.NewClient(opt)
// 	if err := client.Ping(context.Background()).Err(); err != nil {
// 		client.Close()
// 		return nil, fmt.Errorf("ping error: %w", err)
// 	}

// 	return client, nil
// }

// func (rm *RedisManager) createClusterClient(cfg ClusterConfig) (*redis.ClusterClient, error) {
// 	addrs := []string{}

// 	for _, node := range cfg.Nodes {
// 		if node.isNotActive() {
// 			continue
// 		}

// 		if node.URL != "" {
// 			opt, err := redis.ParseURL(node.URL)
// 			if err != nil {
// 				log.Printf("Invalid Redis cluster node URL (%s): %v", node.URL, err)
// 				continue
// 			}
// 			addrs = append(addrs, opt.Addr)
// 		} else {
// 			addrs = append(addrs, fmt.Sprintf("%s:%d", node.Host, node.Port))
// 		}
// 	}

// 	if len(addrs) == 0 {
// 		return nil, fmt.Errorf("no active nodes found")
// 	}

// 	opts := &redis.ClusterOptions{
// 		Addrs:    addrs,
// 		Password: cfg.Password,
// 		Username: cfg.Username,
// 	}

// 	if cfg.PoolSize > 0 {
// 		opts.PoolSize = cfg.PoolSize
// 	}

// 	if cfg.Timeout != "" {
// 		if d, err := time.ParseDuration(cfg.Timeout); err == nil {
// 			opts.DialTimeout = d
// 			opts.ReadTimeout = d
// 			opts.WriteTimeout = d
// 		}
// 	}

// 	client := redis.NewClusterClient(opts)

// 	if err := client.Ping(context.Background()).Err(); err != nil {
// 		client.Close()
// 		return nil, fmt.Errorf("failed to ping cluster: %w", err)
// 	}

// 	return client, nil
// }

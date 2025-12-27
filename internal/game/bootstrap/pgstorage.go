package bootstrap

import (

	"github.com/nineteen-night/empty-room-game/config"
	"github.com/nineteen-night/empty-room-game/internal/game/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGStorage {
    var shardConfigs []struct {
        Host, Username, Password, DBName, SSLMode string
        Port int
    }
    
    for _, shard := range cfg.GameShards {
        shardConfigs = append(shardConfigs, struct {
            Host, Username, Password, DBName, SSLMode string
            Port int
        }{
            Host:     shard.Host,
            Port:     shard.Port,
            Username: shard.Username,
            Password: shard.Password,
            DBName:   shard.DBName,
            SSLMode:  shard.SSLMode,
        })
    }
    
    storage, err := pgstorage.NewPGStorage(shardConfigs)
    if err != nil { panic(err) }
    
    return storage
}
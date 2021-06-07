package psqlconnector

import (
	"fmt"
	"time"
)

type PsqlConfigurations struct {
	Host                string        `mapstructure:"host"`
	Port                int           `mapstructure:"port"`
	Username            string        `mapstructure:"username"`
	Password            string        `mapstructure:"password"`
	DBname              string        `mapstructure:"dbname"`
	QueryTimeout        time.Duration `mapstructure:"querytimeout"`
	HealthcheckInterval time.Duration `mapstructure:"healthcheck_interval"`
}

func (pc *PsqlConfigurations) GetConfig() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		pc.Host, pc.Port, pc.Username, pc.Password, pc.DBname,
	)
}

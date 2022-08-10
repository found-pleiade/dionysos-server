package utils

import (
	"time"

	"github.com/sony/sonyflake"
)

var UUIDGenerator *sonyflake.Sonyflake = sonyflake.NewSonyflake(sonyflake.Settings{StartTime: time.Now(), MachineID: nil, CheckMachineID: nil})

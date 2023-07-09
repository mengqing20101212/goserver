package common

import "time"

func GetYYYY_MM_DD_HH_mm_ss() string {
	return time.Now().Format("2006_01_01_15_04_05")
}

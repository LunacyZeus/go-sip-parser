package csv_utils

import (
	"fmt"
	"strconv"
	"time"
)

func convertTimeStamp(timestampMilli int64) string {
	// 毫秒级 Unix 时间戳字符串（例如："1633024800000"）
	//timestampMilliStr := "1633024800000"
	timestampMilliStr := fmt.Sprintf("%d", timestampMilli)
	// 将字符串转换为 int64 类型
	timestampMilli, err := strconv.ParseInt(timestampMilliStr, 10, 64)
	if err != nil {
		//fmt.Println("Error converting timestamp:", err)
		return timestampMilliStr
	}

	// 将毫秒级时间戳转换为秒级时间戳
	timestamp := timestampMilli / 1000

	// 将 Unix 时间戳转换为 time.Time 类型，默认是 UTC 时间
	t := time.Unix(timestamp, 0).UTC()

	// 格式化输出为 "年-月-日 时:分:秒" 格式
	formattedTime := t.Format("2006-01-02 15:04:05")

	// 输出转换后的时间
	//fmt.Println(formattedTime)
	return formattedTime
}

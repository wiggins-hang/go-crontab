package common

const (
	// 任务保存目录
	JobSaveDir = "/cron/jobs/"

	// 任务强杀目录
	JobKillerDir = "/cron/killer/"

	// 任务锁目录
	JobLockDir = "/cron/lock/"

	// 服务注册目录
	JobWorkersDir = "/cron/workers/"

	// 保存任务事件
	JobEventType = 1

	// 删除任务事件
	JobEventDelete = 2

	// 强杀任务事件
	JobEventKill = 3
)

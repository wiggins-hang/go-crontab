package db_model

type JobLog struct {
	Id           int64  `gorm:"column:id;type:int(11);primary_key" json:"id"`
	JobName      string `gorm:"column:job_name;type:varchar(255);NOT NULL" json:"jobName"`
	Command      string `gorm:"column:command;type:varchar(1000);comment:脚本命令;NOT NULL" json:"command"`
	Err          string `gorm:"column:err;type:varchar(5000);comment:错误信息;NOT NULL" json:"err"`
	Output       string `gorm:"column:output;type:varchar(5000);comment:输出;NOT NULL" json:"output"`
	PlanTime     int64  `gorm:"column:plan_time;type:bigint(20);default:0;comment:计划开始时间;NOT NULL" json:"planTime"`
	ScheduleTime int64  `gorm:"column:schedule_time;type:bigint(20);default:0;comment:实际调度时间;NOT NULL" json:"scheduleTime"`
	StartTime    int64  `gorm:"column:start_time;type:bigint(20);comment:任务执行开始时间;NOT NULL" json:"startTime"`
	EndTime      int64  `gorm:"column:end_time;type:bigint(20);default:0;comment:任务执行结束时间;NOT NULL" json:"endTime"`
}

func (m *JobLog) TableName() string {
	return "job_log"
}

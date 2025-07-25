package models

import (
	"time"
)

// 数据库连接管理表
type SysDatabases struct {
	Id          uint64    `gorm:"column:id;type:bigint UNSIGNED;comment:主键ID;primaryKey;not null;" json:"id"`                       // 主键ID
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime(3);comment:创建时间;default:NULL;" json:"created_at"`                  // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime(3);comment:更新时间;default:NULL;" json:"updated_at"`                  // 更新时间
	DeletedAt   time.Time `gorm:"column:deleted_at;type:datetime(3);comment:删除时间;default:NULL;" json:"deleted_at"`                  // 删除时间
	Name        string    `gorm:"column:name;type:varchar(255);comment:数据库名称;not null;" json:"name"`                                // 数据库名称
	Type        string    `gorm:"column:type;type:varchar(50);comment:数据库类型(mysql/postgresql/redis/mongodb);not null;" json:"type"` // 数据库类型(mysql/postgresql/redis/mongodb)
	Host        string    `gorm:"column:host;type:varchar(255);comment:主机地址;not null;" json:"host"`                                 // 主机地址
	Port        int32     `gorm:"column:port;type:int;comment:端口号;not null;" json:"port"`                                           // 端口号
	Username    string    `gorm:"column:username;type:varchar(255);comment:用户名;not null;" json:"username"`                          // 用户名
	Password    string    `gorm:"column:password;type:varchar(255);comment:密码;not null;" json:"password"`                           // 密码
	Database    string    `gorm:"column:database;type:varchar(255);comment:数据库名;default:NULL;" json:"database"`                     // 数据库名
	Charset     string    `gorm:"column:charset;type:varchar(50);comment:字符集;default:utf8mb4;" json:"charset"`                      // 字符集
	Status      string    `gorm:"column:status;type:varchar(20);comment:状态(active/inactive/error);default:active;" json:"status"`   // 状态(active/inactive/error)
	Description string    `gorm:"column:description;type:text;comment:描述;" json:"description"`                                      // 描述
	Config      string    `gorm:"column:config;type:text;comment:配置信息(JSON);" json:"config"`                                        // 配置信息(JSON)
	LastConnect time.Time `gorm:"column:last_connect;type:datetime(3);comment:最后连接时间;default:NULL;" json:"last_connect"`            // 最后连接时间
	CreatedBy   uint64    `gorm:"column:created_by;type:bigint UNSIGNED;comment:创建者ID;default:NULL;" json:"created_by"`             // 创建者ID
	UpdatedBy   uint64    `gorm:"column:updated_by;type:bigint UNSIGNED;comment:更新者ID;default:NULL;" json:"updated_by"`             // 更新者ID
}

// 数据库用户表
type SysDatabaseUsers struct {
	Id         uint64    `json:"id"`          // 主键ID
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
	DeletedAt  time.Time `json:"deleted_at"`  // 删除时间
	DatabaseId uint64    `json:"database_id"` // 数据库ID
	Username   string    `json:"username"`    // 用户名
	Password   string    `json:"password"`    // 密码
	Host       string    `json:"host"`        // 允许的主机
	Privileges string    `json:"privileges"`  // 权限
	Status     string    `json:"status"`      // 状态
	CreatedBy  uint64    `json:"created_by"`  // 创建者ID
}

// 数据库表信息
type SysDatabaseTables struct {
	Id          uint64    `gorm:"column:id;type:bigint UNSIGNED;comment:主键ID;primaryKey;not null;" json:"id"`         // 主键ID
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime(3);comment:创建时间;default:NULL;" json:"created_at"`    // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime(3);comment:更新时间;default:NULL;" json:"updated_at"`    // 更新时间
	DeletedAt   time.Time `gorm:"column:deleted_at;type:datetime(3);comment:删除时间;default:NULL;" json:"deleted_at"`    // 删除时间
	DatabaseId  uint64    `gorm:"column:database_id;type:bigint UNSIGNED;comment:数据库ID;not null;" json:"database_id"` // 数据库ID
	TableName   string    `gorm:"column:table_name;type:varchar(255);comment:表名;not null;" json:"table_name"`         // 表名
	Engine      string    `gorm:"column:engine;type:varchar(50);comment:存储引擎;default:InnoDB;" json:"engine"`          // 存储引擎
	Rows        int64     `gorm:"column:rows;type:bigint;comment:行数;default:0;" json:"rows"`                          // 行数
	DataLength  int64     `gorm:"column:data_length;type:bigint;comment:数据长度;default:0;" json:"data_length"`          // 数据长度
	IndexLength int64     `gorm:"column:index_length;type:bigint;comment:索引长度;default:0;" json:"index_length"`        // 索引长度
	Comment     string    `gorm:"column:comment;type:text;comment:表注释;" json:"comment"`                               // 表注释
	Status      string    `gorm:"column:status;type:varchar(20);comment:状态;default:active;" json:"status"`            // 状态
}

// 数据库查询历史表
type SysDatabaseQueries struct {
	Id            uint64    `gorm:"column:id;type:bigint UNSIGNED;comment:主键ID;primaryKey;not null;" json:"id"`               // 主键ID
	CreatedAt     time.Time `gorm:"column:created_at;type:datetime(3);comment:创建时间;default:NULL;" json:"created_at"`          // 创建时间
	UpdatedAt     time.Time `gorm:"column:updated_at;type:datetime(3);comment:更新时间;default:NULL;" json:"updated_at"`          // 更新时间
	DeletedAt     time.Time `gorm:"column:deleted_at;type:datetime(3);comment:删除时间;default:NULL;" json:"deleted_at"`          // 删除时间
	DatabaseId    uint64    `gorm:"column:database_id;type:bigint UNSIGNED;comment:数据库ID;not null;" json:"database_id"`       // 数据库ID
	Sql           string    `gorm:"column:sql;type:text;comment:SQL语句;not null;" json:"sql"`                                  // SQL语句
	ExecutionTime int32     `gorm:"column:execution_time;type:int;comment:执行时间(毫秒);default:0;" json:"execution_time"`         // 执行时间(毫秒)
	RowsAffected  int32     `gorm:"column:rows_affected;type:int;comment:影响行数;default:0;" json:"rows_affected"`               // 影响行数
	Status        string    `gorm:"column:status;type:varchar(20);comment:状态(success/failed);default:success;" json:"status"` // 状态(success/failed)
	ErrorMessage  string    `gorm:"column:error_message;type:text;comment:错误信息;" json:"error_message"`                        // 错误信息
	CreatedBy     uint64    `gorm:"column:created_by;type:bigint UNSIGNED;comment:执行者ID;default:NULL;" json:"created_by"`     // 执行者ID
}

// 数据库监控日志表
type SysDatabaseMonitors struct {
	Id             uint64    `gorm:"column:id;type:bigint UNSIGNED;comment:主键ID;primaryKey;not null;" json:"id"`             // 主键ID
	CreatedAt      time.Time `gorm:"column:created_at;type:datetime(3);comment:创建时间;default:NULL;" json:"created_at"`        // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at;type:datetime(3);comment:更新时间;default:NULL;" json:"updated_at"`        // 更新时间
	DeletedAt      time.Time `gorm:"column:deleted_at;type:datetime(3);comment:删除时间;default:NULL;" json:"deleted_at"`        // 删除时间
	DatabaseId     uint64    `gorm:"column:database_id;type:bigint UNSIGNED;comment:数据库ID;not null;" json:"database_id"`     // 数据库ID
	Connections    int32     `gorm:"column:connections;type:int;comment:当前连接数;default:0;" json:"connections"`                // 当前连接数
	MaxConnections int32     `gorm:"column:max_connections;type:int;comment:最大连接数;default:0;" json:"max_connections"`        // 最大连接数
	Queries        int64     `gorm:"column:queries;type:bigint;comment:查询次数;default:0;" json:"queries"`                      // 查询次数
	SlowQueries    int64     `gorm:"column:slow_queries;type:bigint;comment:慢查询次数;default:0;" json:"slow_queries"`           // 慢查询次数
	Uptime         int64     `gorm:"column:uptime;type:bigint;comment:运行时间(秒);default:0;" json:"uptime"`                     // 运行时间(秒)
	BytesReceived  int64     `gorm:"column:bytes_received;type:bigint;comment:接收字节数;default:0;" json:"bytes_received"`       // 接收字节数
	BytesSent      int64     `gorm:"column:bytes_sent;type:bigint;comment:发送字节数;default:0;" json:"bytes_sent"`               // 发送字节数
	CpuUsage       float64   `gorm:"column:cpu_usage;type:decimal(5, 2);comment:CPU使用率;default:0.00;" json:"cpu_usage"`      // CPU使用率
	MemoryUsage    float64   `gorm:"column:memory_usage;type:decimal(5, 2);comment:内存使用率;default:0.00;" json:"memory_usage"` // 内存使用率
	DiskUsage      float64   `gorm:"column:disk_usage;type:decimal(5, 2);comment:磁盘使用率;default:0.00;" json:"disk_usage"`     // 磁盘使用率
}

// 数据库备份表
type SysDatabaseBackups struct {
	Id         uint64    `gorm:"column:id;type:bigint UNSIGNED;comment:主键ID;primaryKey;not null;" json:"id"`                       // 主键ID
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime(3);comment:创建时间;default:NULL;" json:"created_at"`                  // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;type:datetime(3);comment:更新时间;default:NULL;" json:"updated_at"`                  // 更新时间
	DeletedAt  time.Time `gorm:"column:deleted_at;type:datetime(3);comment:删除时间;default:NULL;" json:"deleted_at"`                  // 删除时间
	DatabaseId uint64    `gorm:"column:database_id;type:bigint UNSIGNED;comment:数据库ID;not null;" json:"database_id"`               // 数据库ID
	Name       string    `gorm:"column:name;type:varchar(255);comment:备份名称;not null;" json:"name"`                                 // 备份名称
	Path       string    `gorm:"column:path;type:varchar(500);comment:备份文件路径;not null;" json:"path"`                               // 备份文件路径
	Size       int64     `gorm:"column:size;type:bigint;comment:文件大小(字节);default:0;" json:"size"`                                  // 文件大小(字节)
	Type       string    `gorm:"column:type;type:varchar(20);comment:备份类型(full/incremental);default:full;" json:"type"`            // 备份类型(full/incremental)
	Status     string    `gorm:"column:status;type:varchar(20);comment:状态(success/failed/running);default:success;" json:"status"` // 状态(success/failed/running)
	Message    string    `gorm:"column:message;type:text;comment:备份信息;" json:"message"`                                            // 备份信息
	BackupTime time.Time `gorm:"column:backup_time;type:datetime(3);comment:备份时间;default:NULL;" json:"backup_time"`                // 备份时间
	CreatedBy  uint64    `gorm:"column:created_by;type:bigint UNSIGNED;comment:创建者ID;default:NULL;" json:"created_by"`             // 创建者ID
}

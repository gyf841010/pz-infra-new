package database

import (
	"context"
	"errors"
	"time"

	"github.com/gyf841010/pz-infra-new/log"
	"github.com/gyf841010/pz-infra-new/logging"
	"github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"gorm.io/plugin/dbresolver"
)

var (
	globalDB                 *gorm.DB
	ErrRecordNotFound        = gorm.ErrRecordNotFound
	ErrInvalidTransaction    = gorm.ErrInvalidTransaction
	ErrNotImplemented        = gorm.ErrNotImplemented
	ErrMissingWhereClause    = gorm.ErrMissingWhereClause
	ErrUnsupportedRelation   = gorm.ErrUnsupportedRelation
	ErrPrimaryKeyRequired    = gorm.ErrPrimaryKeyRequired
	ErrModelValueRequired    = gorm.ErrModelValueRequired
	ErrInvalidData           = gorm.ErrInvalidData
	ErrUnsupportedDriver     = gorm.ErrUnsupportedDriver
	ErrRegistered            = gorm.ErrRegistered
	ErrInvalidField          = gorm.ErrInvalidField
	ErrEmptySlice            = gorm.ErrEmptySlice
	ErrDryRunModeUnsupported = gorm.ErrDryRunModeUnsupported
)

// gorm V2 日志输出需要实现logger.Interface接口
type GormLogger struct {
	logger        logging.Logger // logrus封装,日志记录
	SlowThreshold time.Duration  // 慢查询阈值,用于慢查询日志记录
	SourceField   string         //
}

func (l *GormLogger) LogMode(logger.LogLevel) logger.Interface {
	newlogger := *l
	return &newlogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Infof(msg, data...)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Warnf(msg, data...)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Errorf(msg, data...)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}

	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	// 错误日志
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound)) {
		fields[logrus.ErrorKey] = err
		l.logger.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	// 慢日志
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	// 方便定位问题
	// gorm 执行sql,打印 修改为info级别,原为debug
	// 关闭线上日志 2022=03-03
	l.logger.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}

/** 初始化数据库连接,支持传入从库地址,实现读写分离
 * @description:
 * @param {string} connectString 主库地址
 * @param {logging.Logger} infraLogger log组件实现
 * @param {...string} slaveDSN 从库地址
 * @return {*}
 */
func InitDB(connectString string, infraLogger logging.Logger, slaveDSN ...string) error {
	var gormlog logger.Interface
	if infraLogger != nil {
		gormlog = &GormLogger{
			SlowThreshold: 200 * time.Millisecond,
			logger:        infraLogger,
		}
	} else {
		gormlog = logger.Default
	}

	gormDB, err := gorm.Open(mysql.Open(connectString), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormlog,
	})
	if err != nil {
		log.Errorf("init db error with url %s failed: %s", connectString, err.Error())
		panic(err)
	}

	if len(slaveDSN) > 0 {
		dbs := []gorm.Dialector{}
		for _, item := range slaveDSN {
			dbs = append(dbs, mysql.Open(item))
		}

		err = gormDB.Use(dbresolver.Register(
			dbresolver.Config{
				//Sources:  []gorm.Dialector{mysql.Open(...)}, // 主库
				Replicas: dbs, //从库
				// sources/replicas load balancing policy //负载均衡策略
				Policy: dbresolver.RandomPolicy{},
			},
		))
		if err != nil {
			return err
		}
	}

	gormDB.Logger.LogMode(logger.Info)
	SetDB(gormDB)
	return err
}

func LogMode(level logger.LogLevel) {
	if globalDB != nil {
		globalDB.Logger.LogMode(level)
	}
}

// 通过 NewDB 选项创建一个不带之前条件的新 DB,不受上一条执行语句携带的Statement影响
// 参考: https://gorm.io/zh_CN/docs/session.html#NewDB
func GetDB() *gorm.DB {
	return globalDB.Session(&gorm.Session{
		NewDB: true,
	})
}

// 通过指定session获取不带之前条件的新 DB,在session中可指定一些条件,
// 参考: https://gorm.io/zh_CN/docs/session.html#NewDB
func GetDBWithSession(session *gorm.Session) *gorm.DB {
	return globalDB.Session(session)
}

func SetDB(db *gorm.DB) {
	globalDB = db
}

// 显示指定,在从库查询数据
func Read(dbs ...*gorm.DB) *gorm.DB {
	if len(dbs) > 0 {
		return GetNonTransactionDatabases(dbs).Clauses(dbresolver.Read)
	}
	return GetDB().Clauses(dbresolver.Read)
}

// 显示指定,在主库读写数据
func Write(dbs ...*gorm.DB) *gorm.DB {
	if len(dbs) > 0 {
		return GetNonTransactionDatabases(dbs).Clauses(dbresolver.Write)
	}
	return GetDB().Clauses(dbresolver.Write)
}

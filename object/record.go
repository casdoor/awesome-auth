package object

import (
	"strings"

	"github.com/beego/beego/context"
	"github.com/casbin/casvisor/conf"
	"github.com/casbin/casvisor/util"
)

var logPostOnly bool

func init() {
	logPostOnly = conf.GetConfigBool("logPostOnly")
}

type Record struct {
	Id int `xorm:"int notnull pk autoincr" json:"id"`

	Owner       string `xorm:"varchar(100) index" json:"owner"`
	Name        string `xorm:"varchar(100) index" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Organization string `xorm:"varchar(100)" json:"organization"`
	ClientIp     string `xorm:"varchar(100)" json:"clientIp"`
	User         string `xorm:"varchar(100)" json:"user"`
	Method       string `xorm:"varchar(100)" json:"method"`
	RequestUri   string `xorm:"varchar(1000)" json:"requestUri"`
	Action       string `xorm:"varchar(1000)" json:"action"`

	Object string `xorm:"-" json:"object"`
	// ExtendedUser *User  `xorm:"-" json:"extendedUser"`

	IsTriggered bool `json:"isTriggered"`
}

func NewRecord(ctx *context.Context) *Record {
	ip := strings.Replace(util.GetIPFromRequest(ctx.Request), ": ", "", -1)
	action := strings.Replace(ctx.Request.URL.Path, "/api/", "", -1)
	requestUri := util.FilterQuery(ctx.Request.RequestURI, []string{"accessToken"})
	if len(requestUri) > 1000 {
		requestUri = requestUri[0:1000]
	}

	object := ""
	if ctx.Input.RequestBody != nil && len(ctx.Input.RequestBody) != 0 {
		object = string(ctx.Input.RequestBody)
	}

	record := Record{
		Name:        util.GenerateId(),
		CreatedTime: util.GetCurrentTime(),
		ClientIp:    ip,
		User:        "",
		Method:      ctx.Request.Method,
		RequestUri:  requestUri,
		Action:      action,
		Object:      object,
		IsTriggered: false,
	}
	return &record
}

func AddRecord(record *Record) bool {
	if logPostOnly {
		if record.Method == "GET" {
			return false
		}
	}

	if record.Organization == "app" {
		return false
	}

	record.Owner = record.Organization

	affected, err := adapter.engine.Insert(record)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetRecordCount(field, value string, filterRecord *Record) (int64, error) {
	session := GetSession("", -1, -1, field, value, "", "")
	return session.Count(filterRecord)
}

func GetRecords() ([]*Record, error) {
	records := []*Record{}
	err := adapter.engine.Desc("id").Find(&records)
	if err != nil {
		return records, err
	}

	return records, nil
}

func GetPaginationRecords(offset, limit int, field, value, sortField, sortOrder string, filterRecord *Record) ([]*Record, error) {
	records := []*Record{}
	session := GetSession("", offset, limit, field, value, sortField, sortOrder)
	err := session.Find(&records, filterRecord)
	if err != nil {
		return records, err
	}

	return records, nil
}

func GetRecordsByField(record *Record) ([]*Record, error) {
	records := []*Record{}
	err := adapter.engine.Find(&records, record)
	if err != nil {
		return records, err
	}

	return records, nil
}

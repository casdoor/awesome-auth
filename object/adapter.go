// Copyright 2023 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/beego/beego"
	"github.com/hanzoai/vm/conf"
	"github.com/hanzoai/vm/util"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

var adapter *Adapter

func InitConfig() {
	err := beego.LoadAppConfig("ini", "../conf/app.conf")
	if err != nil {
		panic(err)
	}

	InitAdapter()
}

func InitAdapter() {
	adapter = NewAdapter(conf.GetConfigString("driverName"), strings.TrimSpace(conf.GetConfigDataSourceName()))
}

// Adapter represents the MySQL adapter for policy storage.
type Adapter struct {
	driverName     string
	dataSourceName string
	engine         *xorm.Engine
}

// finalizer is the destructor for Adapter.
func finalizer(a *Adapter) {
	err := a.engine.Close()
	if err != nil {
		panic(err)
	}
}

// NewAdapter is the constructor for Adapter.
func NewAdapter(driverName string, dataSourceName string) *Adapter {
	a := &Adapter{}
	a.driverName = driverName
	a.dataSourceName = dataSourceName

	// Open the DB, create it if not existed.
	a.open()

	// Call the destructor when the object is released.
	runtime.SetFinalizer(a, finalizer)

	return a
}

func (a *Adapter) createDatabase() error {
	dbName := beego.AppConfig.String("dbName")

	if a.driverName == "postgres" {
		// For postgres, connect to the default "postgres" database to create the target DB.
		connStr := a.dataSourceName
		if strings.Contains(connStr, "dbname=") {
			// Replace existing dbname with "postgres" for the admin connection.
			connStr = strings.Replace(connStr, fmt.Sprintf("dbname=%s", dbName), "dbname=postgres", 1)
		} else {
			connStr += " dbname=postgres"
		}
		engine, err := xorm.NewEngine(a.driverName, connStr)
		if err != nil {
			return err
		}
		defer engine.Close()

		// Check if database exists before creating.
		var count int64
		_, err = engine.SQL("SELECT count(*) FROM pg_database WHERE datname = ?", dbName).Get(&count)
		if err != nil {
			return err
		}
		if count == 0 {
			_, err = engine.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
			return err
		}
		return nil
	}

	// MySQL path.
	engine, err := xorm.NewEngine(a.driverName, a.dataSourceName)
	if err != nil {
		return err
	}
	defer engine.Close()

	_, err = engine.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8 COLLATE utf8_general_ci", dbName))
	return err
}

func (a *Adapter) open() {
	if err := a.createDatabase(); err != nil {
		panic(err)
	}

	var dsn string
	dbName := beego.AppConfig.String("dbName")
	if a.driverName == "postgres" {
		// For postgres, set dbname in the connection string.
		dsn = a.dataSourceName
		if strings.Contains(dsn, "dbname=") {
			dsn = strings.Replace(dsn, "dbname=postgres", fmt.Sprintf("dbname=%s", dbName), 1)
		} else {
			dsn += fmt.Sprintf(" dbname=%s", dbName)
		}
	} else {
		// MySQL: append dbName to DSN.
		dsn = a.dataSourceName + dbName
	}

	engine, err := xorm.NewEngine(a.driverName, dsn)
	if err != nil {
		panic(err)
	}

	a.engine = engine
	a.createTable()
}

func (a *Adapter) close() {
	a.engine.Close()
	a.engine = nil
}

func (a *Adapter) createTable() {
	err := a.engine.Sync2(new(Asset))
	if err != nil {
		panic(err)
	}

	err = a.engine.Sync2(new(Provider))
	if err != nil {
		panic(err)
	}

	err = a.engine.Sync2(new(Machine))
	if err != nil {
		panic(err)
	}

	err = a.engine.Sync2(new(Record))
	if err != nil {
		panic(err)
	}

	err = a.engine.Sync2(new(Session))
	if err != nil {
		panic(err)
	}
}

func GetSession(owner string, offset, limit int, field, value, sortField, sortOrder string) *xorm.Session {
	session := adapter.engine.Prepare()
	if offset != -1 && limit != -1 {
		session.Limit(limit, offset)
	}
	if owner != "" {
		session = session.And("owner=?", owner)
	}
	if field != "" && value != "" {
		if util.FilterField(field) {
			session = session.And(fmt.Sprintf("%s like ?", util.SnakeString(field)), fmt.Sprintf("%%%s%%", value))
		}
	}
	if sortField == "" || sortOrder == "" {
		sortField = "created_time"
	}
	if sortOrder == "ascend" {
		session = session.Asc(util.SnakeString(sortField))
	} else {
		session = session.Desc(util.SnakeString(sortField))
	}
	return session
}

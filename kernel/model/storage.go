// SiYuan - Build Your Eternal Digital Garden
// Copyright (c) 2020-present, b3log.org
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package model

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/88250/gulu"
	"github.com/siyuan-note/filelock"
	"github.com/siyuan-note/logging"
	"github.com/siyuan-note/siyuan/kernel/util"
)

type Criterion struct {
	Name        string          `json:"name"`
	Sort        int             `json:"sort"`       //  0：按块类型（默认），1：按创建时间升序，2：按创建时间降序，3：按更新时间升序，4：按更新时间降序，5：按内容顺序（仅在按文档分组时）
	Group       int             `json:"group"`      // 0：不分组，1：按文档分组
	Layout      int             `json:"layout"`     // 0：上下，1：左右
	HasReplace  bool            `json:"hasReplace"` // 是否有替换
	Method      int             `json:"method"`     //  0：文本，1：查询语法，2：SQL，3：正则表达式
	HPath       string          `json:"hPath"`
	IDPath      []string        `json:"idPath"`
	K           string          `json:"k"`           // 搜索关键字
	R           string          `json:"r"`           // 替换关键字
	ReplaceList []string        `json:"replaceList"` // 替换候选列表
	List        []string        `json:"list"`        // 搜索候选列表
	Types       *CriterionTypes `json:"types"`       // 类型过滤选项
}

type CriterionTypes struct {
	MathBlock  bool `json:"mathBlock"`
	Table      bool `json:"table"`
	Blockquote bool `json:"blockquote"`
	SuperBlock bool `json:"superBlock"`
	Paragraph  bool `json:"paragraph"`
	Document   bool `json:"document"`
	Heading    bool `json:"heading"`
	List       bool `json:"list"`
	ListItem   bool `json:"listItem"`
	CodeBlock  bool `json:"codeBlock"`
	HtmlBlock  bool `json:"htmlBlock"`
}

var criteriaLock = sync.Mutex{}

func RemoveCriterion(name string) (err error) {
	criteriaLock.Lock()
	defer criteriaLock.Unlock()

	criteria, err := getCriteria()
	if nil != err {
		return
	}

	for i, c := range criteria {
		if c.Name == name {
			criteria = append(criteria[:i], criteria[i+1:]...)
			break
		}
	}

	err = setCriteria(criteria)
	return
}

func SetCriterion(criterion *Criterion) (err error) {
	criteriaLock.Lock()
	defer criteriaLock.Unlock()

	criteria, err := getCriteria()
	if nil != err {
		return
	}

	update := false
	for i, c := range criteria {
		if c.Name == criterion.Name {
			criteria[i] = criterion
			update = true
			break
		}
	}
	if !update {
		criteria = append(criteria, criterion)
	}

	err = setCriteria(criteria)
	return
}

func GetCriteria() (ret []*Criterion, err error) {
	criteriaLock.Lock()
	defer criteriaLock.Unlock()
	return getCriteria()
}

func setCriteria(criteria []*Criterion) (err error) {
	dirPath := filepath.Join(util.DataDir, "storage")
	if err = os.MkdirAll(dirPath, 0755); nil != err {
		logging.LogErrorf("create storage [criteria] dir failed: %s", err)
		return
	}

	data, err := gulu.JSON.MarshalIndentJSON(criteria, "", "  ")
	if nil != err {
		logging.LogErrorf("marshal storage [criteria] failed: %s", err)
		return
	}

	lsPath := filepath.Join(dirPath, "criteria.json")
	err = filelock.WriteFile(lsPath, data)
	if nil != err {
		logging.LogErrorf("write storage [criteria] failed: %s", err)
		return
	}
	return
}

func getCriteria() (ret []*Criterion, err error) {
	ret = []*Criterion{}
	dataPath := filepath.Join(util.DataDir, "storage/criteria.json")
	if !gulu.File.IsExist(dataPath) {
		return
	}

	data, err := filelock.ReadFile(dataPath)
	if nil != err {
		logging.LogErrorf("read storage [criteria] failed: %s", err)
		return
	}

	if err = gulu.JSON.UnmarshalJSON(data, &ret); nil != err {
		logging.LogErrorf("unmarshal storage [criteria] failed: %s", err)
		return
	}
	return
}

var localStorageLock = sync.Mutex{}

func RemoveLocalStorageVal(key string) (err error) {
	localStorageLock.Lock()
	defer localStorageLock.Unlock()

	localStorage, err := getLocalStorage()
	if nil != err {
		return
	}

	delete(localStorage, key)
	return setLocalStorage(localStorage)
}

func SetLocalStorageVal(key string, val interface{}) (err error) {
	localStorageLock.Lock()
	defer localStorageLock.Unlock()

	localStorage, err := getLocalStorage()
	if nil != err {
		return
	}

	localStorage[key] = val
	return setLocalStorage(localStorage)
}

func SetLocalStorage(val interface{}) (err error) {
	localStorageLock.Lock()
	defer localStorageLock.Unlock()
	return setLocalStorage(val)
}

func GetLocalStorage() (ret map[string]interface{}, err error) {
	localStorageLock.Lock()
	defer localStorageLock.Unlock()
	return getLocalStorage()
}

func setLocalStorage(val interface{}) (err error) {
	dirPath := filepath.Join(util.DataDir, "storage")
	if err = os.MkdirAll(dirPath, 0755); nil != err {
		logging.LogErrorf("create storage [local] dir failed: %s", err)
		return
	}

	data, err := gulu.JSON.MarshalIndentJSON(val, "", "  ")
	if nil != err {
		logging.LogErrorf("marshal storage [local] failed: %s", err)
		return
	}

	lsPath := filepath.Join(dirPath, "local.json")
	err = filelock.WriteFile(lsPath, data)
	if nil != err {
		logging.LogErrorf("write storage [local] failed: %s", err)
		return
	}
	return
}

func getLocalStorage() (ret map[string]interface{}, err error) {
	lsPath := filepath.Join(util.DataDir, "storage/local.json")
	if !gulu.File.IsExist(lsPath) {
		return
	}

	data, err := filelock.ReadFile(lsPath)
	if nil != err {
		logging.LogErrorf("read storage [local] failed: %s", err)
		return
	}

	ret = map[string]interface{}{}
	if err = gulu.JSON.UnmarshalJSON(data, &ret); nil != err {
		logging.LogErrorf("unmarshal storage [local] failed: %s", err)
		return
	}
	return
}

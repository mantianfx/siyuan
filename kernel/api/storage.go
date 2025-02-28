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

package api

import (
	"net/http"

	"github.com/88250/gulu"
	"github.com/gin-gonic/gin"
	"github.com/siyuan-note/siyuan/kernel/model"
	"github.com/siyuan-note/siyuan/kernel/util"
)

func removeCriterion(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}

	name := arg["name"].(string)
	err := model.RemoveCriterion(name)
	if nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

func setCriterion(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}

	param, err := gulu.JSON.MarshalJSON(arg["criterion"])
	if nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	criterion := &model.Criterion{}
	if err = gulu.JSON.UnmarshalJSON(param, criterion); nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}

	err = model.SetCriterion(criterion)
	if nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

func getCriteria(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	data, err := model.GetCriteria()
	if nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
	ret.Data = data
}

func removeLocalStorageVal(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}

	key := arg["key"].(string)

	err := model.RemoveLocalStorageVal(key)
	if nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

func setLocalStorageVal(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}

	key := arg["key"].(string)
	val := arg["val"].(interface{})

	err := model.SetLocalStorageVal(key, val)
	if nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

func getLocalStorage(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	data, err := model.GetLocalStorage()
	if nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
	ret.Data = data
}

func setLocalStorage(c *gin.Context) {
	ret := gulu.Ret.NewResult()
	defer c.JSON(http.StatusOK, ret)

	arg, ok := util.JsonArg(c, ret)
	if !ok {
		return
	}

	val := arg["val"].(interface{})

	err := model.SetLocalStorage(val)
	if nil != err {
		ret.Code = -1
		ret.Msg = err.Error()
		return
	}
}

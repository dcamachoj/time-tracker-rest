package rest

import (
	"context"
	"dcamachoj/time-tracker-rest/common"
	"dcamachoj/time-tracker-rest/dbx"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func cfgEmployee(r *gin.RouterGroup) {
	r.GET("/", ToRest(lstEmployee))
	r.GET("/:id", ToRest(getEmployee))
	r.POST("/", ToRest(addEmployee))
	r.PUT("/:id", ToRest(updEmployee))
	r.DELETE("/:id", ToRest(delEmployee))
	cfgTimeData(r.Group("/:id/timeData"))
}

func lstEmployee(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var req common.RequestPage
	var res common.ResponsePage
	var lst []*dbx.Employee
	var err error
	err = c.ShouldBindQuery(&req)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}
	if req.Size <= 0 {
		req.Size = common.DefaultPage
	}
	res.Assign(&req)
	err = employeeEntity.Page(ctx, &lst, &res)
	if err != nil {
		return common.WrapResponse(ctx, err)
	}

	if lst == nil {
		lst = []*dbx.Employee{}
	}
	return common.ResponseOK(ctx).
		SetData(common.ResponseMap{
			"list": lst,
			"page": res,
		})
}

func getEmployee(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var req common.RequestId
	var res *dbx.Employee
	var err error

	err = c.ShouldBindUri(&req)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = dbx.ExecTx(ctx, func(ctx context.Context) error {
		err = employeeEntity.GetSingle(ctx, req.ID, &res)
		if err != nil {
			return err
		}

		var timeData = &(res.TimeData)
		var timeDataPage = &common.ResponsePage{}
		err = timeDataEntity.Page(ctx, res.ID, timeData, timeDataPage)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return common.WrapResponse(ctx, err)
	}

	return common.ResponseOK(ctx).
		SetData(res)
}

func addEmployee(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var req *dbx.Employee
	var err error
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = dbx.ExecTx(ctx, func(ctx context.Context) error {
		err = employeeEntity.Insert(ctx, req)
		if err != nil {
			return err
		}

		for _, timeData := range req.TimeData {
			timeData.ParentID = req.ID
			err = timeDataEntity.Insert(ctx, timeData)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return common.WrapResponse(ctx, err)
	}

	return common.ResponseOK(ctx).
		SetData(req)
}

func updEmployee(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var rid common.RequestId
	var req *dbx.Employee
	var res *dbx.Employee
	var err error

	err = c.ShouldBindUri(&rid)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	if req.ID > 0 && rid.ID != req.ID {
		err = errors.Errorf("IDs don't match")
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = dbx.ExecTx(ctx, func(ctx context.Context) error {
		err = employeeEntity.Update(ctx, req)
		if err != nil {
			return err
		}

		err = employeeEntity.GetSingle(ctx, req.ID, &res)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return common.WrapResponse(ctx, err)
	}
	return common.ResponseOK(ctx).
		SetData(res)
}

func delEmployee(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var rid common.RequestId
	var err error

	err = c.ShouldBindUri(&rid)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = dbx.ExecTx(ctx, func(ctx context.Context) error {
		err = employeeEntity.Delete(ctx, rid.ID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return common.WrapResponse(ctx, err)
	}
	return common.ResponseOK(ctx)
}

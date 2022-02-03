package rest

import (
	"context"
	"dcamachoj/time-tracker-rest/common"
	"dcamachoj/time-tracker-rest/dbx"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func cfgTimeData(r *gin.RouterGroup) {
	r.GET("/", ToRest(lstTimeData))
	r.GET("/:idChild", ToRest(getTimeData))
	r.POST("/", ToRest(addTimeData))
	r.PUT("/:idChild", ToRest(updTimeData))
	r.DELETE("/:idChild", ToRest(delTimeData))
}

func lstTimeData(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var rid common.RequestId
	var req common.RequestPage
	var res common.ResponsePage
	var lst []*dbx.TimeData
	var err error

	err = c.ShouldBindQuery(&req)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = c.ShouldBindUri(&rid)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	res.Assign(&req)
	err = timeDataEntity.Page(ctx, rid.ID, &lst, &res)
	if err != nil {
		return common.WrapResponse(ctx, err)
	}

	if lst == nil {
		lst = []*dbx.TimeData{}
	}
	return common.ResponseOK(ctx).
		SetData(common.ResponseMap{
			"list": lst,
			"page": res,
		})
}

func getTimeData(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var req common.RequestChildId
	var res *dbx.TimeData
	var err error

	err = c.ShouldBindUri(&req)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = timeDataEntity.GetSingle(ctx, req.ChildID, &res)
	if err != nil {
		return common.WrapResponse(ctx, err)
	}

	return common.ResponseOK(ctx).
		SetData(res)
}
func addTimeData(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var rid common.RequestId
	var req *dbx.TimeData
	var err error

	err = c.ShouldBindUri(&rid)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	req.ParentID = rid.ID

	err = dbx.ExecTx(ctx, func(ctx context.Context) error {
		err = timeDataEntity.Insert(ctx, req)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return common.WrapResponse(ctx, err)
	}

	return common.ResponseOK(ctx).
		SetData(req)
}

func updTimeData(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var rid common.RequestChildId
	var req *dbx.TimeData
	var res *dbx.TimeData
	var err error

	err = c.ShouldBindUri(&rid)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = c.ShouldBindJSON(&req)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	if req.ID > 0 && rid.ChildID != req.ID {
		err = errors.Errorf("IDs don't match")
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	if req.ParentID > 0 && rid.ID != req.ParentID {
		err = errors.Errorf("Parent IDs don't match")
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	req.ID = rid.ChildID
	req.ParentID = rid.ID

	err = dbx.ExecTx(ctx, func(ctx context.Context) error {
		err = timeDataEntity.Update(ctx, req)
		if err != nil {
			return err
		}

		err = timeDataEntity.GetSingle(ctx, req.ID, &res)
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

func delTimeData(ctx context.Context) *common.Response {
	var c = getGin(ctx)
	var rid common.RequestChildId
	var err error

	err = c.ShouldBindUri(&rid)
	if err != nil {
		return common.WrapResponse2(ctx, err, common.ResponseBadRequest)
	}

	err = dbx.ExecTx(ctx, func(ctx context.Context) error {
		err = timeDataEntity.Delete(ctx, rid.ChildID)
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

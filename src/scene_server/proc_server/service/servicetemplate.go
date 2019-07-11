/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProcServer) CreateServiceTemplate(ctx *rest.Contexts) {
	template := new(metadata.ServiceTemplate)
	if err := ctx.DecodeInto(template); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(template.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "create service template, but get business id failed, err: %v", err)
		return
	}

	tpl, err := ps.CoreAPI.CoreService().Process().CreateServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, template)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "create service template failed, err: %v", err)
		return
	}

	if err := ps.AuthManager.RegisterServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header, *tpl); err != nil {
		blog.Errorf("create service template success, but register to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		err := ctx.Kit.CCError.CCError(common.CCErrCommRegistResourceToIAMFailed)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(tpl)
}

func (ps *ProcServer) GetServiceTemplate(ctx *rest.Contexts) {
	templateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	templateID, err := util.GetInt64ByInterface(templateIDStr)
	if err != nil {
		ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "create service template failed, err: %v", common.BKServiceTemplateIDField, err)
		return
	}
	temp, err := ps.CoreAPI.CoreService().Process().GetServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, templateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "get service template failed, err: %v", err)
		return
	}

	ctx.RespEntity(temp)
}

// GetServiceTemplateDetail return more info than GetServiceTemplate
func (ps *ProcServer) GetServiceTemplateDetail(ctx *rest.Contexts) {
	templateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	templateID, err := util.GetInt64ByInterface(templateIDStr)
	if err != nil {
		ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "create service template failed, err: %v", common.BKServiceTemplateIDField, err)
		return
	}
	temp, err := ps.CoreAPI.CoreService().Process().GetServiceTemplateDetail(ctx.Kit.Ctx, ctx.Kit.Header, templateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "get service template failed, err: %v", err)
		return
	}

	ctx.RespEntity(temp)
}

func (ps *ProcServer) UpdateServiceTemplate(ctx *rest.Contexts) {
	template := new(metadata.ServiceTemplate)
	if err := ctx.DecodeInto(template); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(template.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "update service template, but get business id failed, err: %v", err)
		return
	}

	tpl, err := ps.CoreAPI.CoreService().Process().UpdateServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, template.ID, template)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "update service template failed, err: %v", err)
		return
	}

	if err := ps.AuthManager.UpdateRegisteredServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header, *tpl); err != nil {
		blog.Errorf("create service template success, but register to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		err := ctx.Kit.CCError.CCError(common.CCErrCommRegistResourceToIAMFailed)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(tpl)
}

func (ps *ProcServer) ListServiceTemplates(ctx *rest.Contexts) {
	input := new(metadata.ListServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list service template, but get business id failed, err: %v", err)
		return
	}

	if input.Page.Limit >= common.BKMaxPageLimit {
		ctx.RespErrorCodeOnly(common.CCErrCommPageLimitIsExceeded, "list service template, but page limit:%d is over limited.", input.Page.Limit)
		return
	}

	option := metadata.ListServiceTemplateOption{
		BusinessID:        bizID,
		Page:              input.Page,
		ServiceCategoryID: &input.ServiceCategoryID,
	}
	temp, err := ps.CoreAPI.CoreService().Process().ListServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "list service template failed, input: %+v", input)
		return
	}

	ctx.RespEntity(temp)
}

func (ps *ProcServer) ListServiceTemplatesWithDetails(ctx *rest.Contexts) {
	input := new(metadata.ListServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list service template, but get business id failed, err: %v", err)
		return
	}

	if input.Page.Limit >= common.BKMaxPageLimit {
		ctx.RespErrorCodeOnly(common.CCErrCommPageLimitIsExceeded, "list service template, but page limit:%d is over limited.", input.Page.Limit)
		return
	}

	option := metadata.ListServiceTemplateOption{
		BusinessID:        bizID,
		Page:              input.Page,
		ServiceCategoryID: &input.ServiceCategoryID,
	}
	temp, err := ps.CoreAPI.CoreService().Process().ListServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "list service template failed, input: %+v", input)
		return
	}

	details := make([]metadata.ListServiceTemplateWithDetailResult, 0)
	for _, serviceTemplate := range temp.Info {
		// process templates reference count
		option := &metadata.ListProcessTemplatesOption{
			BusinessID:        bizID,
			ServiceTemplateID: serviceTemplate.ID,
		}
		processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed,
				"list service template: %d detail, but list process template failed.", serviceTemplate.ID)
			return
		}

		// module reference
		listModuleOption := &metadata.QueryCondition{
			Condition: mapstr.MapStr(map[string]interface{}{
				common.BKServiceTemplateIDField: serviceTemplate.ID,
			}),
		}
		moduleRst, e := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, listModuleOption)
		if e != nil {
			ctx.RespWithError(e, common.CCErrTopoModuleSelectFailed, "list service template: %d detail, but module failed.", serviceTemplate.ID)
			return
		}

		// service instance reference count
		serviceOption := &metadata.ListServiceInstanceOption{
			BusinessID:        bizID,
			ServiceTemplateID: serviceTemplate.ID,
		}
		serviceInstances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceOption)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed,
				"list service template: %d detail, but list service instance failed.", serviceTemplate.ID)
			return
		}

		details = append(details, metadata.ListServiceTemplateWithDetailResult{
			ServiceTemplate:      serviceTemplate,
			ProcessTemplateCount: int64(processTemplates.Count),
			ServiceInstanceCount: int64(serviceInstances.Count),
			ModuleCount:          int64(moduleRst.Data.Count),
		})
	}

	ctx.RespEntityWithCount(int64(temp.Count), details)
}

// a service template can be delete only when it is not be used any more,
// which means that no process instance belongs to it.
func (ps *ProcServer) DeleteServiceTemplate(ctx *rest.Contexts) {
	input := new(metadata.DeleteServiceTemplatesInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete service template, but get business id failed, err: %v", err)
		return
	}

	iamResources, err := ps.AuthManager.MakeResourcesByServiceTemplateIDs(ctx.Kit.Ctx, ctx.Kit.Header, meta.Delete, bizID, input.ServiceTemplateID)
	if err != nil {
		blog.Errorf("make iam resource by service template failed, templateID: %d, err: %+v, rid: %s", input.ServiceTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	err = ps.CoreAPI.CoreService().Process().DeleteServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceTemplateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcDeleteServiceTemplateFailed, "delete service template: %d failed", input.ServiceTemplateID)
		return
	}

	if err := ps.AuthManager.Authorize.DeregisterResource(ctx.Kit.Ctx, iamResources...); err != nil {
		blog.Errorf("delete service template success, but deregister from iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		err := ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

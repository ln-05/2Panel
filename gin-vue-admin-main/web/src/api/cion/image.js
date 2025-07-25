import service from '@/utils/request'
// @Tags Image
// @Summary 创建image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Image true "创建image表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /image/createImage [post]
export const createImage = (data) => {
  return service({
    url: '/image/createImage',
    method: 'post',
    data
  })
}

// @Tags Image
// @Summary 删除image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Image true "删除image表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /image/deleteImage [delete]
export const deleteImage = (params) => {
  return service({
    url: '/image/deleteImage',
    method: 'delete',
    params
  })
}

// @Tags Image
// @Summary 批量删除image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除image表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /image/deleteImage [delete]
export const deleteImageByIds = (params) => {
  return service({
    url: '/image/deleteImageByIds',
    method: 'delete',
    params
  })
}

// @Tags Image
// @Summary 更新image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.Image true "更新image表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /image/updateImage [put]
export const updateImage = (data) => {
  return service({
    url: '/image/updateImage',
    method: 'put',
    data
  })
}

// @Tags Image
// @Summary 用id查询image表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.Image true "用id查询image表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /image/findImage [get]
export const findImage = (params) => {
  return service({
    url: '/image/findImage',
    method: 'get',
    params
  })
}

// @Tags Image
// @Summary 分页获取image表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取image表列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /image/getImageList [get]
export const getImageList = (params) => {
  return service({
    url: '/image/getImageList',
    method: 'get',
    params
  })
}

// @Tags Image
// @Summary 不需要鉴权的image表接口
// @Accept application/json
// @Produce application/json
// @Param data query cionReq.ImageSearch true "分页获取image表列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /image/getImagePublic [get]
export const getImagePublic = () => {
  return service({
    url: '/image/getImagePublic',
    method: 'get',
  })
}

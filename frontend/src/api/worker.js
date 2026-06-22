import request from './request'

export function getMyDispatches(workerId) {
  return request({
    url: `/workers/${workerId}/dispatches`,
    method: 'get'
  })
}

export function acceptDispatch(id) {
  return request({
    url: `/dispatches/${id}/accept`,
    method: 'put'
  })
}

export function startDispatch(id) {
  return request({
    url: `/dispatches/${id}/start`,
    method: 'put'
  })
}

export function completeDispatch(id, data) {
  return request({
    url: `/dispatches/${id}/complete`,
    method: 'put',
    data
  })
}

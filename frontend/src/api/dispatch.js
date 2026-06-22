import request from './request'

export function checkDispatch(params) {
  const queryParams = { ...params }
  if (Array.isArray(queryParams.worker_ids)) {
    queryParams.worker_ids = queryParams.worker_ids.join(',')
  }
  return request({
    url: '/dispatches/check-conflict',
    method: 'get',
    params: queryParams
  })
}

export function createDispatch(data) {
  return request({
    url: '/dispatches',
    method: 'post',
    data
  })
}

export function getDispatches() {
  return request({
    url: '/dispatches',
    method: 'get'
  })
}

export function getDispatch(id) {
  return request({
    url: `/dispatches/${id}`,
    method: 'get'
  })
}

export function getSchedules(start_date, end_date) {
  return request({
    url: '/schedules',
    method: 'get',
    params: { start_date, end_date }
  })
}

export function createSchedule(data) {
  return request({
    url: '/schedules',
    method: 'post',
    data
  })
}

export function getWorkers() {
  return request({
    url: '/workers',
    method: 'get'
  })
}

export function createWorker(data) {
  return request({
    url: '/workers',
    method: 'post',
    data
  })
}

export function getVehicles() {
  return request({
    url: '/vehicles',
    method: 'get'
  })
}

export function createVehicle(data) {
  return request({
    url: '/vehicles',
    method: 'post',
    data
  })
}

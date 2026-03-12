import request from './request'

// 用户注册
export function register(data) {
  return request.post('/user/register', data)
}

// 用户登录
export function login(data) {
  return request.post('/user/login', data)
}

// 更新用户头像
export function updateHeader(formData) {
  return request.post('/user/header', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 获取用户头像
export function getHeader(username) {
  return request.get(`/user/header/${username}`)
}

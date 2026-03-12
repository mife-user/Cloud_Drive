import request from './request'

// 获取文件列表
export function getFiles(userName) {
  return request.get('/file/view', { params: { user_name: userName } })
}

// 获取单个文件
export function getFile(fileId) {
  return request.get(`/file/view/${fileId}`)
}

// 上传文件
export function uploadFile(formData, onProgress) {
  return request.post('/file/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: onProgress
  })
}

// 删除文件（移入回收站）
export function deleteFile(fileId) {
  return request.delete(`/file/delete/${fileId}`)
}

// 永久删除文件
export function deleteFileForever(fileId) {
  return request.delete(`/file/delete/${fileId}/forever`)
}

// 获取回收站文件
export function getDeletedFiles() {
  return request.get('/file/view/deleted')
}

// 创建分享
export function createShare(fileId) {
  return request.post('/file/share', { file_id: fileId })
}

// 访问分享
export function accessShare(shareId, accessKey) {
  return request.get(`/file/share/${shareId}`, { params: { access_key: accessKey } })
}

// 更新文件权限
export function updatePermissions(fileId, permissions) {
  return request.put(`/file/${fileId}/permissions`, { permissions })
}

// 添加收藏
export function addFavorite(fileId, accessKey = '') {
  return request.post('/file/favorite', { file_id: fileId, access_key: accessKey })
}

// 取消收藏
export function removeFavorite(fileId) {
  return request.delete(`/file/favorite/${fileId}`)
}

// 获取收藏列表
export function getFavorites() {
  return request.get('/file/favorites')
}

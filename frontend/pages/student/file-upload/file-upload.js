// pages/student/file-upload/file-upload.js
const app = getApp()

Page({
  data: {
    taskId: '',
    task: null,
    selectedFiles: [],
    uploading: false,
    uploadProgress: 0
  },

  onLoad(options) {
    if (options.taskId) {
      this.setData({
        taskId: options.taskId
      })
      this.loadTaskInfo()
    }
  },

  // 加载任务信息
  async loadTaskInfo() {
    try {
      const result = await app.request({
        url: `/tasks/${this.data.taskId}`,
        method: 'GET'
      })

      this.setData({
        task: result.data
      })
    } catch (error) {
      app.showToast(error.message || '加载任务信息失败')
    }
  },

  // 选择文件
  onChooseFile() {
    wx.chooseMessageFile({
      count: 1,
      type: 'file',
      success: (res) => {
        const file = res.tempFiles[0]
        this.setData({
          selectedFiles: [file]
        })
      },
      fail: () => {
        app.showToast('选择文件失败')
      }
    })
  },

  // 提交文件
  async onSubmit() {
    if (this.data.selectedFiles.length === 0) {
      app.showToast('请先选择文件')
      return
    }

    try {
      this.setData({ uploading: true })
      
      // 上传文件
      const uploadResult = await this.uploadFile(this.data.selectedFiles[0])
      
      // 提交任务
      const submitResult = await app.request({
        url: `/tasks/${this.data.taskId}/submit`,
        method: 'POST',
        data: {
          files: [uploadResult]
        }
      })

      app.showToast('提交成功', 'success')
      
      // 返回任务详情页
      setTimeout(() => {
        wx.navigateBack()
      }, 1500)

    } catch (error) {
      app.showToast(error.message || '提交失败')
    }

    this.setData({ uploading: false })
  },

  // 上传文件
  uploadFile(file) {
    return new Promise((resolve, reject) => {
      wx.uploadFile({
        url: app.globalData.baseUrl + '/files/upload',
        filePath: file.path,
        name: 'file',
        header: {
          'Authorization': `Bearer ${app.globalData.token}`
        },
        success: (res) => {
          const data = JSON.parse(res.data)
          if (data.code === 200) {
            resolve(data.data)
          } else {
            reject(new Error(data.msg))
          }
        },
        fail: reject
      })
    })
  },

  // 取消上传
  onCancel() {
    wx.navigateBack()
  }
})

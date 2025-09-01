// pages/teacher/task-detail/task-detail.js
const app = getApp()

Page({
  data: {
    taskId: null,
    taskInfo: null,
    submissions: [],
    loading: false
  },

  onLoad(options) {
    if (options.id) {
      this.setData({
        taskId: options.id
      })
      this.loadTaskDetail()
      this.loadSubmissions()
    }
  },

  // 加载任务详情
  loadTaskDetail() {
    this.setData({ loading: true })
    
    wx.request({
      url: `${app.globalData.baseUrl}/tasks/${this.data.taskId}`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({
            taskInfo: res.data.data
          })
        } else {
          wx.showToast({
            title: res.data.message || '加载失败',
            icon: 'none'
          })
        }
      },
      fail: (err) => {
        console.error('加载任务详情失败:', err)
        wx.showToast({
          title: '网络错误',
          icon: 'none'
        })
      },
      complete: () => {
        this.setData({ loading: false })
      }
    })
  },

  // 加载提交记录
  loadSubmissions() {
    wx.request({
      url: `${app.globalData.baseUrl}/submissions/task/${this.data.taskId}`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({
            submissions: res.data.data.submissions || []
          })
        }
      },
      fail: (err) => {
        console.error('加载提交记录失败:', err)
      }
    })
  },

  // 编辑任务
  onEditTask() {
    wx.navigateTo({
      url: `/pages/teacher/publish-task/publish-task?id=${this.data.taskId}`
    })
  },

  // 查看提交详情
  onViewSubmission(e) {
    const submissionId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/teacher/submission-detail/submission-detail?id=${submissionId}`
    })
  },

  // 返回
  onBack() {
    wx.navigateBack()
  }
})

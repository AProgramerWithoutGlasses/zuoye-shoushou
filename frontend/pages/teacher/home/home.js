// pages/teacher/home/home.js
const app = getApp()

Page({
  data: {
    userInfo: null,
    taskList: [],
    statistics: {
      totalTasks: 0,
      activeTasks: 0,
      totalSubmissions: 0,
      pendingReviews: 0
    },
    loading: false
  },

  onLoad() {
    this.loadUserInfo()
    this.loadTaskList()
    this.loadStatistics()
  },

  onShow() {
    // 页面显示时刷新数据
    this.loadTaskList()
    this.loadStatistics()
  },

  // 加载用户信息
  loadUserInfo() {
    const userInfo = app.globalData.userInfo
    if (userInfo) {
      this.setData({
        userInfo: userInfo
      })
    } else {
      // 没有用户信息，跳转到登录页
      wx.redirectTo({
        url: '/pages/auth/login/login'
      })
    }
  },

  // 加载任务列表
  loadTaskList() {
    this.setData({ loading: true })
    
    wx.request({
      url: `${app.globalData.baseUrl}/tasks/teacher`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({
            taskList: res.data.data.tasks || []
          })
        } else {
          wx.showToast({
            title: res.data.message || '加载失败',
            icon: 'none'
          })
        }
      },
      fail: (err) => {
        console.error('加载任务列表失败:', err)
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

  // 加载统计数据
  loadStatistics() {
    wx.request({
      url: `${app.globalData.baseUrl}/tasks/statistics`,
      method: 'GET',
      header: {
        'Authorization': `Bearer ${app.globalData.token}`
      },
      success: (res) => {
        if (res.data.code === 200) {
          this.setData({
            statistics: res.data.data || this.data.statistics
          })
        }
      },
      fail: (err) => {
        console.error('加载统计数据失败:', err)
      }
    })
  },

  // 查看任务详情
  onViewTask(e) {
    const taskId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/teacher/task-detail/task-detail?id=${taskId}`
    })
  },

  // 查看全部任务
  onViewAllTasks() {
    wx.navigateTo({
      url: '/pages/teacher/task-list/task-list'
    })
  },

  // 发布新任务
  onPublishTask() {
    wx.navigateTo({
      url: '/pages/teacher/publish-task/publish-task'
    })
  },

  // 查看提交记录
  onViewSubmissions() {
    wx.navigateTo({
      url: '/pages/teacher/submissions/submissions'
    })
  },

  // 查看统计报告
  onViewStatistics() {
    wx.navigateTo({
      url: '/pages/teacher/statistics/statistics'
    })
  },

  // 管理文件
  onManageFiles() {
    wx.navigateTo({
      url: '/pages/teacher/file-management/file-management'
    })
  },

  // 退出登录
  onLogout() {
    wx.showModal({
      title: '确认退出',
      content: '确定要退出登录吗？',
      success: (res) => {
        if (res.confirm) {
          // 清除全局数据
          app.globalData.userInfo = null
          app.globalData.token = null
          
          // 清除本地存储
          wx.removeStorageSync('userInfo')
          wx.removeStorageSync('token')
          
          // 跳转到登录页
          wx.reLaunch({
            url: '/pages/auth/login/login'
          })
        }
      }
    })
  }
})

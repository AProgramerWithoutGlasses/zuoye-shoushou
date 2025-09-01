// pages/student/profile/profile.js
const app = getApp()

Page({
  data: {
    userInfo: null,
    statistics: {
      totalTasks: 0,
      completedTasks: 0,
      pendingTasks: 0,
      averageScore: 0
    }
  },

  onLoad() {
    this.loadUserInfo()
    this.loadStatistics()
  },

  onShow() {
    this.loadUserInfo()
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

  // 加载统计数据
  loadStatistics() {
    wx.request({
      url: `${app.globalData.baseUrl}/submissions/statistics`,
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

  // 查看已完成任务
  onViewCompletedTasks() {
    wx.switchTab({
      url: '/pages/student/completed-tasks/completed-tasks'
    })
  },

  // 查看待完成任务
  onViewPendingTasks() {
    wx.switchTab({
      url: '/pages/student/home/home'
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

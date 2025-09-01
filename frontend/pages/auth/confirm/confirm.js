// pages/auth/confirm/confirm.js
const app = getApp()

Page({
  data: {
    userInfo: null
  },

  onLoad() {
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

  // 进入系统
  onEnterSystem() {
    const userInfo = this.data.userInfo
    
    if (userInfo.role === 'teacher') {
      wx.redirectTo({
        url: '/pages/teacher/home/home'
      })
    } else {
      wx.switchTab({
        url: '/pages/student/home/home'
      })
    }
  },

  // 重新登录
  onRelogin() {
    // 清除登录信息
    app.logout()
  }
})

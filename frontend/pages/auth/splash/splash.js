// pages/auth/splash/splash.js
const app = getApp()

Page({
  data: {
    
  },

  onLoad(options) {
    // 检查登录状态
    setTimeout(() => {
      this.checkLoginStatus()
    }, 2000) // 2秒后检查登录状态
  },

  // 检查登录状态
  checkLoginStatus() {
    const userInfo = app.globalData.userInfo
    const token = app.globalData.token
    
    if (token && userInfo) {
      // 已登录，跳转到相应页面
      if (userInfo.role === 'student') {
        wx.switchTab({
          url: '/pages/student/home/home'
        })
      } else if (userInfo.role === 'teacher') {
        wx.redirectTo({
          url: '/pages/teacher/home/home'
        })
      }
    } else {
      // 未登录，跳转到登录页
      wx.redirectTo({
        url: '/pages/auth/login/login'
      })
    }
  },

  // 微信授权登录
  async onWxLogin() {
    try {
      app.showLoading('登录中...')
      
      // 获取微信登录码
      const code = await app.login()
      
      // 调用后端微信登录接口
      const result = await app.request({
        url: '/auth/wx-login',
        method: 'POST',
        data: { code },
        needAuth: false
      })
      
      app.hideLoading()
      
      if (result.data.token) {
        // 登录成功
        app.saveLoginStatus(result.data.token, result.data.user)
        
        // 跳转到相应页面
        if (result.data.user.role === 'student') {
          wx.switchTab({
            url: '/pages/student/home/home'
          })
        } else if (result.data.user.role === 'teacher') {
          wx.redirectTo({
            url: '/pages/teacher/home/home'
          })
        }
      } else {
        // 需要进行账号验证
        wx.redirectTo({
          url: '/pages/auth/login/login'
        })
      }
    } catch (error) {
      app.hideLoading()
      // 微信登录失败，跳转到账号登录
      wx.redirectTo({
        url: '/pages/auth/login/login'
      })
    }
  }
})

// pages/auth/login/login.js
const app = getApp()

Page({
  data: {
    username: '',
    password: '',
    loading: false
  },

  onLoad() {
    // 检查是否已登录
    if (app.globalData.token && app.globalData.userInfo) {
      this.redirectToHome()
    }
  },

  // 输入用户名
  onUsernameInput(e) {
    this.setData({
      username: e.detail.value
    })
  },

  // 输入密码
  onPasswordInput(e) {
    this.setData({
      password: e.detail.value
    })
  },

  // 登录
  onLogin() {
    const { username, password } = this.data

    if (!username.trim()) {
      app.showToast('请输入用户名')
      return
    }

    if (!password.trim()) {
      app.showToast('请输入密码')
      return
    }

    this.setData({ loading: true })

    wx.request({
      url: `${app.globalData.baseUrl}/auth/login`,
      method: 'POST',
      header: {
        'Content-Type': 'application/json'
      },
      data: {
        username: username,
        password: password
      },
      success: (res) => {
        console.log('登录响应:', res.data) // 添加调试日志
        
        if (res.data.code === 200) {
          const { token, user } = res.data.data
          
          // 保存登录信息
          app.saveLoginInfo(token, user)
          
          // 跳转到身份确认页面
          wx.redirectTo({
            url: '/pages/auth/confirm/confirm'
          })
        } else {
          // 使用 msg 字段而不是 message
          const errorMsg = res.data.msg || '登录失败'
          app.showToast(errorMsg)
        }
      },
      fail: (err) => {
        console.error('登录失败:', err)
        app.showToast('网络错误，请重试')
      },
      complete: () => {
        this.setData({ loading: false })
      }
    })
  },

  // 跳转到首页
  redirectToHome() {
    const userInfo = app.globalData.userInfo
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

  // 返回上一页
  onBack() {
    wx.navigateBack({
      delta: 1
    })
  }
})

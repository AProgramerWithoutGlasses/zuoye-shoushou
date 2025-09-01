// app.js
App({
  globalData: {
    userInfo: null,
    token: null,
    baseUrl: 'http://localhost:8081/api'
  },

  onLaunch() {
    // 展示本地存储能力
    const logs = wx.getStorageSync('logs') || []
    logs.unshift(Date.now())
    wx.setStorageSync('logs', logs)

    // 检查登录状态
    this.checkLoginStatus()
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token')
    const userInfo = wx.getStorageSync('userInfo')
    
    if (token && userInfo) {
      this.globalData.token = token
      this.globalData.userInfo = userInfo
    }
  },

  // 保存登录信息
  saveLoginInfo(token, userInfo) {
    this.globalData.token = token
    this.globalData.userInfo = userInfo
    
    wx.setStorageSync('token', token)
    wx.setStorageSync('userInfo', userInfo)
  },

  // 退出登录
  logout() {
    // 清除全局数据
    this.globalData.userInfo = null
    this.globalData.token = null
    
    // 清除本地存储
    wx.removeStorageSync('userInfo')
    wx.removeStorageSync('token')
    
    // 跳转到登录页
    wx.reLaunch({
      url: '/pages/auth/login/login'
    })
  },

  // 显示提示
  showToast(title, icon = 'none') {
    wx.showToast({
      title: title,
      icon: icon
    })
  },

  // 显示加载
  showLoading(title = '加载中...') {
    wx.showLoading({
      title: title
    })
  },

  // 隐藏加载
  hideLoading() {
    wx.hideLoading()
  },

  // 网络请求封装
  request(options) {
    return new Promise((resolve, reject) => {
      wx.request({
        url: this.globalData.baseUrl + options.url,
        method: options.method || 'GET',
        data: options.data || {},
        header: {
          'Content-Type': 'application/json',
          'Authorization': this.globalData.token ? `Bearer ${this.globalData.token}` : '',
          ...options.header
        },
        success: (res) => {
          if (res.statusCode === 200) {
            if (res.data.code === 200) {
              resolve(res.data)
            } else if (res.data.code === 401) {
              // token过期，退出登录
              this.logout()
              reject(new Error('登录已过期，请重新登录'))
            } else {
              // 使用 msg 字段而不是 message
              reject(new Error(res.data.msg || '请求失败'))
            }
          } else {
            reject(new Error('网络请求失败'))
          }
        },
        fail: (err) => {
          reject(err)
        }
      })
    })
  }
})


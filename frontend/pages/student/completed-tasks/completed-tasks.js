// pages/student/completed-tasks/completed-tasks.js
const app = getApp()

Page({
  data: {
    tasks: [],
    loading: false,
    page: 1,
    size: 10,
    hasMore: true
  },

  onLoad(options) {
    this.loadTasks()
  },

  onShow() {
    // 刷新数据
    this.refreshTasks()
  },

  onPullDownRefresh() {
    this.refreshTasks()
  },

  onReachBottom() {
    this.loadMoreTasks()
  },

  // 刷新任务列表
  async refreshTasks() {
    this.setData({
      page: 1,
      hasMore: true
    })
    await this.loadTasks()
    wx.stopPullDownRefresh()
  },

  // 加载任务列表
  async loadTasks() {
    if (this.data.loading) return

    this.setData({ loading: true })

    try {
      const result = await app.request({
        url: '/submissions',
        method: 'GET',
        data: {
          page: this.data.page,
          size: this.data.size
        }
      })

      const newTasks = result.data.submissions || []

      this.setData({
        tasks: this.data.page === 1 ? newTasks : [...this.data.tasks, ...newTasks],
        hasMore: newTasks.length === this.data.size
      })

    } catch (error) {
      app.showToast(error.message || '加载失败')
    }

    this.setData({ loading: false })
  },

  // 加载更多
  async loadMoreTasks() {
    if (!this.data.hasMore || this.data.loading) return

    this.setData({
      page: this.data.page + 1
    })
    await this.loadTasks()
  },

  // 查看任务详情
  onTaskTap(e) {
    const taskId = e.currentTarget.dataset.taskId
    wx.navigateTo({
      url: `/pages/student/task-detail/task-detail?id=${taskId}`
    })
  },

  // 格式化时间
  formatTime(timeStr) {
    const time = new Date(timeStr)
    const month = time.getMonth() + 1
    const day = time.getDate()
    return `${month}-${day.toString().padStart(2, '0')}`
  }
})

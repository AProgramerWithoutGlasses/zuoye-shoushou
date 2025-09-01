// pages/student/home/home.js
const app = getApp()

Page({
  data: {
    tasks: [],
    loading: false,
    refreshing: false,
    hasMore: true,
    page: 1,
    size: 10,
    searchKeyword: '',
    selectedTeacher: '',
    selectedPeriod: '',
    teacherList: [],
    pendingCount: 0
  },

  onLoad(options) {
    this.loadTasks()
  },

  onShow() {
    // 每次显示页面时刷新数据
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
      hasMore: true,
      refreshing: true
    })
    await this.loadTasks()
    this.setData({
      refreshing: false
    })
    wx.stopPullDownRefresh()
  },

  // 加载任务列表
  async loadTasks() {
    if (this.data.loading) return

    this.setData({ loading: true })

    try {
      const result = await app.request({
        url: '/tasks/student',
        method: 'GET',
        data: {
          page: this.data.page,
          size: this.data.size,
          status: 'active' // 只获取进行中的任务
        }
      })

      const newTasks = result.data.tasks || []
      
      // 统计未完成数量
      const pendingCount = newTasks.filter(task => !task.submitted).length

      this.setData({
        tasks: this.data.page === 1 ? newTasks : [...this.data.tasks, ...newTasks],
        hasMore: newTasks.length === this.data.size,
        pendingCount: pendingCount
      })

      // 如果是第一页，更新教师列表
      if (this.data.page === 1) {
        this.updateTeacherList()
      }

    } catch (error) {
      app.showToast(error.message || '加载失败')
    }

    this.setData({ loading: false })
  },

  // 更新教师列表
  updateTeacherList() {
    try {
      // 从任务中提取教师信息
      const teachers = [...new Set(this.data.tasks.map(task => task.teacher?.name).filter(Boolean))]
      this.setData({
        teacherList: teachers
      })
    } catch (error) {
      console.error('更新教师列表失败', error)
    }
  },

  // 加载更多任务
  async loadMoreTasks() {
    if (!this.data.hasMore || this.data.loading) return

    this.setData({
      page: this.data.page + 1
    })
    await this.loadTasks()
  },



  // 搜索
  onSearchInput(e) {
    this.setData({
      searchKeyword: e.detail.value
    })
  },

  onSearch() {
    this.refreshTasks()
  },

  // 筛选
  onTeacherFilter() {
    wx.showActionSheet({
      itemList: ['全部教师', ...this.data.teacherList],
      success: (res) => {
        const selected = res.tapIndex === 0 ? '' : this.data.teacherList[res.tapIndex - 1]
        this.setData({
          selectedTeacher: selected
        })
        this.refreshTasks()
      }
    })
  },

  onPeriodFilter() {
    wx.showActionSheet({
      itemList: ['全部时间', '本周', '本月', '即将截止'],
      success: (res) => {
        const periods = ['', 'week', 'month', 'deadline']
        this.setData({
          selectedPeriod: periods[res.tapIndex]
        })
        this.refreshTasks()
      }
    })
  },

  // 跳转到任务详情
  onTaskTap(e) {
    const taskId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/student/task-detail/task-detail?id=${taskId}`
    })
  },

  // 快速提交
  onQuickSubmit(e) {
    const taskId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/student/file-upload/file-upload?taskId=${taskId}`
    })
  },

  // 获取任务状态样式
  getTaskStatusClass(task) {
    if (task.submitted) {
      return 'status-submitted'
    }
    
    const now = new Date()
    const endTime = new Date(task.end_time)
    const timeDiff = endTime - now
    
    if (timeDiff < 0) {
      return 'status-expired'
    } else if (timeDiff < 24 * 60 * 60 * 1000) { // 24小时内
      return 'status-urgent'
    } else {
      return 'status-pending'
    }
  },

  // 格式化时间
  formatTime(timeStr) {
    const time = new Date(timeStr)
    const month = time.getMonth() + 1
    const day = time.getDate()
    const hour = time.getHours()
    const minute = time.getMinutes()
    
    return `${month}-${day.toString().padStart(2, '0')} ${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
  }
})

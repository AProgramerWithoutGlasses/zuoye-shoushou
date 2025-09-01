// pages/student/task-detail/task-detail.js
const app = getApp()

Page({
  data: {
    taskId: '',
    task: {},
    mySubmission: null,
    loading: true,
    submittedStudents: [],
    unsubmittedStudents: [],
    allStudents: [], // 新增：存储所有学生的提交状态
    submitRate: '0.0', // 新增：提交率
    progressWidth: 0, // 新增：进度条宽度
    unsubmittedCount: 0, // 新增：未提交人数
    allowedFormatsText: '' // 新增：允许格式文本
  },

  // 页面加载
  async onLoad(options) {
    if (options.id) {
      this.setData({ taskId: options.id })
      await this.loadTaskDetail()
      await this.loadSubmissions()
      await this.loadMySubmission()
      this.setData({ loading: false })
    }
  },

  onShow() {
    // 每次显示时刷新数据
    if (this.data.taskId) {
      this.loadMySubmission()
    }
  },

  // 加载任务详情
  async loadTaskDetail() {
    try {
      const result = await app.request({
        url: `/tasks/${this.data.taskId}`,
        method: 'GET'
      })

      const task = result.data
      
      // 计算提交率
      let submitRate = '0.0'
      let progressWidth = 0
      let unsubmittedCount = 0
      let allowedFormatsText = ''
      if (task.total_students > 0) {
        submitRate = ((task.submitted_count / task.total_students) * 100).toFixed(1)
        progressWidth = Math.round((task.submitted_count / task.total_students) * 100)
        unsubmittedCount = task.total_students - task.submitted_count
      }
      
      // 处理允许格式文本
      if (task.allowed_formats && Array.isArray(task.allowed_formats)) {
        allowedFormatsText = task.allowed_formats.join(', ')
      }

      this.setData({
        task: task,
        submitRate: submitRate,
        progressWidth: progressWidth,
        unsubmittedCount: unsubmittedCount,
        allowedFormatsText: allowedFormatsText
      })

    } catch (error) {
      console.error('加载任务详情失败', error)
      app.showToast('加载任务详情失败')
    }
  },

  // 加载提交情况
  async loadSubmissions() {
    try {
      const result = await app.request({
        url: `/tasks/${this.data.taskId}/students-status`,
        method: 'GET'
      })

      const studentsData = result.data.students || []
      const submitted = []
      const unsubmitted = []

      // 按学号排序
      studentsData.sort((a, b) => {
        const idA = a.student.student_id || a.student.username
        const idB = b.student.student_id || b.student.username
        return idA.localeCompare(idB)
      })

      // 分类已提交和未提交的学生
      studentsData.forEach(studentData => {
        if (studentData.has_submitted && studentData.submission) {
          submitted.push({
            ...studentData.submission,
            student: studentData.student
          })
        } else {
          unsubmitted.push({
            student: studentData.student
          })
        }
      })

      this.setData({
        allStudents: studentsData,
        submittedStudents: submitted,
        unsubmittedStudents: unsubmitted
      })

    } catch (error) {
      console.error('加载提交情况失败', error)
    }
  },

  // 加载我的提交记录
  async loadMySubmission() {
    try {
      const result = await app.request({
        url: `/tasks/${this.data.taskId}/submission`,
        method: 'GET'
      })

      this.setData({
        mySubmission: result.data
      })

    } catch (error) {
      // 可能还没有提交记录
      console.log('暂无提交记录')
    }
  },

  // 立即提交
  onSubmit() {
    wx.navigateTo({
      url: `/pages/student/file-upload/file-upload?taskId=${this.data.taskId}`
    })
  },

  // 修改提交
  onResubmit() {
    wx.showModal({
      title: '确认修改',
      content: '确定要修改已提交的文件吗？',
      success: (res) => {
        if (res.confirm) {
          wx.navigateTo({
            url: `/pages/student/file-upload/file-upload?taskId=${this.data.taskId}&resubmit=true`
          })
        }
      }
    })
  },

  // 预览文件
  onPreviewFile(e) {
    const fileUrl = e.currentTarget.dataset.url
    const fileName = e.currentTarget.dataset.name
    
    // 如果是图片，直接预览
    if (this.isImageFile(fileName)) {
      wx.previewImage({
        urls: [fileUrl],
        current: fileUrl
      })
    } else {
      // 其他文件类型下载后打开
      wx.downloadFile({
        url: fileUrl,
        success: (res) => {
          wx.openDocument({
            filePath: res.tempFilePath,
            success: () => {
              console.log('打开文档成功')
            },
            fail: () => {
              app.showToast('暂不支持此文件类型的预览')
            }
          })
        },
        fail: () => {
          app.showToast('文件下载失败')
        }
      })
    }
  },

  // 下载文件
  onDownloadFile(e) {
    const fileUrl = e.currentTarget.dataset.url
    const fileName = e.currentTarget.dataset.name
    
    wx.downloadFile({
      url: fileUrl,
      success: (res) => {
        app.showToast('文件已下载到本地', 'success')
      },
      fail: () => {
        app.showToast('下载失败')
      }
    })
  },

  // 判断是否为图片文件
  isImageFile(fileName) {
    const imageExts = ['.jpg', '.jpeg', '.png', '.gif', '.bmp', '.webp']
    const ext = fileName.toLowerCase().substring(fileName.lastIndexOf('.'))
    return imageExts.includes(ext)
  },

  // 格式化时间
  formatTime(timeStr) {
    const time = new Date(timeStr)
    const month = time.getMonth() + 1
    const day = time.getDate()
    const hour = time.getHours()
    const minute = time.getMinutes()
    
    return `${month}-${day.toString().padStart(2, '0')} ${hour.toString().padStart(2, '0')}:${minute.toString().padStart(2, '0')}`
  },

  // 格式化文件大小
  formatFileSize(bytes) {
    if (bytes < 1024) return bytes + ' B'
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB'
  },

  // 获取提交状态文本
  getSubmissionStatusText(submission) {
    if (!submission) return '未提交'
    
    switch (submission.status) {
      case 'submitted':
        return submission.is_on_time ? '已提交' : '迟交'
      case 'late':
        return '迟交'
      case 'reviewed':
        return '已批阅'
      default:
        return '未提交'
    }
  },

  // 获取提交状态样式
  getSubmissionStatusClass(submission) {
    if (!submission) return 'status-pending'
    
    switch (submission.status) {
      case 'submitted':
        return submission.is_on_time ? 'status-submitted' : 'status-late'
      case 'late':
        return 'status-late'
      case 'reviewed':
        return 'status-reviewed'
      default:
        return 'status-pending'
    }
  },

  // 返回上一页
  onBack() {
    wx.navigateBack({
      delta: 1
    })
  }
})

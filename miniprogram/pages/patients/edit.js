const api = require('../../utils/api')

Page({
  data: {
    isEdit: false,
    patientId: '',
    form: {
      name: '',
      id_card: '',
      phone: '',
      gender: 1,
      birthday: '',
      is_default: false
    },
    submitting: false
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ isEdit: true, patientId: options.id })
      this.loadPatient()
    }
  },

  async loadPatient() {
    try {
      const patients = await api.getPatients()
      const patient = patients.find(p => p.id == this.data.patientId)
      if (patient) {
        this.setData({
          form: {
            name: patient.name,
            id_card: patient.id_card,
            phone: patient.phone,
            gender: patient.gender,
            birthday: patient.birthday,
            is_default: patient.is_default
          }
        })
      }
    } catch (e) {
      console.error('Failed to load patient:', e)
    }
  },

  onInput(e) {
    const field = e.currentTarget.dataset.field
    this.setData({ [`form.${field}`]: e.detail.value })
  },

  selectGender(e) {
    const gender = parseInt(e.currentTarget.dataset.gender)
    this.setData({ 'form.gender': gender })
  },

  onBirthdayChange(e) {
    this.setData({ 'form.birthday': e.detail.value })
  },

  onDefaultChange(e) {
    this.setData({ 'form.is_default': e.detail.value })
  },

  async submit() {
    const { form, isEdit, patientId, submitting } = this.data
    if (submitting) return

    if (!form.name) {
      wx.showToast({ title: '请输入姓名', icon: 'none' })
      return
    }
    if (!form.id_card) {
      wx.showToast({ title: '请输入身份证号', icon: 'none' })
      return
    }
    if (!form.phone) {
      wx.showToast({ title: '请输入手机号', icon: 'none' })
      return
    }

    this.setData({ submitting: true })
    try {
      if (isEdit) {
        await api.updatePatient(patientId, form)
      } else {
        await api.createPatient(form)
      }
      wx.showToast({ title: '保存成功', icon: 'success' })
      setTimeout(() => wx.navigateBack(), 1500)
    } catch (e) {
      wx.showToast({ title: e || '保存失败', icon: 'none' })
    } finally {
      this.setData({ submitting: false })
    }
  }
})

<template>
  <div class="worker-view">
    <el-card class="mb15">
      <div class="worker-select-bar">
        <el-icon><User /></el-icon>
        <span class="label">当前师傅：</span>
        <el-select
          v-model="selectedWorkerId"
          placeholder="请选择师傅（模拟登录）"
          style="width: 250px"
          @change="loadDispatches"
        >
          <el-option
            v-for="worker in workerList"
            :key="worker.id"
            :label="worker.name"
            :value="worker.id"
          />
        </el-select>
        <el-button type="primary" @click="loadDispatches" style="margin-left: 10px">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </el-card>

    <template v-if="selectedWorkerId">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card class="order-group-card">
            <template #header>
              <div class="group-header">
                <el-icon><Clock /></el-icon>
                <span>待接单</span>
                <el-tag type="warning" effect="dark">{{ groupedDispatches.pending?.length || 0 }}</el-tag>
              </div>
            </template>
            <el-empty v-if="!groupedDispatches.pending?.length" description="暂无待接单" />
            <div v-for="dispatch in groupedDispatches.pending" :key="dispatch.id" class="dispatch-card">
              <div class="dispatch-card-header">
                <span class="dispatch-id">#{{ dispatch.id }}</span>
                <el-tag type="warning">待接单</el-tag>
              </div>
              <div class="dispatch-card-body">
                <div class="row"><span class="label">客户:</span> {{ dispatch.order?.customer_name }}</div>
                <div class="row"><span class="label">起始:</span> {{ dispatch.order?.start_address }}</div>
                <div class="row"><span class="label">目的:</span> {{ dispatch.order?.end_address }}</div>
                <div class="row">
                  <span class="label">计划:</span> {{ formatDateTime(dispatch.scheduled_start_time) }} ~ {{ formatDateTime(dispatch.scheduled_end_time) }}
                </div>
              </div>
              <div class="dispatch-card-footer">
                <el-button type="primary" size="small" @click="acceptDispatchHandler(dispatch)" :loading="dispatch.loading">
                  <el-icon><Check /></el-icon>
                  确认接单
                </el-button>
                <el-button size="small" @click="viewDispatchDetail(dispatch)">查看详情</el-button>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card class="order-group-card">
            <template #header>
              <div class="group-header">
                <el-icon><Operation /></el-icon>
                <span>进行中</span>
                <el-tag type="success" effect="dark">{{ workingCount }}</el-tag>
              </div>
            </template>
            <el-empty v-if="workingCount === 0" description="暂无进行中派单" />
            <div v-for="dispatch in workingDispatches" :key="dispatch.id" class="dispatch-card">
              <div class="dispatch-card-header">
                <span class="dispatch-id">#{{ dispatch.id }}</span>
                <el-tag :type="dispatch.status === 'accepted' ? 'primary' : 'success'">
                  {{ dispatch.status === 'accepted' ? '待开始' : '作业中' }}
                </el-tag>
              </div>
              <div class="dispatch-card-body">
                <div class="row"><span class="label">客户:</span> {{ dispatch.order?.customer_name }}</div>
                <div class="row"><span class="label">起始:</span> {{ dispatch.order?.start_address }}</div>
                <div class="row"><span class="label">目的:</span> {{ dispatch.order?.end_address }}</div>
                <div class="row">
                  <span class="label">计划:</span> {{ formatDateTime(dispatch.scheduled_start_time) }} ~ {{ formatDateTime(dispatch.scheduled_end_time) }}
                </div>
              </div>
              <div class="dispatch-card-footer">
                <template v-if="dispatch.status === 'accepted'">
                  <el-button type="primary" size="small" @click="startDispatchHandler(dispatch)" :loading="dispatch.loading">
                    <el-icon><VideoPlay /></el-icon>
                    开始作业
                  </el-button>
                </template>
                <template v-else>
                  <el-button type="primary" size="small" @click="openFinishDialog(dispatch)">
                    <el-icon><CircleCheck /></el-icon>
                    完成作业
                  </el-button>
                </template>
                <el-button size="small" @click="viewDispatchDetail(dispatch)">查看详情</el-button>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :span="8">
          <el-card class="order-group-card">
            <template #header>
              <div class="group-header">
                <el-icon><Finished /></el-icon>
                <span>已完成</span>
                <el-tag type="info" effect="dark">{{ groupedDispatches.completed?.length || 0 }}</el-tag>
              </div>
            </template>
            <el-empty v-if="!groupedDispatches.completed?.length" description="暂无已完成派单" />
            <div v-for="dispatch in groupedDispatches.completed" :key="dispatch.id" class="dispatch-card">
              <div class="dispatch-card-header">
                <span class="dispatch-id">#{{ dispatch.id }}</span>
                <el-tag type="success">已完成</el-tag>
              </div>
              <div class="dispatch-card-body">
                <div class="row"><span class="label">客户:</span> {{ dispatch.order?.customer_name }}</div>
                <div class="row"><span class="label">起始:</span> {{ dispatch.order?.start_address }}</div>
                <div class="row"><span class="label">目的:</span> {{ dispatch.order?.end_address }}</div>
                <div class="row">
                  <span class="label">实际体积:</span> {{ dispatch.actual_items_volume || dispatch.order?.items_volume }} m³
                  <span class="label" style="margin-left: 15px">实际人数:</span> {{ dispatch.actual_workers_count || '-' }}
                </div>
              </div>
              <div class="dispatch-card-footer">
                <el-button size="small" @click="viewDispatchDetail(dispatch)">查看详情</el-button>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </template>

    <el-dialog
      v-model="finishDialogVisible"
      title="完成作业"
      width="500px"
    >
      <el-form :model="finishForm" :rules="finishRules" ref="finishFormRef" label-width="120px">
        <el-form-item label="派单ID">
          <el-tag type="primary">#{{ finishForm.dispatch_id }}</el-tag>
        </el-form-item>
        <el-form-item label="实际开始时间" prop="actual_start_time">
          <el-date-picker
            v-model="finishForm.actual_start_time"
            type="datetime"
            placeholder="选择开始时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="实际结束时间" prop="actual_end_time">
          <el-date-picker
            v-model="finishForm.actual_end_time"
            type="datetime"
            placeholder="选择结束时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="实际人数" prop="actual_workers_count">
          <el-input-number v-model="finishForm.actual_workers_count" :min="1" :max="20" style="width: 100%" />
        </el-form-item>
        <el-form-item label="实际体积" prop="actual_items_volume">
          <el-input-number v-model="finishForm.actual_items_volume" :min="0" :max="100" :step="0.5" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="finishForm.remark" type="textarea" :rows="3" placeholder="请输入备注（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="finishDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmFinish" :loading="finishing">
          <el-icon><Check /></el-icon>
          确认完成
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getMyDispatches, acceptDispatch, startDispatch, completeDispatch } from '@/api/worker'
import { getWorkers } from '@/api/dispatch'

const selectedWorkerId = ref('')
const workerList = ref([])
const loading = ref(false)
const finishDialogVisible = ref(false)
const finishing = ref(false)
const finishFormRef = ref(null)
const currentDispatches = ref([])

const finishForm = reactive({
  dispatch_id: null,
  actual_start_time: '',
  actual_end_time: '',
  actual_workers_count: 2,
  actual_items_volume: 0,
  remark: ''
})

const finishRules = {
  actual_start_time: [{ required: true, message: '请选择实际开始时间', trigger: 'change' }],
  actual_end_time: [{ required: true, message: '请选择实际结束时间', trigger: 'change' }],
  actual_workers_count: [{ required: true, message: '请输入实际人数', trigger: 'blur' }],
  actual_items_volume: [{ required: true, message: '请输入实际体积', trigger: 'blur' }]
}

const mockWorkers = [
  { id: 1, name: '张师傅', phone: '13800138001' },
  { id: 2, name: '李师傅', phone: '13800138002' },
  { id: 3, name: '王师傅', phone: '13800138003' },
  { id: 4, name: '赵师傅', phone: '13800138004' },
  { id: 5, name: '刘师傅', phone: '13800138005' }
]

const mockDispatches = [
  {
    id: 1,
    order_id: 2,
    status: 'pending',
    scheduled_start_time: '2025-06-22T09:00:00',
    scheduled_end_time: '2025-06-22T12:00:00',
    order: {
      id: 2,
      customer_name: '李四',
      customer_phone: '13900139000',
      start_address: '北京市西城区金融街15号',
      end_address: '北京市东城区王府井大街88号',
      items_volume: 25,
      estimated_workers: 5
    }
  },
  {
    id: 2,
    order_id: 5,
    status: 'accepted',
    scheduled_start_time: '2025-06-22T14:00:00',
    scheduled_end_time: '2025-06-22T17:00:00',
    order: {
      id: 5,
      customer_name: '钱七',
      customer_phone: '13500135000',
      start_address: '北京市大兴区黄村西大街1号',
      end_address: '北京市顺义区府前东街5号',
      items_volume: 18,
      estimated_workers: 3
    }
  },
  {
    id: 3,
    order_id: 6,
    status: 'in_progress',
    scheduled_start_time: '2025-06-22T08:00:00',
    scheduled_end_time: '2025-06-22T11:00:00',
    actual_start_time: '2025-06-22T08:15:00',
    order: {
      id: 6,
      customer_name: '孙八',
      customer_phone: '13400134000',
      start_address: '北京市房山区良乡西路10号',
      end_address: '北京市门头沟区新桥大街20号',
      items_volume: 10,
      estimated_workers: 2
    }
  },
  {
    id: 4,
    order_id: 3,
    status: 'completed',
    scheduled_start_time: '2025-06-20T09:00:00',
    scheduled_end_time: '2025-06-20T12:00:00',
    actual_start_time: '2025-06-20T09:10:00',
    actual_end_time: '2025-06-20T11:45:00',
    actual_workers_count: 2,
    actual_items_volume: 9,
    order: {
      id: 3,
      customer_name: '王五',
      customer_phone: '13700137000',
      start_address: '北京市丰台区南三环西路16号',
      end_address: '北京市通州区新华大街99号',
      items_volume: 8,
      estimated_workers: 2
    }
  }
]

const groupedDispatches = computed(() => {
  const groups = { pending: [], accepted: [], in_progress: [], completed: [], rejected: [] }
  currentDispatches.value.forEach((dispatch) => {
    if (groups[dispatch.status]) {
      groups[dispatch.status].push(dispatch)
    }
  })
  return groups
})

const workingCount = computed(() => {
  return (groupedDispatches.value.accepted?.length || 0) + (groupedDispatches.value.in_progress?.length || 0)
})

const workingDispatches = computed(() => {
  return [...(groupedDispatches.value.accepted || []), ...(groupedDispatches.value.in_progress || [])]
})

const statusTextMap = {
  pending: '待接单',
  accepted: '已接单',
  in_progress: '作业中',
  completed: '已完成',
  rejected: '已拒绝'
}

function formatDate(date) {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function formatDateTime(date) {
  if (!date) return ''
  const d = new Date(date)
  const pad = (n) => String(n).padStart(2, '0')
  return `${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadWorkers = async () => {
  try {
    const data = await getWorkers()
    workerList.value = data.list || data.records || data || mockWorkers
  } catch (e) {
    workerList.value = mockWorkers
  }
  if (workerList.value.length > 0 && !selectedWorkerId.value) {
    selectedWorkerId.value = workerList.value[0].id
  }
}

const loadDispatches = async () => {
  if (!selectedWorkerId.value) return
  loading.value = true
  try {
    const data = await getMyDispatches(selectedWorkerId.value)
    currentDispatches.value = data.list || data.records || data || mockDispatches
  } catch (e) {
    currentDispatches.value = mockDispatches
  } finally {
    loading.value = false
  }
}

const acceptDispatchHandler = async (dispatch) => {
  dispatch.loading = true
  try {
    await acceptDispatch(dispatch.id)
    ElMessage.success('接单成功')
    dispatch.status = 'accepted'
  } catch (e) {
    ElMessage.error('接单失败')
  } finally {
    dispatch.loading = false
  }
}

const startDispatchHandler = async (dispatch) => {
  dispatch.loading = true
  try {
    await startDispatch(dispatch.id)
    ElMessage.success('已开始作业')
    dispatch.status = 'in_progress'
  } catch (e) {
    ElMessage.error('操作失败')
  } finally {
    dispatch.loading = false
  }
}

const openFinishDialog = (dispatch) => {
  finishForm.dispatch_id = dispatch.id
  finishForm.actual_start_time = dispatch.actual_start_time ? new Date(dispatch.actual_start_time) : new Date()
  finishForm.actual_end_time = new Date()
  finishForm.actual_workers_count = dispatch.order?.estimated_workers || 2
  finishForm.actual_items_volume = dispatch.order?.items_volume || 0
  finishForm.remark = ''
  finishDialogVisible.value = true
}

const confirmFinish = async () => {
  if (!finishFormRef.value) return
  await finishFormRef.value.validate(async (valid) => {
    if (valid) {
      finishing.value = true
      try {
        const pad = (n) => String(n).padStart(2, '0')
        const fmt = (d) => {
          const dt = new Date(d)
          return `${dt.getFullYear()}-${pad(dt.getMonth() + 1)}-${pad(dt.getDate())} ${pad(dt.getHours())}:${pad(dt.getMinutes())}:${pad(dt.getSeconds())}`
        }
        await completeDispatch(finishForm.dispatch_id, {
          actual_start_time: fmt(finishForm.actual_start_time),
          actual_end_time: fmt(finishForm.actual_end_time),
          actual_workers_count: finishForm.actual_workers_count,
          actual_items_volume: finishForm.actual_items_volume,
          remark: finishForm.remark
        })
        ElMessage.success('作业已完成')
        finishDialogVisible.value = false
        loadDispatches()
      } catch (e) {
        ElMessage.error('操作失败')
      } finally {
        finishing.value = false
      }
    }
  })
}

const viewDispatchDetail = (dispatch) => {
  const order = dispatch.order || {}
  ElMessageBox.alert(
    `<div style="text-align: left">
      <p><strong>派单ID:</strong> #${dispatch.id}</p>
      <p><strong>客户:</strong> ${order.customer_name || '-'} (${order.customer_phone || '-'})</p>
      <p><strong>起始地址:</strong> ${order.start_address || '-'}</p>
      <p><strong>目的地址:</strong> ${order.end_address || '-'}</p>
      <p><strong>预估体积:</strong> ${order.items_volume || '-'} m³${dispatch.actual_items_volume ? `（实际: ${dispatch.actual_items_volume} m³）` : ''}</p>
      <p><strong>预估人数:</strong> ${order.estimated_workers || '-'}人${dispatch.actual_workers_count ? `（实际: ${dispatch.actual_workers_count}人）` : ''}</p>
      <p><strong>计划开始:</strong> ${formatDateTime(dispatch.scheduled_start_time)}</p>
      <p><strong>计划结束:</strong> ${formatDateTime(dispatch.scheduled_end_time)}</p>
      ${dispatch.actual_start_time ? `<p><strong>实际开始:</strong> ${formatDateTime(dispatch.actual_start_time)}</p>` : ''}
      ${dispatch.actual_end_time ? `<p><strong>实际结束:</strong> ${formatDateTime(dispatch.actual_end_time)}</p>` : ''}
      ${dispatch.remark ? `<p><strong>备注:</strong> ${dispatch.remark}</p>` : ''}
      <p><strong>状态:</strong> ${statusTextMap[dispatch.status] || dispatch.status}</p>
    </div>`,
    '派单详情',
    { dangerouslyUseHTMLString: true, confirmButtonText: '确定' }
  )
}

onMounted(async () => {
  await loadWorkers()
  if (selectedWorkerId.value) {
    loadDispatches()
  }
})
</script>

<style scoped>
.worker-view {
  padding: 0;
}

.mb15 {
  margin-bottom: 15px;
}

.worker-select-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.worker-select-bar .label {
  font-weight: bold;
}

.group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
  font-size: 16px;
}

.order-group-card {
  min-height: 500px;
}

.dispatch-card {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;
  background-color: #fff;
  transition: all 0.3s;
}

.dispatch-card:hover {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.dispatch-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 8px;
  border-bottom: 1px solid #ebeef5;
  margin-bottom: 8px;
}

.dispatch-id {
  font-weight: bold;
  color: #409eff;
}

.dispatch-card-body {
  font-size: 13px;
  line-height: 1.8;
}

.dispatch-card-body .row {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dispatch-card-body .label {
  color: #909399;
  margin-right: 4px;
}

.dispatch-card-footer {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px solid #ebeef5;
  display: flex;
  gap: 8px;
}
</style>

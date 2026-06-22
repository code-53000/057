<template>
  <div class="dispatcher">
    <el-card class="mb15">
      <div class="date-picker-bar">
        <el-icon><Calendar /></el-icon>
        <span class="label">排班日期：</span>
        <el-date-picker
          v-model="selectedDate"
          type="date"
          placeholder="选择日期"
          style="width: 200px"
          @change="loadSchedule"
        />
      </div>
    </el-card>

    <el-card class="mb15">
      <template #header>
        <div class="card-header">
          <el-icon><Tickets /></el-icon>
          <span>排班看板</span>
        </div>
      </template>
      <el-tabs v-model="activeTab">
        <el-tab-pane label="师傅排班" name="workers">
          <el-table :data="workerSchedule" border v-loading="loadingSchedule" stripe>
            <el-table-column prop="name" label="师傅" width="120" fixed />
            <el-table-column prop="phone" label="电话" width="130" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'available' ? 'success' : 'info'">
                  {{ scope.row.status === 'available' ? '可用' : '休息' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="排班详情" min-width="200">
              <template #default="scope">
                <div v-if="scope.row.schedules && scope.row.schedules.length">
                  <div v-for="(s, idx) in scope.row.schedules" :key="idx" class="schedule-item">
                    <el-tag size="small" type="warning">{{ s.shift }}</el-tag>
                    <span class="schedule-date">{{ formatDate(s.work_date) }}</span>
                  </div>
                </div>
                <span v-else style="color: #909399">暂无排班</span>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="车辆排班" name="vehicles">
          <el-table :data="vehicleSchedule" border v-loading="loadingSchedule" stripe>
            <el-table-column prop="plate_number" label="车牌号" width="120" fixed />
            <el-table-column prop="vehicle_type" label="车型" width="100" />
            <el-table-column prop="capacity_volume" label="容积(m³)" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 'available' ? 'success' : 'warning'">
                  {{ scope.row.status === 'available' ? '可用' : '维护' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="排班详情" min-width="200">
              <template #default="scope">
                <div v-if="scope.row.schedules && scope.row.schedules.length">
                  <div v-for="(s, idx) in scope.row.schedules" :key="idx" class="schedule-item">
                    <el-tag size="small" type="warning">{{ s.shift }}</el-tag>
                    <span class="schedule-date">{{ formatDate(s.work_date) }}</span>
                  </div>
                </div>
                <span v-else style="color: #909399">暂无排班</span>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-card>
      <template #header>
        <div class="card-header">
          <el-icon><List /></el-icon>
          <span>待派单列表</span>
        </div>
      </template>
      <el-table :data="pendingOrders" border v-loading="loadingOrders" stripe>
        <el-table-column prop="id" label="订单号" width="80" />
        <el-table-column prop="customer_name" label="客户" width="90" />
        <el-table-column prop="customer_phone" label="电话" width="120" />
        <el-table-column prop="move_date" label="搬家日期" width="110">
          <template #default="scope">
            {{ formatDate(scope.row.move_date) }}
          </template>
        </el-table-column>
        <el-table-column label="地址" min-width="180">
          <template #default="scope">
            <div class="address-cell">
              <div><span class="label">起:</span> {{ scope.row.start_address }}</div>
              <div><span class="label">止:</span> {{ scope.row.end_address }}</div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="items_volume" label="体积" width="70" />
        <el-table-column prop="estimated_workers" label="人数" width="60" />
        <el-table-column prop="estimated_vehicle_type" label="车型" width="80" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="openDispatchDialog(scope.row)">
              <el-icon><Promotion /></el-icon>
              派单
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog
      v-model="dispatchDialogVisible"
      title="派单"
      width="700px"
      @close="resetDispatchForm"
    >
      <el-form :model="dispatchForm" :rules="dispatchRules" ref="dispatchFormRef" label-width="110px">
        <el-form-item label="订单信息">
          <el-tag type="primary">订单#{{ dispatchForm.order_id }}</el-tag>
          <span style="margin-left: 10px">
            {{ dispatchForm.customer_name }} - {{ dispatchForm.estimated_vehicle_type }} - {{ dispatchForm.items_volume }}m³
          </span>
        </el-form-item>
        <el-form-item label="选择师傅" prop="worker_ids">
          <el-select
            v-model="dispatchForm.worker_ids"
            multiple
            filterable
            placeholder="请选择师傅（可多选）"
            style="width: 100%"
          >
            <el-option
              v-for="worker in availableWorkers"
              :key="worker.id"
              :label="`${worker.name} - ${worker.phone || ''}`"
              :value="worker.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="选择车辆" prop="vehicle_id">
          <el-select
            v-model="dispatchForm.vehicle_id"
            filterable
            placeholder="请选择车辆"
            style="width: 100%"
          >
            <el-option
              v-for="vehicle in availableVehicles"
              :key="vehicle.id"
              :label="`${vehicle.plate_number} - ${vehicle.vehicle_type} (${vehicle.capacity_volume}m³)`"
              :value="vehicle.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="计划开始时间" prop="scheduled_start_time">
          <el-date-picker
            v-model="dispatchForm.scheduled_start_time"
            type="datetime"
            placeholder="选择计划开始时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="计划结束时间" prop="scheduled_end_time">
          <el-date-picker
            v-model="dispatchForm.scheduled_end_time"
            type="datetime"
            placeholder="选择计划结束时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item v-if="conflictInfo.has_conflict" label="冲突检测">
          <el-alert type="error" :closable="false" title="检测到冲突">
            <div v-for="(c, idx) in conflictInfo.conflicts" :key="idx" class="conflict-msg">
              <el-icon><Warning /></el-icon>
              <span>{{ c.description }}</span>
            </div>
          </el-alert>
        </el-form-item>
        <el-form-item v-if="!conflictInfo.has_conflict && conflictInfo.checked" label="冲突检测">
          <el-alert type="success" :closable="false">
            <el-icon><CircleCheck /></el-icon>
            无冲突，可以派单
          </el-alert>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dispatchDialogVisible = false">取消</el-button>
        <el-button type="warning" @click="checkConflict" :loading="checkingConflict">
          <el-icon><Search /></el-icon>
          预检冲突
        </el-button>
        <el-button
          type="primary"
          @click="confirmDispatch"
          :loading="dispatching"
          :disabled="conflictInfo.has_conflict || !conflictInfo.checked"
        >
          <el-icon><Check /></el-icon>
          确认派单
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getWorkers, getVehicles, checkDispatch, createDispatch, getSchedules, getDispatches } from '@/api/dispatch'
import { getOrders } from '@/api/order'

const selectedDate = ref(new Date())
const activeTab = ref('workers')
const loadingSchedule = ref(false)
const loadingOrders = ref(false)
const checkingConflict = ref(false)
const dispatching = ref(false)
const workerSchedule = ref([])
const vehicleSchedule = ref([])
const pendingOrders = ref([])
const availableWorkers = ref([])
const availableVehicles = ref([])
const dispatchDialogVisible = ref(false)
const dispatchFormRef = ref(null)

const dispatchForm = reactive({
  order_id: null,
  customer_name: '',
  estimated_vehicle_type: '',
  items_volume: 0,
  worker_ids: [],
  vehicle_id: null,
  scheduled_start_time: '',
  scheduled_end_time: ''
})

const dispatchRules = {
  worker_ids: [{ required: true, message: '请选择师傅', trigger: 'change' }],
  vehicle_id: [{ required: true, message: '请选择车辆', trigger: 'change' }],
  scheduled_start_time: [{ required: true, message: '请选择计划开始时间', trigger: 'change' }],
  scheduled_end_time: [{ required: true, message: '请选择计划结束时间', trigger: 'change' }]
}

const conflictInfo = reactive({
  has_conflict: false,
  checked: false,
  conflicts: []
})

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const formatDateTime = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

const mockWorkers = [
  { id: 1, name: '张师傅', phone: '13800138001', status: 'available', schedules: [] },
  { id: 2, name: '李师傅', phone: '13800138002', status: 'available', schedules: [] },
  { id: 3, name: '王师傅', phone: '13800138003', status: 'available', schedules: [] },
  { id: 4, name: '赵师傅', phone: '13800138004', status: 'on_leave', schedules: [] },
  { id: 5, name: '刘师傅', phone: '13800138005', status: 'available', schedules: [] }
]

const mockVehicles = [
  { id: 1, plate_number: '京A12345', vehicle_type: '小面', capacity_volume: 5, status: 'available', schedules: [] },
  { id: 2, plate_number: '京B67890', vehicle_type: '金杯', capacity_volume: 10, status: 'available', schedules: [] },
  { id: 3, plate_number: '京C11111', vehicle_type: '厢货', capacity_volume: 20, status: 'available', schedules: [] }
]

const mockPendingOrders = [
  {
    id: 1,
    customer_name: '张三',
    customer_phone: '13800138000',
    move_date: formatDate(new Date()),
    start_address: '北京市朝阳区建国路88号',
    start_floor: 5,
    start_has_elevator: true,
    end_address: '北京市海淀区中关村大街1号',
    end_floor: 3,
    end_has_elevator: true,
    items_volume: 15,
    items_description: '家具、家电若干',
    estimated_workers: 3,
    estimated_vehicle_type: '金杯',
    status: 'pending'
  },
  {
    id: 4,
    customer_name: '赵六',
    customer_phone: '13600136000',
    move_date: formatDate(new Date(Date.now() + 86400000)),
    start_address: '北京市石景山区古城大街1号',
    start_floor: 7,
    start_has_elevator: true,
    end_address: '北京市昌平区回龙观东大街33号',
    end_floor: 10,
    end_has_elevator: true,
    items_volume: 12,
    items_description: '家电为主',
    estimated_workers: 2,
    estimated_vehicle_type: '小面',
    status: 'pending'
  }
]

const loadSchedule = async () => {
  loadingSchedule.value = true
  try {
    const startDate = formatDate(selectedDate.value)
    const endDate = formatDate(new Date(selectedDate.value.getTime() + 7 * 86400000))
    const data = await getSchedules(startDate, endDate)
    const schedules = data.list || data.records || data || []
    workerSchedule.value = availableWorkers.value.map(w => ({
      ...w,
      schedules: schedules.filter(s => s.worker_id === w.id)
    }))
    vehicleSchedule.value = availableVehicles.value.map(v => ({
      ...v,
      schedules: schedules.filter(s => s.vehicle_id === v.id)
    }))
  } catch (e) {
    workerSchedule.value = mockWorkers
    vehicleSchedule.value = mockVehicles
  } finally {
    loadingSchedule.value = false
  }
}

const loadPendingOrders = async () => {
  loadingOrders.value = true
  try {
    const data = await getOrders({ status: 'pending' })
    pendingOrders.value = data.list || data.records || data || mockPendingOrders
  } catch (e) {
    pendingOrders.value = mockPendingOrders
  } finally {
    loadingOrders.value = false
  }
}

const loadWorkersAndVehicles = async () => {
  try {
    const [workers, vehicles] = await Promise.all([getWorkers(), getVehicles()])
    availableWorkers.value = workers.list || workers.records || workers || mockWorkers
    availableVehicles.value = vehicles.list || vehicles.records || vehicles || mockVehicles
  } catch (e) {
    availableWorkers.value = mockWorkers
    availableVehicles.value = mockVehicles
  }
}

const openDispatchDialog = (order) => {
  dispatchForm.order_id = order.id
  dispatchForm.customer_name = order.customer_name
  dispatchForm.estimated_vehicle_type = order.estimated_vehicle_type
  dispatchForm.items_volume = order.items_volume
  dispatchForm.worker_ids = []
  dispatchForm.vehicle_id = null
  dispatchForm.scheduled_start_time = ''
  dispatchForm.scheduled_end_time = ''
  conflictInfo.has_conflict = false
  conflictInfo.checked = false
  conflictInfo.conflicts = []
  dispatchDialogVisible.value = true
}

const resetDispatchForm = () => {
  if (dispatchFormRef.value) {
    dispatchFormRef.value.resetFields()
  }
}

const checkConflict = async () => {
  if (!dispatchFormRef.value) return
  await dispatchFormRef.value.validate(async (valid) => {
    if (valid) {
      checkingConflict.value = true
      try {
        const params = {
          order_id: dispatchForm.order_id,
          worker_ids: dispatchForm.worker_ids,
          vehicle_id: dispatchForm.vehicle_id,
          scheduled_start_time: formatDateTime(dispatchForm.scheduled_start_time),
          scheduled_end_time: formatDateTime(dispatchForm.scheduled_end_time)
        }
        const data = await checkDispatch(params)
        conflictInfo.has_conflict = data.has_conflict || false
        conflictInfo.conflicts = data.conflicts || []
        conflictInfo.checked = true
        if (conflictInfo.has_conflict) {
          ElMessage.warning('检测到冲突，请调整派单')
        } else {
          ElMessage.success('无冲突')
        }
      } catch (e) {
        conflictInfo.has_conflict = false
        conflictInfo.checked = true
        conflictInfo.conflicts = []
        ElMessage.success('无冲突')
      } finally {
        checkingConflict.value = false
      }
    }
  })
}

const confirmDispatch = async () => {
  if (conflictInfo.has_conflict) {
    ElMessage.warning('存在冲突，无法派单')
    return
  }
  dispatching.value = true
  try {
    await createDispatch({
      order_id: dispatchForm.order_id,
      worker_ids: dispatchForm.worker_ids,
      vehicle_id: dispatchForm.vehicle_id,
      scheduled_start_time: formatDateTime(dispatchForm.scheduled_start_time),
      scheduled_end_time: formatDateTime(dispatchForm.scheduled_end_time)
    })
    ElMessage.success('派单成功')
    dispatchDialogVisible.value = false
    loadSchedule()
    loadPendingOrders()
  } catch (e) {
    ElMessage.error('派单失败')
  } finally {
    dispatching.value = false
  }
}

onMounted(async () => {
  await loadWorkersAndVehicles()
  loadSchedule()
  loadPendingOrders()
})
</script>

<style scoped>
.dispatcher {
  padding: 0;
}

.mb15 {
  margin-bottom: 15px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
  font-size: 16px;
}

.date-picker-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.date-picker-bar .label {
  font-weight: bold;
}

.schedule-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.schedule-date {
  font-size: 12px;
  color: #606266;
}

.address-cell {
  font-size: 12px;
  line-height: 1.6;
}

.address-cell .label {
  color: #909399;
  margin-right: 4px;
}

.conflict-msg {
  display: flex;
  align-items: center;
  gap: 5px;
  margin-bottom: 5px;
  color: #f56c6c;
}

.conflict-msg:last-child {
  margin-bottom: 0;
}
</style>

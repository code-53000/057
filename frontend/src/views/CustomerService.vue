<template>
  <div class="customer-service">
    <el-row :gutter="20">
      <el-col :span="10">
        <el-card>
          <template #header>
            <div class="card-header">
              <el-icon><EditPen /></el-icon>
              <span>订单录入</span>
            </div>
          </template>
          <el-form :model="orderForm" :rules="formRules" ref="orderFormRef" label-width="100px">
            <el-form-item label="客户姓名" prop="customer_name">
              <el-input v-model="orderForm.customer_name" placeholder="请输入客户姓名" />
            </el-form-item>
            <el-form-item label="联系电话" prop="customer_phone">
              <el-input v-model="orderForm.customer_phone" placeholder="请输入联系电话" />
            </el-form-item>
            <el-form-item label="搬家日期" prop="move_date">
              <el-date-picker
                v-model="orderForm.move_date"
                type="date"
                placeholder="选择搬家日期"
                style="width: 100%"
              />
            </el-form-item>
            <el-divider>起始地址</el-divider>
            <el-form-item label="起始地址" prop="start_address">
              <el-input v-model="orderForm.start_address" placeholder="请输入起始详细地址" />
            </el-form-item>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-form-item label="楼层" prop="start_floor">
                  <el-input-number v-model="orderForm.start_floor" :min="1" :max="50" style="width: 100%" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="有无电梯" prop="start_has_elevator">
                  <el-switch v-model="orderForm.start_has_elevator" active-text="有" inactive-text="无" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-divider>目的地址</el-divider>
            <el-form-item label="目的地址" prop="end_address">
              <el-input v-model="orderForm.end_address" placeholder="请输入目的详细地址" />
            </el-form-item>
            <el-row :gutter="10">
              <el-col :span="12">
                <el-form-item label="楼层" prop="end_floor">
                  <el-input-number v-model="orderForm.end_floor" :min="1" :max="50" style="width: 100%" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="有无电梯" prop="end_has_elevator">
                  <el-switch v-model="orderForm.end_has_elevator" active-text="有" inactive-text="无" />
                </el-form-item>
              </el-col>
            </el-row>
            <el-divider>物品信息</el-divider>
            <el-form-item label="体积(立方米)" prop="items_volume">
              <el-input-number v-model="orderForm.items_volume" :min="0" :max="100" :step="0.5" style="width: 100%" />
            </el-form-item>
            <el-form-item label="物品描述" prop="items_description">
              <el-input
                v-model="orderForm.items_description"
                type="textarea"
                :rows="3"
                placeholder="请描述主要物品，如家具、家电等"
              />
            </el-form-item>
            <el-form-item label="预估人数" prop="estimated_workers">
              <el-input-number v-model="orderForm.estimated_workers" :min="1" :max="20" style="width: 100%" />
            </el-form-item>
            <el-form-item label="预估车型" prop="estimated_vehicle_type">
              <el-select v-model="orderForm.estimated_vehicle_type" placeholder="请选择车型" style="width: 100%">
                <el-option label="小面" value="小面" />
                <el-option label="金杯" value="金杯" />
                <el-option label="厢货" value="厢货" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitOrder" :loading="submitting">
                <el-icon><Check /></el-icon>
                提交订单
              </el-button>
              <el-button @click="resetForm">
                <el-icon><RefreshLeft /></el-icon>
                重置
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      <el-col :span="14">
        <el-card>
          <template #header>
            <div class="card-header">
              <el-icon><List /></el-icon>
              <span>订单列表</span>
            </div>
          </template>
          <div class="filter-bar">
            <el-select v-model="filterStatus" placeholder="状态筛选" clearable style="width: 150px; margin-right: 10px">
              <el-option label="待派单" value="pending" />
              <el-option label="已派单" value="dispatched" />
              <el-option label="作业中" value="in_progress" />
              <el-option label="已完成" value="completed" />
              <el-option label="已取消" value="cancelled" />
            </el-select>
            <el-date-picker
              v-model="filterDate"
              type="date"
              placeholder="日期筛选"
              clearable
              style="width: 180px; margin-right: 10px"
            />
            <el-button type="primary" @click="loadOrders" :loading="loading">
              <el-icon><Search /></el-icon>
              查询
            </el-button>
            <el-button @click="resetFilter">
              <el-icon><RefreshLeft /></el-icon>
              重置
            </el-button>
          </div>
          <el-table :data="orderList" style="width: 100%" v-loading="loading" stripe>
            <el-table-column prop="id" label="订单号" width="80" />
            <el-table-column prop="customer_name" label="客户" width="100" />
            <el-table-column prop="customer_phone" label="电话" width="130" />
            <el-table-column prop="move_date" label="搬家日期" width="120">
              <template #default="scope">
                {{ formatDate(scope.row.move_date) }}
              </template>
            </el-table-column>
            <el-table-column label="地址" min-width="200">
              <template #default="scope">
                <div class="address-cell">
                  <div><span class="label">起:</span> {{ scope.row.start_address }} ({{ scope.row.start_floor }}楼{{ scope.row.start_has_elevator ? '有梯' : '无梯' }})</div>
                  <div><span class="label">止:</span> {{ scope.row.end_address }} ({{ scope.row.end_floor }}楼{{ scope.row.end_has_elevator ? '有梯' : '无梯' }})</div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="items_volume" label="体积" width="70" />
            <el-table-column prop="estimated_vehicle_type" label="车型" width="80" />
            <el-table-column prop="status" label="状态" width="90">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ getStatusText(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="scope">
                <el-button type="primary" link size="small" @click="viewOrder(scope.row)">查看</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            class="pagination"
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="loadOrders"
            @current-change="loadOrders"
          />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createOrder, getOrders } from '@/api/order'

const orderFormRef = ref(null)
const submitting = ref(false)
const loading = ref(false)
const orderList = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const filterStatus = ref('')
const filterDate = ref('')

const orderForm = reactive({
  customer_name: '',
  customer_phone: '',
  move_date: '',
  start_address: '',
  start_floor: 1,
  start_has_elevator: false,
  end_address: '',
  end_floor: 1,
  end_has_elevator: false,
  items_volume: 0,
  items_description: '',
  estimated_workers: 2,
  estimated_vehicle_type: ''
})

const formRules = {
  customer_name: [{ required: true, message: '请输入客户姓名', trigger: 'blur' }],
  customer_phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  move_date: [{ required: true, message: '请选择搬家日期', trigger: 'change' }],
  start_address: [{ required: true, message: '请输入起始地址', trigger: 'blur' }],
  end_address: [{ required: true, message: '请输入目的地址', trigger: 'blur' }],
  estimated_vehicle_type: [{ required: true, message: '请选择车型', trigger: 'change' }]
}

const statusMap = {
  pending: { text: '待派单', type: 'info' },
  dispatched: { text: '已派单', type: 'warning' },
  in_progress: { text: '作业中', type: 'primary' },
  completed: { text: '已完成', type: 'success' },
  cancelled: { text: '已取消', type: 'danger' }
}

const getStatusText = (status) => statusMap[status]?.text || status
const getStatusType = (status) => statusMap[status]?.type || 'info'

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const submitOrder = async () => {
  if (!orderFormRef.value) return
  await orderFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        await createOrder(orderForm)
        ElMessage.success('订单创建成功')
        resetForm()
        loadOrders()
      } catch (e) {
        ElMessage.error('订单创建失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

const resetForm = () => {
  if (orderFormRef.value) {
    orderFormRef.value.resetFields()
  }
  orderForm.start_floor = 1
  orderForm.start_has_elevator = false
  orderForm.end_floor = 1
  orderForm.end_has_elevator = false
  orderForm.items_volume = 0
  orderForm.estimated_workers = 2
}

const resetFilter = () => {
  filterStatus.value = ''
  filterDate.value = ''
  currentPage.value = 1
  loadOrders()
}

const loadOrders = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      size: pageSize.value
    }
    if (filterStatus.value) params.status = filterStatus.value
    if (filterDate.value) params.move_date = formatDate(filterDate.value)
    const data = await getOrders(params)
    orderList.value = data.list || data.records || data || []
    total.value = data.total || orderList.value.length
  } catch (e) {
    orderList.value = mockOrders
    total.value = mockOrders.length
  } finally {
    loading.value = false
  }
}

const viewOrder = (order) => {
  ElMessageBox.alert(
    `<div style="text-align: left">
      <p><strong>订单号:</strong> ${order.id}</p>
      <p><strong>客户:</strong> ${order.customer_name} (${order.customer_phone})</p>
      <p><strong>搬家日期:</strong> ${formatDate(order.move_date)}</p>
      <p><strong>起始地址:</strong> ${order.start_address} (${order.start_floor}楼${order.start_has_elevator ? '有电梯' : '无电梯'})</p>
      <p><strong>目的地址:</strong> ${order.end_address} (${order.end_floor}楼${order.end_has_elevator ? '有电梯' : '无电梯'})</p>
      <p><strong>体积:</strong> ${order.items_volume} m³</p>
      <p><strong>物品描述:</strong> ${order.items_description || '-'}</p>
      <p><strong>预估人数:</strong> ${order.estimated_workers}人</p>
      <p><strong>预估车型:</strong> ${order.estimated_vehicle_type}</p>
      <p><strong>状态:</strong> ${getStatusText(order.status)}</p>
    </div>`,
    '订单详情',
    { dangerouslyUseHTMLString: true, confirmButtonText: '确定' }
  )
}

const mockOrders = [
  {
    id: 1,
    customer_name: '张三',
    customer_phone: '13800138000',
    move_date: '2025-06-25',
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
    id: 2,
    customer_name: '李四',
    customer_phone: '13900139000',
    move_date: '2025-06-26',
    start_address: '北京市西城区金融街15号',
    start_floor: 12,
    start_has_elevator: true,
    end_address: '北京市东城区王府井大街88号',
    end_floor: 6,
    end_has_elevator: false,
    items_volume: 25,
    items_description: '钢琴、衣柜、沙发等大件',
    estimated_workers: 5,
    estimated_vehicle_type: '厢货',
    status: 'dispatched'
  },
  {
    id: 3,
    customer_name: '王五',
    customer_phone: '13700137000',
    move_date: '2025-06-22',
    start_address: '北京市丰台区南三环西路16号',
    start_floor: 2,
    start_has_elevator: false,
    end_address: '北京市通州区新华大街99号',
    end_floor: 8,
    end_has_elevator: true,
    items_volume: 8,
    items_description: '少量包裹和小家电',
    estimated_workers: 2,
    estimated_vehicle_type: '小面',
    status: 'completed'
  }
]

onMounted(() => {
  loadOrders()
})
</script>

<style scoped>
.customer-service {
  padding: 0;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: bold;
  font-size: 16px;
}

.filter-bar {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
}

.address-cell {
  font-size: 12px;
  line-height: 1.6;
}

.address-cell .label {
  color: #909399;
  margin-right: 4px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>

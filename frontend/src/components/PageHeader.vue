<template>
  <div class="page-header">
    <el-row :gutter="20" align="middle">
      <el-col :span="12">
        <div class="header-left">
          <el-icon v-if="icon"><component :is="icon" /></el-icon>
          <h2 class="page-title">{{ title }}</h2>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="header-right">
          <slot name="extra" />
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    default: ''
  },
  icon: {
    type: String,
    default: ''
  }
})
</script>

<style scoped>
.page-header {
  padding: 16px 0;
  margin-bottom: 20px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.page-title {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #303133;
}

.header-right {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 10px;
}
</style>

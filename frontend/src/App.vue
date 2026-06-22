<template>
  <div id="app">
    <el-container>
      <el-header class="app-header">
        <div class="header-content">
          <h1 class="app-title">搬家管理系统</h1>
          <el-menu
            :default-active="activeMenu"
            mode="horizontal"
            router
            class="nav-menu"
            background-color="#409EFF"
            text-color="#fff"
            active-text-color="#ffd04b"
          >
            <el-menu-item index="/customer-service">
              <el-icon><Service /></el-icon>
              <span>客服端</span>
            </el-menu-item>
            <el-menu-item index="/dispatcher">
              <el-icon><Tickets /></el-icon>
              <span>调度端</span>
            </el-menu-item>
            <el-menu-item index="/worker">
              <el-icon><User /></el-icon>
              <span>师傅端</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style>
html, body, #app {
  margin: 0;
  padding: 0;
  height: 100%;
  width: 100%;
}

.app-header {
  background-color: #409EFF;
  padding: 0;
  display: flex;
  align-items: center;
}

.header-content {
  display: flex;
  align-items: center;
  width: 100%;
}

.app-title {
  color: #fff;
  margin: 0 40px 0 20px;
  font-size: 20px;
  font-weight: bold;
}

.nav-menu {
  flex: 1;
  border-bottom: none;
}

.app-main {
  background-color: #f5f7fa;
  min-height: calc(100vh - 60px);
}
</style>

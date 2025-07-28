<template>
  <div class="orchestration-list">
    <el-card>
      <div class="header">
        <el-button type="primary" @click="onCreate">创建编排</el-button>
        <div class="actions">
          <el-input v-model="search" placeholder="搜索" size="small" clearable @input="onSearch" />
          <el-button icon="el-icon-refresh" @click="fetchList" size="small">刷新</el-button>
        </div>
      </div>
      <el-table :data="list" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="source" label="来源" />
        <el-table-column prop="dir" label="编排目录">
          <template #default="{ row }">
            <el-link type="primary" :underline="false">{{ row.dir }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="容器状态" />
        <el-table-column prop="count" label="容器数量" />
        <el-table-column prop="createdAt" label="创建时间" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-link type="primary" @click="onEdit(row)">编辑</el-link>
            <el-divider direction="vertical" />
            <el-link type="danger" @click="onDelete(row)">删除</el-link>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="fetchList"
        @size-change="fetchList"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getOrchestrationList, deleteOrchestration } from '@/api/orchestration'

const list = ref([])
const loading = ref(false)
const search = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getOrchestrationList({ search: search.value, page: page.value, pageSize: pageSize.value })
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    list.value = []
    total.value = 0
  }
  loading.value = false
}

const onSearch = () => {
  page.value = 1
  fetchList()
}
const onCreate = () => {
  // 跳转到创建页面
}
const onEdit = (row) => {
  // 跳转到编辑页面
}
const onDelete = async (row) => {
  try {
    await deleteOrchestration(row.id)
    fetchList()
  } catch (e) {}
}

onMounted(fetchList)
</script>

<style scoped>
.orchestration-list .header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.orchestration-list .actions {
  display: flex;
  gap: 8px;
}
</style>

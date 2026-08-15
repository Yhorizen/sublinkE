<template>
  <el-card class="page-card">
    <template #header>
      <div class="header-row">
        <span>邀请码管理</span>
        <el-button type="primary" @click="openCreate">创建邀请码</el-button>
      </div>
    </template>
    <el-table :data="invites" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="code" label="邀请码" />
      <el-table-column prop="description" label="说明" />
      <el-table-column label="使用次数" width="110">
        <template #default="scope">
          {{ scope.row.usedCount }}{{ scope.row.maxUses > 0 ? ` / ${scope.row.maxUses}` : '' }}
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="170" />
      <el-table-column label="状态" width="90">
        <template #default="scope">
          <el-switch :model-value="scope.row.enabled" @change="(val) => toggleInvite(scope.row, Boolean(val))" />
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="scope">
          <el-button type="danger" link @click="removeInvite(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>

  <el-dialog v-model="dialogVisible" title="创建邀请码" width="420px">
    <el-form :model="form" label-width="90px">
      <el-form-item label="邀请码">
        <el-input v-model="form.code" placeholder="留空自动生成" />
      </el-form-item>
      <el-form-item label="说明">
        <el-input v-model="form.description" placeholder="备注(可选)" />
      </el-form-item>
      <el-form-item label="使用上限">
        <el-input-number v-model="form.maxUses" :min="0" />
        <span class="cfg-hint">0 表示不限次数</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="submitCreate">创建</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { addInvite, deleteInvite, getInvites, updateInvite, type InviteItem } from '@/api/admin'

const invites = ref<InviteItem[]>([])
const dialogVisible = ref(false)
const form = ref({ code: '', description: '', maxUses: 0 })

const fetchInvites = async () => {
  const res = await getInvites()
  invites.value = res.data || []
}

const openCreate = () => {
  form.value = { code: '', description: '', maxUses: 0 }
  dialogVisible.value = true
}

const submitCreate = async () => {
  await addInvite({
    code: form.value.code || undefined,
    description: form.value.description || undefined,
    maxUses: form.value.maxUses
  })
  ElMessage.success('创建成功')
  dialogVisible.value = false
  await fetchInvites()
}

const toggleInvite = async (row: InviteItem, enabled: boolean) => {
  await updateInvite({ id: row.id, enabled })
  ElMessage.success('更新成功')
  await fetchInvites()
}

const removeInvite = (row: InviteItem) => {
  ElMessageBox.confirm(`确认删除邀请码 ${row.code} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    await deleteInvite(row.id)
    ElMessage.success('删除成功')
    await fetchInvites()
  })
}

onMounted(fetchInvites)
</script>

<style scoped>
.page-card { margin: 10px; }
.header-row { display: flex; justify-content: space-between; align-items: center; }
.cfg-hint { margin-left: 8px; font-size: 12px; color: var(--el-text-color-secondary); }
</style>

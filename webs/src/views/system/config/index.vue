<template>
  <el-card class="page-card">
    <template #header>
      <div class="header-row">
        <span>注册配置</span>
      </div>
    </template>
    <el-form label-width="150px">
      <el-form-item label="默认订阅ID">
        <el-select v-model="defaultSubscriptionId" placeholder="请选择默认订阅">
          <el-option label="不分配" :value="0" />
          <el-option v-for="item in subscriptions" :key="item.ID" :label="`${item.ID} - ${item.Name}`" :value="item.ID" />
        </el-select>
      </el-form-item>
      <el-form-item label="启用邀请码注册">
        <el-switch v-model="inviteRequired" />
      </el-form-item>
      <el-form-item label="User-Agent检测">
        <el-switch v-model="uaCheckEnabled" />
        <span class="cfg-hint">开启后订阅拉取仅允许代理软件，浏览器/脚本直接访问将被拦截</span>
      </el-form-item>
      <el-form-item v-if="uaCheckEnabled" label="UA关键字">
        <el-input v-model="uaCheckKeywords" placeholder="逗号分隔，命中任一即视为代理软件" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="saveConfig">保存</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { getAdminConfig, setAdminConfig } from '@/api/admin'
import { getSubs } from '@/api/subcription/subs'

const defaultSubscriptionId = ref<number>(0)
const inviteRequired = ref(false)
const uaCheckEnabled = ref(false)
const uaCheckKeywords = ref('')
const subscriptions = ref<{ ID: number; Name: string }[]>([])

const fetchConfig = async () => {
  const res = await getAdminConfig()
  defaultSubscriptionId.value = Number(res.data.default_subscription_id || 0)
  inviteRequired.value = Boolean(res.data.invite_required)
  uaCheckEnabled.value = Boolean(res.data.ua_check_enabled)
  uaCheckKeywords.value = res.data.ua_check_keywords || ''
}

const fetchSubscriptions = async () => {
  const res = await getSubs()
  subscriptions.value = (res.data || []).map((item: any) => ({ ID: item.ID, Name: item.Name }))
}

const saveConfig = async () => {
  await setAdminConfig({
    defaultSubscriptionId: defaultSubscriptionId.value,
    inviteRequired: inviteRequired.value,
    uaCheckEnabled: uaCheckEnabled.value,
    uaCheckKeywords: uaCheckKeywords.value
  })
  ElMessage.success('保存成功')
}

onMounted(async () => {
  await fetchSubscriptions()
  await fetchConfig()
})
</script>

<style scoped>
.page-card { margin: 10px; }
.header-row { display: flex; justify-content: space-between; align-items: center; }
.cfg-hint { margin-left: 12px; font-size: 12px; color: var(--el-text-color-secondary); }
</style>

<template>
  <div class="red-packet-record">
    <ProTable
      ref="proTableRef"
      :columns="columns"
      :request-api="getRedPacketRecordList"
      :search-columns="searchColumns"
      :show-pagination="true"
    >
      <!-- 红包类型列 -->
      <template #packetType="{ row }">
        <el-tag :type="row.packetType === 1 ? '' : 'warning'">
          {{ row.packetType === 1 ? '普通红包' : '手气红包' }}
        </el-tag>
      </template>

      <!-- 状态列 -->
      <template #status="{ row }">
        <el-tag :type="getStatusType(row.status)">
          {{ getStatusText(row.status) }}
        </el-tag>
      </template>

      <!-- 进度列 -->
      <template #progress="{ row }">
        <span>{{ row.grabbedCount }} / {{ row.totalCount }}</span>
      </template>

      <!-- 金额列 -->
      <template #totalAmount="{ row }">
        {{ row.totalAmount }} {{ row.symbol }}
      </template>

      <!-- 已抢金额列 -->
      <template #grabbedAmount="{ row }">
        {{ row.grabbedAmount }} {{ row.symbol }}
      </template>

      <!-- 创建时间列 -->
      <template #createdAt="{ row }">
        {{ formatTimestamp(row.createdAt) }}
      </template>

      <!-- 操作列 -->
      <template #actions="{ row }">
        <el-button link type="primary" @click="handleView(row)">查看</el-button>
        <el-button link type="primary" @click="handleViewGrabRecords(row)">领取明细</el-button>
      </template>
    </ProTable>

    <!-- 查看红包详情对话框 -->
    <el-dialog v-model="viewDialogVisible" title="红包详情" width="650px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ viewData.id }}</el-descriptions-item>
        <el-descriptions-item label="红包编号">{{ viewData.packetNo }}</el-descriptions-item>
        <el-descriptions-item label="群组ID">{{ viewData.groupId }}</el-descriptions-item>
        <el-descriptions-item label="红包类型">
          <el-tag :type="viewData.packetType === 1 ? '' : 'warning'">
            {{ viewData.packetType === 1 ? '普通红包' : '手气红包' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="总金额">{{ viewData.totalAmount }} {{ viewData.symbol }}</el-descriptions-item>
        <el-descriptions-item label="总个数">{{ viewData.totalCount }}</el-descriptions-item>
        <el-descriptions-item label="已抢金额">{{ viewData.grabbedAmount }} {{ viewData.symbol }}</el-descriptions-item>
        <el-descriptions-item label="已抢个数">{{ viewData.grabbedCount }}</el-descriptions-item>
        <el-descriptions-item label="剩余金额">{{ viewData.remainAmount }} {{ viewData.symbol }}</el-descriptions-item>
        <el-descriptions-item label="剩余个数">{{ viewData.remainCount }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(viewData.status)">{{ getStatusText(viewData.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="消息ID">{{ viewData.messageId || '-' }}</el-descriptions-item>
        <el-descriptions-item label="祝福语" :span="2">{{ viewData.blessingWords || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ formatTimestamp(viewData.createdAt) }}</el-descriptions-item>
        <el-descriptions-item label="过期时间">{{ formatTime(viewData.expireAt) }}</el-descriptions-item>
        <el-descriptions-item label="完成时间">{{ formatTime(viewData.completedAt) }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <!-- 领取明细对话框 -->
    <el-dialog v-model="grabDialogVisible" title="领取明细" width="800px">
      <div style="margin-bottom: 10px; color: #909399;">
        红包编号: {{ currentPacketNo }} | 总额: {{ currentPacketAmount }} {{ currentPacketSymbol }}
      </div>
      <el-table :data="grabRecords" border style="width: 100%" v-loading="grabLoading">
        <el-table-column prop="sequence" label="序号" width="60" />
        <el-table-column prop="telegramUsername" label="Telegram用户名" min-width="140">
          <template #default="{ row }">
            {{ row.telegramUsername || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="userId" label="平台用户ID" width="110" />
        <el-table-column prop="amount" label="抢到金额" width="150">
          <template #default="{ row }">
            <span :style="{ color: row.isBest === 1 ? '#E6A23C' : '' , fontWeight: row.isBest === 1 ? 'bold' : '' }">
              {{ row.amount }}
            </span>
            <el-tag v-if="row.isBest === 1" type="warning" size="small" style="margin-left: 4px;">手气最佳</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="grabbedAt" label="抢红包时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.grabbedAt) }}
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import ProTable from '@/components/ProTable/index.vue';
import { getRedPacketRecordList, getGrabRecordList } from '@/api/modules/redpacket';

const proTableRef = ref();
const viewDialogVisible = ref(false);
const grabDialogVisible = ref(false);
const grabLoading = ref(false);
const viewData = ref({});
const grabRecords = ref([]);
const currentPacketNo = ref('');
const currentPacketAmount = ref('');
const currentPacketSymbol = ref('');

// 表格列配置
const columns = [
  { prop: 'id', label: 'ID', width: 70 },
  { prop: 'packetNo', label: '红包编号', minWidth: 160 },
  { prop: 'groupId', label: '群组ID', width: 140 },
  { prop: 'packetType', label: '红包类型', width: 100, slot: 'packetType' },
  { prop: 'totalAmount', label: '总金额', width: 150, slot: 'totalAmount' },
  { prop: 'progress', label: '领取进度', width: 110, slot: 'progress' },
  { prop: 'grabbedAmount', label: '已抢金额', width: 150, slot: 'grabbedAmount' },
  { prop: 'status', label: '状态', width: 90, slot: 'status' },
  { prop: 'createdAt', label: '创建时间', width: 180, slot: 'createdAt' }
];

// 搜索列配置
const searchColumns = [
  { prop: 'packet_no', label: '红包编号', type: 'input' },
  { prop: 'group_id', label: '群组ID', type: 'input' },
  {
    prop: 'status',
    label: '状态',
    type: 'select',
    options: [
      { label: '进行中', value: 1 },
      { label: '已抢完', value: 2 },
      { label: '已过期', value: 3 }
    ]
  }
];

// 状态文本
const getStatusText = status => {
  const map = { 1: '进行中', 2: '已抢完', 3: '已过期' };
  return map[status] || '未知';
};

// 状态类型
const getStatusType = status => {
  const map = { 1: 'success', 2: 'info', 3: 'danger' };
  return map[status] || '';
};

// 格式化时间戳（秒级）
const formatTimestamp = ts => {
  if (!ts) return '-';
  return new Date(ts * 1000).toLocaleString('zh-CN');
};

// 格式化时间（ISO 字符串）
const formatTime = t => {
  if (!t) return '-';
  const d = new Date(t);
  if (isNaN(d.getTime())) return '-';
  return d.toLocaleString('zh-CN');
};

// 查看详情
function handleView(row) {
  viewData.value = { ...row };
  viewDialogVisible.value = true;
}

// 查看领取明细
async function handleViewGrabRecords(row) {
  currentPacketNo.value = row.packetNo;
  currentPacketAmount.value = row.totalAmount;
  currentPacketSymbol.value = row.symbol;
  grabRecords.value = [];
  grabDialogVisible.value = true;
  grabLoading.value = true;

  try {
    const res = await getGrabRecordList({ packet_id: row.id, pageSize: 100 });
    grabRecords.value = res.data?.list || res.list || res.data || [];
  } catch (e) {
    console.error('查询领取明细失败', e);
  } finally {
    grabLoading.value = false;
  }
}
</script>

<style scoped lang="scss">
.red-packet-record {
  padding: 20px;
}
</style>

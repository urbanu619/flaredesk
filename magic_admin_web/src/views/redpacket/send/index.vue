<template>
  <div class="red-packet-send">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>手动发送红包</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="120px"
        style="max-width: 600px"
      >
        <el-form-item label="群组选择" prop="groupId">
          <el-select
            v-model="formData.groupId"
            placeholder="请选择Telegram群组"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="group in groupList"
              :key="group.chatId"
              :label="group.title"
              :value="group.chatId"
            />
          </el-select>
          <div class="form-tip">如果没有群组，请先将机器人添加到Telegram群组中</div>
        </el-form-item>

        <el-form-item label="群组名称" prop="groupName">
          <el-input v-model="formData.groupName" placeholder="请输入群组名称（选填）" />
        </el-form-item>

        <el-form-item label="红包类型" prop="packetType">
          <el-radio-group v-model="formData.packetType">
            <el-radio :label="1">
              <span>普通红包</span>
              <span class="radio-desc">（每人金额相同）</span>
            </el-radio>
            <el-radio :label="2">
              <span>手气红包</span>
              <span class="radio-desc">（每人金额随机）</span>
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="消息语言" prop="lang">
          <el-select v-model="formData.lang" style="width: 200px">
            <el-option label="🇻🇳 Tiếng Việt (越南语)" value="vi" />
            <el-option label="🇮🇩 Bahasa Indonesia (印尼语)" value="id" />
            <el-option label="🇺🇸 English (英语)" value="en" />
            <el-option label="🇨🇳 中文" value="zh" />
          </el-select>
          <div class="form-tip">发送到群组的消息将使用所选语言</div>
        </el-form-item>

        <el-form-item label="红包金额" prop="totalAmount">
          <el-input-number
            v-model="formData.totalAmount"
            :min="0"
            :precision="2"
            :step="1000"
            placeholder="请输入红包总金额"
            style="width: 200px"
          />
          <span class="ml-2">{{ formData.symbol }}</span>
          <div v-if="formData.totalCount > 0 && formData.totalAmount > 0" class="form-tip">
            {{ formData.packetType === 1 ? '每人固定' : '每人平均' }}：{{
              (formData.totalAmount / formData.totalCount).toFixed(2)
            }}
            {{ formData.symbol }}
          </div>
        </el-form-item>

        <el-form-item label="红包个数" prop="totalCount">
          <el-input-number
            v-model="formData.totalCount"
            :min="1"
            :max="100"
            placeholder="请输入红包个数"
            style="width: 200px"
          />
          <span class="ml-2">个</span>
        </el-form-item>

        <el-form-item label="过期时间" prop="expireMinutes">
          <el-input-number
            v-model="formData.expireMinutes"
            :min="1"
            :max="1440"
            placeholder="红包过期时间"
            style="width: 200px"
          />
          <span class="ml-2">分钟</span>
          <div class="form-tip">红包发出后多久过期，默认10分钟</div>
        </el-form-item>

        <el-form-item v-if="formData.packetType === 2" label="单人最大金额" prop="maxGrabAmount">
          <el-input-number
            v-model="formData.maxGrabAmount"
            :min="0"
            :precision="2"
            :step="1000"
            placeholder="0=不限制"
            style="width: 200px"
          />
          <span class="ml-2">{{ formData.symbol }}（0=不限制）</span>
          <div class="form-tip">手气红包每人最多可领取的金额，0表示不限制</div>
        </el-form-item>

        <el-form-item label="祝福语" prop="blessingWords">
          <el-input
            v-model="formData.blessingWords"
            type="textarea"
            :rows="4"
            placeholder="请输入祝福语，例如：恭喜发财，大吉大利！"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="币种" prop="symbol">
          <el-select v-model="formData.symbol" style="width: 200px">
            <el-option label="VND (越南盾)" value="VND" />
            <el-option label="IDR (印尼盾)" value="IDR" />
            <el-option label="USD (美元)" value="USD" />
            <el-option label="CNY (人民币)" value="CNY" />
          </el-select>
          <div class="form-tip">语言切换时自动联动币种，也可手动选择</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            立即发送
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- Telegram 消息预览 -->
      <el-divider />
      <div class="preview-box">
        <h3>📱 Telegram 消息预览</h3>
        <div class="telegram-preview">
          <div class="telegram-message">
            <div class="message-header">
              <span class="sender-name">🤖 {{ t.adminName }}</span>
            </div>
            <div class="message-content">
              <div class="red-packet-icon">🧧</div>
              <div class="red-packet-text">{{ t.sendTitle }}</div>
              <div class="blessing-words">💰 {{ formData.blessingWords || (formData.packetType === 1 ? t.defaultBlessing1 : t.defaultBlessing2) }}</div>
              <div class="divider"></div>
              <div class="packet-info">
                <div class="info-row">
                  <span class="info-label">{{ t.totalAmount }}:</span>
                  <span class="info-value">{{ formData.totalAmount.toFixed(2) }} {{ formData.symbol }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ t.count }}:</span>
                  <span class="info-value">{{ formData.totalCount }} {{ t.countUnit }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ t.type }}:</span>
                  <span class="info-value">{{ formData.packetType === 1 ? t.normalPacket : t.luckyPacket }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">{{ t.validFor }}:</span>
                  <span class="info-value">{{ formData.expireMinutes }} {{ t.minutes }}</span>
                </div>
                <div v-if="formData.packetType === 2 && formData.maxGrabAmount > 0" class="info-row">
                  <span class="info-label">{{ t.maxGrab }}:</span>
                  <span class="info-value">{{ formData.maxGrabAmount }} {{ formData.symbol }}</span>
                </div>
              </div>
              <div class="call-to-action">{{ t.comeGrab }}</div>
              <div class="grab-button">
                <el-button type="danger" size="large" style="width: 100%">
                  {{ t.grabBtn }}
                </el-button>
              </div>
            </div>
          </div>
          <div class="preview-tip">
            {{ t.previewTip }}
          </div>
        </div>
      </div>
    </el-card>

    <!-- 发送历史 -->
    <el-card class="box-card mt-4">
      <template #header>
        <div class="card-header">
          <span>最近发送记录</span>
        </div>
      </template>

      <el-table :data="recentRecords" border stripe>
        <el-table-column prop="packetNo" label="红包编号" width="180" />
        <el-table-column prop="groupId" label="群组" />
        <el-table-column prop="totalAmount" label="金额" />
        <el-table-column prop="totalCount" label="个数" />
        <el-table-column prop="grabbedCount" label="已抢" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1" type="success">进行中</el-tag>
            <el-tag v-else-if="row.status === 2" type="info">已抢完</el-tag>
            <el-tag v-else type="warning">已过期</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="发送时间">
          <template #default="{ row }">
            {{ formatTimestamp(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { sendRedPacketDirect, getTelegramGroups, getRedPacketRecordList } from '@/api/modules/redpacket';

// 语言 → 默认币种映射
const langSymbolMap = { vi: 'VND', id: 'IDR', en: 'USD', zh: 'CNY' };

// 表单实例
const formRef = ref();

// 加载状态
const loading = ref(false);

// 群组列表
const groupList = ref([
  // 示例数据，实际应该从API获取
  { id: '-1001234567890', name: '测试群组1' },
  { id: '-1009876543210', name: '测试群组2' }
]);

// 最近发送记录
const recentRecords = ref([]);

// 表单数据
const formData = reactive({
  groupId: '',
  groupName: '',
  packetType: 2,
  totalAmount: 10000,
  totalCount: 10,
  symbol: 'VND',
  lang: 'vi',
  expireMinutes: 10,
  maxGrabAmount: 0,
  blessingWords: '恭喜发财，大吉大利！'
});

// 多语言预览文本
const langTexts = {
  vi: {
    adminName: 'Quản trị viên',
    sendTitle: 'đã gửi một lì xì',
    defaultBlessing1: 'Chúc mừng phát tài, đại cát đại lợi',
    defaultBlessing2: 'Lì xì may mắn, nhanh tay nhận nào!',
    totalAmount: 'Tổng số tiền',
    count: 'Số lượng',
    countUnit: 'cái',
    type: 'Loại',
    normalPacket: 'Lì xì thường',
    luckyPacket: 'Lì xì may mắn',
    validFor: 'Thời hạn',
    minutes: 'phút',
    maxGrab: 'Giới hạn/người',
    comeGrab: 'Nhanh tay nhận lì xì nào!',
    grabBtn: '🎁 Nhận lì xì',
    previewTip: '💡 Đây là bản xem trước, tin nhắn thực tế gửi đến nhóm Telegram có thể hơi khác'
  },
  id: {
    adminName: 'Admin',
    sendTitle: 'mengirim angpao',
    defaultBlessing1: 'Gong Xi Fa Cai, semoga beruntung',
    defaultBlessing2: 'Angpao keberuntungan, ayo ambil!',
    totalAmount: 'Total',
    count: 'Jumlah',
    countUnit: 'buah',
    type: 'Tipe',
    normalPacket: 'Angpao biasa',
    luckyPacket: 'Angpao keberuntungan',
    validFor: 'Berlaku',
    minutes: 'menit',
    maxGrab: 'Maks/orang',
    comeGrab: 'Ayo ambil angpao!',
    grabBtn: '🎁 Ambil angpao',
    previewTip: '💡 Ini adalah pratinjau, pesan sebenarnya yang dikirim ke grup Telegram mungkin sedikit berbeda'
  },
  en: {
    adminName: 'System Admin',
    sendTitle: 'sent a red packet',
    defaultBlessing1: 'Wishing you prosperity and good fortune',
    defaultBlessing2: 'Lucky red packet, grab it now!',
    totalAmount: 'Total',
    count: 'Count',
    countUnit: '',
    type: 'Type',
    normalPacket: 'Fixed red packet',
    luckyPacket: 'Lucky red packet',
    validFor: 'Valid for',
    minutes: 'min',
    maxGrab: 'Max/person',
    comeGrab: 'Grab the red packet now!',
    grabBtn: '🎁 Grab',
    previewTip: '💡 This is a preview, the actual message sent to Telegram group may look slightly different'
  },
  zh: {
    adminName: '系统管理员',
    sendTitle: '发了一个红包',
    defaultBlessing1: '恭喜发财，大吉大利',
    defaultBlessing2: '拼手气红包，快来抢！',
    totalAmount: '总金额',
    count: '个数',
    countUnit: '个',
    type: '类型',
    normalPacket: '普通红包',
    luckyPacket: '手气红包',
    validFor: '有效期',
    minutes: '分钟',
    maxGrab: '单人上限',
    comeGrab: '快来抢红包吧！',
    grabBtn: '🎁 抢红包',
    previewTip: '💡 这是预览效果，实际发送到 Telegram 群组的消息样式可能略有不同'
  }
};

const t = computed(() => langTexts[formData.lang] || langTexts.vi);

// 语言切换时自动联动币种
watch(() => formData.lang, (lang) => {
  formData.symbol = langSymbolMap[lang] || 'VND';
});

// 表单验证规则
const rules = {
  groupId: [{ required: true, message: '请选择群组', trigger: 'change' }],
  packetType: [{ required: true, message: '请选择红包类型', trigger: 'change' }],
  totalAmount: [
    { required: true, message: '请输入红包金额', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '金额必须大于0', trigger: 'blur' }
  ],
  totalCount: [
    { required: true, message: '请输入红包个数', trigger: 'blur' },
    { type: 'number', min: 1, max: 100, message: '个数必须在1-100之间', trigger: 'blur' }
  ]
};

// 加载群组列表
const loadGroups = async () => {
  try {
    const res = await getTelegramGroups();
    if (res.data && res.data.groups) {
      groupList.value = res.data.groups;
    }
  } catch (error) {
    console.error('加载群组列表失败:', error);
  }
};

// 加载最近发送记录
const loadRecentRecords = async () => {
  try {
    const res = await getRedPacketRecordList({ page: 1, limit: 5 });
    if (res.data && res.data.list) {
      recentRecords.value = res.data.list;
    }
  } catch (error) {
    console.error('加载发送记录失败:', error);
  }
};

// 提交表单
const handleSubmit = () => {
  formRef.value.validate(async valid => {
    if (!valid) return;

    ElMessageBox.confirm(
      `确定要向群组发送 ${formData.totalAmount} ${formData.symbol} 的${
        formData.packetType === 1 ? '普通' : '手气'
      }红包吗？`,
      '确认发送',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(async () => {
      try {
        loading.value = true;
        const res = await sendRedPacketDirect(formData);
        ElMessage.success('红包发送成功！');

        // 刷新最近记录
        loadRecentRecords();

        // 可选：重置表单
        // handleReset();
      } catch (error) {
        ElMessage.error(error.message || '发送失败');
      } finally {
        loading.value = false;
      }
    });
  });
};

// 重置表单
const handleReset = () => {
  formRef.value?.resetFields();
  Object.assign(formData, {
    groupId: '',
    groupName: '',
    packetType: 2,
    totalAmount: 10000,
    totalCount: 10,
    symbol: 'VND',
    lang: 'vi',
    expireMinutes: 10,
    maxGrabAmount: 0,
    blessingWords: '恭喜发财，大吉大利！'
  });
};

// 格式化时间戳
const formatTimestamp = ts => {
  if (!ts) return '-';
  return new Date(ts * 1000).toLocaleString('zh-CN');
};

// 页面加载
onMounted(() => {
  loadGroups();
  loadRecentRecords();
});
</script>

<style scoped lang="scss">
.red-packet-send {
  height: 100%;
  padding: 20px;

  .box-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: bold;
    }
  }

  .form-tip {
    font-size: 12px;
    color: #999;
    margin-top: 5px;
    line-height: 1.5;
  }

  .radio-desc {
    font-size: 12px;
    color: #999;
    margin-left: 4px;
  }

  .ml-2 {
    margin-left: 8px;
  }

  .mt-4 {
    margin-top: 20px;
  }

  .preview-box {
    margin-top: 20px;

    h3 {
      margin-bottom: 15px;
      font-size: 16px;
      color: #303133;
    }

    .telegram-preview {
      .telegram-message {
        max-width: 400px;
        background: #ffffff;
        border: 1px solid #e4e7ed;
        border-radius: 8px;
        padding: 16px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

        .message-header {
          margin-bottom: 12px;

          .sender-name {
            font-weight: 600;
            color: #409eff;
            font-size: 14px;
          }
        }

        .message-content {
          .red-packet-icon {
            font-size: 32px;
            text-align: center;
            margin-bottom: 8px;
          }

          .red-packet-text {
            text-align: center;
            font-size: 16px;
            color: #303133;
            margin-bottom: 12px;
          }

          .blessing-words {
            text-align: center;
            font-size: 14px;
            color: #606266;
            margin-bottom: 16px;
            font-style: italic;
          }

          .divider {
            height: 1px;
            background: #ebeef5;
            margin: 16px 0;
          }

          .packet-info {
            margin-bottom: 16px;

            .info-row {
              display: flex;
              justify-content: space-between;
              margin-bottom: 8px;
              font-size: 14px;

              &:last-child {
                margin-bottom: 0;
              }

              .info-label {
                color: #909399;
              }

              .info-value {
                color: #303133;
                font-weight: 500;
              }
            }
          }

          .call-to-action {
            text-align: center;
            font-size: 14px;
            color: #67c23a;
            margin-bottom: 12px;
            font-weight: 500;
          }

          .grab-button {
            margin-top: 12px;
          }
        }
      }

      .preview-tip {
        margin-top: 12px;
        padding: 8px 12px;
        background: #f4f4f5;
        border-radius: 4px;
        font-size: 12px;
        color: #909399;
        text-align: center;
      }
    }
  }
}
</style>

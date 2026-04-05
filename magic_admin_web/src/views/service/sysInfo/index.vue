<template>
  <div class="table-box">

    <div class="card table-main">

      <el-card style="">
        <template #header>
          <div class="card-header">
            <div class="left">产品配置</div>
            <div class="right">
              <el-button type="success" @click="refresh">刷新 </el-button>
              <el-button type="info" @click="edit">编辑</el-button>
            </div>
          </div>
        </template>
        
        <el-descriptions :column="3">
          <el-descriptions-item label="领取收益手续费比例:">{{tableData.claimFeeRatio}}</el-descriptions-item>
         </el-descriptions>
         <el-descriptions :column="3">
          <el-descriptions-item label="提现免审额度:">{{tableData.withdrawNoauditLimit}}</el-descriptions-item>
         </el-descriptions>
         <el-descriptions :column="3">
          <el-descriptions-item label="提现最小限额:">{{tableData.minWithdrawAmount}} &nbsp;&nbsp;&nbsp;
            <!-- <span class="red">注：0表示无限制</span> -->
          </el-descriptions-item>
         </el-descriptions>
         <el-descriptions :column="3">
          <!-- -1表示无限制，0表示不能提现 -->
          <el-descriptions-item label="每笔提现最大限额:">{{tableData.maxWithdrawAmount}} &nbsp;&nbsp;&nbsp;
            <!-- <span class="red">注：-1表示无限制，0表示不能提现</span> -->
          </el-descriptions-item>
         </el-descriptions>
         <el-descriptions :column="3">
          <!-- 0=禁止，1=允许 -->
          <el-descriptions-item label="是否允许提现（全局）:">
            <!-- <span :class=" tableData.withdrawStatus == 1 ? 'blue' : 'red'">{{ tableData.withdrawStatus == 1 ? "允许" : "禁止" }}</span> -->
            <el-tag
              :type="tableData.withdrawStatus ? 'primary' : 'danger'"
              disable-transitions
              >
              {{ tableData.withdrawStatus ? '允许' : '禁止' }}
            </el-tag>
          </el-descriptions-item>
         </el-descriptions>
         <el-descriptions :column="3">
          <el-descriptions-item label="提现手续费比例:">{{tableData.withdrawRatio}}</el-descriptions-item>
         </el-descriptions>
      </el-card>

      
    </div>


    <!-- 编辑 -->
    <el-dialog v-model="dialogSet" title="编辑"  destroy-on-close center width="70%">
      <el-form ref="form" :model="formData" label-width="180" :rules="rules">
        <el-form-item label="领取收益手续费比例" prop="claimFeeRatio">
          <el-input v-model="formData.claimFeeRatio" type="number" clearable/>
        </el-form-item>
        <el-form-item label="提现免审额度" prop="withdrawNoauditLimit">
          <el-input v-model="formData.withdrawNoauditLimit" type="number" clearable/>
          <!-- <div class="red-tip">
            0表示都要审核
          </div> -->
        </el-form-item>
        <el-form-item label="提现最小限额" prop="minWithdrawAmount">
          <el-input v-model="formData.minWithdrawAmount" type="number" clearable/>
          <!-- <div class="red-tip">
            0表示无限制
          </div> -->
        </el-form-item>
        <el-form-item label="每笔提现最大限额" prop="maxWithdrawAmount">
          <el-input v-model="formData.maxWithdrawAmount" type="number" clearable/>
          <!-- <div class="red-tip">
            -1表示无限制，0表示不能提现
          </div> -->
        </el-form-item>
        <el-form-item label="是否允许提现（全局）" prop="withdrawStatus">
          <el-select v-model="formData.withdrawStatus" placeholder="是否允许提现（全局）" clearable class="form-item">
            <!-- ：0=禁止，1=允许 -->
            <el-option label="允许" :value="1"/>
            <el-option label="禁止" :value="0"/>
          </el-select>
        </el-form-item>
        <el-form-item label="提现手续费比例" prop="withdrawRatio">
          <el-input v-model="formData.withdrawRatio" type="number" clearable/>
          <!-- <div class="red-tip">
           小于等于0 则是使用的实时价格
          </div> -->
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogSet = false">取消</el-button>
          <el-button type="primary" @click="setForm(form)" :loading="loading">
            设置
          </el-button>
        </div>
      </template>
    </el-dialog>

  </div>
  

</template>

<script setup name="question">
import {onMounted, ref} from "vue";
import { ElNotification, ElMessage } from "element-plus";
import { appSysInfo, setSysInfo } from "@/api/modules/service.js";

const tableData = ref([])

const loading = ref(false)
const dialogSet = ref(false) // 设置弹窗
const form = ref(null) // 设置表单
const rules = ref({ // 表单规则
  claimFeeRatio: [{required: true, message: '请选择', trigger: 'blur'}],
  withdrawNoauditLimit: [{required: true, message: '请输入', trigger: 'blur'}],
  minWithdrawAmount: [{required: true, message: '请输入', trigger: 'blur'}],
  maxWithdrawAmount: [{required: true, message: '请输入', trigger: 'blur'}],
  withdrawStatus: [{required: true, message: '请输入', trigger: 'blur'}],
  withdrawRatio: [{required: true, message: '请输入', trigger: 'blur'}],
})

const formData = ref({ // 表单数据
  claimFeeRatio: null,
  withdrawNoauditLimit: null,
  minWithdrawAmount: null,
  maxWithdrawAmount: null, 
  withdrawStatus: null,
  withdrawRatio: null,
})


onMounted(() => { 
  getList()
})

// 刷新
const refresh = () => {
  getList()
}


// 获取列表
const getList = async () => {
  tableData.value = []
  let params = {
    order: "id desc",
  }
  let res = await appSysInfo(params)
  if (res.code === 200) {
    tableData.value = res.data?.list[0]
  } else {
    ElNotification.error(res.msg)
  }
};


// 编辑
const edit = () => { 
  formData.value ={ ...tableData.value }
  dialogSet.value = true
}

// 设置
const setForm = (formEl) => { 
  if (!formEl) return;
  formEl.validate(async valid => {
    if (!valid) return;
    loading.value = true;
    
    try {
      const params = {
        ...formData.value,
        claimFeeRatio: Number(formData.value.claimFeeRatio),
        withdrawNoauditLimit: Number(formData.value.withdrawNoauditLimit),
        minWithdrawAmount: Number(formData.value.minWithdrawAmount),
        maxWithdrawAmount: Number(formData.value.maxWithdrawAmount),
        withdrawStatus: Boolean(formData.value.withdrawStatus),
        withdrawRatio: Number(formData.value.withdrawRatio),
      }
      const res = await setSysInfo(params)
      if (res.code === 200) { 
        ElMessage({
          message: '设置成功',
          type: 'success',
        })
        dialogSet.value = false
        getList()
      }
    } finally {
      loading.value = false;
    }
  });
}

</script>

<style scoped lang="scss">
.card-header {
  display: flex;
  justify-content: space-between;

  .left {
    display: flex;
    align-items: center;
  }

  .right {
    display: flex;
    align-items: center;
  }
}

.red {
  color: var(--el-color-danger);
  text-decoration: underline;
}

.blue {
  color: var(--el-color-success);
}

.red-tip {
  color: var(--el-color-danger);
  font-size: 12px; 
  line-height: 1.5; 
  margin-top: 4px;
}

.el-card {
  overflow-y: auto;
}
</style>

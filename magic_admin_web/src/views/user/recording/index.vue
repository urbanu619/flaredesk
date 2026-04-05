<template>
  <div class="table-box">
    <!--    头部搜索-->
    <div class="card table-search" v-show="isShowSearch">
      <el-form ref="elSearchFormRef" :inline="true" size="small" :model="searchForm" class="demo-form-inline"
               @keyup.enter="onSubmit">
        <template v-if="showAllQuery">

        </template>
        <div class="operation">
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onResetSearch">重置</el-button>
          <el-button link type="primary" icon="arrow-down" @click="showAllQuery = true" v-if="!showAllQuery">展开</el-button>
          <el-button link type="primary" icon="arrow-up" @click="showAllQuery = false" v-else>收起</el-button>
        </div>
      </el-form>
    </div>

    <div class="card table-main">
      <!-- 表格头部 操作按钮 -->
      <div class="table-header">
        <div class="header-button-lf">
          <el-button type="success" @click="exportTable">导出</el-button>
          <!-- <el-button type="success" @click="mutilHandler">导入</el-button> -->
        </div>
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refresh"/>
          <el-button :icon="isShowSearch ? 'ArrowUpBold' :'ArrowDownBold'" circle @click="isShowSearch = !isShowSearch"/>
        </div>
      </div>

      <!-- 表格主体 -->
      <el-table ref="myTable" :data="tableData" border size="small" :cell-style="cellStyle">
        <el-table-column prop="id" label="用户ID" align="center" min-width="120">
          <template #default="scope">
            <span v-address-format="scope.row.id"></span>
          </template>
        </el-table-column>
        <el-table-column prop="parentId" label="上级ID" align="center" min-width="120">
          <template #default="scope">
            <span v-address-format="scope.row.parentId"></span>
          </template>
        </el-table-column>
        <!-- <el-table-column prop="lockLevel" label="团队锁定等级" width="120" align="center"/> -->
        <el-table-column prop="level" label="等级/锁定等级" align="center" min-width="120">
          <!-- <template #header>
            <span>业绩等级/锁定业绩等级</span>
            <el-tooltip content="业绩等级：通过用户小区业绩升级；锁定业绩等级：锁定后保底升级到该等级" placement="top">
              <el-icon style="color: var(--el-color-primary);" :size="16">
                <QuestionFilled/>
              </el-icon>
            </el-tooltip>
          </template> -->
          <template #default="scope">
            <span >{{ valueToLabel(levelList, scope.row.level) }} /</span> 
            <el-icon class="copy-icon locked" v-if=" scope.row.lockLevel > 0"><Lock /></el-icon>
            <el-icon class="copy-icon unlock" v-else><Unlock /></el-icon>
            {{ scope.row.lockLevel }}
          </template>
        </el-table-column>
        <el-table-column prop="inviteCount" label="直推人数" align="center" min-width="80" />
        <el-table-column prop="teamCount" label="团队人数" align="center" min-width="80" />
        <el-table-column prop="fundUsdAmount" label="联创基金USD" align="center" min-width="140" />
        <el-table-column prop="personAchievement" label="个人业绩" align="center" min-width="140" />
        <el-table-column prop="teamTodayAchievement" label="当日团队业绩" align="center" min-width="140" />
        <el-table-column prop="teamAchievement" label="团队总业绩" align="center" min-width="140" />

<!--         
        <el-table-column fixed="right" label="操作" align="center" min-width="80">
          <template #default="scope">
            <el-button type="primary" link @click="openDialog(scope.row)">编辑</el-button>
          </template>
        </el-table-column> -->

      </el-table>

      <!-- 分页组件 -->
      <slot name="pagination">
        <Pagination
          :pageable="pageable"
          @handleCurrent="handleCurrent"
        />
      </slot>
    </div>


    <!-- 锁定等级弹窗 -->
    <el-dialog v-model="dialogLockLevel" title="锁定业绩等级"  destroy-on-close center width="50%">
      <el-form ref="formRef" :model="lockData" :rules="rules" label-width="auto">
        <el-form-item label="用户ID" prop="userId">
          <el-input v-model="lockData.userId" placeholder="用户ID" disabled/>
        </el-form-item>
        <el-form-item label="当前等级" prop="level">
          <el-input v-model="lockData.level" placeholder="金额" disabled/>
        </el-form-item>
        <el-form-item label="业绩等级" prop="lockLevel">
          <el-select v-model="lockData.lockLevel" placeholder="请选择等级" clearable style="width: 100%;">
            <el-option v-for="(item , index) in levelList " :key="item.value" :label="item.label" :value="item.value"/>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogLockLevel = false">取消</el-button>
          <el-button type="primary" @click="setLockLevel" :loading="loading">
            锁定
          </el-button>
        </div>
      </template>
    </el-dialog>


    <!-- 设置奖励等级弹窗 -->
    <el-dialog v-model="dialogRewardLevel" title="设置奖励等级"  destroy-on-close center width="50%">
      <el-form ref="formRefRewardLevel" :model="rewardLevelData" :rules="rulesRewardLevel" label-width="auto">
        <el-form-item label="用户ID" prop="userId">
          <el-input v-model="rewardLevelData.userId" placeholder="用户ID" disabled/>
        </el-form-item>
        <el-form-item label="奖励等级" prop="level">
          <el-input v-model="rewardLevelData.rewardLevel" placeholder="请输入奖励等级" clearable @input="onInput" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogRewardLevel = false">取消</el-button>
          <el-button type="primary" @click="setRewardLevelHandler" :loading="loading">
            设置
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 直推人数 -->
    <el-dialog 
      style="background: #ffffff; width: 650px"
      v-model="DirectPushormVisible"
      :title="DirectPushormTitle"
      :close-on-click-modal="false"
      @close="closeOnClickModal"
    >
      <div class="DirectPushorView">
        <el-tree
            ref="treeRef"
            style="width: 600px; max-height: 70vh; overflow-x: scroll;overflow-y: scroll;"
            :data="DirectPushormList"
            node-key="id"
            :default-expanded-keys="defaultExpandedKeys"
            @node-click="handleNodeClick"

        >
          <template #default="{ node, data }">
            <div class="custom-tree-node" :class="data.inviteCount > 0 ? 'has-leef' : ''">
              <div class="level">Lv{{ data.level }}</div>
              <div class="line">|</div>
              <div class="uid">uid：{{ data.id }}</div>
              <div class="line">|</div>
              <div class="personAchievementUsdAmount">个人业绩：{{ data.personAchievementUsdAmount }}</div>
              <div class="line">|</div>
              <div class="teamAchievementQuantity">团队业绩：{{ data.teamAchievementQuantity }}</div>
              <div class="line">|</div>
              <div class="inviteCount">直推人数：{{ data.inviteCount }}</div>
              <div class="line">|</div>
              <div class="teamCount">团队人数：{{ data.teamCount }}</div>
            </div>
          </template>


          <!--        等级  level-->
          <!--        uid   id-->
          <!--        个人业绩  personAchievementUsdAmount-->
          <!--        团队业绩  teamAchievementQuantity-->
          <!--        直推人数   inviteCount-->
          <!--        团队人数  teamCount-->


        </el-tree>
      </div>

    </el-dialog>

      <!-- 编辑弹窗 -->
      <el-dialog v-model="dialogEditData" title="编辑用户指标"  destroy-on-close center width="50%">
      <el-form ref="formRefEdit" :model="editData" :rules="rulesEdit" label-width="auto">
        <el-form-item label="用户ID" prop="id">
          <el-input v-model="editData.id" placeholder="用户ID" disabled/>
        </el-form-item>
        <el-form-item label="活期每日限额" prop="dailyCurrentLimit">
          <el-input v-model="editData.dailyCurrentLimit" placeholder="活期每日限额" clearable type="number" class="form-item"/>
        </el-form-item>
        <el-form-item label="锁定等级" prop="lockLevel">
          <el-input v-model="editData.lockLevel" placeholder="锁定等级" clearable type="number" class="form-item"/>
        </el-form-item>
        
      </el-form>
    
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogEditData = false">取消</el-button>
          <el-button type="primary" @click="editConfirm" :loading="loading">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>


  <!-- 批量导入弹窗 -->
  <el-dialog v-model="dialogTransferMutil" title="批量导入" destroy-on-close center width="50%">
    <el-form ref="formRefMutil" :model="transferMutil" :rules="rulesMutil" label-width="auto">
      <!-- <el-form-item label="转移者编号" prop="systemNo">
        <el-input v-model="transferMutil.systemNo" placeholder="最大30位" clearable maxlength="30" class="form-item"/>
      </el-form-item> -->
      <el-form-item label="表格数据" prop="batch">
        <el-upload
          class="upload-demo"
          accept=".csv"
          :show-file-list="false"
          :before-upload="handleFileUpload"
        >
          <el-button type="primary">导入 CSV</el-button>
        </el-upload>
        <el-button type="success" @click="downloadTemplate" style="margin-left: 10px;">
          下载模板
        </el-button>
        <div class="red-tip">字段备注: uid为交易所id, orderNo导入订单编号, quantity为股票数量.</div>
        <span class="red" style="color: red">所有字段必填</span>
      </el-form-item>

    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="transferCancelMutil">取消</el-button>
        <el-button type="primary" @click="transferConfirmMutil" :loading="loading"> 确定</el-button>
      </div>
    </template>
  </el-dialog>

  </div>
  
  <ExportFieldFiltering 
    ref="ExportFieldFilteringView"
    :derived_field="derived_field"
    @filesselectedfiles="filesselectedfiles">
  </ExportFieldFiltering>

</template>

<script setup>
import { onMounted, ref, reactive } from "vue";
import Papa from 'papaparse'
import { h } from 'vue'
import { Lock } from '@element-plus/icons-vue'
import Pagination from "@/components/Pangination/Pagination.vue"
import { Refresh } from "@element-plus/icons-vue";
import { ElNotification, ElMessageBox } from "element-plus";
import { appUserQuota, batchImportStock, updateUserQuota, setLevel, setRewardLevel } from "@/api/modules/user.js";
import { valueToLabel } from '@/utils/index.ts'
import { levelList, levelListOption } from '@/utils/dict'
import { downLoadFile } from "@/api/downloadFile.js"
import ExportFieldFiltering from "@/components/ExportFieldFiltering/index.vue"

// ***************搜索框相关*****************************
// 表单规则校验
const rulesEdit = ref({
  id: [{required: true, message: '不能为空', trigger: 'change',}],
  // dailyCurrentLimit: [{required: true, message: '不能为空', trigger: 'change',}],
  // lockLevel: [{required: true, message: '不能为空', trigger: 'change',}],
})

const rules = ref({
  lockLevel: [{required: true, message: '请选择', trigger: ['blur', 'change']}],
})
const rulesRewardLevel = ref({
  rewardLevel: [
    {
      required: true,
      message: '请输入占股',
      trigger: 'blur'
    },
    {
      validator: (_, value, callback) => {
        const num = Number(value)
        const pattern = /^[1-9]\d?$|^12$|^0$/  // 0~100 的正整数（含0）
        if (pattern.test(value) && num >= 0 && num <= 12) {
          callback()
        } else {
          callback(new Error('请输入0~12的正整数'))
        }
      },
      trigger: 'blur'
    }
  ],
})

const showAllQuery = ref(true)
const isShowSearch = ref(true)

const searchForm = ref({
  pid: null,
  id: null,
  levels: null,
})
const tableData = ref([])
const myTable = ref(null) // 表格ref
const dialogLockLevel = ref(false) // 锁定钱包弹窗
const formRef = ref(null) // 锁定表单
const lockData = ref({  // 锁定的等级
  userId: null, // 用户ID
  level: null, // 当前等级
  lockLevel: null, // 锁定等级
})

const rewardLevelData = ref({  
  userId: null, // 用户ID
  rewardLevel: null, // 奖励等级
})
const dialogRewardLevel = ref(false) // 设置奖励等级的弹窗
const formRefRewardLevel = ref(null) // 奖励等级的弹窗表单


// 编辑对象
const editData = reactive({
  id: null, // userId
  dailyCurrentLimit: null, // 用户活期每日限额
  lockLevel: null, // 锁定等级
})

const dialogEditData = ref(false) // 编辑弹窗开启/关闭

const formRefEdit = ref(null) // 编辑表单

const loading = ref(false) // loading


// 查询参数
const pageable = reactive({
  pageNum: 1,
  pageSize: 30,
  total: 0
})

const ExportFieldFilteringView = ref(null)
const derived_field = ref([])

onMounted(() => { 
  getList()
})

// 设置表格列背景颜色
const cellStyle = ({ column }) => {
  if (column.property === 'largeRegionAchievement' || column.property === 'fewTeamAchievement' || column.property === 'personAchievement' || column.property === 'teamAchievement') {
    return { background: '#f0f9ff' }
  }
  if (column.property === 'level' ) {
    return { background: 'rgba(240, 255, 251, 1)' }
  }
  if (column.property === 'nodeEquityLevel' ) {
    return { background: 'rgba(244, 255, 240, 1)' }
  }
  if (column.property === 'rewardLevel' ) {
    return { background: 'rgba(255, 249, 240, 1)' }
  }
  
  return {}
}

// 头部查询按钮
const onSubmit = () => {
  console.log('submit!')
  getList()
}
// 刷新
const refresh = () => {
  getList()
}

// 重置
const onResetSearch = () => {
  searchForm.value.pid = null
  searchForm.value.id = null
  searchForm.value.levels = null

  getList()
}

// 分页回调
const handleCurrent = (data) => {
  pageable.pageNum = data.current;
  pageable.pageSize = data.pageSize;

  getList()
}

// 获取列表
const getList = async () => {
  tableData.value = []
  let params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: "id desc",
    ...searchForm.value
  }

  let res = await appUserQuota(params)
  if (res.code === 200) {
    tableData.value = res.data.list
    pageable.total = res.data.paging.total

    let arrList = []
    // derived_field.value = res.data.cols
    // 选择可显示的字段集合
    res.data.cols?.map(item => { 
      if (item?.json == 'id' || item?.json == 'parentId' || item?.json == 'level' || item?.json == 'lockLevel' || item?.json == 'inviteCount' ||
        item?.json == 'teamCount' || item?.json == 'fundUsdAmount' || item?.json == 'personAchievement' || item?.json == 'teamTodayAchievement' || item?.json == 'teamAchievement' 
      ) { 
        arrList.push(item)
      }
    })
    derived_field.value = arrList

  } else {
    ElNotification.error(res.msg)
  }
};


// 打开设置奖励等级弹窗
const openRewardLevel = (row) => { 
  rewardLevelData.value.userId = row.id
  rewardLevelData.value.rewardLevel = row.rewardLevel
  dialogRewardLevel.value = true
}

// 关闭设置奖励等级弹窗
const closeRewardLevel= () => { 
  rewardLevelData.value.userId = ''
  rewardLevelData.value.rewardLevel = ''

  dialogRewardLevel.value = false
}


// 奖励等级输入校验
const onInput = (val) => {
  
  // 仅保留数字字符
  let filtered = val.replace(/[^0-9]/g, '')

 // 限制最大值为12
 if (filtered !== '') {
   const num = parseInt(filtered, 10)
   filtered = num > 12 ? '12' : String(num)
 }

 rewardLevelData.value.rewardLevel = filtered
}

// 设置奖励等级
const setRewardLevelHandler = (row) => { 
  formRefRewardLevel.value?.validate( async (valid) => {
    if (valid) {
      const params = {
        userId: rewardLevelData.value.userId,
        rewardLevel: Number(rewardLevelData.value.rewardLevel),
      }
      let res = await setRewardLevel(params)
      if (res.code === 200) {
        ElNotification.success('设置成功')
        closeRewardLevel()
        getList()
      } else {
        ElNotification.error(res.msg)
      }
    } else {
      console.log('error submit!', fields)
    }
  })
}

// 打开锁定等级弹窗
const openLock = (row) => { 
  lockData.value.userId = row.id
  lockData.value.level = row.level
  dialogLockLevel.value = true
}

// 关闭锁定等级弹窗
const closeLock = (row) => { 
  lockData.value.userId = ''
  lockData.value.level = ''
  lockData.value.lockLevel = ''

  dialogLockLevel.value = false
}

// 锁定等级
const setLockLevel = (row) => { 
  formRef.value?.validate( async (valid) => {
    if (valid) {
      const params = {
        userId: lockData.value.userId,
        level: lockData.value.lockLevel,
      }
      let res = await setLevel(params)
      if (res.code === 200) {
        ElNotification.success(res.msg)
        closeLock()
        getList()
      } else {
        // ElNotification.error(res.msg)
      }
    } else {
      console.log('error submit!', fields)
    }
  })
}

// 排序变化时触发
const onSortChange = ({ prop, order }) => {
  sortState.prop = prop;
  sortState.order = order;
  pageable.pageNum = 1; // 排序改变时重置为第一页
  getList()
}


// #endregion 锁定等级

const DirectPushormVisible = ref(false)
const DirectPushormTitle = ref("")
const treeRef = ref()
const DirectPushormList = ref([])
const defaultExpandedKeys = ref([]) //默认展开

//关闭直推弹框
const closeOnClickModal = () => {
  DirectPushormTitle.value = ""
  DirectPushormList.value = []
  defaultExpandedKeys.value = []
}

// 表格获取直推
const DirectPush = async (row) => {
  console.log("row:", row)
  if (row.inviteCount == 0) {
    ElNotification("UID:" + row.id + ",无直推下级")

    return
  }
  DirectPushormTitle.value = "UID:" + row.id + "的团队成员"
  /*
  等级  level
  uid   id
  个人业绩  personAchievementUsdAmount
  团队业绩  teamAchievementQuantity


  直推人数   inviteCount
  团队人数  teamCount
  * */

  let res = await getData_dirPush(row.id, row.inviteCount)
  if (res.length != 0) {
    DirectPushormList.value = sliceDate(res)
  }
  DirectPushormVisible.value = true
}

// 弹框 点击直推列表 获取下级
const handleNodeClick = async (data) => {
  if (data.inviteCount == 0) {
    ElNotification("UID:" + data.id + ",无直推下级")
    return
  }
  if (data.children.length != 0) {
    return
  }
  let res = await getData_dirPush(data.id, data.inviteCount)
  if (res.length != 0) {
    addChildrenById(DirectPushormList.value, data.id, sliceDate(res));
    console.log(DirectPushormList.value)
    defaultExpandedKeys.value = [data.id]

  }
}


// 获取团队人数
const getData_dirPush = async (pid, inviteCount) => {
  return new Promise((resolve, reject) => {
    let params = {
      parentId: pid,
      pageNum: 1,
      pageSize: 30000,
    }
    if (inviteCount) { 
      params.pageSize = inviteCount
    }
    appUserQuota(params).then((res) => {
      if (res.code == 200) {
        resolve(res.data.list)
      } else {
        reject([])
      }
    })
  })
}
// 数据处理
const sliceDate = (data) => {
  let arr = []
  data.forEach((v) => {
    let obj = {
      level: v.level,
      id: v.id,
      personAchievementUsdAmount: v.personAchievement, // 个人业绩
      teamAchievementQuantity: v.teamAchievement, // 团队业绩
      inviteCount: v.inviteCount,  //直推人数
      teamCount: v.teamCount, // 团队人数
      children: [],
    }
    arr.push(obj)
  })
  return arr
}
// 递归插入
const addChildrenById = (data, targetId, childrenToAdd) => {
  for (let i = 0; i < data.length; i++) {
    const item = data[i];

    // 如果找到匹配的ID
    if (item.id === targetId) {
      // 如果children不存在则创建
      if (!item.children) {
        item.children = [];
      }
      // 添加子数组
      item.children.push(...childrenToAdd);
      return true;
    }

    // 如果有children则递归查找
    if (item.children && item.children.length > 0) {
      const found = addChildrenById(item.children, targetId, childrenToAdd);
      if (found) return true;
    }
  }

  return false;
}

// 编辑弹窗打开
const openDialog = (row) => {
  editData.id = null
  editData.dailyCurrentLimit = null
  editData.lockLevel = null
  editData.id = row.id
  if (row.dailyCurrentLimit) {
    editData.dailyCurrentLimit = row.dailyCurrentLimit
  }
  if (row.lockLevel > 0) {
    editData.lockLevel = row.lockLevel
  }
  

  dialogEditData.value = true
}

// 确定编辑
const editConfirm = () => {
  formRefEdit.value?.validate(async (valid) => {
    if (valid) {
      const params = {
        id: editData.id,
      }
      if (editData.dailyCurrentLimit !== null) { 
        params.dailyCurrentLimit = Number(editData.dailyCurrentLimit)
      }
      if (editData.lockLevel !== null) { 
        params.lockLevel = Number(editData.lockLevel)
      }
      let res = await updateUserQuota(params)
      if (res.code === 200) {
        ElNotification.success('操作成功')
        dialogEditData.value = false
        getList()
      } else {
        // ElNotification.error(res.msg)
      }
      console.log(params, 'submit!')
    } else {
      console.log('error submit!', fields)
    }
  })
}


// 字段弹框打开
const exportTable = () => {
  ExportFieldFilteringView.value.show(derived_field.value)
}

// 字段选择回调   导出 start
const filesselectedfiles = async (str) => {
  let params = {
    current: pageable.pageNum,
    pageSize: pageable.pageSize,
    order: "id desc",
    isExport: true,
    fields: str,
    ...searchForm.value
  }

  let res = await appUserQuota(params)

  try {
    let downFile = {fileName: "用户指标", fileUrl: res.data.url}
    console.log(downFile)
    if (res.data.url) {
      await downLoadFile(downFile)
      ElNotification.success("导出成功")
    } else { 
      ElNotification.success("暂无数据")
    }
  } catch { 
    ElNotification.error("导出失败")
  }
}

// # 批量导入模版
const dialogTransferMutil = ref(false) // 批量划转弹窗开启/关闭
const formRefMutil = ref(null) // 批量划转表单
const uploadTableData = ref([])
const transferMutil = reactive({ // 批量划转对象
  // systemNo: null,// 操作者
  batch: null,// 数据
})
const rulesMutil = ref({
  // systemNo: [{required: true, message: '请输入', trigger: ['blur', 'change']}],
  // batch: [{required: true, message: '请输入', trigger: ['blur', 'change']}],
})


// 批量导入股票算力弹窗打开
const mutilHandler = (row) => {
  // transferMutil.systemNo = null
  transferMutil.batch = null
  dialogTransferMutil.value = true
}

// 关闭批量导入股票算力弹窗
const transferCancelMutil = () => {
  transferMutil.batch = null
  // transferMutil.systemNo = null
  uploadTableData.value = []
  dialogTransferMutil.value = false
}

// 批量上传
const handleFileUpload = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = (e) => {
      const csvText = e.target.result

      // 解析 CSV → 数组
      const result = Papa.parse(csvText, {
        header: true, // 第一行当表头
        skipEmptyLines: true
      })

      uploadTableData.value = result.data.map(
        item => {
          console.log(item, 123)
          const obj = {
            uid: item.uid.trim(),
            orderNo: item.orderNo.trim(),
            quantity: item.quantity.trim(),
          }
          return obj
        }
      ) // 数据行
      console.log(uploadTableData.value, 111)

      resolve(false) // 阻止 el-upload 自动上传
    }
    reader.onerror = reject
    reader.readAsText(file, 'utf-8')
  })
}


// 确定批量导入股票算力
const transferConfirmMutil = () => {
  formRefMutil.value?.validate(async (valid) => {
    console.log("valid", valid)

    if (valid) {
      if (uploadTableData.value?.length == 0) {
        ElNotification.error('请上传文件')
        return
      }
      const params = {
        // systemNo: transferMutil.systemNo,
        batch: uploadTableData.value,
        target: "x_staking" //+ new Date().getTime()
      }
      // if (transferMutil.systemNo) {
      //   params.systemNo = transferMutil.systemNo
      // }
      let res = await batchImportStock(params)
      if (res.code === 200) {
        // ElNotification.success('成功金额' + res.data?.realAmount)
        ElMessageBox({
          title: '提示',
          message: h('div', null, [
            h('p', {style: 'max-height: 60vh;  overflow-y: auto'}, res.data?.msg),
            h('p', {style: 'color: teal'}, '成功金额' + res.data?.realAmount),
          ]),
        })
        transferCancelMutil()
        getList()
      } else {
        ElNotification.error(res.msg)
      }
      console.log(params, 'submit!')
    } else {
      console.log('error submit!', fields)
    }
  })
}

// 下载模版
const downloadTemplate = () => {
  // 模板 CSV 内容
  const csvContent = `uid,orderNo,quantity\n30011,100,1,1,1`

  // 创建 Blob
  const blob = new Blob([csvContent], {type: 'text/csv;charset=utf-8;'})

  // 创建下载链接
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', '批量导入股票算力(uid=交易所id,orderNo=订单号,,quantity=股票数量).csv')
  link.style.visibility = 'hidden'

  // 触发下载
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}
const resetOrder = (row) => {
  ElMessageBox.confirm("重新触发保险股生效", '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    const params = {}
    resetValidBachelor(params).then(res => {
      if (res.code === 200) {
        ElNotification.success('操作成功')
        getList()
      } else {
        ElNotification.error(res.msg)
      }
    })
  }).catch(() => {
    ElMessage({
      type: 'info',
      message: 'Delete canceled',
    })
  })
}


</script>

<style scoped lang="scss">
.copy-icon{
  margin-left: 10px;
  cursor: pointer;
}

.locked{
  color: red;
}
.unlock{
  color: green;
}

.custom-tree-node {
  display: flex;
  padding: 0 10px;
  font-size: 14px;
  line-height: 20px;

  .level {
    background: #c9cdd0;
    border-radius: 2px;
    padding: 0 10px;
    height: 20px;
    text-align: center;
    line-height: 20px;
    font-size: 14px;
  }

  .line {
    margin: 0 10px;
  }

}

::v-deep {
  .el-icon.el-tree-node__expand-icon.is-leaf:has(+ .has-leef) {
    visibility: visible;
    color: var(--el-tree-expand-icon-color);
  }
}

.no-margin {
  margin-left: 0;
}

</style>

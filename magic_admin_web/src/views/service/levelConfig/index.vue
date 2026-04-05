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
        <div class="header-button-ri">
          <el-button :icon="Refresh" circle @click="refresh"/>
          <el-button :icon="isShowSearch ? 'ArrowUpBold' :'ArrowDownBold'" circle @click="isShowSearch = !isShowSearch"/>
        </div>
      </div>

      <!-- 表格主体 -->
      <el-table :data="tableData" border size="small">
        <!-- <el-table-column prop="id" label="ID" align="center" min-width="100"/> -->
        <el-table-column prop="level" label="用户等级" align="center" min-width="80"/>
        <!-- <el-table-column prop="icon" label="等级图标" align="center"> 
          <template #default="scope">
            <el-image class="img" :src="scope.row.icon"/>
          </template>
        </el-table-column> -->
        <el-table-column prop="levelName" label="等级名称" align="center" min-width="80"/>
        <!-- <el-table-column prop="inviteTeamLevel" label="升级所需等级" align="center" min-width="120">
          <template #header>
            <span>升级所需等级</span>
            <el-tooltip content="升级所需下属团队等级" placement="top">
              <el-icon style="color: var(--el-color-primary);" :size="16">
                <QuestionFilled/>
              </el-icon>
            </el-tooltip>
          </template>
          <template #default="scope">
            {{ scope.row.inviteTeamLevel }}
          </template>
        </el-table-column>
        <el-table-column prop="inviteTeamNums" label="升级所需个数" align="center" min-width="120">
          <template #header>
            <span>升级所需个数</span>
            <el-tooltip content="升级所需下属团队等级个数" placement="top">
              <el-icon style="color: var(--el-color-primary);" :size="16">
                <QuestionFilled/>
              </el-icon>
            </el-tooltip>
          </template>
          <template #default="scope">
            {{ scope.row.inviteTeamNums }}
          </template>
        </el-table-column> -->
        <el-table-column prop="ratio" label="团队奖励比例" align="center" min-width="120"/>
        <el-table-column prop="personAmount" label="个人业绩要求" align="center" min-width="120"/>
        <el-table-column prop="teamAicAmount" label="团队业绩要求" align="center" min-width="120"/>
        <!-- <el-table-column prop="isGetEqualLevelReward" label="是否获得奖励" align="center" min-width="120">
          <template #header>
            <span>是否获得奖励</span>
            <el-tooltip content="是否可以获得直推平级奖励" placement="top">
              <el-icon style="color: var(--el-color-primary);" :size="16">
                <QuestionFilled/>
              </el-icon>
            </el-tooltip>
          </template>
          <template #default="scope">
            <el-tag
              :type="scope.row.isGetEqualLevelReward ? 'primary' : 'danger'"
              disable-transitions
              >
              {{ scope.row.isGetEqualLevelReward ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="equalLevelReward" label="平级奖励比例" align="center" min-width="120"/> -->
        <!-- <el-table-column prop="equalLevelReward" label="平级奖励" align="center" min-width="120">
          <template #header>
            <span>是否获得奖励</span>
            <el-tooltip content="是否可以获得直推平级奖励" placement="top">
              <el-icon style="color: var(--el-color-primary);" :size="16">
                <QuestionFilled/>
              </el-icon>
            </el-tooltip>
          </template>
          <template #default="scope">
            <el-tag
              :type="scope.row.equalLevelReward ? 'primary' : 'danger'"
              disable-transitions
              >
              {{ scope.row.equalLevelReward ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column> -->
        <!-- <el-table-column prop="createdAt" label="创建时间" align="center" min-width="100">
          <template #default="scope">
            {{ scope.row.createdAt ?dayjs(scope.row.createdAt * 1000).format("YYYY-MM-DD HH:mm:ss") : ''}}
          </template>
        </el-table-column> -->
        <!-- <el-table-column prop="updatedAt" label="更新时间" align="center" min-width="100">
          <template #default="scope">
            {{dayjs(scope.row.updatedAt * 1000).format("YYYY-MM-DD HH:mm:ss")}}
          </template>
        </el-table-column> -->
        <!-- <el-table-column fixed="right" label="操作" align="center" min-width="120">
          <template #default="scope">
            <el-button type="primary" :icon="Lock" size="small" @click="openEdit(scope.row)">
              编辑
            </el-button>
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



    <!-- 编辑 -->
    <el-dialog v-model="dialogSet" title="编辑"  destroy-on-close center width="70%">
      <el-form ref="form" :model="formData" label-width="180" :rules="rules">
        <el-form-item label="id" prop="id">
          <el-input v-model="formData.id" disabled/>
        </el-form-item>
        <!-- <el-form-item label="等级" prop="level">
          <el-input v-model="formData.level" disabled/>
        </el-form-item> -->
        <!-- <el-form-item label="等级图标" prop="icon">
          <upImage :avatar="formData.icon" @upImgSuccess="upSucess" @upDel="upDel" ></upImage>
        </el-form-item> -->
        <el-form-item label="奖励比例" prop="ratio">
          <el-input v-model="formData.ratio" />
        </el-form-item>
        <el-form-item label="个人业绩要求" prop="personalRegularAmount">
          <el-input v-model="formData.personalRegularAmount" />
        </el-form-item>
        <el-form-item label="团队业绩要求" prop="teamRegularAmount">
          <el-input v-model="formData.teamRegularAmount" />
        </el-form-item>
        <el-form-item label="是否获得直推平级奖励" prop="isGetEqualLevelReward">
          <el-select v-model="formData.isGetEqualLevelReward" placeholder="请选择" clearable style="width: 100%;" >
            <el-option label="是" :value="true"/>
            <el-option label="否" :value="false"/>
          </el-select>
        </el-form-item>
        <el-form-item label="平级奖励比例" prop="equalLevelReward">
          <el-input v-model="formData.equalLevelReward" />
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
import { onMounted, ref, reactive } from "vue";
import dayjs from "dayjs";
import Pagination from "@/components/Pangination/Pagination.vue"
import upImage from "@/components/Upload/Img.vue"
import { Refresh } from "@element-plus/icons-vue";
import { ElMessage, ElNotification} from "element-plus";
import { appLevelConfig, setLevelConfig } from "@/api/modules/service.js";
import { Lock } from '@element-plus/icons-vue'
 

// ***************搜索框相关*****************************
const showAllQuery = ref(true)
const isShowSearch = ref(true)

const searchForm = ref({
 
})
const tableData = ref([])


const pageable = reactive({
  pageNum: 1,
  pageSize: 30,
  total: 0
})

const loading = ref(false)
const dialogSet = ref(false) // 设置弹窗
const form = ref(null) // 设置表单
const rules = ref({ // 表单规则
  // level: [{required: true, message: '请输入', trigger: 'blur'}],
  icon: [{required: true, message: '请上传', trigger: 'blur'}],
  ratio: [{required: true, message: '请输入', trigger: 'blur'}],
  personalRegularAmount: [{required: true, message: '请输入', trigger: 'blur'}],
  teamRegularAmount: [{required: true, message: '请输入', trigger: 'blur'}],
  isGetEqualLevelReward: [{ required: true, message: '请选择', trigger: ['blur', 'change'] }],
  equalLevelReward: [{required: true, message: '请输入', trigger: 'blur'}],
})

const formData = ref({ // 表单数据
  id: null, // id
  // level: null, // 等级
  icon: null, // 图标
  ratio: null, // 团队奖励比例
  personalRegularAmount: null, // 个人业绩要求
  teamRegularAmount: null, // 团队业绩要求
  isGetEqualLevelReward: null, // 是否获得奖励
  equalLevelReward: null, // 平级奖励比例
})


onMounted(() => { 
  getList()
})

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
  searchForm.value.id = null
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
    ...searchForm.value
  }
  let res = await appLevelConfig(params)
  if (res.code === 200) {
    tableData.value = res.data.list
    pageable.total = res.data.paging.total

  } else {
    ElNotification.error(res.msg)
  }
};

// 编辑
const openEdit = (row) => { 
  formData.value.id = row.id
  // formData.value.level = row.level
  formData.value.icon = row.icon
  formData.value.ratio = row.ratio
  formData.value.personalRegularAmount = row.personalRegularAmount
  formData.value.teamRegularAmount = row.teamRegularAmount
  formData.value.isGetEqualLevelReward = row.isGetEqualLevelReward
  formData.value.equalLevelReward = row.equalLevelReward
  
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
        level: formData.value.level,
      }
      const res = await setLevelConfig(params)
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

// 图片上传完成回调
const upSucess = (data) => {
  formData.value.icon = data
}

// 图片删除回调
const upDel = () => {
  formData.value.icon = ""
}

</script>

<style scoped lang="scss">


</style>

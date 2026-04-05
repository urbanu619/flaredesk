<template>
  <!--设置谷歌验证码-->
  <el-dialog v-model="dialogFormVisible_setGoogle" title="设置谷歌验证码" :close-on-click-modal="false" :show-close="false">
    <el-form ref="formUserRef_set" :model="form_user" label-width="140" :rules="rules">

      <el-form-item label="原Google Codey" prop="originalGoogleCode" v-if="!userInfo.isSetGoogleAuth">
        <el-input v-model="form_user.originalGoogleCode" clearable/>
      </el-form-item>
      <el-form-item label="Google Key">
        <div>
          <div class="codeView" style="display: block;margin-bottom: 10px" @click="getKey">
            <qrcode-vue :value="googleQrcode" :size="80"></qrcode-vue>
          </div>
          <div class="iconView copy-btn1" style="cursor: pointer;color: green" @click="copy(googleKeyCode)" >
            {{ googleKeyCode }}
          </div>
        </div>

      </el-form-item>
      <el-form-item label="Google Codey" prop="googleCode">
        <el-input v-model="form_user.googleCode" clearable/>
      </el-form-item>


    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="popupUserCancal_set(formUserRef_set)">取消</el-button>
        <el-button type="primary" @click="popupUserConfirm_set(formUserRef_set)">确定</el-button>
      </span>
    </template>
  </el-dialog>


  <!--取消谷歌验证码-->
  <el-dialog v-model="dialogFormVisible_closeGoogle" title="取消谷歌验证码" :close-on-click-modal="false" :show-close="false">
    <el-form ref="formUserRef_del" :model="form_user" label-width="120" :rules="rules">
      <el-form-item label="Google Codey" prop="googleCode">
        <el-input v-model="form_user.googleCode" clearable/>
      </el-form-item>

    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="popupUserCancal_del(formUserRef_del)">取消</el-button>
        <el-button type="primary" @click="popupUserConfirm_del(formUserRef_del)">确定</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { ElNotification } from "element-plus";
import QrcodeVue from 'qrcode.vue'
import { copy } from '@/utils/index'
import {useUserStore} from "@/stores/modules/user";
import {ElMessage} from "element-plus";
import { cloearsysGoogle, setsysGoogle, userGetGoogleKey } from "@/api/modules/system";

const userStore = useUserStore();
const userInfo = userStore.userInfo

const dialogFormVisible_setGoogle = ref(false)
const dialogFormVisible_closeGoogle = ref(false)
const formUserRef_set = ref()
const formUserRef_del = ref()

// 表单规则校验
const rules = ref({
  originalGoogleCode: [{required: true, message: '不能为空', trigger: 'change',}],
  googleCode: [{required: true, message: '不能为空', trigger: 'change',}]
})

const googleKeyCode = ref("") //谷歌验证码
const googleQrcode = ref("") //谷歌验证码二维码


const form_user = ref({
  googleKey: "",  //  string  谷歌密钥  必需
  googleCode: "",  //  string  谷歌验证码  必需
  originalGoogleCode: "",  //  string  原谷歌验证码

})

onMounted(() => {
  form_user.value = userStore.userInfo
})

// 查看谷歌验证码
const showSetGoogle = async () => {
  let res = await getKeyCode()
  if (res) {
    googleKeyCode.value = res.googleKey
    googleQrcode.value = res.qrcode

    dialogFormVisible_setGoogle.value = true
  }
}

// 查看谷歌验证码
const getKey = async () => {
  let res = await getKeyCode()
  if (res) {
    googleKeyCode.value = res.googleKey
    googleQrcode.value = res.qrcode
  }

}

// 获取谷歌验证码接口
const getKeyCode = () => {
  return new Promise((resolve, reject) => {
    userGetGoogleKey({}).then(res_ => {
      if (res_.code == 200) {
        resolve(res_.data)
      } else {
        reject()
        ElMessage.error(res.msg)
      }
    })
  })
}

//管理员弹框取消
const popupUserCancal_set = (formEl) => {
  formEl.resetFields()
  dialogFormVisible_setGoogle.value = false

}

// 确定设置谷歌验证码
const popupUserConfirm_set = (formEl) => {
  if (!formEl) return
  formEl.validate((valid) => {
    if (!valid) return
    let {googleKey, googleCode, originalGoogleCode} = form_user.value
    let params = {
      googleKey: googleKeyCode.value,
      googleCode,
      originalGoogleCode
    }
    setsysGoogle(params).then(res => {
      if (res.code == 200) {
        ElNotification.success("操作成功")
        formEl.resetFields()
        dialogFormVisible_setGoogle.value = false

      } else {
        ElMessage.error(res.msg)
      }

    })

  })
}

// 取消谷歌验证码
const popupUserCancal_del = (formEl) => {
  formEl.resetFields()
  dialogFormVisible_closeGoogle.value = false
}

// 确定取消谷歌验证码
const popupUserConfirm_del = (formEl) => {
  if (!formEl) return
  formEl.validate((valid) => {
    if (!valid) return
    let {googleCode} = form_user.value
    let params = {
      googleCode
    }
    cloearsysGoogle(params).then(res => {
      if (res.code == 200) {
        ElNotification.success("操作成功")
        formEl.resetFields()
        dialogFormVisible_closeGoogle.value = false
      } else {
        ElMessage.error(res.msg)
      }
    })
  })
}

// 开启谷歌绑定/修改弹窗
const openDialog = () => {
  dialogFormVisible_setGoogle.value = true;
  showSetGoogle()
};

// 开启删除谷歌验证码弹窗
const closeGoogleAuthDialog = () => {
  dialogFormVisible_closeGoogle.value = true;
};

// 暴露方法
defineExpose({ openDialog , closeGoogleAuthDialog});
</script>

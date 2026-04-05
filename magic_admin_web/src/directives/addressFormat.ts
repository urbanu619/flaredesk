// 自定义指令，加密地址，文字蓝色带下划线，点击可以复制完整内容
import { ElMessage } from 'element-plus'

export default {
  mounted(el: HTMLElement, binding: { value: string }) {
    const copyStr = String(binding.value)
    if (!copyStr || copyStr.length < 10) {
      el.textContent = String(copyStr) || ''
      return
    }

    // 加密后的地址
    const formatted = `${copyStr?.slice(0, 6)}...${copyStr?.slice(-4)}`

    // 设置样式
    el.style.color = '#409EFF' // Element Plus 主色调
    el.style.textDecoration = 'underline'
    el.style.cursor = 'pointer'

    // 创建 tooltip
    el.setAttribute('title', copyStr)

    // 设置内容
    el.textContent = formatted

    // 点击复制
    el.addEventListener('click', () => {
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(copyStr).then(() => {
          ElMessage.success('复制成功')
        }).catch(() => {
          ElMessage.error('复制失败')
        })
      } else {
        console.warn('当前环境不支持 clipboard API');
        // 兼容旧浏览器
        // const textarea = document.createElement('textarea')
        // textarea.value = copyStr
        // textarea.style.position = 'fixed'
        // textarea.style.opacity = '0'
        // document.body.appendChild(textarea)
        // textarea.select()
        // try {
        //   document.execCommand('copy')
        //   ElMessage.success('复制成功')
        // } catch (e) {
        //   ElMessage.error('复制失败')
        // }
        // document.body.removeChild(textarea)
        try {
            console.log(copyStr);
          if (!copyStr) {
            ElMessage.error("复制失败")
            return
          }
          const save = function (e) {
              e.clipboardData.setData("text/plain", copyStr);
              e.preventDefault();                     // 阻止默认行为
            };
            document.addEventListener("copy", save);  // 添加一个copy事件
            document.execCommand("copy");             // 执行copy方法
            // 复制成功提示
            ElMessage.success("复制成功")
          } catch (error) {
            ElMessage.error("复制失败")
        }
      }
    })
  }
}

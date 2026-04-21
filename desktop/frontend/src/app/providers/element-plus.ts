import type { App } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'

export function setupElementPlus(app: App): void {
  app.use(ElementPlus, { locale: zhCn, size: 'small' })
}

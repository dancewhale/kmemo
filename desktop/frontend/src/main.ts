import { createApp } from 'vue'
import App from './App.vue'
import { router } from './app/router'
import { pinia } from './app/store'
import { setupElementPlus } from './app/providers/element-plus'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import { useSettingsStore } from '@/modules/settings/stores/settings.store'
import { ROUTE_NAMES } from '@/shared/constants/routes'
import '@/app/styles/index.scss'

const app = createApp(App)
app.use(pinia)
const workspace = useWorkspaceStore()
workspace.restoreLayoutFromStorage()
const settings = useSettingsStore()
settings.initialize()
app.use(router)
setupElementPlus(app)
app.mount('#app')

void router.isReady().then(async () => {
  if (router.currentRoute.value.name === ROUTE_NAMES.settings) {
    return
  }
  const targetContext = settings.getStartupContext()
  if (router.currentRoute.value.name !== ROUTE_NAMES[targetContext]) {
    await router.replace({ name: ROUTE_NAMES[targetContext] })
  }
  workspace.setContext(targetContext)
})

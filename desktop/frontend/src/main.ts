import { createApp } from 'vue'
import App from './App.vue'
import { router } from './app/router'
import { pinia } from './app/store'
import { setupElementPlus } from './app/providers/element-plus'
import { setupShortcuts } from './app/providers/shortcuts'
import { useWorkspaceStore } from '@/modules/workspace/stores/workspace.store'
import '@/app/styles/index.scss'

const app = createApp(App)
app.use(pinia)
useWorkspaceStore().restoreLayoutFromStorage()
app.use(router)
setupElementPlus(app)
setupShortcuts()
app.mount('#app')

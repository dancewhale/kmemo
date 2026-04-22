import type { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import BlankLayout from '@/layouts/BlankLayout.vue'
import ReadingPage from '@/pages/ReadingPage.vue'
import KnowledgePage from '@/pages/KnowledgePage.vue'
import SettingsPage from '@/pages/SettingsPage.vue'
import { ROUTE_NAMES, ROUTE_PATHS } from '@/shared/constants/routes'

export const routes: RouteRecordRaw[] = [
  {
    path: ROUTE_PATHS.root,
    component: MainLayout,
    children: [
      { path: '', redirect: { name: ROUTE_NAMES.reading } },
      { path: 'reading', name: ROUTE_NAMES.reading, component: ReadingPage },
      { path: 'knowledge', name: ROUTE_NAMES.knowledge, component: KnowledgePage },
      { path: 'inbox', redirect: { name: ROUTE_NAMES.reading } },
      { path: 'review', redirect: { name: ROUTE_NAMES.reading } },
      { path: 'search', redirect: { name: ROUTE_NAMES.reading } },
    ],
  },
  {
    path: ROUTE_PATHS.settings,
    component: BlankLayout,
    children: [{ path: '', name: ROUTE_NAMES.settings, component: SettingsPage }],
  },
]

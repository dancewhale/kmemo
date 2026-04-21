import type { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import BlankLayout from '@/layouts/BlankLayout.vue'
import InboxPage from '@/pages/InboxPage.vue'
import ReadingPage from '@/pages/ReadingPage.vue'
import KnowledgePage from '@/pages/KnowledgePage.vue'
import ReviewPage from '@/pages/ReviewPage.vue'
import SearchPage from '@/pages/SearchPage.vue'
import SettingsPage from '@/pages/SettingsPage.vue'
import { ROUTE_NAMES, ROUTE_PATHS } from '@/shared/constants/routes'

export const routes: RouteRecordRaw[] = [
  {
    path: ROUTE_PATHS.root,
    component: MainLayout,
    children: [
      { path: '', redirect: { name: ROUTE_NAMES.reading } },
      { path: 'inbox', name: ROUTE_NAMES.inbox, component: InboxPage },
      { path: 'reading', name: ROUTE_NAMES.reading, component: ReadingPage },
      { path: 'knowledge', name: ROUTE_NAMES.knowledge, component: KnowledgePage },
      { path: 'review', name: ROUTE_NAMES.review, component: ReviewPage },
      { path: 'search', name: ROUTE_NAMES.search, component: SearchPage },
    ],
  },
  {
    path: ROUTE_PATHS.settings,
    component: BlankLayout,
    children: [{ path: '', name: ROUTE_NAMES.settings, component: SettingsPage }],
  },
]

import { createRouter, createWebHashHistory } from "vue-router";

import WorkspacePage from "../pages/workspace/WorkspacePage.vue";
import ReviewPage from "../pages/review/ReviewPage.vue";
import SettingsPage from "../pages/settings/SettingsPage.vue";

export const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      redirect: "/workspace",
    },
    {
      path: "/workspace",
      name: "workspace",
      component: WorkspacePage,
    },
    {
      path: "/review",
      name: "review",
      component: ReviewPage,
    },
    {
      path: "/settings",
      name: "settings",
      component: SettingsPage,
    },
  ],
});

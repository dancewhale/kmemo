import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import { router } from "./router";
import "./styles/reset.css";
import "./styles/theme.css";
import "./styles/layout.css";

if (import.meta.env.DEV && typeof window !== "undefined" && !window.go?.main?.App) {
  // 仅开发：直连 Vite 时无 Wails 注入，避免用户误以为「已运行 wails」却仍无桥接。
  console.warn(
    "[kmemo] 未检测到 Wails 桥接（window.go.main.App）。若你正在 Chrome 中打开的是 Vite 地址（例如 :9245），请改用 wails 弹出的窗口，或终端里 Wails 开发服务器地址（常见 :34115）。",
  );
}

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.mount("#app");

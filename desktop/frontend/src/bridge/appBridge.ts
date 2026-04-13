type AppBridge = Record<string, (...args: any[]) => Promise<any>>;

declare global {
  interface Window {
    go?: {
      main?: {
        App?: AppBridge;
      };
    };
  }
}

function bridgeUnavailableMessage(): string {
  const base =
    "Wails 桥接不可用：当前页面没有经过 Wails 的注入流程，无法调用 Go 后端。";
  if (import.meta.env.DEV) {
    return (
      base +
      " 开发时请勿在 Chrome 中直接打开 Vite 端口（本项目为 http://127.0.0.1:9245）。" +
      " 请使用「wails dev」或「task run:wails」打开的原生窗口；若要用浏览器调试，请在终端输出里找到 Wails 开发服务器地址（常见为 http://127.0.0.1:34115）并用该地址打开。"
    );
  }
  return base + " 请使用「wails build」产物或「wails dev」窗口运行本应用。";
}

function getBridge(): AppBridge {
  const bridge = window.go?.main?.App;
  if (!bridge) {
    throw new Error(bridgeUnavailableMessage());
  }
  return bridge;
}

export function callApp<T>(method: string, ...args: any[]): Promise<T> {
  const bridge = getBridge();
  const fn = bridge[method];
  if (typeof fn !== "function") {
    throw new Error(`Wails bridge method not found: ${method}`);
  }
  return fn(...args) as Promise<T>;
}

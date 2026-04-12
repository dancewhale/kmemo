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

function getBridge(): AppBridge {
  const bridge = window.go?.main?.App;
  if (!bridge) {
    throw new Error("Wails bridge is unavailable. Please run the app inside Wails.");
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

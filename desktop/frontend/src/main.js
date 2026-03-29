const root = document.getElementById("app");

const title = document.createElement("h1");
title.textContent = "kmemo";

const subtitle = document.createElement("p");
subtitle.textContent = "渐进阅读 / SuperMemo 风格 — 当前为工程骨架（无业务逻辑）";

const lineVersion = document.createElement("p");
const linePython = document.createElement("p");
lineVersion.id = "line-version";
linePython.id = "line-python";

root.append(title, subtitle, lineVersion, linePython);

function bindGo() {
  const bridge = window.go?.main?.App;
  if (!bridge) {
    lineVersion.textContent = "Go 绑定：在 Wails 窗口中打开本页（task run:wails）";
    linePython.textContent = "";
    return;
  }
  bridge
    .GetVersion()
    .then((v) => {
      lineVersion.textContent = `Go 版本字符串：${v}`;
    })
    .catch(() => {
      lineVersion.textContent = "Go 绑定存在，但 GetVersion 调用失败";
    });
  bridge
    .PythonEndpoint()
    .then((addr) => {
      linePython.textContent = `Python gRPC 地址（配置）：${addr}`;
    })
    .catch(() => {
      linePython.textContent = "";
    });
}

bindGo();

以下是可采集的浏览器与操作系统信息清单，按类别以比较格式输出：

---

### **一、浏览器信息采集清单**

| **类别**       | **具体信息项**               | **说明/示例**                     |
| -------------------- | ---------------------------------- | --------------------------------------- |
| **基础信息**   | 浏览器名称                         | Chrome, Firefox, Safari, Edge, Opera 等 |
|                      | 浏览器版本                         | 如 Chrome 120.0.0.0                     |
|                      | 用户代理字符串（User Agent）       | 包含浏览器、操作系统、设备等原始信息    |
|                      | 语言设置                           | 浏览器偏好语言（如 zh-CN, en-US）       |
| **功能支持**   | Cookie 是否启用                    | Boolean 值                              |
|                      | JavaScript 是否启用                | Boolean 值                              |
|                      | LocalStorage / SessionStorage 支持 | Boolean 值                              |
|                      | 屏幕色彩深度                       | 如 24 bit, 30 bit                       |
|                      | 时区                               | 如 Asia/Shanghai                        |
| **屏幕与视口** | 屏幕分辨率                         | 如 1920×1080                           |
|                      | 可用屏幕尺寸                       | 排除任务栏等系统界面                    |
|                      | 视口尺寸（viewport）               | 页面可视区域大小                        |
|                      | 设备像素比（DPR）                  | 物理像素与逻辑像素比例                  |
| **网络与连接** | 在线状态                           | navigator.onLine                        |
|                      | 网络类型（部分支持）               | 4g, wifi（通过 Navigator.connection）   |
| **其他**       | 插件列表（如 Flash, PDF Viewer）   | navigator.plugins（部分浏览器已限制）   |
|                      | 字体列表（需 JS 检测）             | 通过 Canvas 检测字体支持                |
|                      | 浏览器内核版本                     | 如 WebKit 537.36, Gecko 109.0           |

---

### **二、操作系统信息采集清单**

| **类别**       | **具体信息项** | **说明/示例**                                                   |
| -------------------- | -------------------- | --------------------------------------------------------------------- |
| **基础信息**   | 操作系统名称         | Windows, macOS, Linux, Android, iOS, iPadOS 等                        |
|                      | 操作系统版本         | 如 Windows 10, macOS 14.0, Android 13                                 |
|                      | 设备类型             | Desktop, Mobile, Tablet, TV, Console                                  |
| **架构与内核** | 系统架构             | x86, x64, arm, arm64（从 User Agent 或 navigator.userAgentData 解析） |
|                      | 内核版本（部分）     | 如 Linux 内核版本（通过 User Agent）                                  |
| **时区与区域** | 系统时区             | Intl.DateTimeFormat().resolvedOptions().timeZone                      |
|                      | 系统区域设置         | 如 zh-CN, en-GB                                                       |
| **其他**       | 触摸屏支持           | 最大触摸点数（navigator.maxTouchPoints）                              |
|                      | 电池状态（部分支持） | 充电状态、电量（navigator.getBattery，需 HTTPS）                      |
|                      | CPU 核心数           | navigator.hardwareConcurrency                                         |

---

### **三、可补充采集的设备信息**

1. **设备型号**（移动端常见）：
   - 如 iPhone 15 Pro, Samsung Galaxy S24（通过 User Agent 解析）。
2. **GPU 信息**（通过 WebGL 渲染器字符串）：
   - 如 “Apple GPU”, “NVIDIA GeForce RTX 4090”。
3. **输入设备支持**：
   - 鼠标、触摸屏、陀螺仪、摄像头、麦克风等。
4. **浏览器指纹特征**（需组合多维度信息）：
   - Canvas 指纹、WebGL 指纹、音频指纹等。

---

### **注意事项**

- **隐私限制**：现代浏览器（如 Safari、Firefox）会限制部分信息获取（如插件列表、精确字体列表）。
- **User Agent 不可全信**：可被修改或屏蔽，建议结合多维度检测。
- **合规要求**：采集信息需遵守 GDPR、CCPA 等隐私法规，通常需用户知情同意。

---

可根据实际需求选择采集字段，建议优先使用标准化 API（如 `navigator.userAgentData` 替代部分 User Agent 解析）。

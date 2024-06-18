# gpt-sovits-backend
## 使用须知
- 需要掌握 autodl 部署tts服务
- 使用者需要掌握golang的编译运行
- 需要 **gpt-sovits** 的 **autodl** 上的服务

## 操作步骤
### autodl
- 需要上传 `sovits_weight` 和 `gpt_weight `
- **命名格式**: 里面必须包含角色名，不可有重复的角色名，如 `派蒙-平静-e10.ckpt`
- 启动`tts`推理界面

### 代码
- 在项目根目录添加 `wav` 文件夹
- 文件夹内是参考音频 格式为`角色名+说的话.wav`
比如`派蒙+既然罗莎莉亚说足迹上有元素力，用元素视野应该能很清楚地看到吧。`
- 在根目录下添加 `config.yaml` 文件 比如
```
tts:
  GRADIO_URL: xxxx.gradio.live
```

### 编译和运行
- 需要有golang环境 > 1.18
- 根目录下执行 `go mod tidy && go build -o tts main.go` , 然后 `./tts`

### TODO
* [x] 支持多自定义多角色
* [ ] 优化hash逻辑
* [ ] docker一键部署
* [ ] 支持较高并发
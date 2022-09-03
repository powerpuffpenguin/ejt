# ejt

將 jsonnet 轉換爲 yaml/json 的小工具

[English](README.md) 中文

# Background

本喵想在一些小型環境下使用 envoy 替換 nginx，這些環境不適合安裝複雜的 envoy 控制面板因爲這些環境很簡單或輕量，但也需要足夠的靈活性可以隨時動態的改變設定。envoy 提供了 dynamic-resources-from-filesystem 很適合，然而直接使用 yaml 書寫配置有諸多問題，例如 無法將複雜配置分散到單個小的檔案中使用 import/include 之類的技術導入，大部分類似的 listenr/cluster 每次都重複填寫一樣的屬性值，並且 envoy 只在 move 時才會重新加載 dynamic-resources 指定的內容。

jsonnet 可以很好解決這些 yaml 的不足，於是本喵寫了這個小工具用於將 jsonnet 轉換到 yaml，並自動將 yaml 複製到 envoy dynamic-resources 監控的目標路徑。

此外因爲 jsonnet/json/yaml 結構一致所以也順便支持了 jsonnet 到 json 的轉換。

# 如何使用

創建 **ejt.jsonnet** 定義檔案，你可以執行下述指令在當前工作目錄下創建定義檔案：

```
ejt init
```

**ejt.jsonnet** 定義了要從哪裏獲取 jsonnet，將轉換代碼輸出到哪裏，以及將輸出檔案移動或複製到哪裏以便類似 envoy 的監控系統觸發更新。最終看起來像這樣
```
{
  version: 'v0.0.1',
  endpoints: [
    {
      output: './dst',  // redirect output structure to the directory.
      target: './envoy',  // target root directory.
      source: './src',  // source root directory.
      resources: [
        'envoy.jsonnet',
      ],
    },
  ],
}
```

編寫好你的 jsonnet 後執行下述指令生成 yaml 並且移動到目標路徑
```
ejt yaml -m
```
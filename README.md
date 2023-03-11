# Dcard Backend Intern 2023 Assignment
本系統旨在提供一個通用的列表系統，使得各種推薦系統可以方便的呼叫這個系統儲存以及瀏覽排序結果。
## 測試環境
Arch Linux 6.1.2-arch1-1
## 使用方法
### 安裝
安裝套件以及 ProtoBuf 請參考 `make help`
### 環境變數
將 `.env.example` 複製到 `.env`。從 `.env` 載入環境變數請使用以下指令：
```bash
source env.sh
```
### 測試
```bash
make test
```
## APIs
### 查詢用 API
#### `GET /list/:key`
查詢列表的編號，等同於 getHead 的功能。透過列表回傳的 `id` 以及 `version` 可以組合查詢出最新的列表。
##### Parameters
- key
    - string
    - 列表的鍵，唯一值。
##### Response
回傳為 JSON 格式。
- id
    - number (uint)
    - 列表的 ID
- version
    - number (uint32)
    - 列表目前的版本。

#### `GET /list/:id/:version`
##### Parameters
- id
    - number (uint)
    - 列表的 id，使用上述 API 可以由列表 key 查詢。
- version
    - number (uint32)
    - 列表的最新版本編號，使用上述 API 可以由列表 key 查詢。

### 設定用 API
設定用 API 由 gRPC 組成。以下是系統預期的處理流程：
1. 取得列表 ID，可以夠過 `Head` 方法。
2. 產生 version 編號。Version 只要不與存在資料庫中的先前版本重複即可，例如可以使用 `time`，或是使用原先列表的 version + 1。
3. 呼叫 `SetList` 將項目資料，即 `ArticleID` 與排列順序 `NodeOrder` 新增至資料庫中。`NodeOrder` 為遞增排序 (ASC)。不建議在同個版本中使用重複的 `NodeOrder`。這個 API 需要傳入先前取得的 `ListID` 以及 `Version`。
4. 呼叫 `SetListVersion` 將列表更新為正式版本。

其中 SetList 傳送跟接收都使用 stream，可以建立長連線，處理多批次資料，避免重複建立連線。

## 其他程式說明
### Serializer
根據定義
```go
type Serializable[t any] interface {
	Serialize() t
}
```
`Serializable[t]` 需要物件實做 `Serialize` 方法，將自身序列化為 `t` 型態，`t` 就是要呈現在 client 的格式。而 `CursorPagedSerializer` 以及對應的 `SerializeCursorPagedItems` 則可以直接傳入 cursor 以及資料結構體本身（與資料庫溝通的），呼叫 Serialize 方法，獲得最後回傳的結果。
雖然在資料模型很少的情況下這麼設計有點過度設計的嫌疑，不過我的情境是假設有很多的資料模型需要做成列表（就是那個 Article）。`GetListNodes` 方法也可以更改為透過傳入參數的方法，替換成需要使用的模型。

### 資料庫的選用
資料在設計上涉及頻繁的寫入操作，**寫入的操作是預期要比讀取多的**。會這麼假定是因為題目中假設一小時會更新一次推薦列表，而且每個人的推薦清單是個人化的，但是使用者很大機率不會每小時看超過一次的完整清單。基於此，我認為對於資料做反正規劃就比較沒有必要，不用將 Article 內容一併儲存進去 ListNode 中。

基於上述假設，由於不需要儲存資料本身，只需要儲存 Foreign Key，資料是統一格式的，我認為使用關聯式資料庫比較合適。考慮到我比較熟悉 PostgreSQL，所以就使用 Postgres。過期刪除的部份，雖然使用 MongoDB 可以設定 TTL 鍵，但若是在每小時更新，每天刪除過期資料的前提下，刪除的資料應該遠大於不須刪除的資料。這種情況下 index 似乎比較沒有這麼大優勢。
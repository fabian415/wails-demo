# 概要

需求：介面提供 8 個選項（有的是「檔案」、有的是「資料夾」、有的是「盤符」），使用者透過檔案/資料夾/盤符選擇器選定後，把所選項的絕對路徑寫入一個配置檔 build.conf，只寫入被「選中」的項目。這是為後續 build 步驟準備資料。

下面依序說明每個選項的意義、UI 建議、驗證規則、寫入 build.conf 時的細節（格式、轉義、原子寫入）、以及 build 階段的使用檢查。

## 每個選項的詳細說明與驗證建議
### (1) ISO（文件）
含義：指向一個光碟映像檔（常見副檔名 .iso、有時 .img）。
UI：檔案選擇器（只允許檔案），可設定 filter（.iso, .img）。
驗證：
檔案是否存在且可讀（exists && readable）。
副檔名驗證（非強制，但建議提示）。
（可選）檢查檔案大小是否合理（> 1MB）。

寫入建議：單一路徑字串。

### (2) DEB（文件夾）
含義：一個資料夾，內含 .deb 套件或相關套件集合。
UI：資料夾選擇器（只允許資料夾）。
驗證：
資料夾存在且可讀。
可檢查資料夾內是否至少有一個 .deb（若你的流程需要）。
寫入建議：單一路徑（若支援多個套件來源，可用資料夾或支援陣列）。

### (3) LOGO（文件）
含義：單一圖檔（品牌 logo），例如 PNG、SVG、JPG。
UI：檔案選擇器（允許圖片副檔名）。
驗證：
檔案存在且可讀。
副檔名或 MIME 類型檢測（可選）。
（可選）尺寸或解析度檢查。
寫入建議：單一路徑字串。

### (4) 驱动（文件夹）
含義：放驅動檔的資料夾（例如驅動安裝程式、Windows 驅動或 Linux 模組）。
UI：資料夾選擇器。
驗證：
資料夾存在。
檔案權限（可讀）。
（可選）檢查常見檔案類型（.inf/.sys/.ko 等）是否存在，視需求而定。
寫入建議：單一路徑或陣列（視多驅動來源需求）。

### (5) 预执行脚本（文件）
含義：build 前要執行的腳本檔（shell script、python 等）。
UI：檔案選擇器（單檔）。
驗證：
檔案存在且可讀。
對 *nix：應檢查是否可執行（chmod +x 或提醒使用者設定執行權限）。
檢查 shebang（#!）或檔案類型（可選）。
寫入建議：單一路徑字串；build 執行前再檢查並且不要直接執行沒有驗證的腳本（安全性考量）。

### (6) 自启动服务（文件夹）
含義：包含 systemd unit、init script 或可放入 image 的自啟動服務定義的資料夾。
UI：資料夾選擇器。
驗證：
資料夾存在。
可選：檢查是否含 .service 或預期格式。
寫入建議：單一路徑（代表整個 services 資料夾）。如果需要分多個 service，build 階段可掃描該資料夾。

### (7) 其他文件（文件夹）
含義：其它需要一併打包的檔案集合（config、資源、第三方工具等）。
UI：資料夾選擇器。
驗證：資料夾存在、可讀。
寫入建議：單一路徑或陣列（視需求）。

### (8) 启动盘（盘符）
含義：代表要寫入或使用的實體磁碟 / 隨身碟（Windows 表示盤符 E:，Linux 可能是 /dev/sdb 或 mount point /media/xxx）。
UI：
在桌面環境可列出可用磁碟清單（使用系統 API 列出已掛載或未掛載的 block devices）。
或讓使用者輸入裝置路徑/盤符。
驗證：
在 Windows：檢查盤符是否存在。
在 Linux：檢查該路徑是否為 block device（/dev/sdX），或檢查是否為 mount point。

## 強烈建議在 UI 顯示警告（危險操作：將會格式化/寫入）並要求二次確認。
寫入建議：記錄裝置識別（優先用 device node /dev/sdb 或 UUID；若記錄 E:，在不同機器上可能不一致）。

## UI 行為與互動細節（建議）
每個選項呈現方式：checkbox + 「選擇」按鈕 + 只讀顯示欄顯示已選路徑 + clear（清除）按鈕。

Checkbox 為「啟用/選中」，只有勾選時才把該項寫入 build.conf。

按「選擇」開啟檔案/資料夾/盤符選擇器（根據項目類型限定選擇器）。

顯示絕對路徑：選完後立即把經過 realpath 或 abspath 的絕對路徑顯示給使用者。

實時驗證：選取後立刻檢查是否存在、類型是否正確（檔案/資料夾/裝置），並用綠勾/紅叉或文字提示。

提示/tooltip：每一項提供小說明（ex: ISO 必須為映像檔；啟動盤選擇會格式化等）。

多平台：對 Windows 路徑顯示 E:\...，對 Linux 顯示 /home/... 或 /dev/sdb。

## 儲存/更新邏輯：

使用者按「保存設定」或自動儲存時：對應欄位若勾選則寫入；若未勾選則從 build.conf 移除該 key（或保持注釋）。

若使用者變更已存在路徑，覆寫 build.conf 中該 key 的值。

安全確認：對「啟動盤」或「預執行腳本」等危險操作，要二次確認或 require admin privileges。

## build.conf 的格式與範例
選擇格式要點：
若只需簡單 key=value，INI-style（KEY="value"）最直觀。
若需要陣列、多個條目、或跨平台解析，建議 JSON 或 YAML（JSON 可被多種語言快速解析，並支援陣列與物件結構）。
下面提供兩種示例。

### A. INI-like (簡單，易讀)
```cmd!
# UTF-8 編碼
ISO="/home/user/ubuntu-24.04.iso"
DEB_FOLDER="/home/user/debs"
LOGO="/home/user/assets/logo.png"
DRIVER_FOLDER="/home/user/driver_package"
PRE_SCRIPT="/home/user/scripts/pre_build.sh"
AUTOSTART_FOLDER="/home/user/autostart_services"
OTHER_FOLDER="/home/user/extra_files"
BOOT_DEVICE="/dev/sdb"
```

值用雙引號包起，方便包含空白與特殊字元（寫入時要 escape 引號）。
不存在或未選中的項目，可以完全不出現在檔案中。


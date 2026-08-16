# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

对完全相同的一批事件重复跑 detect，incident 数量没变，但同一个 incident 的 count 每跑一次就翻倍增长，而它关联的事件 ID 列表始终只有那几条，看起来计数和真实证据对不上。先不要修改代码。请调查计数为什么会被重复累加，给出可核验证据、完整因果链，并定位具体 Go 文件和符号。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/gogo-49
- 仓库地址：https://github.com/zhanglei10281852-gif/gogo-49.git
- parent SHA：5f2d6d48b024a4152094f41c24ae1bfa43136815

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/gogo-49.git bug-repro
cd bug-repro
git checkout --detach 5f2d6d48b024a4152094f41c24ae1bfa43136815
go test ./internal/detect -run "^TestRepeatedDetectionDoesNotInflateIncidentCount$" -count=1 -v
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/detect -run "^TestRepeatedDetectionDoesNotInflateIncidentCount$" -count=1 -v
=== RUN   TestRepeatedDetectionDoesNotInflateIncidentCount
    repeat_regression_test.go:29: reprocessing unchanged events inflated incident: [{ID:49a89d9e4919ce4c01e7a6c8 RuleID:errors Title:Repeated errors Status:open Severity:ERROR Group:source=api Fingerprint:ff2bf27fc3ddbe88eee031f20972997f16016321366b0fd59edc7a49411d1676 FirstSeen:2025-01-01 00:00:00 +0000 UTC LastSeen:2025-01-01 00:01:00 +0000 UTC OpenedAt:2025-01-01 00:02:00 +0000 UTC UpdatedAt:2025-01-01 00:02:00 +0000 UTC ResolvedAt:<nil> Count:4 EventIDs:[346c71bb9ac29474166aa4bd 6c29bb97da8876bd7b3fb382] Labels:map[] History:[{At:2025-01-01 00:02:00 +0000 UTC From: To:open Reason:threshold 2 reached in 5m}]}]
--- FAIL: TestRepeatedDetectionDoesNotInflateIncidentCount (0.00s)
FAIL
FAIL	LogPilot/internal/detect	0.002s
FAIL

```

stderr：

```text
(empty)
```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/detect -run "^TestRepeatedDetectionDoesNotInflateIncidentCount$" -count=1 -v
=== RUN   TestRepeatedDetectionDoesNotInflateIncidentCount
    repeat_regression_test.go:29: reprocessing unchanged events inflated incident: [{ID:49a89d9e4919ce4c01e7a6c8 RuleID:errors Title:Repeated errors Status:open Severity:ERROR Group:source=api Fingerprint:ff2bf27fc3ddbe88eee031f20972997f16016321366b0fd59edc7a49411d1676 FirstSeen:2025-01-01 00:00:00 +0000 UTC LastSeen:2025-01-01 00:01:00 +0000 UTC OpenedAt:2025-01-01 00:02:00 +0000 UTC UpdatedAt:2025-01-01 00:02:00 +0000 UTC ResolvedAt:<nil> Count:4 EventIDs:[346c71bb9ac29474166aa4bd 6c29bb97da8876bd7b3fb382] Labels:map[] History:[{At:2025-01-01 00:02:00 +0000 UTC From: To:open Reason:threshold 2 reached in 5m}]}]
--- FAIL: TestRepeatedDetectionDoesNotInflateIncidentCount (0.02s)
FAIL
FAIL	LogPilot/internal/detect	0.163s
FAIL

```

stderr：

```text
(empty)
```

## 通过条件

定位 internal/model/model.go 的 (*Incident).Merge，并说明它与 internal/detect/detect.go 中 reconcile 的调用关系；解释 count 被累加而事件 ID 去重后从未回算的完整机制；有证据且目标仓库零改动。

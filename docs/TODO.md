# gofsm-gen 開発タスク管理

vibe-kanbanを使用した開発タスク管理ボードです。

## 📈 プロジェクト進捗状況

**最終更新**: 2025-11-23

### 全体進捗
- **Phase 1**: 🚧 進行中 (約50%完了)
  - ✅ プロジェクト基盤: 完了
  - ✅ 内部モデル実装: 完了（テストカバレッジ97.4%）
  - ✅ ドキュメント: 完了
  - ⏳ パーサー/ジェネレーター: 未着手
  - ⏳ CLI実装: 未着手
- **Phase 2**: ⏸️ 未開始
- **Phase 3**: ⏸️ 未開始
- **Phase 4**: ⏸️ 未開始

### 次のステップ
Phase 1の内部モデル実装とドキュメント作成が完了したため、以下のタスクが**ブロック解除**されました：

**すぐに着手可能なタスク** (優先度順):
1. 🎯 **P0** YAMLパーサーの実装 (pkg/parser/yaml.go)
2. 🎯 **P0** コード生成器インターフェースの定義 (pkg/generator/generator.go)
3. 🎯 **P0** Validator構造体の実装 (pkg/analyzer/validator.go)
4. ⚡ **P0** state_machine.tmplテンプレートの作成
5. ⚡ **P1** exhaustive統合の調査と実装

---

## 凡例

- 🔒 **Blocked by**: このタスクは他のタスクに依存しています
- ⚡ **Parallel OK**: 他のタスクと並列実行可能
- 🎯 **Blocker**: 他のタスクがこのタスクに依存しています
- ⏱️ **Quick Win**: 短時間で完了可能なタスク

---

## 📋 Backlog

### Phase 1: 基本機能実装

#### 🎯 ブロッキングタスク（他のタスクの前提条件）

**プロジェクトセットアップ**
- [x] 🎯 プロジェクト構造の初期化
- [x] 🎯 go.modファイルの作成
- [x] ⚡ CI/CDパイプラインの構築
- [x] ⚡ 開発環境のセットアップガイド作成

**内部モデル実装** (基盤となるデータ構造)
- [x] 🎯 FSMModel構造体の実装 (pkg/model/fsm.go)
- [x] 🎯 State構造体の実装 (pkg/model/state.go)
- [x] 🎯 Event構造体の実装 (pkg/model/event.go)
- [x] 🎯 Transition構造体の実装 (pkg/model/transition.go)
- [x] ⚡ StateGraph構造体の実装 (pkg/model/graph.go)

#### ⚡ 並列実行可能タスク（依存関係が少ない）

**ドキュメント作成** (他の実装と並行可能)
- [x] ⚡⏱️ READMEの作成
- [x] ⚡⏱️ インストールガイドの作成 (docs/installation.md)
- [x] ⚡ 基本的な使い方ガイド (docs/usage.md)
- [x] ⚡ YAML定義リファレンス (docs/yaml-reference.md)
- [x] ⚡ API仕様書 (docs/api.md)
- [x] ⚡⏱️ exhaustive統合の調査 (docs/exhaustive-integration-investigation.md)
- [ ] ⚡⏱️ コントリビューションガイド (CONTRIBUTING.md)

**テンプレート作成** (Generator実装と並行可能)
- [ ] ⚡ state_machine.tmplテンプレートの作成
- [ ] ⚡ test.tmplテンプレートの作成
- [ ] ⚡ mock.tmplテンプレートの作成

#### 🔒 依存関係ありタスク

**パーサー実装** ✅ ブロック解除済み
- [ ] 🎯 YAMLパーサーの実装 (pkg/parser/yaml.go)
  - **優先度**: P0 (Critical)
  - **推定工数**: 4-6時間
  - **開発方針**: TDD（テストファースト）
  - **参考**: docs/yaml-reference.md に仕様あり
- [ ] 🔒 YAML定義構造体の作成
  - Blocked by: YAMLパーサー構造設計
- [ ] 🔒 YAMLパーサーのユニットテスト作成
  - Blocked by: YAMLパーサーの実装（TDDで同時進行）
- [ ] 🔒 パーサーエラーハンドリングの実装
  - Blocked by: YAMLパーサーの実装

**基本コード生成器** ✅ ブロック解除済み
- [ ] 🎯 コード生成器インターフェースの定義 (pkg/generator/generator.go)
  - **優先度**: P0 (Critical)
  - **推定工数**: 2-3時間
  - **開発方針**: インターフェース定義 + TDD
  - **参考**: docs/api.md に設計あり
- [ ] 🔒 基本コード生成器の実装 (pkg/generator/code_generator.go)
  - Blocked by: コード生成器インターフェース, state_machine.tmplテンプレート
  - **推定工数**: 6-8時間
- [ ] 🔒 テンプレートデータ準備ロジックの実装
  - Blocked by: コード生成器インターフェース
  - **推定工数**: 3-4時間
- [ ] 🔒 gofmtによるフォーマット処理の実装
  - Blocked by: 基本コード生成器の実装
  - **推定工数**: 1-2時間
- [ ] 🔒 生成コードのユニットテスト
  - Blocked by: 基本コード生成器の実装（TDDで同時進行）
  - **方針**: Golden fileテスト使用

**静的解析基盤** ✅ 部分的にブロック解除
- [ ] 🎯 網羅性チェッカーの実装 (pkg/analyzer/exhaustive.go)
  - **優先度**: P1 (High)
  - **推定工数**: 4-5時間
  - **参考**: docs/exhaustive-integration-investigation.md に調査結果あり
  - **方針**: exhaustiveツールとの統合実装
- [ ] 🔒 exhaustiveアノテーション自動挿入機能
  - Blocked by: 基本コード生成器の実装
  - **推定工数**: 2-3時間
- [ ] 🔒 静的解析のテスト作成
  - Blocked by: 網羅性チェッカーの実装
  - **推定工数**: 3-4時間

**モデル検証器** ✅ ブロック解除済み
- [ ] 🎯 Validator構造体の実装 (pkg/analyzer/validator.go)
  - **優先度**: P0 (Critical)
  - **推定工数**: 3-4時間
  - **開発方針**: TDD
  - **依存**: pkg/model（完了済み）
- [ ] 🔒 状態の検証ロジック
  - Blocked by: Validator構造体の実装
  - **推定工数**: 2時間
- [ ] 🔒 イベントの検証ロジック
  - Blocked by: Validator構造体の実装
  - **推定工数**: 2時間
- [ ] 🔒 遷移の検証ロジック
  - Blocked by: Validator構造体の実装
  - **推定工数**: 2-3時間
- [ ] 🔒 到達可能性の検証 (reachability analysis)
  - Blocked by: Validator構造体の実装
  - **依存**: StateGraph実装（完了済み）
  - **推定工数**: 3-4時間
- [ ] 🔒 決定性の検証 (determinism check)
  - Blocked by: Validator構造体の実装
  - **推定工数**: 2-3時間
- [ ] 🔒 重複遷移チェック
  - Blocked by: Validator構造体の実装
  - **推定工数**: 1-2時間
- [ ] 🔒 検証エラーレポート機能
  - Blocked by: Validator構造体の実装
  - **推定工数**: 2時間

**CLI実装**
- [ ] 🔒 CLIエントリポイントの実装 (cmd/gofsm-gen/main.go)
  - Blocked by: パーサー実装, コード生成器実装
- [ ] 🔒 コマンドライン引数のパース
  - Blocked by: CLIエントリポイント
- [ ] 🔒 ヘルプメッセージの実装
  - Blocked by: CLIエントリポイント
- [ ] 🔒 エラーメッセージの実装
  - Blocked by: CLIエントリポイント
- [ ] 🔒 詳細出力モード (-v flag)
  - Blocked by: CLIエントリポイント

**サンプルコード**
- [ ] 🔒 シンプルなドアロック状態機械のサンプル
  - Blocked by: CLI実装完了
- [ ] 🔒 Order管理状態機械のサンプル
  - Blocked by: CLI実装完了
- [ ] 🔒 examples/ディレクトリの構築
  - Blocked by: CLI実装完了

### Phase 2: 高度な機能実装

#### 🎯 ブロッキングタスク

**ガード/アクション機能** (Phase 1完了が前提)
- [ ] 🎯 Guard構造体の実装
  - Blocked by: Phase 1 内部モデル実装
- [ ] 🎯 Action構造体の実装
  - Blocked by: Phase 1 内部モデル実装
- [ ] 🔒 ガード条件評価ロジック
  - Blocked by: Guard構造体, 基本コード生成器
- [ ] 🔒 アクション実行ロジック
  - Blocked by: Action構造体, 基本コード生成器
- [ ] 🔒 エントリー/イグジットアクションの実装
  - Blocked by: Action構造体
- [ ] 🔒 非同期アクションのサポート
  - Blocked by: アクション実行ロジック
- [ ] 🔒 ガード/アクションのテスト
  - Blocked by: ガード/アクション実装

**Go DSLサポート**
- [ ] 🎯 Go ASTパーサーの実装 (pkg/parser/ast.go)
  - Blocked by: Phase 1 内部モデル実装
- [ ] 🔒 DSLパーサーの実装 (pkg/parser/dsl.go)
  - Blocked by: Go ASTパーサー
- [ ] 🔒 DSL用APIの設計
  - Blocked by: 内部モデル実装
- [ ] 🔒 Fluent API実装
  - Blocked by: DSL用API設計
- [ ] 🔒 DSLのサンプルコード作成
  - Blocked by: Fluent API実装
- [ ] 🔒 DSLパーサーのテスト
  - Blocked by: DSLパーサー実装

#### ⚡ 並列実行可能タスク

**視覚化機能** (基本モデルがあれば独立実装可能)
- [ ] ⚡ Mermaid生成器の実装 (pkg/visualizer/mermaid.go)
  - Blocked by: Phase 1 内部モデル実装のみ
- [ ] ⚡ Graphviz生成器の実装 (pkg/visualizer/graphviz.go)
  - Blocked by: Phase 1 内部モデル実装のみ
- [ ] 🔒 状態遷移図の自動生成
  - Blocked by: Mermaid/Graphviz実装
- [ ] 🔒 視覚化オプションの実装
  - Blocked by: CLI実装

**ランタイムサポート** (独立実装可能)
- [ ] ⚡ Loggerインターフェースの実装 (pkg/runtime/logger.go)
- [ ] ⚡ TransitionLoggerの実装
  - Blocked by: Loggerインターフェース
- [ ] ⚡ RuntimeValidatorの実装 (pkg/runtime/validator.go)
- [ ] ⚡ 実行時統計収集機能
  - Blocked by: RuntimeValidator
- [ ] ⚡ ガード競合検出ロジック
  - Blocked by: RuntimeValidator

#### 🔒 依存関係ありタスク

**テストコード自動生成**
- [ ] 🔒 テストジェネレーターの実装 (pkg/generator/test_generator.go)
  - Blocked by: Phase 1 基本コード生成器
- [ ] 🔒 test.tmplテンプレートの作成
  - Blocked by: Phase 1 テンプレート基盤
- [ ] 🔒 状態遷移テストの自動生成
  - Blocked by: テストジェネレーター
- [ ] 🔒 ガード条件テストの自動生成
  - Blocked by: ガード機能実装
- [ ] 🔒 エッジケーステストの自動生成
  - Blocked by: テストジェネレーター
- [ ] 🔒 生成されたテストコードの検証
  - Blocked by: テスト自動生成完了

**モック自動生成**
- [ ] 🔒 モックジェネレーターの実装 (pkg/generator/mock_generator.go)
  - Blocked by: Phase 1 基本コード生成器
- [ ] 🔒 mock.tmplテンプレートの作成
  - Blocked by: Phase 1 テンプレート基盤
- [ ] 🔒 Guard/Actionモックの生成
  - Blocked by: ガード/アクション機能
- [ ] 🔒 モックの使用例ドキュメント
  - Blocked by: モック生成完了

**パフォーマンス最適化**
- [ ] 🔒 ゼロアロケーションモードの実装
  - Blocked by: Phase 1 基本コード生成器
- [ ] ⚡ パフォーマンスベンチマークの作成 (benchmarks/)
- [ ] 🔒 状態遷移のベンチマーク
  - Blocked by: Phase 1 CLI完成
- [ ] 🔒 ガード付き遷移のベンチマーク
  - Blocked by: ガード機能実装
- [ ] 🔒 メモリプロファイリング
  - Blocked by: ベンチマーク作成
- [ ] 🔒 パフォーマンス目標の達成確認 (<50ns/transition)
  - Blocked by: 全ベンチマーク完了

### Phase 3: ツール統合

#### 🎯 ブロッキングタスク

**VSCode拡張機能**
- [ ] 🎯 VSCode拡張のプロジェクトセットアップ
- [ ] ⚡ YAML定義のシンタックスハイライト
  - Blocked by: プロジェクトセットアップのみ
- [ ] 🔒 定義ファイルの自動補完
  - Blocked by: Phase 1 YAMLパーサー
- [ ] 🔒 リアルタイムバリデーション
  - Blocked by: Phase 1 Validator
- [ ] 🔒 状態遷移図のプレビュー機能
  - Blocked by: Phase 2 視覚化機能
- [ ] 🔒 コード生成のショートカット
  - Blocked by: Phase 1 CLI完成
- [ ] 🔒 拡張機能のテスト
  - Blocked by: 拡張機能の全機能実装
- [ ] 🔒 VSCode Marketplaceへの公開
  - Blocked by: 拡張機能テスト完了

#### ⚡ 並列実行可能タスク

**gopls統合** (LSP調査は独立可能)
- [ ] ⚡⏱️ LSP統合の調査
- [ ] 🔒 型情報の提供
  - Blocked by: Phase 1 基本実装
- [ ] 🔒 定義へのジャンプ機能
  - Blocked by: Phase 1 基本実装
- [ ] 🔒 リファクタリングサポート
  - Blocked by: Phase 1 基本実装

**開発ツール** (独立実装可能)
- [ ] ⚡ ホットリロード機能
  - Blocked by: Phase 1 CLI完成のみ
- [ ] ⚡ ファイルウォッチャーの実装
- [ ] ⚡ デバッグモードの実装
- [ ] ⚡ 開発サーバーの実装
  - Blocked by: Phase 1 CLI完成のみ

**移行ツール** (独立実装可能)
- [ ] ⚡ looplab/fsm → gofsm-gen 移行ツール
  - Blocked by: Phase 1 YAMLパーサーのみ
- [ ] ⚡ qmuntal/stateless → gofsm-gen 移行ツール
  - Blocked by: Phase 1 YAMLパーサーのみ
- [ ] ⚡⏱️ 移行ガイドの作成
- [ ] ⚡ 移行スクリプトの作成
  - Blocked by: 移行ツール実装

### Phase 4: 高度な状態機械機能

#### 🎯 ブロッキングタスク

**階層的ステート**
- [ ] 🎯 親子状態モデルの実装
  - Blocked by: Phase 1 内部モデル実装
- [ ] 🔒 複合状態のサポート
  - Blocked by: 親子状態モデル
- [ ] 🔒 初期子状態の実装
  - Blocked by: 親子状態モデル
- [ ] 🔒 状態の階層的遷移ロジック
  - Blocked by: 複合状態サポート
- [ ] 🔒 階層的状態のテスト
  - Blocked by: 階層的遷移ロジック

**ヒストリーステート**
- [ ] 🎯 ShallowHistoryの実装
  - Blocked by: 階層的ステート実装
- [ ] 🔒 DeepHistoryの実装
  - Blocked by: ShallowHistory実装
- [ ] 🔒 ヒストリーステートのテスト
  - Blocked by: History実装

**パラレルステート**
- [ ] 🎯 並行状態モデルの設計
  - Blocked by: Phase 1 内部モデル実装
- [ ] 🔒 パラレル状態遷移の実装
  - Blocked by: 並行状態モデル
- [ ] 🔒 同期ポイントの実装
  - Blocked by: パラレル状態遷移
- [ ] 🔒 パラレルステートのテスト
  - Blocked by: パラレル状態実装

**内部遷移**
- [ ] ⚡ InternalTransitionの実装
  - Blocked by: Phase 1 内部モデル実装のみ
- [ ] 🔒 内部遷移のテスト
  - Blocked by: InternalTransition実装

## 🔄 In Progress

現在進行中のタスクはありません

## ✅ Done

### Phase 0: プロジェクト準備
- [x] プロジェクト概要設計書の作成 (docs/overview-design.md)
- [x] プロジェクト詳細設計書の作成 (docs/detailed-design.md)
- [x] CLAUDE.md の作成
- [x] TODO管理ボードの作成 (docs/TODO.md)
- [x] TDD methodology and test quality guidelines の追加

### Phase 1: 基本機能実装（進行中）

#### ブロッキングタスク - 完了 ✓
- [x] プロジェクト構造の初期化
- [x] go.modファイルの作成
- [x] CI/CDパイプラインの構築 (.github/workflows/ci.yml, .github/workflows/release.yml)
- [x] 開発環境のセットアップガイド作成 (SETUP.md)
- [x] FSMModel構造体の実装 (pkg/model/fsm.go + テスト)
- [x] State構造体の実装 (pkg/model/state.go + テスト)
- [x] Event構造体の実装 (pkg/model/event.go + テスト)
- [x] Transition構造体の実装 (pkg/model/transition.go + テスト)
- [x] StateGraph構造体の実装 (pkg/model/graph.go + テスト)

#### 並列実行可能タスク - 完了 ✓
- [x] READMEの作成
- [x] インストールガイドの作成 (docs/installation.md)
- [x] 基本的な使い方ガイド (docs/usage.md)
- [x] YAML定義リファレンス (docs/yaml-reference.md)
- [x] API仕様書 (docs/api.md)
- [x] exhaustive統合の調査 (docs/exhaustive-integration-investigation.md)
- [x] Makefileの作成（ビルド自動化）
- [x] golangci-lintの設定 (.golangci.yml)
- [x] MITライセンスの追加 (LICENSE)

## 📊 マイルストーン

### M1: Phase 1 完了 (基本機能) - 🚧 進行中 (約50%完了)
**完了済み**:
- ✅ プロジェクト基盤構築（構造、go.mod、CI/CD、Makefile、golangci-lint）
- ✅ 内部モデル実装（FSM, State, Event, Transition, Graph）- テストカバレッジ97.4%
- ✅ ドキュメント作成（README, installation.md, usage.md, yaml-reference.md, api.md, exhaustive調査）

**残作業** (クリティカルパス):
- ⏳ YAMLパーサーの実装 (pkg/parser/)
- ⏳ 基本コード生成器の実装 (pkg/generator/)
- ⏳ テンプレート作成 (templates/)
- ⏳ モデル検証器の実装 (pkg/analyzer/validator.go)
- ⏳ 静的解析基盤の実装 (pkg/analyzer/exhaustive.go)
- ⏳ CLI実装 (cmd/gofsm-gen/)
- ⏳ サンプルコード作成 (examples/)

**目標日**: TBD

**推奨される次のアクション**:
1. YAMLパーサーの実装（TDD方式で進める）
2. コード生成器インターフェースの定義
3. state_machine.tmplテンプレートの作成

### M2: Phase 2 完了 (高度な機能)
- ガード/アクション機能の完全実装
- Go DSLサポート
- テスト/モック自動生成
- パフォーマンス目標達成
- **目標日**: TBD

### M3: Phase 3 完了 (ツール統合)
- VSCode拡張機能リリース
- 既存ライブラリからの移行ツール
- **目標日**: TBD

### M4: Phase 4 完了 (高度な状態機械)
- 階層的ステートサポート
- ヒストリーステート実装
- パラレルステート実装
- **目標日**: TBD

## 🎯 優先度ラベル

タスクには以下の優先度を設定：

- **P0 (Critical)**: Phase 1の基本機能、プロジェクト成功に必須
- **P1 (High)**: Phase 2の主要機能
- **P2 (Medium)**: Phase 3のツール統合
- **P3 (Low)**: Phase 4の高度な機能、将来的な拡張

## 📝 メモ

### 技術的な決定事項
- **Go Version**: Go 1.18+ をサポート（ジェネリクスを活用）
- **コード生成**: `text/template` を使用
- **YAMLパース**: `gopkg.in/yaml.v3`
- **静的解析**: `golang.org/x/tools/go/analysis` + `exhaustive` ツール統合
- **テスト戦略**: TDD (Test-Driven Development) を厳格に適用
- **コードスタイル**: golangci-lint による品質管理

### 開発方針
1. **TDD優先**: 全ての新機能は「テストファースト」で実装
2. **高品質テスト**: 意味のあるテストのみ作成（マジックナンバー禁止）
3. **テーブル駆動テスト**: 複数のシナリオを網羅的にテスト
4. **カバレッジ目標**: > 90%（現在pkg/modelは97.4%達成）
5. **ドキュメント駆動**: 実装前にドキュメントで仕様を明確化

### リスクと課題
- Go言語の制約による機能制限の可能性
- パフォーマンス目標（<50ns/transition）の達成
- 既存ライブラリとの差別化の明確化
- コミュニティの採用促進
- exhaustiveツールの適切な統合とCI/CDでの実行

### 成功指標
- **GitHub Stars**: 1年で1000+
- **採用プロジェクト数**: 100+
- **テストカバレッジ**: > 90% (現在: 97.4% for pkg/model ✅)
- **パフォーマンス**: < 50ns/transition
- **ビルド時間**: < 1秒 for 1000 states
- **ドキュメント品質**: 全APIの完全なgodocカバレッジ

### 最近の成果 (2025-11-23)
- ✅ 内部モデル実装完了（97.4%カバレッジ）
- ✅ 包括的なドキュメント作成完了
- ✅ CI/CD パイプライン構築完了
- ✅ exhaustive統合調査完了
- ✅ プロジェクト基盤の確立

## 🚀 Phase 1 実装戦略

### 推奨される実装順序

**Week 1-2: パーサーとバリデーター基盤**
1. YAMLパーサーの実装（TDD）→ 4-6時間
2. Validator構造体の実装 → 3-4時間
3. 基本的な検証ロジック（状態、イベント、遷移）→ 6時間
4. 到達可能性・決定性検証 → 5-7時間

**Week 3-4: コード生成基盤**
1. state_machine.tmplテンプレート作成 → 4-6時間
2. コード生成器インターフェース定義 → 2-3時間
3. 基本コード生成器実装 → 6-8時間
4. Golden fileテスト作成 → 3-4時間

**Week 5: CLI実装とサンプル**
1. CLIエントリポイント実装 → 3-4時間
2. コマンドライン引数パース → 2-3時間
3. サンプルコード作成（ドアロック、注文管理）→ 4-5時間
4. End-to-endテスト → 3-4時間

**Week 6: 静的解析とポリッシュ**
1. exhaustive統合実装 → 4-5時間
2. アノテーション自動挿入 → 2-3時間
3. ドキュメント最終調整 → 2-3時間
4. パフォーマンステスト → 2-3時間

**合計推定工数**: 約60-80時間（1.5-2ヶ月、週20時間ペース）

### 並列作業の機会

以下のタスクは依存関係が少なく、並列実行可能：
- テンプレート作成（templates/*.tmpl）
- コントリビューションガイド作成
- ベンチマーク基盤の準備
- Phase 2の設計ドキュメント作成

### クリティカルパス

```
内部モデル ✅
    ↓
YAMLパーサー + Validator (並列可能)
    ↓
テンプレート + コード生成器インターフェース
    ↓
基本コード生成器
    ↓
CLI実装
    ↓
サンプルコード + E2Eテスト
    ↓
Phase 1 完了 🎉
```

### 品質チェックリスト

Phase 1完了前に確認すべき項目：
- [ ] 全パッケージのテストカバレッジ > 90%
- [ ] golangci-lint クリーン
- [ ] CI/CD グリーン（全テストパス）
- [ ] サンプルコードが動作
- [ ] ドキュメント完全性チェック
- [ ] パフォーマンステスト実施
- [ ] README.mdのクイックスタートが動作
- [ ] 生成コードがexhaustiveチェック通過

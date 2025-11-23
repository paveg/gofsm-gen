# gofsm-gen 開発タスク管理

vibe-kanbanを使用した開発タスク管理ボードです。

## 📋 Backlog

### Phase 1: 基本機能実装

#### プロジェクトセットアップ
- [ ] プロジェクト構造の初期化
- [ ] go.modファイルの作成
- [ ] CI/CDパイプラインの構築
- [ ] 開発環境のセットアップガイド作成

#### パーサー実装
- [ ] YAMLパーサーの実装 (pkg/parser/yaml.go)
- [ ] YAML定義構造体の作成
- [ ] YAMLパーサーのユニットテスト作成
- [ ] パーサーエラーハンドリングの実装

#### 内部モデル実装
- [ ] FSMModel構造体の実装 (pkg/model/fsm.go)
- [ ] State構造体の実装 (pkg/model/state.go)
- [ ] Event構造体の実装 (pkg/model/event.go)
- [ ] Transition構造体の実装 (pkg/model/transition.go)
- [ ] StateGraph構造体の実装 (pkg/model/graph.go)

#### 基本コード生成器
- [ ] コード生成器インターフェースの定義 (pkg/generator/generator.go)
- [ ] 基本コード生成器の実装 (pkg/generator/code_generator.go)
- [ ] state_machine.tmplテンプレートの作成
- [ ] テンプレートデータ準備ロジックの実装
- [ ] gofmtによるフォーマット処理の実装
- [ ] 生成コードのユニットテスト

#### 静的解析基盤
- [ ] exhaustive統合の調査
- [ ] 網羅性チェッカーの実装 (pkg/analyzer/exhaustive.go)
- [ ] exhaustiveアノテーション自動挿入機能
- [ ] 静的解析のテスト作成

#### モデル検証器
- [ ] Validator構造体の実装 (pkg/analyzer/validator.go)
- [ ] 状態の検証ロジック
- [ ] イベントの検証ロジック
- [ ] 遷移の検証ロジック
- [ ] 到達可能性の検証 (reachability analysis)
- [ ] 決定性の検証 (determinism check)
- [ ] 重複遷移チェック
- [ ] 検証エラーレポート機能

#### CLI実装
- [ ] CLIエントリポイントの実装 (cmd/gofsm-gen/main.go)
- [ ] コマンドライン引数のパース
- [ ] ヘルプメッセージの実装
- [ ] エラーメッセージの実装
- [ ] 詳細出力モード (-v flag)

#### ドキュメント作成
- [ ] READMEの作成
- [ ] インストールガイドの作成
- [ ] 基本的な使い方ガイド
- [ ] YAML定義リファレンス
- [ ] API仕様書
- [ ] コントリビューションガイド

#### サンプルコード
- [ ] シンプルなドアロック状態機械のサンプル
- [ ] Order管理状態機械のサンプル
- [ ] examples/ディレクトリの構築

### Phase 2: 高度な機能実装

#### ガード/アクション機能
- [ ] Guard構造体の実装
- [ ] Action構造体の実装
- [ ] ガード条件評価ロジック
- [ ] アクション実行ロジック
- [ ] エントリー/イグジットアクションの実装
- [ ] 非同期アクションのサポート
- [ ] ガード/アクションのテスト

#### Go DSLサポート
- [ ] Go ASTパーサーの実装 (pkg/parser/ast.go)
- [ ] DSLパーサーの実装 (pkg/parser/dsl.go)
- [ ] DSL用APIの設計
- [ ] Fluent API実装
- [ ] DSLのサンプルコード作成
- [ ] DSLパーサーのテスト

#### テストコード自動生成
- [ ] テストジェネレーターの実装 (pkg/generator/test_generator.go)
- [ ] test.tmplテンプレートの作成
- [ ] 状態遷移テストの自動生成
- [ ] ガード条件テストの自動生成
- [ ] エッジケーステストの自動生成
- [ ] 生成されたテストコードの検証

#### モック自動生成
- [ ] モックジェネレーターの実装 (pkg/generator/mock_generator.go)
- [ ] mock.tmplテンプレートの作成
- [ ] Guard/Actionモックの生成
- [ ] モックの使用例ドキュメント

#### 視覚化機能
- [ ] Mermaid生成器の実装 (pkg/visualizer/mermaid.go)
- [ ] Graphviz生成器の実装 (pkg/visualizer/graphviz.go)
- [ ] 状態遷移図の自動生成
- [ ] 視覚化オプションの実装

#### ランタイムサポート
- [ ] Loggerインターフェースの実装 (pkg/runtime/logger.go)
- [ ] TransitionLoggerの実装
- [ ] RuntimeValidatorの実装 (pkg/runtime/validator.go)
- [ ] 実行時統計収集機能
- [ ] ガード競合検出ロジック

#### パフォーマンス最適化
- [ ] ゼロアロケーションモードの実装
- [ ] パフォーマンスベンチマークの作成 (benchmarks/)
- [ ] 状態遷移のベンチマーク
- [ ] ガード付き遷移のベンチマーク
- [ ] メモリプロファイリング
- [ ] パフォーマンス目標の達成確認 (<50ns/transition)

### Phase 3: ツール統合

#### VSCode拡張機能
- [ ] VSCode拡張のプロジェクトセットアップ
- [ ] YAML定義のシンタックスハイライト
- [ ] 定義ファイルの自動補完
- [ ] リアルタイムバリデーション
- [ ] 状態遷移図のプレビュー機能
- [ ] コード生成のショートカット
- [ ] 拡張機能のテスト
- [ ] VSCode Marketplaceへの公開

#### gopls統合
- [ ] LSP統合の調査
- [ ] 型情報の提供
- [ ] 定義へのジャンプ機能
- [ ] リファクタリングサポート

#### 開発ツール
- [ ] ホットリロード機能
- [ ] ファイルウォッチャーの実装
- [ ] デバッグモードの実装
- [ ] 開発サーバーの実装

#### 移行ツール
- [ ] looplab/fsm → gofsm-gen 移行ツール
- [ ] qmuntal/stateless → gofsm-gen 移行ツール
- [ ] 移行ガイドの作成
- [ ] 移行スクリプトの作成

### Phase 4: 高度な状態機械機能

#### 階層的ステート
- [ ] 親子状態モデルの実装
- [ ] 複合状態のサポート
- [ ] 初期子状態の実装
- [ ] 状態の階層的遷移ロジック
- [ ] 階層的状態のテスト

#### ヒストリーステート
- [ ] ShallowHistoryの実装
- [ ] DeepHistoryの実装
- [ ] ヒストリーステートのテスト

#### パラレルステート
- [ ] 並行状態モデルの設計
- [ ] パラレル状態遷移の実装
- [ ] 同期ポイントの実装
- [ ] パラレルステートのテスト

#### 内部遷移
- [ ] InternalTransitionの実装
- [ ] 内部遷移のテスト

## 🔄 In Progress

現在進行中のタスクはありません

## ✅ Done

- [x] プロジェクト概要設計書の作成 (docs/overview-design.md)
- [x] プロジェクト詳細設計書の作成 (docs/detailed-design.md)
- [x] CLAUDE.md の作成
- [x] TODO管理ボードの作成 (docs/TODO.md)

## 📊 マイルストーン

### M1: Phase 1 完了 (基本機能)
- YAML定義からの基本的なコード生成が動作
- exhaustive統合による網羅性チェック
- 基本的なCLIツールの完成
- **目標日**: TBD

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
- Go 1.18+ をサポート
- `text/template` を使用したコード生成
- `gopkg.in/yaml.v3` をYAMLパーサーに使用
- `golang.org/x/tools/go/analysis` を静的解析に使用
- `exhaustive` ツールとの統合

### リスクと課題
- Go言語の制約による機能制限の可能性
- パフォーマンス目標（<50ns/transition）の達成
- 既存ライブラリとの差別化の明確化
- コミュニティの採用促進

### 成功指標
- GitHub Stars: 1年で1000+
- 採用プロジェクト数: 100+
- テストカバレッジ: > 90%
- パフォーマンス: 既存ライブラリと同等以上

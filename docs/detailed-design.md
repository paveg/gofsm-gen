# gofsm-gen 詳細設計書

## 1. モジュール構成

### 1.1 パッケージ構造

```
github.com/yourusername/gofsm-gen/
├── cmd/
│   └── gofsm-gen/
│       └── main.go                 # CLIエントリポイント
├── pkg/
│   ├── parser/                     # 定義パーサー
│   │   ├── yaml.go                 # YAML パーサー
│   │   ├── dsl.go                  # Go DSL パーサー
│   │   └── ast.go                  # Go AST パーサー
│   ├── model/                      # 内部データモデル
│   │   ├── fsm.go                  # FSMモデル定義
│   │   ├── state.go                # 状態モデル
│   │   ├── event.go                # イベントモデル
│   │   └── transition.go           # 遷移モデル
│   ├── generator/                  # コード生成器
│   │   ├── generator.go            # 生成器インターフェース
│   │   ├── code_generator.go       # メインコード生成
│   │   ├── test_generator.go       # テストコード生成
│   │   └── mock_generator.go       # モック生成
│   ├── analyzer/                   # 静的解析
│   │   ├── exhaustive.go           # 網羅性チェック
│   │   ├── validator.go            # モデル検証
│   │   └── graph.go                # グラフ解析
│   ├── visualizer/                 # 視覚化
│   │   ├── mermaid.go              # Mermaid生成
│   │   └── graphviz.go             # Graphviz生成
│   └── runtime/                    # ランタイムサポート
│       ├── logger.go               # ロギングインターフェース
│       ├── validator.go            # 実行時検証
│       └── context.go              # コンテキスト管理
├── templates/                       # コード生成テンプレート
│   ├── state_machine.tmpl
│   ├── test.tmpl
│   └── mock.tmpl
├── examples/                        # サンプルコード
└── tools/                          # 開発ツール
    └── vscode-extension/           # VSCode拡張
```

## 2. データモデル詳細

### 2.1 内部モデル定義

```go
// pkg/model/fsm.go

package model

import "time"

// FSMModel はステートマシンの完全な定義
type FSMModel struct {
    Name        string                 // マシン名
    Package     string                 // 生成先パッケージ名
    Initial     string                 // 初期状態
    States      map[string]*State      // 状態定義
    Events      map[string]*Event      // イベント定義
    Transitions []*Transition          // 遷移定義
    Metadata    *Metadata              // メタデータ
    
    // 内部フィールド（解析結果）
    stateGraph  *StateGraph           // 状態グラフ
    validated   bool                   // 検証済みフラグ
    errors      []ValidationError      // 検証エラー
}

// State は状態の定義
type State struct {
    Name        string                 // 状態名
    Type        StateType              // 状態タイプ
    Entry       *Action                // エントリーアクション
    Exit        *Action                // イグジットアクション
    Internal    []*InternalTransition  // 内部遷移
    Properties  map[string]interface{} // カスタムプロパティ
    
    // 階層的ステート用（Phase 4）
    Parent      *State                 // 親状態
    Children    []*State               // 子状態
    Initial     *State                 // 初期子状態
    History     HistoryType            // 履歴タイプ
}

type StateType int
const (
    StateTypeSimple StateType = iota  // 単純状態
    StateTypeComposite                 // 複合状態
    StateTypeFinal                     // 最終状態
    StateTypePseudo                    // 擬似状態
)

// Event はイベントの定義
type Event struct {
    Name       string                  // イベント名
    Parameters []*Parameter            // パラメータ定義
    Metadata   map[string]interface{}  // メタデータ
}

// Parameter はイベントパラメータ
type Parameter struct {
    Name     string                    // パラメータ名
    Type     string                    // 型名（Go型）
    Required bool                      // 必須フラグ
    Default  interface{}               // デフォルト値
}

// Transition は状態遷移の定義
type Transition struct {
    From       string                  // 遷移元状態
    To         string                  // 遷移先状態
    Event      string                  // トリガーイベント
    Guard      *Guard                  // ガード条件
    Action     *Action                 // 遷移アクション
    Priority   int                     // 優先度（複数遷移時）
    Internal   bool                    // 内部遷移フラグ
}

// Guard はガード条件
type Guard struct {
    Name       string                  // ガード関数名
    Expression string                  // ガード式（オプション）
    Parameters []*Parameter            // パラメータ
}

// Action はアクション定義
type Action struct {
    Name       string                  // アクション関数名
    Async      bool                    // 非同期実行フラグ
    Timeout    time.Duration           // タイムアウト
    Parameters []*Parameter            // パラメータ
}

// Metadata はFSMのメタデータ
type Metadata struct {
    Version     string                 // バージョン
    Author      string                 // 作者
    Description string                 // 説明
    Tags        []string               // タグ
    Generated   time.Time              // 生成日時
}
```

### 2.2 グラフモデル

```go
// pkg/model/graph.go

package model

// StateGraph は状態遷移グラフ
type StateGraph struct {
    Nodes map[string]*StateNode        // 状態ノード
    Edges []*TransitionEdge            // 遷移エッジ
}

// StateNode はグラフのノード
type StateNode struct {
    State      *State                  // 状態
    InEdges    []*TransitionEdge       // 入力エッジ
    OutEdges   []*TransitionEdge       // 出力エッジ
    Reachable  bool                    // 到達可能フラグ
    Depth      int                     // 初期状態からの深さ
}

// TransitionEdge はグラフのエッジ
type TransitionEdge struct {
    Transition *Transition             // 遷移
    From       *StateNode              // 開始ノード
    To         *StateNode              // 終了ノード
    Weight     int                     // エッジの重み
}
```

## 3. パーサー実装

### 3.1 YAMLパーサー

```go
// pkg/parser/yaml.go

package parser

import (
    "fmt"
    "io"
    "gopkg.in/yaml.v3"
    "github.com/yourusername/gofsm-gen/pkg/model"
)

// YAMLParser はYAML形式の定義をパースする
type YAMLParser struct {
    strict bool  // 厳密モード
}

// YAMLDefinition はYAMLの構造定義
type YAMLDefinition struct {
    Machine struct {
        Name    string `yaml:"name"`
        Initial string `yaml:"initial"`
    } `yaml:"machine"`
    
    States []struct {
        Name  string `yaml:"name"`
        Entry string `yaml:"entry,omitempty"`
        Exit  string `yaml:"exit,omitempty"`
    } `yaml:"states"`
    
    Events []string `yaml:"events"`
    
    Transitions []struct {
        From   string `yaml:"from"`
        To     string `yaml:"to"`
        On     string `yaml:"on"`
        Guard  string `yaml:"guard,omitempty"`
        Action string `yaml:"action,omitempty"`
    } `yaml:"transitions"`
}

func (p *YAMLParser) Parse(r io.Reader) (*model.FSMModel, error) {
    var def YAMLDefinition
    decoder := yaml.NewDecoder(r)
    if p.strict {
        decoder.KnownFields(true)
    }
    
    if err := decoder.Decode(&def); err != nil {
        return nil, fmt.Errorf("failed to decode YAML: %w", err)
    }
    
    return p.buildModel(&def)
}

func (p *YAMLParser) buildModel(def *YAMLDefinition) (*model.FSMModel, error) {
    fsm := &model.FSMModel{
        Name:    def.Machine.Name,
        Initial: def.Machine.Initial,
        States:  make(map[string]*model.State),
        Events:  make(map[string]*model.Event),
    }
    
    // 状態を構築
    for _, s := range def.States {
        state := &model.State{
            Name: s.Name,
            Type: model.StateTypeSimple,
        }
        
        if s.Entry != "" {
            state.Entry = &model.Action{Name: s.Entry}
        }
        if s.Exit != "" {
            state.Exit = &model.Action{Name: s.Exit}
        }
        
        fsm.States[s.Name] = state
    }
    
    // イベントを構築
    for _, e := range def.Events {
        fsm.Events[e] = &model.Event{Name: e}
    }
    
    // 遷移を構築
    for _, t := range def.Transitions {
        transition := &model.Transition{
            From:  t.From,
            To:    t.To,
            Event: t.On,
        }
        
        if t.Guard != "" {
            transition.Guard = &model.Guard{Name: t.Guard}
        }
        if t.Action != "" {
            transition.Action = &model.Action{Name: t.Action}
        }
        
        fsm.Transitions = append(fsm.Transitions, transition)
    }
    
    return fsm, nil
}
```

### 3.2 Go DSLパーサー

```go
// pkg/parser/dsl.go

package parser

import (
    "go/ast"
    "go/parser"
    "go/token"
    "github.com/yourusername/gofsm-gen/pkg/model"
)

// DSLParser はGo DSL形式の定義をパースする
type DSLParser struct {
    packageName string
}

func (p *DSLParser) Parse(filename string) (*model.FSMModel, error) {
    fset := token.NewFileSet()
    file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
    if err != nil {
        return nil, err
    }
    
    // ASTを走査してDSL呼び出しを検出
    var dslCalls []*ast.CallExpr
    ast.Inspect(file, func(n ast.Node) bool {
        if call, ok := n.(*ast.CallExpr); ok {
            if p.isDSLCall(call) {
                dslCalls = append(dslCalls, call)
            }
        }
        return true
    })
    
    return p.buildModelFromDSL(dslCalls)
}

func (p *DSLParser) isDSLCall(call *ast.CallExpr) bool {
    // dsl.Machine() 呼び出しを検出
    if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
        if ident, ok := sel.X.(*ast.Ident); ok {
            return ident.Name == "dsl" && sel.Sel.Name == "Machine"
        }
    }
    return false
}

func (p *DSLParser) buildModelFromDSL(calls []*ast.CallExpr) (*model.FSMModel, error) {
    // DSL呼び出しチェーンを解析してモデルを構築
    // 実装省略（流れるAPIの解析）
    return nil, nil
}
```

## 4. コード生成器実装

### 4.1 メインコード生成器

```go
// pkg/generator/code_generator.go

package generator

import (
    "bytes"
    "fmt"
    "go/format"
    "text/template"
    "github.com/yourusername/gofsm-gen/pkg/model"
)

// CodeGenerator はメインのコード生成器
type CodeGenerator struct {
    model     *model.FSMModel
    template  *template.Template
    options   *Options
}

type Options struct {
    Package         string
    GenerateTests   bool
    GenerateMocks   bool
    Exhaustive      bool
    ZeroAllocation  bool
}

func (g *CodeGenerator) Generate() ([]byte, error) {
    // テンプレート用データを準備
    data := g.prepareTemplateData()
    
    // テンプレート実行
    var buf bytes.Buffer
    if err := g.template.Execute(&buf, data); err != nil {
        return nil, fmt.Errorf("template execution failed: %w", err)
    }
    
    // gofmtでフォーマット
    formatted, err := format.Source(buf.Bytes())
    if err != nil {
        return nil, fmt.Errorf("gofmt failed: %w", err)
    }
    
    return formatted, nil
}

func (g *CodeGenerator) prepareTemplateData() map[string]interface{} {
    return map[string]interface{}{
        "Package":      g.options.Package,
        "MachineName":  g.model.Name,
        "States":       g.generateStatesData(),
        "Events":       g.generateEventsData(),
        "Transitions":  g.generateTransitionsData(),
        "Exhaustive":   g.options.Exhaustive,
        "Imports":      g.generateImports(),
    }
}

func (g *CodeGenerator) generateStatesData() []map[string]interface{} {
    var states []map[string]interface{}
    
    for _, state := range g.model.States {
        states = append(states, map[string]interface{}{
            "Name":       state.Name,
            "ConstName":  g.toConstName(state.Name),
            "HasEntry":   state.Entry != nil,
            "HasExit":    state.Exit != nil,
            "EntryFunc":  state.Entry,
            "ExitFunc":   state.Exit,
        })
    }
    
    return states
}

func (g *CodeGenerator) generateTransitionsTable() string {
    // 網羅的なswitch文を生成
    var buf bytes.Buffer
    
    // exhaustiveアノテーション
    if g.options.Exhaustive {
        buf.WriteString("//exhaustive:enforce\n")
    }
    
    buf.WriteString("switch m.state {\n")
    
    // 各状態ごとに処理
    for stateName, state := range g.model.States {
        buf.WriteString(fmt.Sprintf("case %s:\n", g.toConstName(stateName)))
        buf.WriteString("    switch event {\n")
        
        // この状態からの遷移を収集
        transitions := g.collectTransitionsFrom(stateName)
        
        // 各イベントごとに処理
        for eventName := range g.model.Events {
            if trans := g.findTransition(transitions, eventName); trans != nil {
                g.generateTransitionCase(&buf, trans)
            } else {
                g.generateInvalidTransitionCase(&buf, eventName)
            }
        }
        
        buf.WriteString("    }\n")
    }
    
    buf.WriteString("}\n")
    
    return buf.String()
}

func (g *CodeGenerator) generateTransitionCase(buf *bytes.Buffer, trans *model.Transition) {
    buf.WriteString(fmt.Sprintf("    case %s:\n", g.toConstName(trans.Event)))
    
    // ガードチェック
    if trans.Guard != nil {
        buf.WriteString(fmt.Sprintf(`
        if m.guards.%s != nil {
            if !m.guards.%s(ctx, m.context) {
                return ErrGuardFailed
            }
        }
`, trans.Guard.Name, trans.Guard.Name))
    }
    
    // 状態遷移
    buf.WriteString(fmt.Sprintf("        m.state = %s\n", g.toConstName(trans.To)))
    
    // アクション実行
    if trans.Action != nil {
        buf.WriteString(fmt.Sprintf(`
        if m.actions.%s != nil {
            if err := m.actions.%s(ctx, oldState, m.state, m.context); err != nil {
                return err
            }
        }
`, trans.Action.Name, trans.Action.Name))
    }
    
    buf.WriteString("        return nil\n")
}
```

### 4.2 テンプレート定義

```go
// templates/state_machine.tmpl

const stateMachineTemplate = `// Code generated by gofsm-gen. DO NOT EDIT.
package {{.Package}}

import (
    "context"
    "errors"
    "fmt"
    {{range .Imports}}
    "{{.}}"
    {{end}}
)

// エラー定義
var (
    ErrInvalidTransition = errors.New("invalid transition")
    ErrGuardFailed      = errors.New("guard condition failed")
    ErrInvalidState     = errors.New("invalid state")
)

// 状態の列挙型
type {{.MachineName}}State int

const (
    {{range $i, $state := .States}}
    {{if eq $i 0}}{{$state.ConstName}} {{$.MachineName}}State = iota{{else}}{{$state.ConstName}}{{end}}
    {{end}}
)

// String は状態を文字列化する
func (s {{.MachineName}}State) String() string {
    switch s {
    {{range .States}}
    case {{.ConstName}}: return "{{.Name}}"
    {{end}}
    default: return fmt.Sprintf("{{.MachineName}}State(%d)", s)
    }
}

// イベントの列挙型
type {{.MachineName}}Event int

const (
    {{range $i, $event := .Events}}
    {{if eq $i 0}}{{$event.ConstName}} {{$.MachineName}}Event = iota{{else}}{{$event.ConstName}}{{end}}
    {{end}}
)

// コンテキスト型
type {{.MachineName}}Context struct {
    // ユーザー定義フィールド
}

// ガード関数型
type {{.MachineName}}Guards struct {
    {{range .Guards}}
    {{.Name}} func(ctx context.Context, c *{{$.MachineName}}Context) bool
    {{end}}
}

// アクション関数型
type {{.MachineName}}Actions struct {
    {{range .Actions}}
    {{.Name}} func(ctx context.Context, from, to {{$.MachineName}}State, c *{{$.MachineName}}Context) error
    {{end}}
}

// ステートマシン本体
type {{.MachineName}}Machine struct {
    state   {{.MachineName}}State
    context *{{.MachineName}}Context
    guards  {{.MachineName}}Guards
    actions {{.MachineName}}Actions
}

// New{{.MachineName}}Machine は新しいステートマシンを作成する
func New{{.MachineName}}Machine(
    guards {{.MachineName}}Guards,
    actions {{.MachineName}}Actions,
) *{{.MachineName}}Machine {
    return &{{.MachineName}}Machine{
        state:   {{.InitialState}},
        context: &{{.MachineName}}Context{},
        guards:  guards,
        actions: actions,
    }
}

// State は現在の状態を返す
func (m *{{.MachineName}}Machine) State() {{.MachineName}}State {
    return m.state
}

// Transition は状態遷移を実行する
func (m *{{.MachineName}}Machine) Transition(ctx context.Context, event {{.MachineName}}Event) error {
    oldState := m.state
    
    {{.TransitionsTable}}
    
    return nil
}

// PermittedEvents は現在の状態で許可されたイベントを返す
func (m *{{.MachineName}}Machine) PermittedEvents() []{{.MachineName}}Event {
    {{if .Exhaustive}}//exhaustive:enforce{{end}}
    switch m.state {
    {{range .States}}
    case {{.ConstName}}:
        return []{{$.MachineName}}Event{ {{range .PermittedEvents}}{{.}}, {{end}} }
    {{end}}
    }
    return nil
}
`
```

## 5. 静的解析実装

### 5.1 網羅性チェッカー

```go
// pkg/analyzer/exhaustive.go

package analyzer

import (
    "fmt"
    "go/ast"
    "go/token"
    "go/types"
    "golang.org/x/tools/go/analysis"
    "github.com/yourusername/gofsm-gen/pkg/model"
)

// ExhaustiveAnalyzer は網羅性をチェックする
type ExhaustiveAnalyzer struct {
    model *model.FSMModel
    pass  *analysis.Pass
}

// Analyzer はanalysisパッケージ用のAnalyzer定義
var Analyzer = &analysis.Analyzer{
    Name: "gofsm-exhaustive",
    Doc:  "check exhaustiveness of state machine switch statements",
    Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
    for _, file := range pass.Files {
        ast.Inspect(file, func(n ast.Node) bool {
            switch stmt := n.(type) {
            case *ast.SwitchStmt:
                checkSwitch(pass, stmt)
            }
            return true
        })
    }
    return nil, nil
}

func checkSwitch(pass *analysis.Pass, stmt *ast.SwitchStmt) {
    // exhaustive:enforceコメントを確認
    if !hasExhaustiveComment(stmt) {
        return
    }
    
    // switch対象の型を取得
    tagType := pass.TypesInfo.TypeOf(stmt.Tag)
    if tagType == nil {
        return
    }
    
    // 列挙型の全値を取得
    allValues := collectEnumValues(pass, tagType)
    if len(allValues) == 0 {
        return
    }
    
    // カバーされている値を収集
    covered := make(map[string]bool)
    hasDefault := false
    
    for _, clause := range stmt.Body.List {
        cc, ok := clause.(*ast.CaseClause)
        if !ok {
            continue
        }
        
        if cc.List == nil {
            hasDefault = true
            continue
        }
        
        for _, expr := range cc.List {
            if val := getConstValue(pass, expr); val != "" {
                covered[val] = true
            }
        }
    }
    
    // defaultがある場合は網羅的
    if hasDefault {
        return
    }
    
    // 欠けている値を報告
    var missing []string
    for _, val := range allValues {
        if !covered[val] {
            missing = append(missing, val)
        }
    }
    
    if len(missing) > 0 {
        pass.Reportf(stmt.Pos(), "missing cases in switch: %v", missing)
    }
}

func hasExhaustiveComment(stmt *ast.SwitchStmt) bool {
    // コメントをチェック
    // 実装省略
    return false
}

func collectEnumValues(pass *analysis.Pass, t types.Type) []string {
    // 型の全定数値を収集
    // 実装省略
    return nil
}

func getConstValue(pass *analysis.Pass, expr ast.Expr) string {
    // 式から定数値を取得
    // 実装省略
    return ""
}
```

### 5.2 モデル検証器

```go
// pkg/analyzer/validator.go

package analyzer

import (
    "fmt"
    "github.com/yourusername/gofsm-gen/pkg/model"
)

// Validator はFSMモデルを検証する
type Validator struct {
    model  *model.FSMModel
    errors []ValidationError
}

type ValidationError struct {
    Type     ErrorType
    Message  string
    Location string
}

type ErrorType string
const (
    ErrorTypeUnreachableState   ErrorType = "unreachable_state"
    ErrorTypeDuplicateTransition ErrorType = "duplicate_transition"
    ErrorTypeInvalidInitial      ErrorType = "invalid_initial"
    ErrorTypeMissingState        ErrorType = "missing_state"
    ErrorTypeMissingEvent        ErrorType = "missing_event"
    ErrorTypeConflictingGuards   ErrorType = "conflicting_guards"
)

func (v *Validator) Validate() []ValidationError {
    v.errors = nil
    
    v.validateStates()
    v.validateEvents()
    v.validateTransitions()
    v.validateReachability()
    v.validateDeterminism()
    
    return v.errors
}

func (v *Validator) validateStates() {
    // 初期状態の存在確認
    if v.model.Initial != "" {
        if _, exists := v.model.States[v.model.Initial]; !exists {
            v.addError(ErrorTypeInvalidInitial, 
                fmt.Sprintf("initial state '%s' not found", v.model.Initial))
        }
    }
    
    // 状態名の重複チェック
    seen := make(map[string]bool)
    for name := range v.model.States {
        if seen[name] {
            v.addError(ErrorTypeDuplicateTransition,
                fmt.Sprintf("duplicate state name: %s", name))
        }
        seen[name] = true
    }
}

func (v *Validator) validateTransitions() {
    // 遷移の検証
    for _, trans := range v.model.Transitions {
        // From状態の存在確認
        if _, exists := v.model.States[trans.From]; !exists {
            v.addError(ErrorTypeMissingState,
                fmt.Sprintf("transition from unknown state: %s", trans.From))
        }
        
        // To状態の存在確認
        if _, exists := v.model.States[trans.To]; !exists {
            v.addError(ErrorTypeMissingState,
                fmt.Sprintf("transition to unknown state: %s", trans.To))
        }
        
        // イベントの存在確認
        if _, exists := v.model.Events[trans.Event]; !exists {
            v.addError(ErrorTypeMissingEvent,
                fmt.Sprintf("transition with unknown event: %s", trans.Event))
        }
    }
    
    // 重複遷移のチェック
    v.checkDuplicateTransitions()
}

func (v *Validator) checkDuplicateTransitions() {
    type key struct {
        from  string
        event string
    }
    
    transitions := make(map[key][]*model.Transition)
    
    for _, trans := range v.model.Transitions {
        k := key{from: trans.From, event: trans.Event}
        transitions[k] = append(transitions[k], trans)
    }
    
    // 同じ(from, event)で複数の遷移がある場合
    for k, transList := range transitions {
        if len(transList) > 1 {
            // ガードがない遷移が複数ある場合はエラー
            noGuardCount := 0
            for _, trans := range transList {
                if trans.Guard == nil {
                    noGuardCount++
                }
            }
            
            if noGuardCount > 1 {
                v.addError(ErrorTypeConflictingGuards,
                    fmt.Sprintf("multiple unguarded transitions from %s on %s", 
                        k.from, k.event))
            }
        }
    }
}

func (v *Validator) validateReachability() {
    // グラフ構築
    graph := v.buildStateGraph()
    
    // 初期状態から到達可能性を検証
    visited := make(map[string]bool)
    v.dfs(v.model.Initial, visited, graph)
    
    // 到達不可能な状態を報告
    for stateName := range v.model.States {
        if !visited[stateName] {
            v.addError(ErrorTypeUnreachableState,
                fmt.Sprintf("state '%s' is unreachable from initial state", stateName))
        }
    }
}

func (v *Validator) buildStateGraph() map[string][]string {
    graph := make(map[string][]string)
    
    for _, trans := range v.model.Transitions {
        graph[trans.From] = append(graph[trans.From], trans.To)
    }
    
    return graph
}

func (v *Validator) dfs(state string, visited map[string]bool, graph map[string][]string) {
    if visited[state] {
        return
    }
    
    visited[state] = true
    
    for _, next := range graph[state] {
        v.dfs(next, visited, graph)
    }
}

func (v *Validator) addError(errType ErrorType, message string) {
    v.errors = append(v.errors, ValidationError{
        Type:    errType,
        Message: message,
    })
}
```

## 6. ランタイムサポート実装

### 6.1 ロギングインターフェース

```go
// pkg/runtime/logger.go

package runtime

import (
    "context"
    "fmt"
)

// Logger はロギングインターフェース
type Logger interface {
    Debug(format string, args ...interface{})
    Info(format string, args ...interface{})
    Warn(format string, args ...interface{})
    Error(format string, args ...interface{})
}

// StructuredLogger は構造化ログ用インターフェース
type StructuredLogger interface {
    Logger
    WithFields(fields map[string]interface{}) Logger
    WithContext(ctx context.Context) Logger
}

// NoopLogger は何もしないロガー
type NoopLogger struct{}

func (l *NoopLogger) Debug(format string, args ...interface{}) {}
func (l *NoopLogger) Info(format string, args ...interface{})  {}
func (l *NoopLogger) Warn(format string, args ...interface{})  {}
func (l *NoopLogger) Error(format string, args ...interface{}) {}

// TransitionLogger は遷移専用ロガー
type TransitionLogger struct {
    logger Logger
}

func (l *TransitionLogger) LogTransition(from, to string, event string, success bool) {
    if success {
        l.logger.Info("transition succeeded: %s -> %s on %s", from, to, event)
    } else {
        l.logger.Warn("transition failed: %s on %s", from, event)
    }
}

func (l *TransitionLogger) LogGuardEvaluation(guard string, result bool) {
    l.logger.Debug("guard '%s' evaluated to %v", guard, result)
}

func (l *TransitionLogger) LogActionExecution(action string, err error) {
    if err != nil {
        l.logger.Error("action '%s' failed: %v", action, err)
    } else {
        l.logger.Debug("action '%s' executed successfully", action)
    }
}
```

### 6.2 実行時検証器

```go
// pkg/runtime/validator.go

package runtime

import (
    "context"
    "fmt"
    "sync"
)

// RuntimeValidator は実行時検証を行う
type RuntimeValidator struct {
    enabled bool
    logger  Logger
    mu      sync.Mutex
    
    // 統計情報
    transitionCount map[string]int
    guardEvalCount  map[string]int
    conflicts       []ConflictInfo
}

type ConflictInfo struct {
    State   string
    Event   string
    Guards  []string
    Results map[string]bool
}

func NewRuntimeValidator(logger Logger) *RuntimeValidator {
    return &RuntimeValidator{
        enabled:         true,
        logger:          logger,
        transitionCount: make(map[string]int),
        guardEvalCount:  make(map[string]int),
    }
}

func (v *RuntimeValidator) ValidateTransition(
    ctx context.Context,
    from, to string,
    event string,
    guards map[string]func() bool,
) error {
    if !v.enabled {
        return nil
    }
    
    v.mu.Lock()
    defer v.mu.Unlock()
    
    // 遷移カウント
    key := fmt.Sprintf("%s->%s:%s", from, to, event)
    v.transitionCount[key]++
    
    // 複数ガードの相互排他性チェック
    if len(guards) > 1 {
        results := make(map[string]bool)
        trueCount := 0
        
        for name, fn := range guards {
            result := fn()
            results[name] = result
            if result {
                trueCount++
            }
            v.guardEvalCount[name]++
        }
        
        // 複数のガードがtrueの場合は警告
        if trueCount > 1 {
            v.conflicts = append(v.conflicts, ConflictInfo{
                State:   from,
                Event:   event,
                Guards:  getKeys(guards),
                Results: results,
            })
            
            v.logger.Warn("multiple guards returned true for %s->%s on %s: %v",
                from, to, event, results)
        }
    }
    
    return nil
}

func (v *RuntimeValidator) GetStatistics() map[string]interface{} {
    v.mu.Lock()
    defer v.mu.Unlock()
    
    return map[string]interface{}{
        "transitionCount": v.transitionCount,
        "guardEvalCount":  v.guardEvalCount,
        "conflicts":       v.conflicts,
    }
}

func getKeys(m map[string]func() bool) []string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

## 7. CLI実装

### 7.1 メインエントリポイント

```go
// cmd/gofsm-gen/main.go

package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "path/filepath"
    
    "github.com/yourusername/gofsm-gen/pkg/generator"
    "github.com/yourusername/gofsm-gen/pkg/parser"
)

func main() {
    var (
        spec          = flag.String("spec", "", "FSM specification file (YAML/HCL)")
        typeName      = flag.String("type", "", "Type name for Go DSL mode")
        infer         = flag.Bool("infer", false, "Infer FSM from Go type")
        out           = flag.String("out", "", "Output file path")
        pkg           = flag.String("package", "", "Package name")
        generateTests = flag.Bool("generate-tests", false, "Generate test code")
        generateMocks = flag.Bool("generate-mocks", false, "Generate mock code")
        visualize     = flag.String("visualize", "", "Generate visualization (mermaid/graphviz)")
        exhaustive    = flag.Bool("exhaustive", true, "Add exhaustive annotations")
        verbose       = flag.Bool("v", false, "Verbose output")
    )
    
    flag.Parse()
    
    // 入力検証
    if *spec == "" && *typeName == "" && !*infer {
        log.Fatal("must specify -spec, -type, or -infer")
    }
    
    // FSMモデルを構築
    model, err := buildModel(*spec, *typeName, *infer)
    if err != nil {
        log.Fatalf("failed to build model: %v", err)
    }
    
    // 出力ファイル名を決定
    if *out == "" {
        *out = fmt.Sprintf("%s_fsm.gen.go", toSnakeCase(model.Name))
    }
    
    // パッケージ名を決定
    if *pkg == "" {
        *pkg = inferPackageName(*out)
    }
    
    // コード生成オプション
    opts := &generator.Options{
        Package:        *pkg,
        GenerateTests:  *generateTests,
        GenerateMocks:  *generateMocks,
        Exhaustive:     *exhaustive,
    }
    
    // メインコード生成
    if err := generateMainCode(model, *out, opts); err != nil {
        log.Fatalf("failed to generate main code: %v", err)
    }
    
    if *verbose {
        log.Printf("Generated %s", *out)
    }
    
    // テストコード生成
    if *generateTests {
        testFile := strings.TrimSuffix(*out, ".go") + "_test.go"
        if err := generateTestCode(model, testFile, opts); err != nil {
            log.Fatalf("failed to generate test code: %v", err)
        }
        if *verbose {
            log.Printf("Generated %s", testFile)
        }
    }
    
    // モック生成
    if *generateMocks {
        mockFile := strings.TrimSuffix(*out, ".go") + "_mock.go"
        if err := generateMockCode(model, mockFile, opts); err != nil {
            log.Fatalf("failed to generate mock code: %v", err)
        }
        if *verbose {
            log.Printf("Generated %s", mockFile)
        }
    }
    
    // 視覚化
    if *visualize != "" {
        if err := generateVisualization(model, *visualize); err != nil {
            log.Fatalf("failed to generate visualization: %v", err)
        }
        if *verbose {
            log.Printf("Generated visualization")
        }
    }
}

func buildModel(spec, typeName string, infer bool) (*model.FSMModel, error) {
    if spec != "" {
        // YAMLまたはHCLから構築
        return parseSpecFile(spec)
    }
    
    if typeName != "" {
        // Go DSLから構築
        return parseDSL(typeName)
    }
    
    if infer {
        // Go型から推論
        return inferFromType()
    }
    
    return nil, fmt.Errorf("no input source specified")
}
```

## 8. テスト実装

### 8.1 単体テスト例

```go
// pkg/generator/code_generator_test.go

package generator

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/yourusername/gofsm-gen/pkg/model"
)

func TestCodeGenerator_Generate(t *testing.T) {
    tests := []struct {
        name    string
        model   *model.FSMModel
        options *Options
        wantErr bool
    }{
        {
            name: "simple state machine",
            model: &model.FSMModel{
                Name:    "Door",
                Initial: "locked",
                States: map[string]*model.State{
                    "locked":   {Name: "locked"},
                    "unlocked": {Name: "unlocked"},
                },
                Events: map[string]*model.Event{
                    "unlock": {Name: "unlock"},
                    "lock":   {Name: "lock"},
                },
                Transitions: []*model.Transition{
                    {From: "locked", To: "unlocked", Event: "unlock"},
                    {From: "unlocked", To: "locked", Event: "lock"},
                },
            },
            options: &Options{
                Package:    "door",
                Exhaustive: true,
            },
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gen := NewCodeGenerator(tt.model, tt.options)
            
            code, err := gen.Generate()
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, code)
                
                // 生成されたコードの検証
                codeStr := string(code)
                assert.Contains(t, codeStr, "type DoorState int")
                assert.Contains(t, codeStr, "type DoorEvent int")
                assert.Contains(t, codeStr, "//exhaustive:enforce")
            }
        })
    }
}
```

## 9. パフォーマンスベンチマーク

```go
// benchmarks/performance_test.go

package benchmarks

import (
    "context"
    "testing"
)

func BenchmarkStateMachine_Transition(b *testing.B) {
    // 生成されたコードのベンチマーク
    sm := NewDoorLockMachine(
        DoorLockGuards{},
        DoorLockActions{},
    )
    
    ctx := context.Background()
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        sm.state = DoorLockStateLocked
        sm.Transition(ctx, DoorLockEventUnlock)
    }
}

func BenchmarkStateMachine_GuardedTransition(b *testing.B) {
    guards := DoorLockGuards{
        ValidKey: func(ctx context.Context, c *DoorLockContext) bool {
            return true
        },
    }
    
    sm := NewDoorLockMachine(guards, DoorLockActions{})
    ctx := context.Background()
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        sm.state = DoorLockStateLocked
        sm.Transition(ctx, DoorLockEventUnlock)
    }
}
```

## 10. 設定ファイル

### 10.1 プロジェクト設定

```yaml
# .gofsm-gen.yml
version: 1.0
defaults:
  package: fsm
  exhaustive: true
  generate_tests: true
  
generators:
  - name: main
    template: templates/state_machine.tmpl
  - name: test
    template: templates/test.tmpl
    
analyzers:
  - exhaustive
  - reachability
  - determinism
  
output:
  format: gofmt
  header: |
    // Code generated by gofsm-gen. DO NOT EDIT.
```

このように、概要設計では全体的なアーキテクチャと方針を定義し、詳細設計では具体的な実装レベルのコード構造とアルゴリズムを定義しています。この設計に基づいて実装を進めることで、Rustレベルの網羅性チェックに近い安全性を持つGoのステートマシンライブラリを実現できます。

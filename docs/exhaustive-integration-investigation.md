# exhaustive Tool Integration Investigation

**Date**: 2025-11-23
**Status**: Complete
**Phase**: Phase 1 - 静的解析基盤

## Executive Summary

This document presents findings from investigating the `exhaustive` static analysis tool for integration into gofsm-gen. The exhaustive tool provides Rust-like exhaustiveness checking for switch statements on enum-like constants in Go, which aligns perfectly with our goal of compile-time safety for state machine code generation.

**Key Recommendation**: Integrate exhaustive tool through both code generation annotations and optional static analysis integration.

## 1. Overview of exhaustive Tool

### 1.1 Tool Information

- **Repository**: https://github.com/nishanths/exhaustive
- **Package**: `github.com/nishanths/exhaustive`
- **Purpose**: Check exhaustiveness of switch statements of enum-like constants in Go source code
- **Analysis Framework**: `golang.org/x/tools/go/analysis` compliant

### 1.2 Installation

```bash
# Command line tool
go install github.com/nishanths/exhaustive/cmd/exhaustive@latest

# As a library/analyzer
go get github.com/nishanths/exhaustive
```

### 1.3 Integration Status

The tool is widely adopted and integrated into:
- golangci-lint (popular Go linting framework)
- Custom analysis pipelines via `analysis.Analyzer` interface

## 2. How exhaustive Identifies Enum Types

### 2.1 Enum Type Definition

The exhaustive tool identifies a type as an enum when it satisfies:

1. **Named type** with underlying type of:
   - Integer types (int, int8, int16, int32, int64, uint, uint8/byte, uint16, uint32, uint64)
   - Float types (float32, float64)
   - String type

2. **At least one constant** of that type defined in the same block

3. **Constants declared in the same block** as the type definition

### 2.2 Example Valid Enum

```go
// Valid enum - type and constants in same file-level block
type OrderState int

const (
    OrderStatePending  OrderState = iota
    OrderStateApproved
    OrderStateRejected
    OrderStateShipped
)
```

### 2.3 Scope Considerations

**Default behavior**: Discovers enums in all scopes (file-level and function-level blocks)

**With `-package-scope-only` flag**: Only discovers enums at file-level blocks (recommended for generated code)

## 3. Annotation Syntax and Usage

### 3.1 Comment Directives

#### `//exhaustive:ignore`

Instructs exhaustive to skip checking a specific switch statement or map literal.

**Syntax**:
```go
//exhaustive:ignore [optional explanation]
switch state {
case StatePending:
    // handle pending
// Other cases not required
}
```

**Use cases**:
- Partial switch statements where default case handles remaining cases intentionally
- Switch statements that don't require full coverage by design

#### `//exhaustive:enforce`

Explicitly marks a switch statement or map literal for exhaustiveness checking.

**Syntax**:
```go
//exhaustive:enforce
switch state {
case StatePending:
    // handle pending
case StateApproved:
    // handle approved
case StateRejected:
    // handle rejected
case StateShipped:
    // handle shipped
}
```

**Use cases**:
- When using `-explicit-exhaustive-switch` or `-explicit-exhaustive-map` flags
- Explicitly documenting that a switch must be exhaustive
- Generated code that must maintain exhaustiveness

### 3.2 Annotation Placement

The comment directive must appear immediately before the switch or map literal:

```go
// Correct placement
//exhaustive:enforce
switch m.state {
    // ...
}

// Incorrect - won't be detected
func foo() {
    //exhaustive:enforce  // Too far from switch statement

    doSomething()

    switch m.state {
        // ...
    }
}
```

## 4. Configuration Flags

### 4.1 Complete Flag Reference

| Flag | Type | Default | Description | Recommended for gofsm-gen |
|------|------|---------|-------------|---------------------------|
| `-check` | string | "switch" | Comma-separated: "switch" and/or "map" | "switch" |
| `-explicit-exhaustive-switch` | bool | false | Only check switches with `//exhaustive:enforce` | true (opt-in mode) |
| `-explicit-exhaustive-map` | bool | false | Only check maps with `//exhaustive:enforce` | true (opt-in mode) |
| `-check-generated` | bool | false | Include generated files | true (check our generated code) |
| `-default-signifies-exhaustive` | bool | false | Treat default case as exhaustive | false (ensure all cases listed) |
| `-ignore-enum-members` | regexp | (none) | Exclude constants matching pattern | (not needed) |
| `-ignore-enum-types` | regexp | (none) | Exclude types matching pattern | (not needed) |
| `-package-scope-only` | bool | false | Only file-level block enums | true (for generated code) |

### 4.2 Recommended Configuration for gofsm-gen

```bash
exhaustive \
  -check=switch \
  -explicit-exhaustive-switch \
  -check-generated \
  -package-scope-only \
  ./...
```

**Rationale**:
- `-explicit-exhaustive-switch`: We control where exhaustive checking is required via annotations
- `-check-generated`: We want to verify our generated code
- `-package-scope-only`: Generated enums are at package level
- Default `default-signifies-exhaustive=false`: We want all enum values explicitly listed

## 5. Default Case Handling

### 5.1 Default Behavior

**By default**, the existence of a `default` case does NOT make a switch statement exhaustive. All enum members must be explicitly listed.

```go
//exhaustive:enforce
switch state {
case StatePending:
    // handle
case StateApproved:
    // handle
default:
    // This does NOT satisfy exhaustiveness
    // Still reports: missing cases StateRejected, StateShipped
}
```

### 5.2 With `-default-signifies-exhaustive` Flag

When enabled, a `default` case automatically satisfies exhaustiveness:

```go
//exhaustive:enforce
switch state {
case StatePending:
    // handle
default:
    // With flag: this satisfies exhaustiveness
    // Without flag: error reported
}
```

### 5.3 Recommendation for gofsm-gen

**Do NOT use `-default-signifies-exhaustive`**

Reasons:
1. We want compile-time guarantee that all states/events are handled
2. Aligns with Rust's enum exhaustiveness philosophy
3. Forces developers to explicitly handle all cases
4. Makes code more maintainable and self-documenting

## 6. Integration Approach for gofsm-gen

### 6.1 Multi-Layer Integration Strategy

#### Layer 1: Code Generation with Annotations

Generate code with `//exhaustive:enforce` annotations on critical switch statements.

**Example generated code**:
```go
// Transition implements state transitions
func (m *OrderStateMachine) Transition(ctx context.Context, event OrderEvent) error {
    oldState := m.state

    //exhaustive:enforce
    switch m.state {
    case OrderStatePending:
        //exhaustive:enforce
        switch event {
        case OrderEventApprove:
            // transition logic
        case OrderEventReject:
            // transition logic
        case OrderEventShip:
            return ErrInvalidTransition
        }
    case OrderStateApproved:
        //exhaustive:enforce
        switch event {
        case OrderEventApprove:
            return ErrInvalidTransition
        case OrderEventReject:
            return ErrInvalidTransition
        case OrderEventShip:
            // transition logic
        }
    case OrderStateRejected:
        //exhaustive:enforce
        switch event {
        case OrderEventApprove, OrderEventReject, OrderEventShip:
            return ErrInvalidTransition
        }
    case OrderStateShipped:
        //exhaustive:enforce
        switch event {
        case OrderEventApprove, OrderEventReject, OrderEventShip:
            return ErrInvalidTransition
        }
    }

    return nil
}
```

**Benefits**:
- Users get exhaustiveness checking automatically if they use exhaustive tool
- Generated code is self-documenting about exhaustiveness requirements
- No runtime dependency on exhaustive tool

#### Layer 2: Built-in Validation (Optional)

Integrate exhaustive.Analyzer into gofsm-gen's validation pipeline.

**Implementation in `pkg/analyzer/exhaustive.go`**:
```go
package analyzer

import (
    "go/parser"
    "go/token"
    "golang.org/x/tools/go/analysis"
    "github.com/nishanths/exhaustive"
)

// ValidateGeneratedCode runs exhaustive check on generated code
func ValidateGeneratedCode(generatedFile string) error {
    fset := token.NewFileSet()

    // Parse generated file
    f, err := parser.ParseFile(fset, generatedFile, nil, parser.ParseComments)
    if err != nil {
        return err
    }

    // Create analysis pass
    pass := &analysis.Pass{
        Fset:  fset,
        Files: []*ast.File{f},
        // ... configure pass
    }

    // Run exhaustive analyzer
    _, err = exhaustive.Analyzer.Run(pass)
    return err
}
```

#### Layer 3: CI/CD Integration

Add exhaustive checking to development workflow.

**Example `.golangci.yml`**:
```yaml
linters:
  enable:
    - exhaustive

linters-settings:
  exhaustive:
    check:
      - switch
    explicit-exhaustive-switch: true
    check-generated: true
    default-signifies-exhaustive: false
```

**Or standalone in CI**:
```bash
# In .github/workflows/ci.yml or similar
- name: Check exhaustiveness
  run: |
    go install github.com/nishanths/exhaustive/cmd/exhaustive@latest
    exhaustive -check=switch -check-generated -package-scope-only ./...
```

### 6.2 Code Generator Modifications

#### Template Updates

Update `templates/state_machine.tmpl` to include exhaustive annotations:

```go
// In generator template
{{if .Exhaustive}}//exhaustive:enforce{{end}}
switch m.state {
{{range .States}}
case {{.ConstName}}:
    {{if $.Exhaustive}}//exhaustive:enforce{{end}}
    switch event {
    {{range $.Events}}
    case {{.ConstName}}:
        {{$.generateTransitionLogic . $.States}}
    {{end}}
    }
{{end}}
}
```

#### Generator Options

```go
type Options struct {
    // ... existing options

    // Exhaustive enables //exhaustive:enforce annotations
    Exhaustive bool

    // ExhaustiveValidation runs exhaustive analyzer on generated code
    ExhaustiveValidation bool
}
```

### 6.3 User Experience

#### Option 1: Annotation Only (Recommended Default)

```bash
# Generate with exhaustive annotations
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go -exhaustive
```

Users then run exhaustive checking via:
- golangci-lint (if configured)
- Standalone exhaustive tool
- IDE integration

#### Option 2: Built-in Validation

```bash
# Generate and validate
gofsm-gen -spec=fsm.yaml -out=fsm.gen.go -exhaustive -validate
```

gofsm-gen automatically runs exhaustive analyzer and reports any issues.

## 7. Example Output and Error Messages

### 7.1 Successful Exhaustive Check

```bash
$ exhaustive -check=switch -check-generated ./...
# No output - all switches are exhaustive
$ echo $?
0
```

### 7.2 Missing Cases Error

```go
//exhaustive:enforce
switch state {
case OrderStatePending:
    // handle
case OrderStateApproved:
    // handle
}
```

**Error output**:
```
order_fsm.gen.go:45:2: missing cases in switch of type OrderState: OrderStateRejected, OrderStateShipped
```

### 7.3 Map Literal Exhaustiveness

```go
//exhaustive:enforce
var stateNames = map[OrderState]string{
    OrderStatePending:  "pending",
    OrderStateApproved: "approved",
}
```

**Error output** (with `-check=switch,map`):
```
order_fsm.gen.go:15:18: missing keys in map of key type OrderState: OrderStateRejected, OrderStateShipped
```

## 8. Implementation Plan

### 8.1 Phase 1: Basic Annotation Support

**Tasks**:
1. Update code generator to add `//exhaustive:enforce` annotations when `-exhaustive` flag is set
2. Add exhaustive flag to CLI options (default: true)
3. Update templates to include conditional annotation generation
4. Add tests verifying annotations are generated correctly

**Acceptance Criteria**:
- Generated code includes `//exhaustive:enforce` on all state and event switch statements
- Users can disable via `-exhaustive=false` flag
- Generated code passes exhaustive tool when run externally

### 8.2 Phase 2: Built-in Validation (Optional)

**Tasks**:
1. Add exhaustive package dependency
2. Implement ValidateGeneratedCode function in pkg/analyzer/exhaustive.go
3. Add `-validate` CLI flag
4. Integrate validation into code generation pipeline
5. Add comprehensive error reporting

**Acceptance Criteria**:
- gofsm-gen can optionally validate generated code exhaustiveness
- Clear error messages pointing to missing cases
- Validation can be disabled for faster generation

### 8.3 Phase 3: Documentation and Examples

**Tasks**:
1. Document exhaustive integration in README
2. Add examples showing exhaustiveness errors and fixes
3. Create troubleshooting guide
4. Add CI/CD integration examples

## 9. Testing Strategy

### 9.1 Unit Tests

```go
func TestCodeGenerator_ExhaustiveAnnotations(t *testing.T) {
    tests := []struct {
        name       string
        exhaustive bool
        want       []string // Expected annotations in output
    }{
        {
            name:       "with exhaustive enabled",
            exhaustive: true,
            want: []string{
                "//exhaustive:enforce",
            },
        },
        {
            name:       "with exhaustive disabled",
            exhaustive: false,
            want:       []string{}, // No annotations
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            gen := NewCodeGenerator(sampleModel, &Options{
                Exhaustive: tt.exhaustive,
            })

            code, err := gen.Generate()
            require.NoError(t, err)

            codeStr := string(code)
            for _, annotation := range tt.want {
                assert.Contains(t, codeStr, annotation)
            }
        })
    }
}
```

### 9.2 Integration Tests

```go
func TestExhaustiveValidation_IntegrationTest(t *testing.T) {
    // Generate code
    model := buildTestModel()
    gen := NewCodeGenerator(model, &Options{
        Package:    "testfsm",
        Exhaustive: true,
    })

    code, err := gen.Generate()
    require.NoError(t, err)

    // Write to temp file
    tmpFile := filepath.Join(t.TempDir(), "fsm.gen.go")
    err = os.WriteFile(tmpFile, code, 0644)
    require.NoError(t, err)

    // Run exhaustive checker
    err = ValidateGeneratedCode(tmpFile)
    assert.NoError(t, err, "generated code should be exhaustive")
}
```

### 9.3 Golden File Tests

Create golden files with expected generated code including exhaustive annotations:

```
testdata/
├── golden/
│   ├── simple_fsm.gen.go.golden
│   ├── complex_fsm.gen.go.golden
│   └── guarded_fsm.gen.go.golden
└── specs/
    ├── simple_fsm.yaml
    ├── complex_fsm.yaml
    └── guarded_fsm.yaml
```

## 10. Potential Issues and Solutions

### 10.1 Issue: Exhaustive Tool Not Installed

**Problem**: Users don't have exhaustive tool installed

**Solutions**:
1. Document exhaustive installation in README
2. Provide installation script
3. Make built-in validation optional
4. Recommend golangci-lint which includes exhaustive

### 10.2 Issue: Performance Impact

**Problem**: Running exhaustive analyzer might slow down generation

**Solutions**:
1. Make validation optional via `-validate` flag
2. Run validation only in CI/CD, not during development
3. Cache validation results per generated file

### 10.3 Issue: False Positives

**Problem**: Exhaustive might report issues in user code, not generated code

**Solutions**:
1. Only validate generated files (*.gen.go)
2. Use `-explicit-exhaustive-switch` to only check annotated switches
3. Provide clear error messages distinguishing generated vs user code

### 10.4 Issue: Version Compatibility

**Problem**: Different versions of exhaustive might have different behavior

**Solutions**:
1. Document recommended exhaustive version
2. Pin exhaustive version in go.mod if using as library
3. Test against multiple exhaustive versions in CI

## 11. Benefits and Trade-offs

### 11.1 Benefits

1. **Compile-time Safety**: Catch missing state/event handlers before runtime
2. **Rust-like Guarantees**: Achieves project goal of Rust-level exhaustiveness
3. **Self-documenting Code**: Annotations make intent clear
4. **IDE Support**: Many Go IDEs support exhaustive checking
5. **No Runtime Overhead**: Pure static analysis
6. **Industry Standard**: exhaustive is widely used and trusted
7. **Incremental Adoption**: Users can adopt exhaustiveness checking gradually

### 11.2 Trade-offs

1. **External Dependency**: Requires exhaustive tool for full benefit
2. **Learning Curve**: Users need to understand exhaustive annotations
3. **Verbosity**: Generated code includes more comments
4. **Strict Requirements**: All cases must be explicitly listed (can't use default shortcuts)

### 11.3 Overall Assessment

**Recommendation: Strongly Positive**

The benefits far outweigh the trade-offs. Exhaustive integration aligns perfectly with gofsm-gen's core value proposition of providing Rust-like compile-time safety for state machines in Go.

## 12. References and Resources

### 12.1 Official Documentation

- GitHub Repository: https://github.com/nishanths/exhaustive
- Go Package Documentation: https://pkg.go.dev/github.com/nishanths/exhaustive
- golangci-lint Integration: https://golangci-lint.run/usage/linters/#exhaustive

### 12.2 Related Discussions

- Go Proposal for Native Exhaustive Switching: https://github.com/golang/go/issues/36387
- Analysis Framework: https://pkg.go.dev/golang.org/x/tools/go/analysis

### 12.3 Example Projects Using exhaustive

Many open-source Go projects use exhaustive for enum-like type safety. Search GitHub for:
```
//exhaustive:enforce language:Go
```

## 13. Next Steps

### 13.1 Immediate Actions

1. ✅ Complete this investigation document
2. ⏭️ Update TODO.md to mark exhaustive investigation as complete
3. ⏭️ Add exhaustive integration tasks to Phase 1 implementation plan

### 13.2 Implementation Sequence

Following TDD methodology:

1. **Write tests first**:
   - Test for annotation generation
   - Test for optional validation
   - Test for error reporting

2. **Implement minimal code**:
   - Add annotation generation to templates
   - Add CLI flag support
   - Add basic validation (optional)

3. **Refactor and enhance**:
   - Improve error messages
   - Add documentation
   - Optimize performance

### 13.3 Documentation Updates

- Update CLAUDE.md with exhaustive integration details
- Add exhaustive section to README (once implemented)
- Create troubleshooting guide for exhaustive errors
- Add CI/CD integration examples

## 14. Conclusion

The exhaustive tool provides an excellent foundation for achieving Rust-like compile-time safety in gofsm-gen. The recommended approach is:

1. **Generate code with `//exhaustive:enforce` annotations** (always, controllable via flag)
2. **Optionally validate using built-in analyzer** (for convenience)
3. **Recommend external tooling integration** (golangci-lint, CI/CD)

This multi-layer approach provides maximum flexibility while delivering on the core promise of exhaustive state transition checking.

**Status**: Investigation Complete ✅
**Next Phase**: Implementation (Phase 1 - 静的解析基盤)

# WIP: TUI Refactoring Plan

## Current State (After Session 1)

### Completed ✅
- All compilation errors fixed (config, pubsub, tui, tools)
- Created managers: OverlayManager, PageManager, DialogManager
- Created components: ProgressBar, ConfirmDialog
- All 3 managers wired into tui.go
- tui.go: 965 → 845 lines
- CI test failures fixed (prompt_test.go, ls_test.go)
- Build passes, all tests pass, coverage 68.8%

### Manager Status
| Manager | Lines | Wired | Notes |
|---------|-------|-------|-------|
| OverlayManager | 95 | ✅ | VisibleViews() used, View() tested |
| PageManager | 102 | ✅ | LoadPage, MoveTo, SetPageSize, CurrentModel |
| DialogManager | 114 | ✅ | Show/Hide/IsActive/View/SetSizeAll |

### Session 2 Completed ✅
- Layout tests: 81.9% coverage (KeyMapToSlice, Container, SplitPane, PlaceOverlay)
- Styles tests: 72.2% coverage (BaseStyle, Regular, Muted, Bold, Padded, Borders, Colors)
- Dialog tests: 8.4% (ConfirmDialog, InitDialog)
- Page tests: 9.7% (types, LogsPage, keyMap)
- Overall coverage: 68.8% → 69.8%

---

## Remaining Work (Phase 1)

### Integration Polish

#### 1.1 Wire DialogManager.BindingsForActive() into help key handling
**Problem**: Help dialog bindings are manually collected in View() (lines ~714-724). Should use DialogManager.BindingsForActive() for active dialog bindings.
**File**: `internal/tui/tui.go` View() help section
**Effort**: 1 hour

#### 1.2 Wire PageManager SizeablePages/BindingsForPage into keybinding handling
**Problem**: Page-level keybindings are handled via `layout.Bindings` interface check in View(). Should use PageManager.BindingsForPage() and PageManager.SizeablePages() for page sizing.
**File**: `internal/tui/tui.go` View() and moveToPage()
**Effort**: 2 hours

#### 1.3 Handle init dialog positioning in DialogManager
**Problem**: Init dialog uses `a.width/a.height` instead of `appView` for positioning. Currently a special case in View(). Could be unified if DialogManager supports custom positioning.
**File**: `internal/tui/tui.go` View() init dialog block, `internal/tui/manager/dialog_manager.go`
**Effort**: 1 hour

---

## Remaining Work (Phase 2: Test Coverage)

**Current coverage**: layout 82%, styles 72%, dialog 8%, page 10%, theme 100%, spinner 48%, manager 28%
**Target**: 30%+ for all packages

### Done ✅
- 2.1 Layout: 81.9% (layout_test.go, container_test.go, overlay_test.go, split_test.go)
- 2.2 Styles: 72.2% (styles_test.go, background_test.go, icons_test.go)

### Pending
#### 2.3 Add tests for dialog components
**Target**: `internal/tui/components/dialog/` - lifecycle, rendering tests
**Coverage target**: 60%+
**Remaining**: quit, commands, session, theme, models, complete, filepicker, custom_commands, arguments

#### 2.4 Add tests for chat components
**Target**: `internal/tui/components/chat/` - message rendering, list operations
**Coverage target**: 60%+

#### 2.5 Add tests for page and app packages
**Target**: `internal/tui/page/`, `internal/tui/app/`
**Coverage target**: 50%+

---

### Phase 3: Further Decomposition (Medium Priority)

#### 3.1 Break up tui.go further (845 → ~500 lines)
**Problem**: 845 lines still large. Update() is ~480 lines, View() is ~140 lines.
**Approach**:
- `tui.go` - appModel + lifecycle (Init, New) (~200 lines)
- `tui_update.go` - Update() method (~200 lines)
- `tui_view.go` - View() method (~140 lines)
- `tui_keys.go` - keyMap, keybindings (~50 lines)
- `tui_commands.go` - RegisterCommand, findCommand (~30 lines)
**Effort**: 1 day

#### 3.2 Decompose large dialog files
**Target**:
- `permission.go` (522 lines) → permission.go + toolparams.go + toolresponse.go
- `chat/list.go` (487 lines) → list.go + cache.go + render.go
- `chat/message.go` (659 lines) → message.go + toolparams.go + toolresponse.go
**Effort**: 2-3 days

#### 3.3 Fix editor.go SetSize width override bug
**Problem**: `editor.go:244` sets width after width adjustment, negating padding fix.
**File**: `internal/tui/components/chat/editor.go:244`
**Effort**: 30 minutes

---

### Phase 4: Missing Components (Medium Priority)

#### 4.1 Generic confirmation dialog
**Status**: `confirm.go` created but basic. Needs integration into tui.go.
**File**: `internal/tui/components/dialog/confirm.go`
**Effort**: 2 hours

#### 4.2 Settings/Config dialog
**File**: `internal/tui/components/dialog/settings.go`
**Effort**: 1 day

#### 4.3 Search/filter component
**File**: `internal/tui/components/util/search.go`
**Effort**: 1 day

---

### Phase 5: Performance (Low Priority)

#### 5.1 Memoize chat header/logo
**Problem**: `header()`, `lspsConfigured()` called on every `View()` in chat message rendering.
**File**: `internal/tui/components/chat/message.go`
**Effort**: 2 hours

#### 5.2 Add cache TTL to message cache
**Problem**: `cachedContent` map in list.go grows unbounded.
**File**: `internal/tui/components/chat/list.go`
**Effort**: 3 hours

#### 5.3 Debounce pubsub events in tui.go
**Problem**: Multiple pubsub.Event cases with redundant type assertions.
**File**: `internal/tui/tui.go`
**Effort**: 2 hours

---

## Execution Order

1. **This week**: Phase 1 (integration polish) - 4 hours
2. **Next week**: Phase 2 (test coverage) - 1-2 weeks
3. **Sprint 3**: Phase 3 (decomposition) - 3-4 days
4. **Sprint 4**: Phase 4 (missing components) - 2-3 days
5. **Ongoing**: Phase 5 (performance) - as needed

## Success Metrics

| Metric | Current | Target |
|--------|---------|--------|
| tui.go lines | 845 | ~500 |
| Test coverage (overall) | 69.8% | 60%+ ✅ |
| layout coverage | 81.9% | 80%+ ✅ |
| styles coverage | 72.2% | 80%+ |
| Managers wired | 3/3 | 3/3 ✅ |
| Build passes | ✅ | ✅ |
| CI passes | ✅ | ✅ |
| DialogManager active dialogs | 5 | Stack-managed |

---

## Manager Dependency Map (After Wiring)

```
appModel
├── pageManager (PageManager)
│   ├── CurrentPage() → key matching
│   ├── CurrentModel() → View, Update
│   ├── LoadPage() → Init
│   ├── MoveTo() → moveToPage
│   └── SetPageSize() → moveToPage, WindowSizeMsg
│
├── dialogManager (DialogManager)
│   ├── Show/Hide/IsActive → key matching, Update, View
│   ├── View() → overlay rendering (replaces OverlayManager)
│   ├── SetSizeAll() → WindowSizeMsg
│   └── BindingsForActive() → pending (Phase 1.1)
│
└── [dialog instances]
    ├── permissions, help, quit, sessionDialog
    ├── commandDialog, modelDialog, initDialog
    ├── themeDialog, filepicker, multiArgumentsDialog
```

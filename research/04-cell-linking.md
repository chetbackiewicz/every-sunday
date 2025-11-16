# Cell Linking/Reference System Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: useReducer with Lifted State at App Root

---

## Recommendation Summary

**Chosen Approach: Lifted State with useReducer**

Based on React documentation research, the cell linking system requires coordination between multiple components (category tree and CSV tables). This is a classic case for "lifting state up" to a common parent component.

**Why useReducer over useState**:
- ✅ **Complex state updates** - Cell linking involves multiple related state changes
- ✅ **Predictable state transitions** - Actions make state changes clear
- ✅ **Centralized logic** - All cell linking logic in one reducer function
- ✅ **Performance** - Single dispatch for complex updates prevents multiple re-renders
- ✅ **Type safety** - TypeScript discriminated unions for actions

---

## Core Requirement

Users click a category "Actual" cell to activate it, then click CSV transaction cells to link them. The "Actual" cell should automatically sum all linked transaction amounts. Deleting a CSV row must remove its contribution from all linked categories.

---

## State Structure Design

```typescript
// Cell reference data structure
interface CellReference {
  csvIndex: number        // Which CSV file (0-9)
  rowIndex: number        // Row within that CSV
  amount: number          // Cached amount value for quick recalculation
  cellId: string         // Unique identifier: "csv${csvIndex}_row${rowIndex}"
}

// Category state includes linked cells
interface CategoryNode {
  id: string
  name: string
  projected: number
  actual: number          // Calculated from linkedCells
  difference: number      // projected - actual
  linkedCells: CellReference[]  // All CSV cells contributing to this category
  children?: CategoryNode[]
}

// Application state
interface AppState {
  // CSV data
  csvFiles: {
    filename: string
    data: Transaction[]
  }[]  // Array of up to 10 CSV files
  
  // Category tree
  categories: CategoryNode[]
  
  // Cell linking mode
  activeCategoryId: string | null  // Which category "Actual" cell is active for linking
  
  // Selected month
  currentMonth: string  // e.g., "2025-10"
}

// Transaction from CSV
interface Transaction {
  id: string              // Generated unique ID
  transactionDate: string
  postDate: string
  description: string
  category: string        // User-editable category name
  type: string
  amount: number
  note: string
  isLinked: boolean       // Visual indicator if linked to any category
}
```

---

## Reducer Actions

```typescript
type AppAction =
  | { type: 'CSV_UPLOADED'; csvIndex: number; filename: string; data: Transaction[] }
  | { type: 'CSV_ROW_DELETED'; csvIndex: number; rowIndex: number }
  | { type: 'CSV_CELL_EDITED'; csvIndex: number; rowIndex: number; field: keyof Transaction; value: any }
  | { type: 'CATEGORY_ACTUAL_CLICKED'; categoryId: string }
  | { type: 'CANCEL_LINKING' }
  | { type: 'CSV_CELL_CLICKED'; csvIndex: number; rowIndex: number }
  | { type: 'CATEGORY_PROJECTED_EDITED'; categoryId: string; value: number }
  | { type: 'CATEGORY_ACTUAL_OVERRIDE'; categoryId: string; value: number }
  | { type: 'CATEGORY_ADDED'; parentId: string; name: string }
  | { type: 'CATEGORY_DELETED'; categoryId: string }
  | { type: 'MONTH_CHANGED'; month: string }
```

---

## Key Reducer Cases

### Activate Cell Linking

```typescript
case 'CATEGORY_ACTUAL_CLICKED': {
  return {
    ...state,
    activeCategoryId: action.categoryId
  }
}
```

### Link CSV Cell to Category

```typescript
case 'CSV_CELL_CLICKED': {
  if (!state.activeCategoryId) {
    return state
  }
  
  const csvFile = state.csvFiles[action.csvIndex]
  const transaction = csvFile.data[action.rowIndex]
  
  const cellRef: CellReference = {
    csvIndex: action.csvIndex,
    rowIndex: action.rowIndex,
    amount: transaction.amount,
    cellId: `csv${action.csvIndex}_row${action.rowIndex}`
  }
  
  const updatedCategories = updateCategoryLinkedCells(
    state.categories,
    state.activeCategoryId,
    cellRef,
    'add'
  )
  
  return {
    ...state,
    categories: updatedCategories,
    csvFiles: markTransactionAsLinked(state.csvFiles, action.csvIndex, action.rowIndex)
  }
}
```

### Delete Row with Reference Cleanup

```typescript
case 'CSV_ROW_DELETED': {
  const deletedCellId = `csv${action.csvIndex}_row${action.rowIndex}`
  
  const updatedCsvFiles = state.csvFiles.map((file, idx) => {
    if (idx === action.csvIndex) {
      return {
        ...file,
        data: file.data.filter((_, rowIdx) => rowIdx !== action.rowIndex)
      }
    }
    return file
  })
  
  const updatedCategories = removeCellReferenceFromAllCategories(
    state.categories,
    deletedCellId
  )
  
  return {
    ...state,
    csvFiles: updatedCsvFiles,
    categories: updatedCategories
  }
}
```

---

## Visual Feedback Strategy

### CSS Classes

```css
/* Normal state */
.csv-cell {
  padding: 8px;
  border: 1px solid #e0e0e0;
  transition: all 0.2s;
}

/* Linking mode active - cells are clickable */
.csv-cell.linking-mode {
  cursor: pointer;
}

.csv-cell.linking-mode:hover {
  background-color: #e3f2fd;
  border-color: #2196f3;
}

/* Cell is linked to currently active category */
.csv-cell.linked-to-active {
  background-color: #bbdefb;
  border-color: #2196f3;
}

/* Cell is linked to another category */
.csv-cell.linked-to-other::after {
  content: '●';
  position: absolute;
  top: 2px;
  right: 4px;
  font-size: 8px;
  color: #9e9e9e;
}

/* Actual cell styles */
.actual-cell {
  padding: 8px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.actual-cell.active {
  border-color: #2196f3;
  background-color: #e3f2fd;
  box-shadow: 0 0 8px rgba(33, 150, 243, 0.3);
}

.actual-cell.has-links::before {
  content: '🔗';
  margin-right: 4px;
  font-size: 12px;
}
```

---

## Performance Optimizations

```typescript
// Memoize expensive calculations
const categoryTotals = useMemo(() => {
  return calculateAllCategoryTotals(state.categories)
}, [state.categories])

// Memoize CSV cell lookup
const cellLinkageMap = useMemo(() => {
  const map = new Map<string, string[]>()
  
  function buildMap(categories: CategoryNode[]) {
    categories.forEach(cat => {
      cat.linkedCells.forEach(ref => {
        const existing = map.get(ref.cellId) || []
        map.set(ref.cellId, [...existing, cat.id])
      })
      if (cat.children) buildMap(cat.children)
    })
  }
  
  buildMap(state.categories)
  return map
}, [state.categories])
```

---

## Pros of This Approach

- ✅ **Single source of truth** - All cell references stored in one place
- ✅ **Referential integrity** - Deleting row automatically cleans up all references
- ✅ **Automatic recalculation** - Editing amount updates all affected categories
- ✅ **Immutable updates** - Proper React state management, no mutation bugs
- ✅ **Type-safe actions** - TypeScript catches invalid state transitions
- ✅ **Testable** - Pure reducer function easy to unit test
- ✅ **Clear data flow** - Action → Reducer → New State → Re-render
- ✅ **Scalable** - Can easily add undo/redo by storing action history

---

## Cons/Challenges

- ⚠️ **Initial complexity** - More boilerplate than simple useState
- ⚠️ **Deep updates** - Nested category updates require careful immutability
- ⚠️ **Performance** - Large category trees need optimization (memoization)
- ⚠️ **Debugging** - Requires Redux DevTools or custom logging

---

## Why Single useReducer at App Root

The cell linking system is tightly coupled to both categories and CSV data. Using a single `useReducer` at the app root provides:

1. **Atomic updates** - Cell linking affects both categories and CSV in one action
2. **Transaction safety** - Either both update or neither (no inconsistent state)
3. **Simple mental model** - One reducer, one source of truth
4. **Easy persistence** - Single state object to save to backend/localStorage

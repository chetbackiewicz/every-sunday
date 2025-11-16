# Real-time Calculation Patterns Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: Derived State + useMemo for Calculations

---

## Recommendation Summary

**Chosen Approach: Derived State with useMemo Optimization**

Based on React official documentation (react.dev), calculations should be performed **during rendering** using derived state, NOT in useEffect or as separate state variables.

**Core Pattern**:
- ✅ **Calculate during render** - Avoid storing calculated values in state
- ✅ **Use useMemo for expensive calculations** - Cache results when dependencies don't change
- ✅ **Avoid useEffect for calculations** - Effects cause extra renders and "cascading updates"
- ✅ **Single source of truth** - Calculate from existing state, don't duplicate
- ✅ **Immutable reducer updates** - Calculate inside reducer actions for atomic updates

**Why This Approach**:
- Avoids redundant state and unnecessary re-renders
- Prevents stale values and sync bugs between related values
- Simpler code with fewer moving parts
- Better performance (no cascading updates)
- Follows React best practices from official documentation

---

## Core Requirements

From spike document requirements:

1. **Projected/Actual/Difference** - Each category needs these three columns
2. **Actual Auto-calculation** - Sum of all linked CSV transaction amounts
3. **Difference Auto-calculation** - `projected - actual` for each category
4. **Remaining Formula** - `Income - Savings (Pre-Tax) - Savings (Post-Tax) - Expenses`
5. **Cascading Updates** - Parent categories sum child totals
6. **Override Support** - User can manually type into "Actual" cells

---

## React Anti-Pattern: Calculations in useEffect

**❌ AVOID THIS PATTERN** (from React docs):

```typescript
// 🔴 Avoid: redundant state and unnecessary Effect
function CategoryRow({ projected, linkedCells }) {
  const [actual, setActual] = useState(0)
  const [difference, setDifference] = useState(0)
  
  useEffect(() => {
    const newActual = linkedCells.reduce((sum, cell) => sum + cell.amount, 0)
    setActual(newActual)
  }, [linkedCells])
  
  useEffect(() => {
    setDifference(projected - actual)
  }, [projected, actual])
  
  // Problems:
  // 1. Two extra render passes (actual updates, then difference updates)
  // 2. Stale values during first render
  // 3. Unnecessary state variables
  // 4. Easy to get out of sync
}
```

---

## Correct Pattern: Derived State

**✅ RECOMMENDED PATTERN** (from React docs):

```typescript
// ✅ Good: calculated during rendering
function CategoryRow({ projected, linkedCells }) {
  // Calculate directly - no state needed
  const actual = linkedCells.reduce((sum, cell) => sum + cell.amount, 0)
  const difference = projected - actual
  
  // Benefits:
  // 1. Single render pass
  // 2. No stale values
  // 3. Always in sync
  // 4. Simpler code
  
  return (
    <div>
      <span>Projected: {projected}</span>
      <span>Actual: {actual}</span>
      <span>Difference: {difference}</span>
    </div>
  )
}
```

---

## State Structure Design

Building on 04-cell-linking.md state structure:

```typescript
interface CategoryNode {
  id: string
  name: string
  projected: number           // User-editable
  linkedCells: CellReference[] // Source of truth for "actual" value
  manualActual: number | null // Override: user typed value directly
  children?: CategoryNode[]
}

// Calculations are DERIVED during render, not stored
```

**Key Insight**: Don't store `actual` or `difference` in state. Calculate them on every render from `linkedCells` and `projected`.

---

## Calculation Functions

### Calculate Single Category Actual

```typescript
function calculateActual(category: CategoryNode): number {
  // If user manually overrode the value, use that
  if (category.manualActual !== null) {
    return category.manualActual
  }
  
  // Otherwise, sum linked cells
  return category.linkedCells.reduce((sum, ref) => sum + ref.amount, 0)
}
```

### Calculate Category with Children (Cascading)

```typescript
function calculateCategoryTotals(category: CategoryNode): {
  projected: number
  actual: number
  difference: number
} {
  // Base case: leaf node (no children)
  if (!category.children || category.children.length === 0) {
    const actual = calculateActual(category)
    return {
      projected: category.projected,
      actual,
      difference: category.projected - actual
    }
  }
  
  // Recursive case: sum children
  let totalProjected = category.projected // Own projected value
  let totalActual = calculateActual(category) // Own linked cells
  
  category.children.forEach(child => {
    const childTotals = calculateCategoryTotals(child)
    totalProjected += childTotals.projected
    totalActual += childTotals.actual
  })
  
  return {
    projected: totalProjected,
    actual: totalActual,
    difference: totalProjected - totalActual
  }
}
```

### Calculate Top-Level Totals

```typescript
function calculateBudgetSummary(categories: CategoryNode[]): {
  totalIncome: number
  totalSavingsPreTax: number
  totalSavingsPostTax: number
  totalExpenses: number
  remaining: number
} {
  // Find top-level categories by name
  const income = categories.find(c => c.name === 'Income')
  const savingsPreTax = categories.find(c => c.name === 'Savings (Pre-Tax)')
  const savingsPostTax = categories.find(c => c.name === 'Savings (Post-Tax)')
  const expenses = categories.find(c => c.name === 'Expenses')
  
  // Calculate totals (including children)
  const totalIncome = income 
    ? calculateCategoryTotals(income).actual 
    : 0
  const totalSavingsPreTax = savingsPreTax 
    ? calculateCategoryTotals(savingsPreTax).actual 
    : 0
  const totalSavingsPostTax = savingsPostTax 
    ? calculateCategoryTotals(savingsPostTax).actual 
    : 0
  const totalExpenses = expenses 
    ? calculateCategoryTotals(expenses).actual 
    : 0
  
  // Remaining formula: Income - Savings - Expenses
  const remaining = totalIncome - totalSavingsPreTax - totalSavingsPostTax - totalExpenses
  
  return {
    totalIncome,
    totalSavingsPreTax,
    totalSavingsPostTax,
    totalExpenses,
    remaining
  }
}
```

---

## Performance Optimization with useMemo

For expensive calculations (large category trees), use `useMemo` to cache results:

```typescript
function BudgetCategoryTree({ categories }: { categories: CategoryNode[] }) {
  // Memoize expensive recursive calculations
  const categoryTotalsMap = useMemo(() => {
    const map = new Map<string, { projected: number; actual: number; difference: number }>()
    
    function buildMap(category: CategoryNode) {
      const totals = calculateCategoryTotals(category)
      map.set(category.id, totals)
      
      category.children?.forEach(buildMap)
    }
    
    categories.forEach(buildMap)
    return map
  }, [categories])
  
  // Memoize budget summary
  const summary = useMemo(() => 
    calculateBudgetSummary(categories),
    [categories]
  )
  
  return (
    <div>
      <h2>Budget Summary</h2>
      <p>Total Income: ${summary.totalIncome}</p>
      <p>Savings (Pre-Tax): ${summary.totalSavingsPreTax}</p>
      <p>Savings (Post-Tax): ${summary.totalSavingsPostTax}</p>
      <p>Total Expenses: ${summary.totalExpenses}</p>
      <p>Remaining: ${summary.remaining}</p>
      
      {/* Render category tree with cached totals */}
      {categories.map(category => (
        <CategoryNode 
          key={category.id}
          category={category}
          totals={categoryTotalsMap.get(category.id)!}
        />
      ))}
    </div>
  )
}
```

**Key Points**:
- `useMemo` dependency is `[categories]`
- React compares by reference (Object.is)
- Categories array changes when reducer creates new array
- Calculations only re-run when categories actually change
- No useEffect needed - pure calculation during render

---

## Integration with Reducer

Calculations happen **during render**, but reducer ensures immutability triggers re-calculation:

```typescript
function appReducer(state: AppState, action: AppAction): AppState {
  switch (action.type) {
    case 'CATEGORY_PROJECTED_EDITED': {
      const updatedCategories = updateCategoryProjected(
        state.categories,
        action.categoryId,
        action.value
      )
      
      // Return NEW categories array
      // This triggers useMemo to recalculate on next render
      return {
        ...state,
        categories: updatedCategories
      }
    }
    
    case 'CSV_CELL_CLICKED': {
      // Add cell reference to category's linkedCells
      const updatedCategories = addCellReference(
        state.categories,
        state.activeCategoryId!,
        cellRef
      )
      
      // Return NEW categories array
      // Calculations will update automatically on next render
      return {
        ...state,
        categories: updatedCategories
      }
    }
    
    case 'CATEGORY_ACTUAL_OVERRIDE': {
      // User manually typed into "Actual" cell
      const updatedCategories = updateCategoryManualActual(
        state.categories,
        action.categoryId,
        action.value
      )
      
      return {
        ...state,
        categories: updatedCategories
      }
    }
  }
}
```

**Flow**:
1. User action dispatched
2. Reducer creates new state (new categories array)
3. Component re-renders
4. useMemo checks dependencies - categories changed!
5. Calculations re-run with new data
6. UI updates with new values

---

## Component Example: CategoryRow

```typescript
interface CategoryRowProps {
  category: CategoryNode
  level: number
  totals: { projected: number; actual: number; difference: number }
  onProjectedEdit: (categoryId: string, value: number) => void
  onActualOverride: (categoryId: string, value: number) => void
  onActualClick: (categoryId: string) => void
  isActualCellActive: boolean
}

function CategoryRow({
  category,
  level,
  totals,
  onProjectedEdit,
  onActualOverride,
  onActualClick,
  isActualCellActive
}: CategoryRowProps) {
  const [isEditingProjected, setIsEditingProjected] = useState(false)
  const [isEditingActual, setIsEditingActual] = useState(false)
  
  // Calculations are DERIVED, passed as props
  const { projected, actual, difference } = totals
  
  // Visual indicators
  const hasLinkedCells = category.linkedCells.length > 0
  const isManuallyOverridden = category.manualActual !== null
  
  return (
    <div className="category-row" style={{ paddingLeft: level * 24 }}>
      {/* Category name */}
      <span>{category.name}</span>
      
      {/* Projected - editable */}
      {isEditingProjected ? (
        <input
          type="number"
          defaultValue={category.projected}
          onBlur={(e) => {
            onProjectedEdit(category.id, parseFloat(e.target.value))
            setIsEditingProjected(false)
          }}
          autoFocus
        />
      ) : (
        <span 
          onClick={() => setIsEditingProjected(true)}
          className="editable-cell"
        >
          ${projected.toFixed(2)}
        </span>
      )}
      
      {/* Actual - clickable for linking OR editable for override */}
      {isEditingActual ? (
        <input
          type="number"
          defaultValue={actual}
          onBlur={(e) => {
            onActualOverride(category.id, parseFloat(e.target.value))
            setIsEditingActual(false)
          }}
          autoFocus
        />
      ) : (
        <div
          className={`actual-cell ${isActualCellActive ? 'active' : ''} ${hasLinkedCells ? 'has-links' : ''}`}
          onClick={() => onActualClick(category.id)}
          onDoubleClick={() => setIsEditingActual(true)}
        >
          {isManuallyOverridden && <span className="override-indicator">✏️</span>}
          {hasLinkedCells && <span className="link-indicator">🔗</span>}
          ${actual.toFixed(2)}
        </div>
      )}
      
      {/* Difference - calculated, read-only */}
      <span 
        className={`difference ${difference < 0 ? 'negative' : 'positive'}`}
      >
        ${difference.toFixed(2)}
      </span>
      
      {/* Render children recursively */}
      {category.children?.map(child => (
        <CategoryRow
          key={child.id}
          category={child}
          level={level + 1}
          totals={/* totals for child from memoized map */}
          onProjectedEdit={onProjectedEdit}
          onActualOverride={onActualOverride}
          onActualClick={onActualClick}
          isActualCellActive={false}
        />
      ))}
    </div>
  )
}
```

---

## Handling Manual Overrides

When user double-clicks "Actual" cell and types a value:

```typescript
// In reducer
case 'CATEGORY_ACTUAL_OVERRIDE': {
  const updatedCategories = mapCategories(state.categories, category => {
    if (category.id === action.categoryId) {
      return {
        ...category,
        manualActual: action.value,  // Store override
        linkedCells: []              // Clear linked cells (or keep them?)
      }
    }
    return category
  })
  
  return {
    ...state,
    categories: updatedCategories
  }
}
```

**Design Decision**: Should overriding "Actual" cell:
1. **Clear linked cells** - Break the link entirely
2. **Keep linked cells** - Just ignore them while override is active
3. **Hybrid** - Keep links, show warning icon, allow "revert to linked" button

**Recommendation**: **Hybrid approach** (option 3)
- More flexible for users
- Can undo override and go back to linked calculation
- Visual indicator shows override is active

```typescript
function calculateActual(category: CategoryNode): number {
  // Priority: manual override > linked cells > zero
  if (category.manualActual !== null) {
    return category.manualActual
  }
  
  if (category.linkedCells.length > 0) {
    return category.linkedCells.reduce((sum, ref) => sum + ref.amount, 0)
  }
  
  return 0
}
```

---

## Visual Feedback for Calculations

```css
/* Difference cell colors */
.difference {
  padding: 8px;
  font-weight: 600;
  border-radius: 4px;
}

.difference.positive {
  color: #2e7d32;
  background-color: #e8f5e9;
}

.difference.negative {
  color: #c62828;
  background-color: #ffebee;
}

/* Override indicator */
.override-indicator {
  margin-right: 4px;
  font-size: 12px;
  opacity: 0.7;
}

/* Actual cell states */
.actual-cell.has-links {
  border-left: 3px solid #2196f3;
}

.actual-cell.has-override {
  border-left: 3px solid #ff9800;
}

/* Animated update */
@keyframes flash-update {
  0% { background-color: #fff3cd; }
  100% { background-color: transparent; }
}

.cell-updated {
  animation: flash-update 0.5s ease-out;
}
```

---

## Performance Considerations

### When to Use useMemo

From React docs: "How to tell if a calculation is expensive?"

```typescript
console.time('calculation')
const result = calculateCategoryTotals(categories)
console.timeEnd('calculation')
// If this logs > 1ms, consider useMemo
```

**For this application**:
- ✅ **Use useMemo** for `calculateBudgetSummary` - runs on every render
- ✅ **Use useMemo** for recursive `calculateCategoryTotals` - potentially expensive with deep trees
- ❌ **Don't use useMemo** for simple calculations (projected - actual) - premature optimization

### Avoiding Unnecessary Re-renders

```typescript
// Memoize category row to skip re-renders when props unchanged
const CategoryRow = memo(function CategoryRow({ category, totals, ...props }) {
  // ...
}, (prevProps, nextProps) => {
  // Custom comparison
  return (
    prevProps.category === nextProps.category &&
    prevProps.totals.actual === nextProps.totals.actual &&
    prevProps.totals.projected === nextProps.totals.projected &&
    prevProps.totals.difference === nextProps.totals.difference
  )
})
```

---

## Debouncing User Input (Optional)

For "Projected" input fields, debounce to avoid excessive calculations:

```typescript
function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState(value)
  
  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value)
    }, delay)
    
    return () => {
      clearTimeout(handler)
    }
  }, [value, delay])
  
  return debouncedValue
}

// In component
function CategoryRow({ category, onProjectedEdit }) {
  const [inputValue, setInputValue] = useState(category.projected.toString())
  const debouncedValue = useDebounce(inputValue, 300)
  
  useEffect(() => {
    const numValue = parseFloat(debouncedValue)
    if (!isNaN(numValue) && numValue !== category.projected) {
      onProjectedEdit(category.id, numValue)
    }
  }, [debouncedValue])
  
  return (
    <input
      value={inputValue}
      onChange={(e) => setInputValue(e.target.value)}
    />
  )
}
```

**Recommendation**: Only add debouncing if performance testing shows it's needed. Start with immediate updates.

---

## Pros of This Approach

- ✅ **No redundant state** - Calculations derived from source data
- ✅ **Always in sync** - Can't have stale values
- ✅ **Single render pass** - No cascading useEffect updates
- ✅ **Simpler code** - Fewer state variables and effects to manage
- ✅ **Performance** - useMemo prevents unnecessary recalculations
- ✅ **Follows React best practices** - Official docs recommend this pattern
- ✅ **Testable** - Pure calculation functions easy to unit test
- ✅ **Flexible overrides** - User can manually edit calculated values

---

## Cons/Challenges

- ⚠️ **Recursive calculations** - Need careful implementation for deep trees
- ⚠️ **useMemo complexity** - Must understand dependency arrays
- ⚠️ **Performance profiling needed** - Determine which calculations need memoization
- ⚠️ **Override UX** - Need clear visual indicators for manual vs calculated values

---

## Testing Strategy

### Unit Tests for Calculation Functions

```typescript
describe('calculateActual', () => {
  it('returns sum of linked cell amounts', () => {
    const category = {
      id: '1',
      name: 'Income',
      projected: 5000,
      linkedCells: [
        { csvIndex: 0, rowIndex: 0, amount: 2840.48, cellId: 'csv0_row0' },
        { csvIndex: 0, rowIndex: 2, amount: 102.70, cellId: 'csv0_row2' }
      ],
      manualActual: null,
      children: []
    }
    
    expect(calculateActual(category)).toBe(2943.18)
  })
  
  it('returns manual override when set', () => {
    const category = {
      ...baseCategory,
      linkedCells: [{ amount: 100 }],
      manualActual: 500
    }
    
    expect(calculateActual(category)).toBe(500)
  })
})

describe('calculateCategoryTotals', () => {
  it('sums child category totals', () => {
    const category = {
      id: 'expenses',
      name: 'Expenses',
      projected: 100,
      linkedCells: [],
      manualActual: null,
      children: [
        {
          id: 'house',
          name: 'House',
          projected: 2000,
          linkedCells: [{ amount: 1800 }],
          manualActual: null,
          children: []
        },
        {
          id: 'car',
          name: 'Car',
          projected: 500,
          linkedCells: [{ amount: 450 }],
          manualActual: null,
          children: []
        }
      ]
    }
    
    const totals = calculateCategoryTotals(category)
    expect(totals.projected).toBe(2600) // 100 + 2000 + 500
    expect(totals.actual).toBe(2250)    // 0 + 1800 + 450
    expect(totals.difference).toBe(350) // 2600 - 2250
  })
})
```

---

## Next Steps

1. **Implement calculation functions** in shared utils file
2. **Add useMemo to category tree component** for performance
3. **Integrate with existing reducer** from 04-cell-linking.md
4. **Add manual override support** to CategoryNode interface
5. **Create visual indicators** for linked vs overridden cells
6. **Performance test** with large dataset (1000+ transactions)
7. **Add debouncing** if input lag observed

---

## References

- [React: You Might Not Need an Effect](https://react.dev/learn/you-might-not-need-an-effect)
- [React: useMemo Reference](https://react.dev/reference/react/useMemo)
- [React: Keeping Components Pure](https://react.dev/learn/keeping-components-pure)
- [React: Thinking in React](https://react.dev/learn/thinking-in-react)

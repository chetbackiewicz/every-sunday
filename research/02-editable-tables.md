# Editable Table/Grid Component Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: TanStack Table (React Table v8) ✅ SELECTED

---

## Recommendation Summary

**TanStack Table (React Table v8)** chosen over AG Grid because:
1. **Complete design control** - Headless architecture for custom budget UI
2. **Lightweight** - 10-15kb vs 500kb (better page load performance)
3. **Free and open source** - No licensing costs or restrictions
4. **Sufficient features** - Has everything needed (editing, selection, deletion)
5. **Custom cell linking** - Can build exact cell reference behavior needed
6. **Modern React patterns** - Hooks-based API aligns with React + TypeScript stack
7. **No vendor lock-in** - Full control over implementation

---

## TanStack Table (React Table v8) - Headless Table Library

**Source**: https://tanstack.com/table/latest, https://github.com/TanStack/table

### Key Features

- ✅ **Headless architecture** - 100% control over markup and styles
- ✅ **Lightweight** - 10-15kb with tree-shaking
- ✅ **Framework agnostic core** - React wrapper available
- ✅ **TypeScript first** - excellent type safety
- ✅ **Column definitions** - flexible column configuration with custom cell renderers
- ✅ **Row selection** - built-in `getIsSelected()`, `toggleSelectedHandler()`
- ✅ **Editable cells** - via custom `cell` property in column definitions
- ✅ **Pagination, sorting, filtering** - built-in features
- ✅ **Table state management** - controlled or uncontrolled
- ✅ **Virtual scrolling support** - works with `@tanstack/react-virtual`
- ✅ **Table meta pattern** - pass custom functions via `table.options.meta`

### Editable Cell Implementation Pattern

```typescript
// Define default column with editable cell
const defaultColumn: Partial<ColumnDef<Transaction>> = {
  cell: ({ getValue, row, column, table }) => {
    const initialValue = getValue()
    const [value, setValue] = useState(initialValue)
    
    // Update table data when user finishes editing
    const onBlur = () => {
      table.options.meta?.updateData(row.index, column.id, value)
    }
    
    return (
      <input
        value={value as string}
        onChange={e => setValue(e.target.value)}
        onBlur={onBlur}
      />
    )
  }
}

// Table setup with meta functions
const table = useReactTable({
  data,
  columns,
  defaultColumn,
  meta: {
    updateData: (rowIndex, columnId, value) => {
      setData(old => 
        old.map((row, index) => {
          if (index === rowIndex) {
            return { ...old[rowIndex], [columnId]: value }
          }
          return row
        })
      )
    }
  },
  getCoreRowModel: getCoreRowModel(),
})
```

### Row Deletion

- Simple array filter in state
- Example: `setData(old => old.filter((_, i) => i !== rowIndex))`
- Can trigger callbacks to update cell references

### Cell Selection (Custom Implementation)

- Track selected cells in separate state: `useState<Set<string>>(new Set())`
- Cell ID format: `${rowIndex}_${columnId}`
- Apply conditional CSS classes based on selection state
- Click handlers to toggle selection

### Pros

- Maximum flexibility - build exactly what's needed
- Small bundle size (important for web app performance)
- No opinionated styling - perfect for custom budget UI design
- Excellent documentation with 30+ examples
- Very active community (27k+ GitHub stars, 170k+ dependents)
- Can implement custom cell linking behavior easily
- Free and open source (MIT license)

### Cons

- More setup required vs opinionated solutions
- Must implement cell-level selection manually (row selection built-in only)
- No built-in editing UI components (must build custom)
- Requires understanding of table concepts (rows, columns, cells, models)

---

## AG Grid - Enterprise Data Grid (Not Selected)

**Source**: https://www.ag-grid.com/react-data-grid/, https://github.com/ag-grid/ag-grid

### Why Not Selected

- **Large bundle size** - Community edition ~500kb minified (vs 10-15kb for TanStack)
- **Opinionated styling** - requires CSS overrides for custom look
- **Enterprise features require paid license** - advanced features cost $999+/developer/year
- **Steeper learning curve** - many configuration options
- **Overkill for this use case** - paying for features not needed

### When AG Grid Would Be Better

- Need enterprise features (row grouping, pivoting, Excel export)
- Want out-of-the-box professional grid with minimal setup
- Budget allows for enterprise license
- Team has AG Grid expertise

---

## Feature Comparison

| Feature         | TanStack Table | AG Grid Community    | AG Grid Enterprise   |
| --------------- | -------------- | -------------------- | -------------------- |
| Bundle Size     | 10-15kb        | ~500kb               | ~500kb               |
| License         | MIT (Free)     | MIT (Free)           | Commercial (Paid)    |
| Cell Editing    | Custom         | Built-in             | Built-in + Advanced  |
| Row Deletion    | Custom         | Built-in             | Built-in             |
| Cell Selection  | Custom         | Built-in             | Built-in + Advanced  |
| Styling Control | 100% Custom    | Theme + CSS Override | Theme + CSS Override |
| Batch Editing   | Custom         | ❌                    | ✅                    |
| Validation      | Custom         | Basic                | Advanced             |
| Learning Curve  | Medium         | Medium-High          | High                 |
